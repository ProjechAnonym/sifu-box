package model

import (
	"fmt"
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
