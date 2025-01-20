package singbox

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sifu-box/models"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)
func fetchFromRemote(provider models.Provider, client *http.Client, logger *zap.Logger) ([]models.Outbound, error) {
	req, err := http.NewRequest("GET", provider.Path, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
		return nil, fmt.Errorf("'%s'出错: 创建请求失败", provider.Name)
	}
	res, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("发送请求失败: [%s]", err.Error()))
		return nil, fmt.Errorf("'%s'出错: 发送请求失败", provider.Name)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("读取'%s'响应失败: [%s]",  provider.Name, err.Error()))
			return nil, fmt.Errorf("读取'%s'响应失败", provider.Name)
		}
		outbounds, err := parseFileContent(content, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("解析'%s'文件失败: [%s]", provider.Name, err.Error()))
			return nil, fmt.Errorf("'%s'出错: %s", provider.Name, err.Error())
		}
		return outbounds, nil
	}
	return nil, fmt.Errorf("'%s'未知响应, 状态码: %d", provider.Name, res.StatusCode)
}

func fetchFromLocal(provider models.Provider, logger *zap.Logger) ([]models.Outbound, error) {
	file, err := os.Open(provider.Path)
	if err != nil {
		logger.Error(fmt.Sprintf("打开'%s'文件失败: [%s]", provider.Name, err.Error()))
		return nil, fmt.Errorf("打开'%s'文件失败", provider.Name)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("读取'%s'文件失败: [%s]", provider.Name, err.Error()))
		return nil, fmt.Errorf("读取'%s'文件失败", provider.Name)
	}
	outbounds, err := parseFileContent(content, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("解析'%s'文件失败: [%s]", provider.Name, err.Error()))
		return nil, fmt.Errorf("'%s'出错: %s", provider.Name, err.Error())
	}
	return outbounds, nil
}

func parseFileContent(content []byte, logger *zap.Logger) ([]models.Outbound, error) {
	var providerInfo map[string]interface{}
	if err := yaml.Unmarshal(content, &providerInfo); err != nil {
		logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
		return nil, fmt.Errorf("解析响应失败")
	}
	var outbounds []models.Outbound
	proxies, ok := providerInfo["proxies"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("'proxies'字段丢失或不正确")
	}
	for _, proxy := range proxies {
		proxyMap, ok := proxy.(map[string]interface{})
		if !ok {
			logger.Error("该节点不是map类型")
			continue
		}
		protocol, ok := proxyMap["type"].(string)
		if !ok {
			logger.Error("该节点没有'type'字段")
			continue
		}
		name, ok := proxyMap["name"].(string)
		if !ok {
			logger.Error("该节点没有'name'字段")
			continue
		}
		switch protocol {
			case "ss":
				shadowSocks := models.ShadowSocks{}
				err := error(nil)
				var outbound models.Outbound = &shadowSocks
				outbound, err = outbound.Transform(proxyMap, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("'%s'节点解析ShadowSocks代理失败: [%s]", name, err.Error()))
					continue
				}
				outbounds = append(outbounds, outbound)
			case "vmess":
				vmess := models.VMess{}
				err := error(nil)
				var outbound models.Outbound = &vmess
				outbound, err = outbound.Transform(proxyMap, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("'%s'节点解析Vmess代理失败: [%s]", name, err.Error()))
					continue
				}
				outbounds = append(outbounds, outbound)
			case "trojan":
				trojan := models.Trojan{}
				err := error(nil)
				var outbound models.Outbound = &trojan
				outbound, err = outbound.Transform(proxyMap, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("'%s'节点解析Trojan代理失败: [%s]", name, err.Error()))
					continue
				}
				outbounds = append(outbounds, outbound)
		}
	}
	return outbounds, nil
}