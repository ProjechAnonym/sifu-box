package generate

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sifu-box/ent"
	"sifu-box/singbox"
	"sifu-box/utils"

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

// Generate 根据模板和出站映射生成配置文件
// 参数:
//
//	dir: 配置文件保存的目录路径
//	template: 配置模板对象, 包含出站组、路由、入站等配置信息
//	outbound_map: 出站配置映射, 用于扩展出站配置
//	logger: 日志记录器, 用于记录错误日志
//
// 返回值:
//
//	error: 配置生成或保存过程中出现的错误
func (c *Config) Generate(dir string, template *ent.Template, outbound_map map[string][]map[string]any, logger *zap.Logger) error {
	// 处理模板中的出站组配置
	outbounds := []map[string]any{}
	for _, outbound_group := range template.OutboundGroups {
		outbound, err := outbound_group.NewOutboundGroup(outbound_map)
		if err != nil {
			logger.Error(fmt.Sprintf(`出站组解析失败: [%s]`, err.Error()))
		}
		outbounds = append(outbounds, outbound)
	}

	// 合并额外的出站配置
	for _, v := range outbound_map {
		outbounds = append(outbounds, v...)
	}

	// 将处理后的配置赋值给当前配置对象
	c.Outbounds = outbounds
	c.Route = template.Route
	c.Inbounds = template.Inbounds
	c.Log = &template.Log
	c.Ntp = &template.Ntp
	c.Experiment = &template.Experiment
	c.DNS = template.DNS

	// 将配置对象序列化为JSON格式
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf(`配置"%s"生成失败: [%s]`, template.Name, err.Error())
	}

	// 将生成的配置保存到文件中, 文件名使用模板名称的MD5哈希值
	if err := utils.WriteFile(path.Join(dir, "config", fmt.Sprintf(`%s.json`, fmt.Sprintf("%x", md5.Sum([]byte(template.Name))))), content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return fmt.Errorf(`配置"%s"保存失败: [%s]`, template.Name, err.Error())
	}
	return nil
}
