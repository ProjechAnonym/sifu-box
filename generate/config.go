package generate

import (
	"encoding/json"
	"fmt"
	"sifu-box/ent"
	"sifu-box/singbox"

	"go.uber.org/zap"
)

type Config struct {
	Experiment *singbox.Experiment `json:"experiment,omitempty" yaml:"experiment,omitempty"`
	Ntp        *singbox.Ntp        `json:"ntp,omitempty" yaml:"ntp,omitempty"`
	Log        *singbox.Log        `json:"log,omitempty" yaml:"log,omitempty"`
	DNS        singbox.DNS         `json:"dns" yaml:"dns"`
	Inbounds   []map[string]any    `json:"inbounds" yaml:"inbounds"`
	Outbounds  []map[string]any    `json:"outbounds" yaml:"outbounds"`
	Route      singbox.Route       `json:"route" yaml:"route"`
}

func (c *Config) Generate(template *ent.Template, outbound_map map[string][]map[string]any, logger *zap.Logger) (map[string]any, error) {
	outbounds := []map[string]any{}
	for _, outbound_group := range template.OutboundGroups {
		outbound, err := outbound_group.NewOutboundGroup(outbound_map)
		if err != nil {
			logger.Error(fmt.Sprintf(`出站组解析失败: [%s]`, err.Error()))
		}
		outbounds = append(outbounds, outbound)
	}
	for _, v := range outbound_map {
		outbounds = append(outbounds, v...)
	}
	c.Outbounds = outbounds
	c.Route = template.Route
	c.Inbounds = template.Inbounds
	c.Log = &template.Log
	c.Ntp = &template.Ntp
	c.Experiment = &template.Experiment
	c.DNS = template.DNS
	content, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println(string(content))
	return nil, nil
}
