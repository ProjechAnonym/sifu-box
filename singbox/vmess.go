package singbox

import (
	"fmt"
	"sifu-box/models"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func marshVmess(vmessMap map[string]interface{}, logger *zap.Logger) (*models.VMess, error) {
	vmessContent, err := yaml.Marshal(vmessMap)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化json字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化json字符串失败")
	}
	var vmess models.VMess
	if err := yaml.Unmarshal(vmessContent, &vmess); err != nil {
		logger.Error(fmt.Sprintf("反序列化json字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化json字符串失败")
	}
	vmess.Network = ""
	network, ok := vmessMap["network"]
	if ok {
		switch network.(string) {
		case "ws":
			wsOptContent, err := yaml.Marshal(vmessMap["ws-opts"])
			if err != nil {
				logger.Error(fmt.Sprintf("'%s' 序列化ws-opts字段失败: [%s]", vmessMap["name"].(string), err.Error()))
				return nil, fmt.Errorf("序列化ws-opts字段失败")
			}
			var transport models.Transport
			if err := yaml.Unmarshal(wsOptContent, &transport); err != nil {
				logger.Error(fmt.Sprintf("'%s' 反序列化ws-opts字段失败: [%s]", vmessMap["name"].(string), err.Error()))
				return nil, fmt.Errorf("序列化ws-opts字段失败")
			}
			transport.Type = "ws"
			vmess.Transport = &transport
		}
	}
	return &vmess, nil
}