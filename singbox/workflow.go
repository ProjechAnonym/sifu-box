package singbox

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/template"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func GenerateConfigFiles(entClient *ent.Client, buntClient *buntdb.DB, specificProvider []string, specificTemplate []string, workDir string, server bool, rwLock *sync.RWMutex, logger *zap.Logger) []error {
	var err error
	var providers []models.Provider
	var rulesets []models.RuleSet
	templateMap := make(map[string]models.Template)
	if server {
		var providerList []*ent.Provider
		var templateList []*ent.Template
		if specificProvider != nil {
			providerList, err = entClient.Provider.Query().Where(provider.NameIn(specificProvider...)).All(context.Background())
		}else{
			providerList, err = entClient.Provider.Query().All(context.Background())
		}
		if err != nil {
			logger.Error(fmt.Sprintf("获取机场信息失败: [%s]", err.Error()))
			return []error{fmt.Errorf("获取机场信息失败")}
		}
		for _, provider := range providerList {
			providers = append(providers, models.Provider{
				Name: provider.Name,
				Path: provider.Path,
				Remote: provider.Remote,
				Detour: provider.Detour,
			})
		}
		
		rulesetsList, err := entClient.RuleSet.Query().All(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取路由规则集信息失败: [%s]", err.Error()))
			return []error{fmt.Errorf("获取路由规则集信息失败")}
		}
		for _, ruleset := range rulesetsList {
			rulesets = append(rulesets, models.RuleSet{
				Type: ruleset.Type,
				Tag: ruleset.Tag,
				Format: ruleset.Format,
				China: ruleset.China,
				NameServer: ruleset.NameServer,
				Label: ruleset.Label,
				Path: ruleset.Path,
				DownloadDetour: ruleset.DownloadDetour,
				UpdateInterval: ruleset.UpdateInterval,
			})
		}

		if specificTemplate != nil {
			templateList, err = entClient.Template.Query().Where(template.NameIn(specificTemplate...)).All(context.Background())
		}else{
			templateList, err = entClient.Template.Query().All(context.Background())
		}
		
		if err != nil {
			logger.Error(fmt.Sprintf("获取路由规则集信息失败: [%s]", err.Error()))
			return []error{fmt.Errorf("获取路由规则集信息失败")}
		}
		for _, template := range templateList {
			templateMap[template.Name] = template.Content
		}
	}else{
		configurationStr, err := utils.GetValue(buntClient, models.SINGBOXSETTINGKEY, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("获取配置信息失败: [%s]", err.Error()))
			return []error{fmt.Errorf("获取配置信息失败")}
		}
		var configuration models.Configuration
		if err := yaml.Unmarshal([]byte(configurationStr), &configuration); err != nil {
			logger.Error(fmt.Sprintf("解析配置信息失败: [%s]", err.Error()))
			return []error{fmt.Errorf("解析配置信息失败")}
		}
		providers = configuration.Providers
		rulesets = configuration.Rulesets
		templateMap = configuration.Templates
	}
	return merge(providers, rulesets, templateMap, workDir, server, rwLock, logger)
}

// TransferConfig 转移并应用配置文件
// 该函数负责从给定的路径备份当前配置，生成新的配置文件，并根据新配置重新加载服务。
// 如果重载服务失败，它将尝试恢复原始配置文件。
// 参数:
//   workDir string: 工作目录路径，用于存储临时文件和备份。
//   singboxSetting models.SingboxEnv: singbox环境设置，包含配置路径、命令等信息。
//   logger *zap.Logger: 日志记录器，用于记录日志信息。
// 返回值:
//   error: 如果过程中发生任何错误，返回该错误。
func ApplyNewConfig(workDir string, singboxSetting models.Singbox, buntClient *buntdb.DB, logger *zap.Logger) error {
    providerName, err := utils.GetValue(buntClient, models.CURRENTPROVIDER, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置机场失败: [%s]", err.Error()))
		return fmt.Errorf("获取当前配置机场失败")
	}
	templateName, err := utils.GetValue(buntClient, models.CURRENTTEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置模板失败: [%s]", err.Error()))
		return fmt.Errorf("获取当前配置模板失败")
	}
	// 备份当前配置文件，以防新配置应用过程中发生错误
    if err := backupConfig(singboxSetting.ConfigPath, workDir, logger); err != nil {
        return err
    }
    
    // 计算机场名称的MD5哈希值，用于生成唯一的配置文件名
    providerHashName, err := utils.EncryptionMd5(providerName)
    if err != nil {
        logger.Error(fmt.Sprintf("计算机场名称哈希值失败: [%s]", err.Error()))
        return fmt.Errorf("计算机场名称哈希值失败")
    }
    
    // 转移新的配置文件到指定路径
    if err := transferConfig(singboxSetting.ConfigPath, filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, templateName, fmt.Sprintf("%s.json", providerHashName)), logger); err != nil {
        return err
    }
    
    // 检查服务状态，如果服务不在运行或检查失败，则尝试启动服务，否则重载服务
    status, err := checkService(false, logger, singboxSetting.Commands[models.CHECKCOMMAND])
    if err != nil || !status {
        if err := bootService(logger, singboxSetting.Commands[models.BOOTCOMMAND]); err != nil {
            return err
        }
    }else{
        if err := reloadService(logger, singboxSetting.Commands[models.RELOADCOMMAND]); err != nil {
            return err
        }        
    }
    
    // 再次检查服务状态，确认服务是否成功应用了新配置
    status, err = checkService(true, logger, singboxSetting.Commands[models.CHECKCOMMAND])
    if status && err == nil {
		logger.Debug(fmt.Sprintf("重载'%s'基于'%s'模板的配置文件成功", providerName, templateName))
        return nil
    }else if err != nil{
        logger.Error(fmt.Sprintf("重载'%s'基于'%s'模板的配置文件失败: [%s]", providerName, templateName, err.Error()))
    }else{
        err = errors.New("未知错误")
        logger.Error(fmt.Sprintf("重载'%s'基于'%s'模板的配置文件失败: [%s]", providerName, templateName, err.Error()))
    }
    
    // 如果服务重载失败，尝试恢复原始配置文件
    if err := recoverConfig(singboxSetting.ConfigPath, workDir, logger); err != nil {
        return fmt.Errorf("恢复原始配置文件失败")
    }
   
    return nil
}