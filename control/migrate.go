package control

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/model"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func Migrate(ent_client *ent.Client, bunt_client *buntdb.DB, logger *zap.Logger) ([]byte, error) {
	providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询机场数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询机场数据失败: [%s]`, err.Error())
	}
	provider_list := []struct {
		model.Provider
		UUID      string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
		Nodes     []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
		Templates []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}
	ruleset_list := []struct {
		model.Ruleset
		Templates []string `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}
	template_list := []model.Template{}
	for _, provider := range providers {
		provider_list = append(provider_list, struct {
			model.Provider
			UUID      string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes     []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		}{
			Provider: model.Provider{
				Name:   provider.Name,
				Path:   provider.Path,
				Remote: provider.Remote},
			Nodes:     provider.Nodes,
			UUID:      provider.UUID,
			Templates: provider.Templates,
		})
	}
	rulesets, err := ent_client.Ruleset.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询规则集数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询规则集数据失败: [%s]`, err.Error())
	}
	for _, ruleset := range rulesets {
		ruleset_list = append(ruleset_list, struct {
			model.Ruleset
			Templates []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		}{
			Templates: ruleset.Templates,
			Ruleset: model.Ruleset{
				Name:           ruleset.Name,
				Path:           ruleset.Path,
				Remote:         ruleset.Remote,
				Binary:         ruleset.Binary,
				UpdateInterval: ruleset.UpdateInterval,
				DownloadDetour: ruleset.DownloadDetour,
			},
		})
	}
	templates, err := ent_client.Template.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询模板数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询模板数据失败: [%s]`, err.Error())
	}
	for _, template := range templates {
		template_list = append(template_list, model.Template{
			Name:           template.Name,
			Ntp:            &template.Ntp,
			Inbounds:       template.Inbounds,
			Providers:      template.Providers,
			Route:          &template.Route,
			OutboundsGroup: template.OutboundGroups,
			DNS:            &template.DNS,
			Experiment:     &template.Experiment,
			Log:            &template.Log,
		})
	}
	content, err := yaml.Marshal(struct {
		Providers []struct {
			model.Provider
			UUID      string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes     []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"providers,omitempty" yaml:"providers,omitempty"`
		Rulesets []struct {
			model.Ruleset
			Templates []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
		Templates []model.Template `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{Templates: template_list, Providers: provider_list, Rulesets: ruleset_list})
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf(`反序列化配置信息失败: [%s]`, err.Error())
	}
	return content, nil
}
