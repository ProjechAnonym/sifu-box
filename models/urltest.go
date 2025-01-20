package models

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type URLTest struct {
	Type                      string   `json:"type" yaml:"type"`
	Tag                       string   `json:"tag" yaml:"tag"`
	Outbounds                 []string `json:"outbounds" yaml:"outbounds"`
	URL                       string   `json:"url,omitempty" yaml:"url,omitempty"`
	Interval                  string   `json:"interval,omitempty" yaml:"interval,omitempty"`
	Tolerance                 int      `json:"tolerance,omitempty" yaml:"tolerance,omitempty"`
	IdleTimeout               string   `json:"idle_timeout,omitempty" yaml:"idle_timeout,omitempty"`
	InterruptExistConnections bool     `json:"interrupt_exist_connections" yaml:"interrupt_exist_connections"`
}

func (u *URLTest) Transform(message map[string]interface{}, logger *zap.Logger) (Outbound, error){
	urlTestContent, err := yaml.Marshal(message)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化yaml字符串失败")
	}
	var urlTest URLTest
	if err := yaml.Unmarshal(urlTestContent, &urlTest); err != nil {
		logger.Error(fmt.Sprintf("反序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化yaml字符串失败")
	}
	return &urlTest, nil
}
func (u *URLTest)  GetTag () string {
	return u.Tag
}