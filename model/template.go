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
	Inbounds       []map[string]any        `json:"inbounds,omitempty" yaml:"inbounds,omitempty"`
	OutboundsGroup []singbox.OutboundGroup `json:"outbounds_group,omitempty" yaml:"outbounds_group,omitempty"`
	DNS            *singbox.DNS            `json:"dns,omitempty" yaml:"dns,omitempty"`
	Experiment     *singbox.Experiment     `json:"experiment,omitempty" yaml:"experiment,omitempty"`
	Log            *singbox.Log            `json:"log,omitempty" yaml:"log,omitempty"`
	Route          *singbox.Route          `json:"route,omitempty" yaml:"route,omitempty"`
	Providers      []string                `json:"providers,omitempty" yaml:"providers,omitempty"`
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

// LinkProvidersTable 将模板与提供商进行关联, 更新提供商的模板列表
//
// 参数:
//   - ent_client: ent客户端, 用于数据库操作
//
// 返回值:
//   - error: 操作过程中发生的错误, 如果操作成功则返回nil
func (t *Template) LinkProvidersTable(ent_client *ent.Client) error {
	// 遍历模板关联的所有提供商
	for _, name := range t.Providers {
		// 查询提供商信息, 获取其当前的模板列表
		provider_msg, err := ent_client.Provider.Query().Where(provider.NameEQ(name)).Select(provider.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询机场"%s"失败: %s`, t.Name, name, err.Error())
		}

		// 将当前模板添加到提供商的模板列表中, 并去重
		provider_msg.Templates = append(provider_msg.Templates, t.Name)
		template_map := make(map[string]bool)
		template_list := []string{}

		// 使用map进行去重操作
		for _, v := range provider_msg.Templates {
			template_map[v] = true
		}

		// 将去重后的模板列表重新构建为切片
		for k := range template_map {
			template_list = append(template_list, k)
		}

		// 更新提供商的模板列表
		if _, err := ent_client.Provider.UpdateOne(provider_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新机场"%s"失败: %s`, t.Name, name, err.Error())
		}
	}
	return nil
}

// LinkRulesetsTable 将模板与规则集进行关联, 将当前模板名称添加到对应规则集的模板列表中
// 参数:
//   - ent_client: ent数据库客户端, 用于执行数据库查询和更新操作
//
// 返回值:
//   - error: 操作过程中出现的错误信息, 如果操作成功则返回nil
func (t *Template) LinkRulesetsTable(ent_client *ent.Client) error {
	// 遍历模板中的所有规则集
	for _, rule_set := range t.Route.Rule_sets {
		// 查询指定名称的规则集, 获取其模板列表
		ruleset_msg, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(rule_set.Tag)).Select(ruleset.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}

		// 将当前模板名称添加到规则集的模板列表中, 并去重
		ruleset_msg.Templates = append(ruleset_msg.Templates, t.Name)
		template_list := []string{}
		template_map := make(map[string]bool)
		// 同一个模板只能作为map的一个键值, 以此实现去重
		for _, v := range ruleset_msg.Templates {
			template_map[v] = true
		}
		// 将去重后的模板名称转换为列表
		for k := range template_map {
			template_list = append(template_list, k)
		}

		// 更新规则集的模板列表
		if _, err := ent_client.Ruleset.UpdateOne(ruleset_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}
	}
	return nil
}

// UnLinkProvidersTable 从模板关联的机场中移除该模板的引用
// 该函数会遍历模板关联的所有机场, 从每个机场的模板列表中移除当前模板的名称
//
// 参数:
//
//	ent_client - ent客户端, 用于数据库操作
//
// 返回值:
//
//	error - 如果在查询或更新过程中发生错误则返回错误信息, 否则返回nil
func (t *Template) UnLinkProvidersTable(ent_client *ent.Client) error {
	// 遍历模板关联的所有机场名称
	for _, name := range t.Providers {
		// 查询指定名称的机场信息, 只选择Templates字段
		provider_msg, err := ent_client.Provider.Query().Where(provider.NameEQ(name)).Select(provider.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询机场"%s"失败: %s`, t.Name, name, err.Error())
		}

		// 构建新的模板列表, 排除当前模板名称
		template_list := []string{}
		for _, v := range provider_msg.Templates {
			if v != t.Name {
				template_list = append(template_list, v)
			}
		}

		// 更新机场的模板列表
		if _, err := ent_client.Provider.UpdateOne(provider_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新机场"%s"失败: %s`, t.Name, name, err.Error())
		}
	}
	return nil
}

