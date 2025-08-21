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

// fetchFromRemote 从远程URL获取数据并解析为出站节点配置
// name: 配置名称, 用于错误信息标识
// url: 远程配置文件的URL地址
// client: HTTP客户端, 用于发送请求
// logger: 日志记录器, 用于记录处理过程中的日志
// 返回值: 解析后的出站节点配置列表和可能的错误信息
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

	// 处理HTTP响应
	if res.StatusCode == 200 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(`读取"%s"响应失败: [%s]`, name, err.Error())
		}

		config := map[string]any{}
		// 尝试按YAML格式解析响应内容
		if err := yaml.Unmarshal(content, &config); err == nil {
			outbounds, err := generateFromYaml(config, logger)
			if err != nil {
				return nil, fmt.Errorf(`生成"%s"出站节点失败: [%s]`, name, err.Error())
			}
			return outbounds, nil
		}
		// 按Base64格式解析响应内容
		outbounds, err := generateFromBase64(content, logger)
		if err != nil {
			return nil, err
		}
		return outbounds, nil
	}

	return nil, fmt.Errorf(`"%s"未知响应, 状态码: %d`, name, res.StatusCode)
}

// fetchFromLocal 从本地YAML文件读取配置并生成出站节点信息
// name: 配置名称, 用于错误提示
// path: 本地文件路径
// logger: 日志记录器
// 返回值: 出站节点信息列表和错误信息
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

	// 根据YAML配置生成出站节点信息
	outbounds, err := generateFromYaml(config, logger)
	if err != nil {
		return nil, fmt.Errorf(`生成"%s"出站节点失败: [%s]`, name, err.Error())
	}

	// 返回解析出的出站信息列表
	return outbounds, nil
}

// generateFromYaml 从YAML配置内容中解析出站代理配置
// 参数:
//   - content: YAML配置文件解析后的map结构, 包含代理配置信息
//   - logger: 用于记录日志的zap.Logger实例
//
// 返回值:
//   - []map[string]interface{}: 解析后的出站代理配置切片, 每个元素代表一个代理配置
//   - error: 解析过程中出现的错误, 如果解析成功则返回nil
func generateFromYaml(content map[string]any, logger *zap.Logger) ([]map[string]interface{}, error) {

	var outbounds []map[string]any
	// 从配置中提取代理列表
	proxies, ok := content[PROXIES_FIELD].([]any)
	if !ok {
		return nil, fmt.Errorf(`"%s"字段丢失或不正确`, PROXIES_FIELD)
	}

	// 遍历所有代理配置并根据协议类型进行解析
	for _, proxy := range proxies {

		// 提取代理名称
		name, ok := proxy.(map[string]any)[NAME].(string)
		if !ok {
			logger.Error(fmt.Sprintf(`获取节点信息出错, 没有找到"%s"字段`, NAME))
			continue
		}
		// 提取代理协议类型
		protocol, ok := proxy.(map[string]any)[TYPE].(string)
		if !ok {
			logger.Error(fmt.Sprintf(`获取节点信息出错, 没有找到"%s"字段`, TYPE))
			continue
		}
		// 验证代理配置是否为有效的map结构
		if _, ok := proxy.(map[string]any); !ok {
			logger.Error(fmt.Sprintf(`无法解析"%s"所含字段`, name))
			continue
		}
		// 根据协议类型调用对应的解析函数
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

// generateFromBase64 从Base64编码的内容中解析出站代理配置
// content: Base64编码的代理配置内容, 每行一个代理链接
// logger: 用于记录错误日志的zap日志器
// 返回解析后的出站代理配置切片和可能的错误
func generateFromBase64(content []byte, logger *zap.Logger) ([]map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return nil, err
	}
	var outbounds []map[string]any

	// 遍历解码后的每一行, 解析不同协议的代理链接
	for line := range strings.SplitSeq(string(data), "\n") {
		link, err := url.Parse(line)
		if err != nil {
			logger.Error(fmt.Sprintf(`"%s"解析URL失败: [%s]`, line, err.Error()))
			continue
		}
		switch link.Scheme {
		case "ss":
			outbounds = append(outbounds, shadowsocksFromBase64(link))
		case "vmess":
			outbounds = append(outbounds, vmessFromBase64(link))
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
