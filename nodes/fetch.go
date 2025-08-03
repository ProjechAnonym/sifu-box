package nodes

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
		if err := yaml.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf(`"%s"响应内容不是有效的YAML格式: [%s]`, name, err.Error())
		}
		// 解析响应内容
		outbounds, err := generateFromYaml(config, logger)
		if err != nil {
			return nil, fmt.Errorf(`生成"%s"出站节点失败: [%s]`, name, err.Error())
		}

		// 返回解析后的出站信息
		return outbounds, nil
	}

	// 如果响应状态码不是200, 返回自定义错误信息
	return nil, fmt.Errorf(`"%s"未知响应, 状态码: %d`, name, res.StatusCode)
}

func fetchFromLocal(name, path string, logger *zap.Logger) ([]map[string]any, error) {
	// 尝试打开本地文件
	file, err := os.Open(path)
	if err != nil {
		// 如果打开文件失败, 记录错误并返回自定义错误信息
		logger.Error(fmt.Sprintf(`打开"%s"文件失败: [%s]`, name, err.Error()))
		return nil, fmt.Errorf(`打开"%s"文件失败: [%s]`, name, err.Error())
	}
	// 确保在函数返回前关闭文件
	defer file.Close()

	// 读取文件的全部内容
	content, err := io.ReadAll(file)
	if err != nil {
		// 如果读取文件失败, 记录错误并返回自定义错误信息
		logger.Error(fmt.Sprintf(`读取"%s"文件失败: [%s]`, name, err.Error()))
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
			logger.Error(fmt.Sprintf(`获取字段出错, 没有找到"%s"字段`, NAME))
			continue
		}
		protocol, ok := proxy.(map[string]any)[TYPE].(string)
		if !ok {
			logger.Error(fmt.Sprintf(`获取字段出错, 没有找到"%s"字段`, TYPE))
			continue
		}
		if _, ok := proxy.(map[string]any); !ok {
			logger.Error(fmt.Sprintf(`无法解析"%s"所含字段`, name))
			continue
		}
		switch protocol {
		case "ss":
			outbounds = append(outbounds, shadowsocksFromYaml(proxy.(map[string]any)))
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