// UnLinkRulesetsTable 从规则集中移除当前模板的关联关系
//
// 该函数会遍历模板的所有路由规则集, 查询对应的规则集记录,
// 然后从规则集的模板列表中移除当前模板名称, 实现解绑操作
//
// 参数:
//
//	ent_client - ent数据库客户端, 用于执行数据库查询和更新操作
//
// 返回值:
//
//	error - 操作过程中发生的错误, 包括查询失败或更新失败的情况
func (t *Template) UnLinkRulesetsTable(ent_client *ent.Client) error {
	// 遍历模板的所有路由规则集
	for _, rule_set := range t.Route.Rule_sets {
		// 查询规则集记录, 只获取Templates字段
		ruleset_msg, err := ent_client.Ruleset.Query().Where(ruleset.NameEQ(rule_set.Tag)).Select(ruleset.FieldTemplates).First(context.Background())
		if err != nil {
			return fmt.Errorf(`模板"%s"查询规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}

		// 构建新的模板列表, 排除当前模板
		template_list := []string{}
		for _, v := range ruleset_msg.Templates {
			if v != t.Name {
				template_list = append(template_list, v)
			}
		}

		// 更新规则集的模板列表
		if _, err := ent_client.Ruleset.UpdateOne(ruleset_msg).SetTemplates(template_list).Save(context.Background()); err != nil {
			return fmt.Errorf(`模板"%s"更新规则集"%s"模板失败: %s`, t.Name, rule_set.Tag, err.Error())
		}
	}
	return nil
}

// EditProviders 处理模板中的提供商配置, 过滤掉不存在的提供商
// 参数: 无
// 返回值:
//
//	error - 处理过程中出现的错误, 如果providers字段为空则返回错误信息
func (t *Template) EditProviders() error {
	if t.Providers == nil {
		return fmt.Errorf(`模板"%s"中"providers"字段为空, 机场不能为空`, t.Name)
	}

	// 创建机场映射表, 用于快速查找存在的机场
	providers_map := make(map[string]bool)
	outbound_group_list := []singbox.OutboundGroup{}
	for _, name := range t.Providers {
		providers_map[name] = true
	}

	// 遍历所有出站组, 过滤掉不在机场列表中的机场
	for _, outbound_group := range t.OutboundsGroup {
		providers_list := []string{}
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

// EditRulesets 用于根据模板中定义的规则集（Rule_sets）来过滤路由和 DNS 规则中的规则集标签
// 该函数会遍历路由和 DNS 的规则列表, 移除那些引用了未定义规则集标签的规则,
// 并更新剩余规则中的规则集标签为已定义的有效标签
//
// 参数：
//
//	无显式参数, 作用于 Template 结构体实例 t
//
// 返回值：
//
//	无返回值修改将直接应用于 t.Route.Rules 和 t.DNS.Rules
func (t *Template) EditRulesets() {
	// 初始化用于存储处理后路由规则和 DNS 规则的切片
	route_rules := []map[string]any{}
	dns_rules := []map[string]any{}

	// 构建一个规则集标签的映射, 用于快速判断标签是否存在
	ruleset_map := make(map[string]bool)
	for _, rule_set := range t.Route.Rule_sets {
		ruleset_map[rule_set.Tag] = true
	}

	// 处理路由规则：过滤并保留引用了有效规则集标签的规则
	for _, rule := range t.Route.Rules {
		filter_tags := []string{}
		ruleset_tags, ok := rule[RULE_SET]
		if !ok {
			// 如果规则中没有规则集字段, 则直接保留该规则
			route_rules = append(route_rules, rule)
			continue
		}
		// 遍历规则中的规则集标签, 筛选出有效的标签
		for _, ruleset_tag := range ruleset_tags.([]any) {
			if _, ok := ruleset_tag.(string); ok {
				if _, exists := ruleset_map[ruleset_tag.(string)]; exists {
					filter_tags = append(filter_tags, ruleset_tag.(string))
				}
			}
		}
		// 如果没有有效的规则集标签, 则丢弃该规则
		if len(filter_tags) == 0 {
			continue
		}
		// 更新规则中的规则集标签为筛选后的有效标签, 并添加到结果中
		rule[RULE_SET] = filter_tags
		route_rules = append(route_rules, rule)
	}

	// 处理 DNS 规则：逻辑与处理路由规则相同
	for _, rule := range t.DNS.Rules {
		filter_tags := []string{}
		ruleset_tags, ok := rule[RULE_SET]
		if !ok {
			dns_rules = append(dns_rules, rule)
			continue
		}
		for _, ruleset_tag := range ruleset_tags.([]any) {
			if _, ok := ruleset_tag.(string); ok {
				if _, exists := ruleset_map[ruleset_tag.(string)]; exists {
					filter_tags = append(filter_tags, ruleset_tag.(string))
				}
			}
		}
		if len(filter_tags) == 0 {
			continue
		}
		rule[RULE_SET] = filter_tags
		dns_rules = append(dns_rules, rule)
	}

	// 将处理后的规则列表赋回模板结构体
	t.Route.Rules = route_rules
	t.DNS.Rules = dns_rules
}
