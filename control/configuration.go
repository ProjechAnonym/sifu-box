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

func Set(name, workDir string, singboxSetting models.Singbox, templateContent models.Template, buntClient *buntdb.DB, entClient *ent.Client, rwLock *sync.RWMutex, logger *zap.Logger) []error{
	exist, err := entClient.Template.Query().Where(template.NameEQ(name)).Exist(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		return []error{fmt.Errorf("数据库查询'%s'数据失败", name)}
	}
	if exist {
		if _, err := entClient.Template.Update().Where(template.NameEQ(name)).SetContent(templateContent).Save(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("更新数据库数据失败: [%s]",err.Error()))
			return []error{fmt.Errorf("数据库更新'%s'数据失败", name)}
		}
	}else{
		if _, err := entClient.Template.Create().SetName(name).SetContent(templateContent).Save(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("保存数据库数据失败: [%s]",err.Error()))
			return []error{fmt.Errorf("数据库保存'%s'数据失败", name)}
		}
	}
	
	if errors := singbox.GenerateConfigFiles(entClient, buntClient, nil, []string{name}, workDir, true, rwLock, logger); errors != nil {
		return errors
	}
	currentTemplate, err := utils.GetValue(buntClient, models.CURRENTTEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置模板失败: [%s]", err.Error()))
		return []error{fmt.Errorf("获取当前配置模板失败")}
	}
	fmt.Println(currentTemplate)
	if currentTemplate == name {
		
		if err := singbox.ApplyNewConfig(workDir, singboxSetting, buntClient, rwLock, logger); err != nil{
			return []error{err}
		}
	}
	return nil
}
func Add(providers []models.Provider, rulesets []models.RuleSet, entClient *ent.Client, buntClient *buntdb.DB, workDir string, rwLock *sync.RWMutex, logger *zap.Logger) []error {
	providerList := make([]*ent.ProviderCreate, len(providers))
	newProviders := make([]string, len(providers))
	if providers != nil {
		for i, provider := range providers {
			newProviders[i] = provider.Name
			providerList[i] = entClient.Provider.Create().SetDetour(provider.Detour).SetName(provider.Name).SetPath(provider.Path).SetRemote(provider.Remote)
		}
		if err := entClient.Provider.CreateBulk(providerList...).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			return []error{fmt.Errorf("保存机场数据失败")}
		}
	}

	rulesetList := make([]*ent.RuleSetCreate, len(rulesets))
	if rulesets != nil {

		for i, ruleset := range rulesets {
			rulesetList[i] = entClient.RuleSet.Create().SetChina(ruleset.China).SetDownloadDetour(ruleset.DownloadDetour).SetFormat(ruleset.Format).SetLabel(ruleset.Label).SetNameServer(ruleset.NameServer).SetPath(ruleset.Path).SetTag(ruleset.Tag).SetType(ruleset.Type).SetUpdateInterval(ruleset.UpdateInterval)
		}
		if err := entClient.RuleSet.CreateBulk(rulesetList...).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			return []error{fmt.Errorf("保存规则集数据失败")}
		}
		newProviders = nil
	}


	errors := singbox.GenerateConfigFiles(entClient, buntClient, newProviders, nil, workDir, true, rwLock, logger)
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
	for _, providerName := range providers {
		if providerName == currentProvider {
			if err := utils.DeleteValue(buntClient, models.CURRENTPROVIDER, logger); err != nil {
				logger.Error(fmt.Sprintf("删除当前机场配置失败: [%s]", err.Error()))
				errors = append(errors, fmt.Errorf("删除当前机场配置失败"))
				return errors
			}
		}
		providerMsg, err := entClient.Provider.Query().Select(provider.FieldRemote, provider.FieldName, provider.FieldPath).Where(provider.NameEQ(providerName)).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库'%s'数据失败: [%s]", providerName, err.Error()))
			errors = append(errors, fmt.Errorf("获取数据库'%s'数据失败", providerName))
			continue
		}
		if !providerMsg.Remote {
			if err := os.Remove(providerMsg.Path); err != nil {
				logger.Error(fmt.Sprintf("删除'%s'机场Yaml文件失败: [%s]", providerMsg.Name, err.Error()))
				errors = append(errors, fmt.Errorf("删除'%s'机场Yaml文件失败", providerMsg.Name))
				continue
			}
		}
		dirs, err := os.ReadDir(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR))
		if err != nil {
			logger.Error(fmt.Sprintf("遍历配置文件夹失败: [%s]", err.Error()))
			errors = append(errors, fmt.Errorf("遍历配置文件夹失败"))
			return errors
		}
		providerHashName, err := utils.EncryptionMd5(providerName)
		if err != nil {
			logger.Error(fmt.Sprintf("计算'%s'哈希值失败: [%s]", providerName, err.Error()))
			errors = append(errors, fmt.Errorf("计算'%s'哈希值失败", providerName))
			continue
		}
		for _, dir := range dirs {
			if !dir.IsDir() {
				logger.Error(fmt.Sprintf("配置文件夹下的模板'%s'不是文件夹", dir.Name()))
				errors = append(errors, fmt.Errorf("配置文件夹下的模板'%s'不是文件夹", dir.Name()))
			}
			
			if err := os.RemoveAll(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, dir.Name(), fmt.Sprintf("%s.json", providerHashName))); err != nil {
				logger.Error(fmt.Sprintf("删除机场'%s'失败: [%s]", providerName, err.Error()))
				errors = append(errors, fmt.Errorf("删除机场'%s'失败", providerName))
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