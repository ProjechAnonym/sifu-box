package control

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/ent/template"
	"sifu-box/models"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func Fetch(entClient *ent.Client, logger *zap.Logger) (*models.Configuration, error) {
	providers, err := entClient.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取机场数据失败: [%s]",err.Error()))
		return nil, fmt.Errorf("获取机场数据失败")
	}
	providerList := make([]models.Provider, len(providers))
	for i, provider := range providers {
		providerList[i] = models.Provider{
			Name: provider.Name,
			Path: provider.Path,
			Remote: provider.Remote,
			Detour: provider.Detour,
		}
	}
	rulesets, err := entClient.RuleSet.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取规则集数据失败: [%s]",err.Error()))
		return nil, fmt.Errorf("获取规则集数据失败")
	}
	rulesetList := make([]models.RuleSet, len(rulesets))
	for i, ruleset := range rulesets {
		rulesetList[i] = models.RuleSet{
			Type: ruleset.Type,
			Tag: ruleset.Tag,
			Format: ruleset.Format,
			China: ruleset.China,
			NameServer: ruleset.NameServer,
			Label: ruleset.Label,
			Path: ruleset.Path,
			DownloadDetour: ruleset.DownloadDetour,
			UpdateInterval: ruleset.UpdateInterval,
		}
	}
	templates, err := entClient.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取模板信息失败")
	}
	templateList := make(map[string]models.Template)
	for _, template := range templates {
		templateList[template.Name] = template.Content
	}
	return &models.Configuration{
		Providers: providerList,
		Rulesets: rulesetList,
		Templates: templateList,
	}, nil
}

func Delete(providers, rulesets, templates []string, workDir string, buntClient *buntdb.DB, entClient *ent.Client, rwLock *sync.RWMutex, logger *zap.Logger) []error{
	var errors []error
	currentProvider, err := utils.GetValue(buntClient, models.CURRENTPROVIDER, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置机场失败: [%s]", err.Error()))
	}
	errors = append(errors, deleteProviders(providers, currentProvider, workDir, buntClient, entClient, rwLock, logger)...)
	
	currentTemplate, err := utils.GetValue(buntClient, models.CURRENTTEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置模板失败: [%s]", err.Error()))
	}
	errors = append(errors, deleteTemplate(entClient, buntClient, templates, currentTemplate, workDir, rwLock, logger)...)
	if rulesets != nil {
		errors = append(errors, deleteRulesets(rulesets, workDir, rwLock, entClient, buntClient, logger)...)
	}
	return errors
}

func deleteRulesets(rulesets []string, workDir string, rwLock *sync.RWMutex, entClient *ent.Client, buntClient *buntdb.DB, logger *zap.Logger) []error {
	var errors []error
	if _, err := entClient.RuleSet.Delete().Where(ruleset.TagIn(rulesets...)).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("连接数据库失败: [%s]", err.Error()))
		errors = append(errors, fmt.Errorf("从数据库删除指定规则集失败"))
	}
	errors = append(errors, singbox.GenerateConfigFiles(entClient, buntClient, nil, nil, workDir, true, rwLock, logger)...)
	return errors
}
func deleteProviders(providers []string, currentProvider, workDir string, buntClient *buntdb.DB, entClient *ent.Client, rwLock *sync.RWMutex, logger *zap.Logger) []error{
	rwLock.Lock()
	defer rwLock.Unlock()
	var errors []error
	for _, provider := range providers {
		if provider == currentProvider {
			if err := utils.DeleteValue(buntClient, models.CURRENTPROVIDER, logger); err != nil {
				logger.Error(fmt.Sprintf("删除当前机场配置失败: [%s]", err.Error()))
				errors = append(errors, fmt.Errorf("删除当前机场配置失败"))
				return errors
			}
		}
		dirs, err := os.ReadDir(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR))
		if err != nil {
			logger.Error(fmt.Sprintf("遍历配置文件夹失败: [%s]", err.Error()))
			errors = append(errors, fmt.Errorf("遍历配置文件夹失败"))
			return errors
		}
		providerHashName, err := utils.EncryptionMd5(provider)
		if err != nil {
			logger.Error(fmt.Sprintf("计算'%s'哈希值失败: [%s]", provider, err.Error()))
			errors = append(errors, fmt.Errorf("计算'%s'哈希值失败", provider))
			continue
		}
		for _, dir := range dirs {
			if !dir.IsDir() {
				logger.Error(fmt.Sprintf("配置文件夹下的模板'%s'不是文件夹", dir.Name()))
				errors = append(errors, fmt.Errorf("配置文件夹下的模板'%s'不是文件夹", dir.Name()))
			}
			
			if err := os.RemoveAll(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, dir.Name(), fmt.Sprintf("%s.json", providerHashName))); err != nil {
				logger.Error(fmt.Sprintf("删除机场'%s'失败: [%s]", provider, err.Error()))
				errors = append(errors, fmt.Errorf("删除机场'%s'失败", provider))
			}
		}
	}
	if _, err := entClient.Provider.Delete().Where(provider.NameIn(providers...)).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("连接数据库失败: [%s]", err.Error()))
		errors = append(errors, fmt.Errorf("从数据库删除指定机场失败"))
	}
	return errors
}

func deleteTemplate(entClient *ent.Client, buntClient *buntdb.DB, templates []string, currentTemplate, workDir string, rwLock *sync.RWMutex, logger *zap.Logger) []error {
	rwLock.Lock()
	defer rwLock.Unlock()
	var errors []error
	
	for _, template := range templates {
		if template == currentTemplate {
			if err := utils.DeleteValue(buntClient, models.CURRENTTEMPLATE, logger); err != nil {
				logger.Error(fmt.Sprintf("删除当前配置模板失败: [%s]", err.Error()))
				errors = append(errors, fmt.Errorf("删除当前配置模板失败"))
				return errors
			}
		}
		if err := os.RemoveAll(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, template)); err != nil {
			logger.Error(fmt.Sprintf("删除模板'%s'失败: [%s]", template, err.Error()))
			errors = append(errors, fmt.Errorf("删除模板'%s'失败", template))
		}
	}
	
	if _, err := entClient.Template.Delete().Where(template.NameIn(templates...)).Exec(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("连接数据库失败: [%s]", err.Error()))
		errors = append(errors, fmt.Errorf("从数据库删除指定模板失败"))
	}
	return errors
}