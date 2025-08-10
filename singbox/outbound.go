package singbox

import "fmt"

type OutboundGroup struct {
	Type      string   `json:"type" yaml:"type"`
	Tag       string   `json:"tag" yaml:"tag"`
	Providers []string `json:"providers" yaml:"providers"`
}

func (o *OutboundGroup) NewOutboundGroup(outbounds map[string][]map[string]any) (map[string]any, error) {
	outbound := map[string]any{}
	outbound["type"] = o.Type
	outbound["tag"] = o.Tag
	if o.Type == "direct" {
		return outbound, nil
	}
	outbound["interrupt_exist_connections"] = false
	outbound_tags := []string{}
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
