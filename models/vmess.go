package models

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type VMess struct {
	Transport           *Transport             `json:"transport,omitempty" yaml:"transport,omitempty"`
	Multiplex           *Multiplex             `json:"multiplex,omitempty" yaml:"multiplex,omitempty"`
	UDP                 bool                   `json:"-" yaml:"udp,omitempty"`
	TLS                 *TLS                   `json:"tls,omitempty" yaml:"tls,omitempty"`
	Type                string                 `json:"type" yaml:"type"`
	Tag                 string                 `json:"tag" yaml:"name"`
	Server              string                 `json:"server" yaml:"server"`
	ServerPort          int                    `json:"server_port" yaml:"port"`
	Network             string                 `json:"network,omitempty" yaml:"network,omitempty"`
	UUID                string                 `json:"uuid" yaml:"uuid"`
	Security            string                 `json:"security,omitempty" yaml:"cipher,omitempty"`
	AlterID             int                    `json:"alter_id" yaml:"alterId"`
	GlobalPadding       bool                   `json:"global_padding,omitempty" yaml:"global_padding,omitempty"`
	AuthenticatedLength bool                   `json:"authenticated_length,omitempty" yaml:"authenticated_length,omitempty"`
	PacketEncoding      string                 `json:"packet_encoding,omitempty" yaml:"packet_encoding,omitempty"`
}

func (v *VMess) Transform(message map[string]interface{}, logger *zap.Logger) (Outbound, error) {
	vmessContent, err := yaml.Marshal(message)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化yaml字符串失败")
	}
	var vmess VMess
	if err := yaml.Unmarshal(vmessContent, &vmess); err != nil {
		logger.Error(fmt.Sprintf("反序列化yaml字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化yaml字符串失败")
	}
	
	network, ok := message["network"]
	if ok {
		switch network.(string) {
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
			vmess.Transport = &transport
		}
	}
	vmess.Network = ""
	return &vmess, nil
}
func (v *VMess) GetTag() string {
	return v.Tag
}