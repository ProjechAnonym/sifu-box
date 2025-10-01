package control

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func Export(ent_client *ent.Client, bunt_client *buntdb.DB, logger *zap.Logger) ([]byte, error) {
	providers, err := ent_client.Provider.Query().All(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf(`查询机场数据失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`查询机场数据失败: [%s]`, err.Error())
	}
	provider_list := []struct {
		model.Provider `json:",inline" yaml:",inline"`
		UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
		Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
		Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}
	ruleset_list := []struct {
		model.Ruleset `json:",inline" yaml:",inline"`
		Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}
	template_list := []model.Template{}
	for _, provider := range providers {
		provider_list = append(provider_list, struct {
			model.Provider `json:",inline" yaml:",inline"`
			UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
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
			model.Ruleset `json:",inline" yaml:",inline"`
			Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
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
			model.Provider `json:",inline" yaml:",inline"`
			UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"providers,omitempty" yaml:"providers,omitempty"`
		Rulesets []struct {
			model.Ruleset `json:",inline" yaml:",inline"`
			Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
		Templates []model.Template `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{Templates: template_list, Providers: provider_list, Rulesets: ruleset_list})
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf(`反序列化配置信息失败: [%s]`, err.Error())
	}
	return content, nil
}
func Import(content []byte, ent_client *ent.Client, logger *zap.Logger) ([]gin.H, error) {
	setting := struct {
		Providers []struct {
			model.Provider `json:",inline" yaml:",inline"`
			UUID           string           `json:"uuid,omitempty" yaml:"uuid,omitempty"`
			Nodes          []map[string]any `json:"nodes,omitempty" yaml:"nodes,omitempty"`
			Templates      []string         `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"providers,omitempty" yaml:"providers,omitempty"`
		Rulesets []struct {
			model.Ruleset `json:",inline" yaml:",inline"`
			Templates     []string `json:"templates,omitempty" yaml:"templates,omitempty"`
		} `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
		Templates []model.Template `json:"templates,omitempty" yaml:"templates,omitempty"`
	}{}
	if err := yaml.Unmarshal(content, &setting); err != nil {
		logger.Error(fmt.Sprintf(`序列化配置文件失败: [%s]`, err.Error()))
		return nil, fmt.Errorf(`序列化配置文件失败: [%s]`, err.Error())
	}
	res := []gin.H{}
	for _, provider := range setting.Providers {
		if err := ent_client.Provider.Create().SetName(provider.Name).SetPath(provider.Path).SetRemote(provider.Remote).SetUUID(provider.UUID).SetNodes(provider.Nodes).SetTemplates(provider.Templates).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加机场"%s"失败: [%s]`, provider.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加机场"%s"失败: [%s]`, provider.Name, err.Error())})
			continue
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加机场"%s"成功`, provider.Name)})
	}
	for _, ruleset := range setting.Rulesets {
		if err := ent_client.Ruleset.Create().SetName(ruleset.Name).SetPath(ruleset.Path).SetRemote(ruleset.Remote).SetBinary(ruleset.Binary).SetUpdateInterval(ruleset.UpdateInterval).SetDownloadDetour(ruleset.DownloadDetour).SetTemplates(ruleset.Templates).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加规则集%s失败: [%s]`, ruleset.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加规则集"%s"失败: [%s]`, ruleset.Name, err.Error())})
			continue
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加规则集"%s"成功`, ruleset.Name)})
	}
	for _, template := range setting.Templates {
		template_instance := ent_client.Template.Create()
		template.CreateFillFields(template_instance)
		if err := template_instance.SetInbounds(template.Inbounds).SetName(template.Name).SetOutboundGroups(template.OutboundsGroup).SetProviders(template.Providers).SetUpdated(true).Exec(context.Background()); err != nil {
			logger.Error(fmt.Sprintf(`添加模板"%s"失败: [%s]`, template.Name, err.Error()))
			res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`添加模板"%s"失败: [%s]`, template.Name, err.Error())})
			continue
		}
		res = append(res, gin.H{"status": true, "message": fmt.Sprintf(`添加模板"%s"成功`, template.Name)})

	}
	return res, nil
}
