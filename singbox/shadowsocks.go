package singbox

import (
	"fmt"
	"sifu-box/models"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func marshShadowSocks(shadowSocksMap map[string]interface{}, logger *zap.Logger) (*models.ShadowSocks, error){
	shadowSocksContent, err := yaml.Marshal(shadowSocksMap)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化json字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化json字符串失败")
	}
	var shadowSocks models.ShadowSocks
	if err := yaml.Unmarshal(shadowSocksContent, &shadowSocks); err != nil {
		logger.Error(fmt.Sprintf("反序列化json字符串失败: [%s]", err.Error()))
		return nil, fmt.Errorf("反序列化json字符串失败")
	}
	shadowSocks.Type = "shadowsocks"
	return &shadowSocks, nil
}