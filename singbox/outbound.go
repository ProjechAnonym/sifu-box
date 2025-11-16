package singbox

import "fmt"

type OutboundGroup struct {
	Type       string   `json:"type" yaml:"type"`
	Tag        string   `json:"tag" yaml:"tag"`
	Providers  []string `json:"providers,omitempty" yaml:"providers,omitempty"`
	Tag_Groups []string `json:"tag_groups,omitempty" yaml:"tag_groups,omitempty"`
}

// NewOutboundGroup 根据提供的出站配置创建一个新的出站组
// 参数:
//   - outbounds: 包含所有可用出站节点的映射,键为提供者名称,值为节点配置列表
//
// 返回值:
//   - map[string]any: 构建的出站组配置
//   - error: 如果构建过程中出现错误则返回相应的错误信息
func (o *OutboundGroup) NewOutboundGroup(outbounds map[string][]map[string]any) (map[string]any, error) {
	outbound := map[string]any{}
	outbound["type"] = o.Type
	outbound["tag"] = o.Tag
	// 如果类型为direct, 直接返回基础配置
	if o.Type == "direct" {
		return outbound, nil
	}
	outbound["interrupt_exist_connections"] = false
	outbound_tags := o.Tag_Groups
	// 遍历所有机场, 收集出站节点标签
	for _, provider := range o.Providers {
		if _, exists := outbounds[provider]; !exists {
			return nil, fmt.Errorf(`该配置不使用"%s"节点`, provider)
		}
		for _, node := range outbounds[provider] {
			if _, ok := node["tag"]; !ok {
				continue
			}
			outbound_tags = append(outbound_tags, node["tag"].(string))
		}

	}
	outbound["outbounds"] = outbound_tags
	return outbound, nil
}
