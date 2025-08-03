package nodes

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func fetchFromRemote(name, url string, client *http.Client, logger *zap.Logger) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(`"%s"出错, 创建请求失败: [%s]`, name, err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(`"%s"出错,发送请求失败: [%s]`, name, err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(`读取"%s"响应失败: [%s]`, name, err.Error())
		}

		config := map[string]any{}
		if err := yaml.Unmarshal(content, &config); err == nil {
			outbounds, err := generateFromYaml(config, logger)
			if err != nil {
				return nil, fmt.Errorf(`生成"%s"出站节点失败: [%s]`, name, err.Error())
			}
			return outbounds, nil
		}
		outbounds, err := generateFromBase64(content, logger)
		if err != nil {
			return nil, err
		}
		return outbounds, nil
	}

	return nil, fmt.Errorf(`"%s"未知响应, 状态码: %d`, name, res.StatusCode)
}

func fetchFromLocal(name, path string, logger *zap.Logger) ([]map[string]any, error) {

	file, err := os.Open(path)
	if err != nil {

		return nil, fmt.Errorf(`打开"%s"文件失败: [%s]`, name, err.Error())
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {

		return nil, fmt.Errorf(`读取"%s"文件失败: [%s]`, name, err.Error())
	}

	config := map[string]any{}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf(`"%s"响应内容不是有效的YAML格式: [%s]`, name, err.Error())
	}
	outbounds, err := generateFromYaml(config, logger)
	if err != nil {
		return nil, fmt.Errorf(`生成"%s"出站节点失败: [%s]`, name, err.Error())
	}

	// 返回解析出的出站信息列表
	return outbounds, nil
}

func generateFromYaml(content map[string]any, logger *zap.Logger) ([]map[string]interface{}, error) {

	var outbounds []map[string]any
	proxies, ok := content[PROXIES_FIELD].([]any)
	if !ok {
		return nil, fmt.Errorf(`"%s"字段丢失或不正确`, PROXIES_FIELD)
	}

	for _, proxy := range proxies {

		name, ok := proxy.(map[string]any)[NAME].(string)
		if !ok {
			logger.Error(fmt.Sprintf(`获取节点信息出错, 没有找到"%s"字段`, NAME))
			continue
		}
		protocol, ok := proxy.(map[string]any)[TYPE].(string)
		if !ok {
			logger.Error(fmt.Sprintf(`获取节点信息出错, 没有找到"%s"字段`, TYPE))
			continue
		}
		if _, ok := proxy.(map[string]any); !ok {
			logger.Error(fmt.Sprintf(`无法解析"%s"所含字段`, name))
			continue
		}
		switch protocol {
		case "ss":
			outbounds = append(outbounds, shadowsocksFromYaml(proxy.(map[string]any)))
		case "vless":
			outbounds = append(outbounds, vlessFromYaml(proxy.(map[string]any)))
		case "vmess":
			outbounds = append(outbounds, vmessFromYaml(proxy.(map[string]any)))
		case "trojan":
			outbounds = append(outbounds, trojanFromYaml(proxy.(map[string]any)))
		default:
			logger.Error(fmt.Sprintf(`"%s"协议暂不支持`, protocol))
			continue
		}
	}
	// 返回解析后的出站代理配置切片
	return outbounds, nil
}
func generateFromBase64(content []byte, logger *zap.Logger) ([]map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return nil, err
	}
	var outbounds []map[string]any
	for line := range strings.SplitSeq(string(data), "\n") {
		link, err := url.Parse(line)
		if err != nil {
			logger.Error(fmt.Sprintf(`"%s"解析URL失败: [%s]`, line, err.Error()))
			continue
		}
		switch link.Scheme {
		case "ss":
			outbound, err := shadowsocksFromBase64(link)
			if err != nil {
				logger.Error(fmt.Sprintf(`"%s"协议解析失败: [%s]`, link.Scheme, err.Error()))
				continue
			}
			outbounds = append(outbounds, outbound)
		case "vmess":
			outbound, err := vmessFromBase64(link)
			if err != nil {
				logger.Error(fmt.Sprintf(`"%s"协议解析失败: [%s]`, link.Scheme, err.Error()))
				continue
			}
			outbounds = append(outbounds, outbound)
		case "trojan":
			outbounds = append(outbounds, trojanFromBase64(link))
		case "vless":
			outbounds = append(outbounds, vlessFromBase64(link))
		default:
			if line != "" {
				logger.Error(fmt.Sprintf(`"%s"协议暂不支持`, strings.Split(line, ":")[0]))
			}
		}
	}

	// 返回解析后的出站代理配置切片
	return outbounds, nil
}
