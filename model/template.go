package model

import (
	"context"
	"fmt"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/singbox"
)

type Template struct {
	Name           string                  `json:"name" yaml:"name"`
	Ntp            *singbox.Ntp            `json:"ntp,omitempty" yaml:"ntp,omitempty"`
	Inbounds       []map[string]any        `json:"inbounds" yaml:"inbounds"`
	OutboundsGroup []singbox.OutboundGroup `json:"outbounds_group,omitempty" yaml:"outbounds_group,omitempty"`
	DNS            *singbox.DNS            `json:"dns" yaml:"dns"`
	Experiment     *singbox.Experiment     `json:"experiment,omitempty" yaml:"experiment,omitempty"`
	Log            *singbox.Log            `json:"log,omitempty" yaml:"log,omitempty"`
	Route          *singbox.Route          `json:"route" yaml:"route"`
	Providers      []string                `json:"providers" yaml:"providers"`
}

func (t *Template) CheckField() error {
	if t.Name == "" {
		return fmt.Errorf(`"name"字段为空, 名称不能为空`)
	}
	if t.Inbounds == nil {
		return fmt.Errorf(`"inbounds"字段为空, 入站配置不能为空`)
	}
	if t.OutboundsGroup == nil {
		return fmt.Errorf(`"outbounds_group"字段为空, 出站组不能为空`)
	}
	if t.DNS == nil {
		return fmt.Errorf(`"dns"字段为空, DNS配置不能为空`)
	}
	if t.Route == nil {
		return fmt.Errorf(`"route"字段为空, 路由配置不能为空`)
	}
	if t.Experiment == nil {
		return fmt.Errorf(`"experiment"字段为空, 实验功能配置不能为空`)
	}
	if t.Providers == nil {
		return fmt.Errorf(`"providers"字段为空, 机场不能为空`)
	}
	return nil
}
func (t *Template) CreateFillFields(template *ent.TemplateCreate) {
	if t.Log != nil {
		template.SetLog(*t.Log)
	}
	if t.Ntp != nil {
		template.SetNtp(*t.Ntp)
	}
	if t.Experiment != nil {
		template.SetExperiment(*t.Experiment)
	}
	if t.DNS != nil {
		template.SetDNS(*t.DNS)
	}
	if t.Route != nil {
		template.SetRoute(*t.Route)
	}
}
func (t *Template) UpdateFillFields(template *ent.TemplateUpdate) {
	if t.Log != nil {
		template.SetLog(*t.Log)
	}
	if t.Ntp != nil {
		template.SetNtp(*t.Ntp)
	}
	if t.Experiment != nil {
		template.SetExperiment(*t.Experiment)
	}
	if t.DNS != nil {
		template.SetDNS(*t.DNS)
	}
	if t.Route != nil {
		template.SetRoute(*t.Route)
	}
}
func (t *Template) LinkProvidersTable(ent_client *ent.Client) error {
	for _, name := range t.Providers {
		provider_msg, err := ent_client.Provider.Query().Where(provider.NameEQ(name)).Select(provider.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询机场"%s"失败: %s`, t.Name, name, err.Error())
		}
		provider_msg.Templates = append(provider_msg.Templates, t.Name)
		template_map := make(map[string]bool)
		template_list := []string{}
		for _, v := range provider_msg.Templates {
			template_map[v] = true
		}
		for k := range template_map {
			template_list = append(template_list, k)
		}
		if _, err := ent_client.Provider.UpdateOne(provider_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新机场"%s"失败: %s`, t.Name, name, err.Error())
		}
	}
	return nil
}
func (t *Template) LinkRulesetsTable(ent_client *ent.Client) error {
	for _, rule_set := range t.Route.Rule_sets {
		ruleset_msg, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(rule_set.Tag)).Select(ruleset.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}
		ruleset_msg.Templates = append(ruleset_msg.Templates, t.Name)
		template_list := []string{}
		template_map := make(map[string]bool)
		for _, v := range ruleset_msg.Templates {
			template_map[v] = true
		}
		for k := range template_map {
			template_list = append(template_list, k)
		}
		if _, err := ent_client.Ruleset.UpdateOne(ruleset_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}
	}
	return nil
}
func (t *Template) UnLinkProvidersTable(ent_client *ent.Client) error {
	for _, name := range t.Providers {
		provider_msg, err := ent_client.Provider.Query().Where(provider.NameEQ(name)).Select(provider.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询机场"%s"失败: %s`, t.Name, name, err.Error())
		}
		template_list := []string{}
		for _, v := range provider_msg.Templates {
			if v != t.Name {
				template_list = append(template_list, v)
			}
		}
		if _, err := ent_client.Provider.UpdateOne(provider_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新机场"%s"失败: %s`, t.Name, name, err.Error())
		}
	}
	return nil
}
func (t *Template) UnLinkRulesetsTable(ent_client *ent.Client) error {
	for _, rule_set := range t.Route.Rule_sets {
		ruleset_msg, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(rule_set.Tag)).Select(ruleset.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}
		template_list := []string{}
		for _, v := range ruleset_msg.Templates {
			if v != t.Name {
				template_list = append(template_list, v)
			}
		}
		if _, err := ent_client.Ruleset.UpdateOne(ruleset_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}
	}
	return nil
}
func (t *Template) EditProviders(ent_client *ent.Client) error {
	if t.Providers == nil {
		return fmt.Errorf(`模板"%s"中"providers"字段为空, 机场不能为空`, t.Name)
	}
	outbound_group_list := []singbox.OutboundGroup{}
	for _, outbound_group := range t.OutboundsGroup {
		providers_map := make(map[string]bool)
		providers_list := []string{}
		for _, name := range t.Providers {
			providers_map[name] = true
		}
		for _, provider := range outbound_group.Providers {
			if _, exists := providers_map[provider]; exists {
				providers_list = append(providers_list, provider)
			}
		}
		outbound_group.Providers = providers_list
		outbound_group_list = append(outbound_group_list, outbound_group)
	}
	t.OutboundsGroup = outbound_group_list
	return nil
}
