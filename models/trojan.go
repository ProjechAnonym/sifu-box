package models

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Trojan struct {
	Type       string     `json:"type" yaml:"type"`
	Tag        string     `json:"tag" yaml:"name"`
	Server     string     `json:"server" yaml:"server"`
	ServerPort int        `json:"server_port" yaml:"port"`
	Network    string     `json:"network,omitempty" yaml:"network,omitempty"`
	Password   string     `json:"password" yaml:"password"`
	UDP        bool       `json:"-" yaml:"udp,omitempty"`
	Multiplex  *Multiplex `json:"multiplex,omitempty" yaml:"multiplex,omitempty"`
	Transport  *Transport `json:"transport,omitempty" yaml:"transport,omitempty"`
	TLS        *TLS       `json:"tls,omitempty" yaml:"tls,omitempty"`
	Dial       `json:",inline" yaml:",inline"`
}

func (t *Trojan) Transform(message map[string]interface{}, logger *zap.Logger) (Outbound, error) {
	trojanContent, err := yaml.Marshal(message)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化yaml字符串失败")
	}
	var trojan Trojan
	if err := yaml.Unmarshal(trojanContent, &trojan); err != nil {
		logger.Error(fmt.Sprintf("反序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化yaml字符串失败")
	}
	network, ok := message["network"].(string)
	if ok {
		switch network {
		case "ws":
			wsOptContent, err := yaml.Marshal(message["ws-opts"])
			if err != nil {
				logger.Error(fmt.Sprintf("序列化ws-opts字段失败: [%s]", err.Error()))
				return nil, fmt.Errorf("序列化ws-opts字段失败")
			}
			var transport Transport
			if err := yaml.Unmarshal(wsOptContent, &transport); err != nil {
				logger.Error(fmt.Sprintf("反序列化ws-opts字段失败: [%s]", err.Error()))
				return nil, fmt.Errorf("序列化ws-opts字段失败")
			}
			transport.Type = "ws"
			trojan.Transport = &transport
		}
	}
	if insecure, ok := message["skip-cert-verify"].(bool); ok {
		trojan.TLS = &TLS{Enabled: true, Insecure: insecure}
	}
	if sni, ok := message["sni"].(string); ok {
		trojan.TLS.ServerName = sni
	}
	trojan.Network = ""
	return &trojan, nil
}
func (t *Trojan) GetTag() string {
	return t.Tag
}