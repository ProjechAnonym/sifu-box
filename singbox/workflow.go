package singbox

import (
	"context"
	"encoding/json"
	"fmt"
	"sifu-box/ent"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func Workflow(entClient *ent.Client, buntClient *buntdb.DB, server bool, logger *zap.Logger) ([]string, error) {
	settingStr, err := utils.GetValue(buntClient, models.SINGBOXSETTINGKEY, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取配置信息失败")
	}
	var setting models.SingboxSetting
	if err := json.Unmarshal([]byte(settingStr), &setting); err != nil {
		logger.Error(fmt.Sprintf("解析配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("解析配置信息失败")
	}
	var providers []models.Provider
	var rulesets []models.RuleSet
	templateMap := make(map[string]models.Template)
	if server {
		providerList, err := entClient.Provider.Query().All(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取机场信息失败: [%s]", err.Error()))
			return nil, fmt.Errorf("获取机场信息失败")
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
			return nil, fmt.Errorf("获取路由规则集信息失败")
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

		templates, err := entClient.Template.Query().All(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取路由规则集信息失败: [%s]", err.Error()))
			return nil, fmt.Errorf("获取路由规则集信息失败")
		}
		for _, template := range templates {
			templateMap[template.Name] = template.Content
		}
	}else{
		providers = setting.Providers
		rulesets = setting.Rulesets
		templateMap = setting.Templates
	}
	merge(providers, rulesets, templateMap, logger)
	// templateMap := setting.Templates

	

	return nil, nil
}