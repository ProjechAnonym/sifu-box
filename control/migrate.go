package control

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/ent/template"
	"sifu-box/models"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func Export(entClient *ent.Client, buntClient *buntdb.DB, logger *zap.Logger) ([]byte, error) {
	providers, err := entClient.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		return nil, fmt.Errorf("获取机场数据失败")
	}
	providerList := make([]models.Provider, len(providers))
	for i, provider := range providers {
		providerList[i] = models.Provider{
			Name: provider.Name,
			Detour: provider.Detour,
			Path: provider.Path,
			Remote: provider.Remote,
		}	
	}

	rulesets, err := entClient.RuleSet.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		return nil, fmt.Errorf("获取规则集数据失败")
	}
	rulesetList := make([]models.RuleSet, len(rulesets))
	for i, ruleset := range rulesets {
		rulesetList[i] = models.RuleSet{
			Tag: ruleset.Tag,
			NameServer: ruleset.NameServer,
			Path: ruleset.Path,
			Type: ruleset.Type,
			Format: ruleset.Format,
			China: ruleset.China,
			Label: ruleset.Label,
			DownloadDetour: ruleset.DownloadDetour,
			UpdateInterval: ruleset.UpdateInterval,
		}
	}

	templates, err := entClient.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取模板信息失败")
	}

	templateMap := make(map[string]models.Template)
	for _, template := range templates {
		templateMap[template.Name] = template.Content
	}
	

	currentProvider, err := utils.GetValue(buntClient, models.CURRENTPROVIDER, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置机场失败: [%s]", err.Error()))
		
	}
	currentTemplate, err := utils.GetValue(buntClient, models.CURRENTTEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置模板失败: [%s]", err.Error()))
		
	}
	currentInterval, err := utils.GetValue(buntClient, models.INTERVAL, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前配置定时任务失败: [%s]", err.Error()))
		
	}
	conf := struct {
		Configuration models.Configuration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
		CurrentApplication map[string]string `json:"current_application,omitempty" yaml:"current_application,omitempty"`
	}{
		Configuration: models.Configuration{Providers: providerList, Rulesets: rulesetList, Templates: templateMap},
		CurrentApplication: map[string]string{models.CURRENTPROVIDER: currentProvider, models.CURRENTTEMPLATE: currentTemplate, models.INTERVAL: currentInterval},
	}
	content, err := yaml.Marshal(conf)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化配置文件失败: [%s]",err.Error()))
		return nil, fmt.Errorf("序列化配置文件失败")
	}
	return content, nil
}

func Import(content []byte, workDir string, singboxSetting models.Singbox, entClient *ent.Client, buntClient *buntdb.DB, scheduler *cron.Cron, jobID *cron.EntryID, execLock *sync.Mutex, rwLock *sync.RWMutex, logger *zap.Logger) error {
	conf := struct {
		Configuration models.Configuration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
		CurrentApplication map[string]string `json:"current_application,omitempty" yaml:"current_application,omitempty"`
	}{}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		logger.Error(fmt.Sprintf("解析配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("解析配置文件失败")
	}
	for key, value := range conf.CurrentApplication {
		if err := utils.SetValue(buntClient, key, value, logger); err != nil {
			logger.Error(fmt.Sprintf("写入配置信息失败: [%s]", err.Error()))
			return fmt.Errorf("写入配置信息失败")
		}
		if key == models.INTERVAL && value != "" {
			scheduler.Remove(*jobID)
			var err error
			*jobID, err = scheduler.AddFunc(value, func(){
				singbox.GenerateConfigFiles(entClient, buntClient, nil, nil, workDir, true, rwLock, logger)
				singbox.ApplyNewConfig(workDir, singboxSetting, buntClient, rwLock, execLock, logger)
			})
			if err != nil {
				logger.Error(fmt.Sprintf("设置定时任务失败: [%s]", err.Error()))
				return err
			}
		}
	}
	for _, supplier := range conf.Configuration.Providers {
		exist, err := entClient.Provider.Query().Where(provider.NameEQ(supplier.Name)).Exist(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		}
		if !exist {
			if _, err := entClient.Provider.Create().SetName(supplier.Name).SetDetour(supplier.Detour).SetPath(supplier.Path).SetRemote(supplier.Remote).Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			}
		}	
	}

	for _, collectionInfo := range conf.Configuration.Rulesets {
		exist, err := entClient.RuleSet.Query().Where(ruleset.TagEQ(collectionInfo.Tag)).Exist(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		}
		if !exist {
			if _, err := entClient.RuleSet.Create().SetTag(collectionInfo.Tag).
													SetNameServer(collectionInfo.NameServer).
													SetPath(collectionInfo.Path).
													SetType(collectionInfo.Type).
													SetFormat(collectionInfo.Format).
													SetChina(collectionInfo.China).
													SetLabel(collectionInfo.Label).
													SetDownloadDetour(collectionInfo.DownloadDetour).
													SetUpdateInterval(collectionInfo.UpdateInterval).
													Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			}
		}
	}
	for key, templateContent := range conf.Configuration.Templates {
		exist, err := entClient.Template.Query().Where(template.NameEQ(key)).Exist(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
		}
		if !exist {
			if _, err := entClient.Template.Create().
											SetName(key).
											SetContent(templateContent).
											Save(context.Background()); err != nil {
				logger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
			}
		}
	}
	return nil
}