package singbox

import "fmt"

type OutboundGroup struct {
	Type      string   `json:"type" yaml:"type"`
	Tag       string   `json:"tag" yaml:"tag"`
	Providers []string `json:"providers" yaml:"providers"`
}

func (o *OutboundGroup) NewOutboundGroup(outbounds map[string][]string) (map[string]any, error) {
	outbound := map[string]any{}
	outbound["type"] = o.Type
	outbound["tag"] = o.Tag
	outbound["interrupt_exist_connections"] = false
	outbound_tags := []string{}
	for _, provider := range o.Providers {
		if _, exists := outbounds[provider]; !exists {
			return nil, fmt.Errorf(`该配置不使用"%s"节点`, provider)
		}
		outbound_tags = append(outbound_tags, outbounds[provider]...)

	}
	outbound["outbounds"] = outbound_tags
	return outbound, nil
}
