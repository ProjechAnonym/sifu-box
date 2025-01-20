package models

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
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
		logger.Error(fmt.Sprintf("序列化json字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化json字符串失败")
	}
	var trojan Trojan
	if err := yaml.Unmarshal(trojanContent, &trojan); err != nil {
		logger.Error(fmt.Sprintf("反序列化json字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化json字符串失败")
	}
	network, ok := message["network"]
	if ok {
		switch network.(string) {
		case "ws":
			wsOptContent, err := yaml.Marshal(message["ws-opts"])
			if err != nil {
				logger.Error(fmt.Sprintf("'%s' 序列化ws-opts字段失败: [%s]", message["name"].(string), err.Error()))
				return nil, fmt.Errorf("序列化ws-opts字段失败")
			}
			var transport Transport
			if err := yaml.Unmarshal(wsOptContent, &transport); err != nil {
				logger.Error(fmt.Sprintf("'%s' 反序列化ws-opts字段失败: [%s]", message["name"].(string), err.Error()))
				return nil, fmt.Errorf("序列化ws-opts字段失败")
			}
			transport.Type = "ws"
			trojan.Transport = &transport
		}
	}
	if _, ok := message["skip-cert-verify"]; ok {
		trojan.TLS = &TLS{Enabled: true, Insecure: message["skip-cert-verify"].(bool)}
	}
	if _, ok := message["sni"]; ok {
		trojan.TLS.ServerName = message["sni"].(string)
	}
	if !trojan.UDP {
		trojan.Network = "tcp"
	}else{
		trojan.Network = ""
	}
	return &trojan, nil
}
func (t *Trojan) GetTag() string {
	return t.Tag
}