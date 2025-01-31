package control

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func Export(entClient *ent.Client, buntClient *buntdb.DB, application *models.Application, logger *zap.Logger) ([]byte, error) {
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
	
	application.Server.SSL = nil

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
		Application models.Application `json:"application,omitempty" yaml:"application,omitempty"`
		Configuration models.Configuration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
		CurrentApplication map[string]string `json:"current_application,omitempty" yaml:"current_application,omitempty"`
	}{
		Application: *application,
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