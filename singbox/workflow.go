package singbox

import (
	"context"
	"encoding/json"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/template"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func Workflow(entClient *ent.Client, buntClient *buntdb.DB, specificProvider []string, specificTemplate []string, workDir string, server bool, logger *zap.Logger) []error {
	settingStr, err := utils.GetValue(buntClient, models.SINGBOXSETTINGKEY, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取配置信息失败: [%s]", err.Error()))
		return []error{fmt.Errorf("获取配置信息失败")}
	}
	var setting models.SingboxSetting
	if err := json.Unmarshal([]byte(settingStr), &setting); err != nil {
		logger.Error(fmt.Sprintf("解析配置信息失败: [%s]", err.Error()))
		return []error{fmt.Errorf("解析配置信息失败")}
	}
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
		providers = setting.Providers
		rulesets = setting.Rulesets
		templateMap = setting.Templates
	}
	return merge(providers, rulesets, templateMap, workDir, server, logger)
}

func TransferConfig(entClient *ent.Client, buntClient *buntdb.DB, workDir string, singboxSetting models.SingboxEnv, logger *zap.Logger) error {
	// if err := backupConfig(singboxSetting.ConfigPath, workDir, logger); err != nil {
	// 	return err
	// }
	// providerHashName, err := utils.EncryptionMd5(singboxSetting.Provider)
	// if err != nil {
	// 	logger.Error(fmt.Sprintf("计算机场名称哈希值失败: [%s]", err.Error()))
	// 	return fmt.Errorf("计算机场名称哈希值失败")
	// }
	
	// if err := transferConfig(singboxSetting.ConfigPath, filepath.Join(workDir, models.TEMPDIR, singboxSetting.Template, providerHashName), workDir, logger); err != nil {
	// 	return err
	// }
	// if err := reloadService(logger, singboxSetting.Command[models.RELOADCOMMAND]); err != nil {
	// 	return err
	// }
	status, err := checkService(true, logger, singboxSetting.Command[models.CHECKCOMMAND])
	if status && err == nil {
		return nil
	}else if err != nil{
		logger.Error(fmt.Sprintf("重载'%s'基于'%s'模板的配置文件失败: [%s]", singboxSetting.Provider, singboxSetting.Template, err.Error()))
	}else{
		logger.Error(fmt.Sprintf("重载'%s'基于'%s'模板的配置文件失败: [%s]", singboxSetting.Provider, singboxSetting.Template, fmt.Errorf("未知错误")))
	}
	// if err := recoverConfig(singboxSetting.ConfigPath, workDir, logger); err != nil {
	// 	return err
	// }

	return nil
}