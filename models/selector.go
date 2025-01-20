package models

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Selector struct {
	Type                      string   `json:"type" yaml:"type"`
	Tag                       string   `json:"tag" yaml:"tag"`
	Outbounds                 []string `json:"outbounds" yaml:"outbounds"`
	Default                   string   `json:"default,omitempty" yaml:"default,omitempty"`
	InterruptExistConnections bool     `json:"interrupt_exist_connections" yaml:"interrupt_exist_connections"`
}

func (s *Selector) Transform(message map[string]interface{}, logger *zap.Logger) (Outbound, error) {
	selectorContent, err := yaml.Marshal(message)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化yaml字符串失败")
	}
	var selector Selector
	if err := yaml.Unmarshal(selectorContent, &selector); err != nil {
		logger.Error(fmt.Sprintf("反序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化yaml字符串失败")
	}
	return &selector, nil
}

func (s *Selector) GetTag() string {
	return s.Tag
}