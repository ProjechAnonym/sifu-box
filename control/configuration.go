package control

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/models"

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