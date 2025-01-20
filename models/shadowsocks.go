package models

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ShadowSocks struct {
	Type       string `json:"type" yaml:"type"`
	Tag        string `json:"tag" yaml:"name"`
	Server     string `json:"server" yaml:"server"`
	ServerPort int    `json:"server_port" yaml:"port"`
	UDP        bool   `json:"-" yaml:"udp,omitempty"`
	Network    string `json:"network" yaml:"network"`
	Method     string `json:"method" yaml:"cipher"`
	Password   string `json:"password" yaml:"password"`
	Plugin     string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	PluginOpts string `json:"plugin_opts,omitempty" yaml:"plugin_opts,omitempty"`
	Dial       `json:",inline" yaml:",inline"`
}

func (s *ShadowSocks) Transform(message map[string]interface{}, logger *zap.Logger) (Outbound, error) {
	shadowSocksContent, err := yaml.Marshal(message)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化yaml字符串失败")
	}
	var shadowSocks ShadowSocks
	if err := yaml.Unmarshal(shadowSocksContent, &shadowSocks); err != nil {
		logger.Error(fmt.Sprintf("反序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化yaml字符串失败")
	}
	shadowSocks.Type = "shadowsocks"
	shadowSocks.Network = ""
	return &shadowSocks, nil
}

func (s *ShadowSocks) GetTag() string {
	return s.Tag
}