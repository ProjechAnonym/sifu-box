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

// fetchFromRemote 从远程获取数据
// 该函数根据提供的信息, 使用指定的HTTP客户端从远程路径获取数据, 并解析为Outbound模型切片
// 参数:
//   provider: 包含远程路径和其他必要信息的提供者模型
//   client: 用于发送HTTP请求的客户端
//   logger: 用于记录日志的实例
// 返回值:
//   []models.Outbound: 解析后的出站信息切片
//   error: 如果在执行过程中遇到任何错误, 返回该错误
func fetchFromRemote(provider models.Provider, client *http.Client, logger *zap.Logger) ([]models.Outbound, error) {
    // 创建GET请求
    req, err := http.NewRequest("GET", provider.Path, nil)
    if err != nil {
        // 如果创建请求失败, 记录错误并返回自定义错误信息
        logger.Error(fmt.Sprintf("创建请求失败: [%s]", err.Error()))
        return nil, fmt.Errorf("'%s'出错: 创建请求失败", provider.Name)
    }
    
    // 发送HTTP请求
    res, err := client.Do(req)
    if err != nil {
        // 如果发送请求失败, 记录错误并返回自定义错误信息
        logger.Error(fmt.Sprintf("发送请求失败: [%s]", err.Error()))
        return nil, fmt.Errorf("'%s'出错: 发送请求失败", provider.Name)
    }
    // 确保在函数返回前关闭响应体
    defer res.Body.Close()
    
    // 检查响应状态码是否为200(成功)
    if res.StatusCode == 200 {
        // 读取响应内容
        content, err := io.ReadAll(res.Body)
        if err != nil {
            // 如果读取响应失败, 记录错误并返回自定义错误信息
            logger.Error(fmt.Sprintf("读取'%s'响应失败: [%s]",  provider.Name, err.Error()))
            return nil, fmt.Errorf("读取'%s'响应失败", provider.Name)
        }
        
        // 解析响应内容
        outbounds, err := parseFileContent(content, logger)
        if err != nil {
            // 如果解析内容失败, 记录错误并返回自定义错误信息
            logger.Error(fmt.Sprintf("解析'%s'文件失败: [%s]", provider.Name, err.Error()))
            return nil, fmt.Errorf("'%s'出错: %s", provider.Name, err.Error())
        }
        
        // 返回解析后的出站信息
        return outbounds, nil
    }
    
    // 如果响应状态码不是200, 返回自定义错误信息
    return nil, fmt.Errorf("'%s'未知响应, 状态码: %d", provider.Name, res.StatusCode)
}

// fetchFromLocal 从本地文件中获取出站信息
// 该函数接收一个包含文件路径的provider对象和一个logger对象用于错误日志记录
// 它返回一个Outbound对象列表, 这些对象是从文件内容解析出来的
// 如果在文件处理或解析过程中遇到错误, 它将记录错误并返回相应的错误信息
func fetchFromLocal(provider models.Provider, logger *zap.Logger) ([]models.Outbound, error) {
    // 尝试打开本地文件
    file, err := os.Open(provider.Path)
    if err != nil {
        // 如果打开文件失败, 记录错误并返回自定义错误信息
        logger.Error(fmt.Sprintf("打开'%s'文件失败: [%s]", provider.Name, err.Error()))
        return nil, fmt.Errorf("打开'%s'文件失败", provider.Name)
    }
    // 确保在函数返回前关闭文件
    defer file.Close()
    
    // 读取文件的全部内容
    content, err := io.ReadAll(file)
    if err != nil {
        // 如果读取文件失败, 记录错误并返回自定义错误信息
        logger.Error(fmt.Sprintf("读取'%s'文件失败: [%s]", provider.Name, err.Error()))
        return nil, fmt.Errorf("读取'%s'文件失败", provider.Name)
    }
    
    // 解析文件内容
    outbounds, err := parseFileContent(content, logger)
    if err != nil {
        // 如果解析文件内容失败, 记录错误并返回自定义错误信息
        logger.Error(fmt.Sprintf("解析'%s'文件失败: [%s]", provider.Name, err.Error()))
        return nil, fmt.Errorf("'%s'出错: %s", provider.Name, err.Error())
    }
    
    // 返回解析出的出站信息列表
    return outbounds, nil
}

// parseFileContent 解析文件内容并返回一组出站代理配置
// 该函数接受文件内容的字节切片和一个日志记录器作为参数
// 它尝试解析内容, 根据解析结果创建并返回一个Outbound对象切片
// 如果解析过程中遇到错误, 它会记录错误信息并返回错误
func parseFileContent(content []byte, logger *zap.Logger) ([]models.Outbound, error) {
    // 初始化一个空的providerInfo映射, 用于存储解析后的文件内容
    var providerInfo map[string]interface{}
    // 使用yaml.Unmarshal解析文件内容如果解析失败, 记录错误并返回错误
    if err := yaml.Unmarshal(content, &providerInfo); err != nil {
        logger.Error(fmt.Sprintf("解析响应失败: [%s]", err.Error()))
        return nil, fmt.Errorf("解析响应失败")
    }

    // 初始化一个空的outbounds切片, 用于存储解析后的出站代理配置
    var outbounds []models.Outbound
    // 从providerInfo中获取'proxies'字段, 如果获取失败或类型不正确, 返回错误
    proxies, ok := providerInfo["proxies"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("'proxies'字段丢失或不正确")
    }

    // 遍历proxies中的每个代理配置
    for _, proxy := range proxies {
        // 将代理配置转换为map[string]interface{}类型, 如果转换失败, 记录错误并跳过当前配置
        proxyMap, ok := proxy.(map[string]interface{})
        if !ok {
            logger.Error("该节点不是map类型")
            continue
        }

        // 从proxyMap中获取'protocol'字段, 如果获取失败或类型不正确, 记录错误并跳过当前配置
        protocol, ok := proxyMap["type"].(string)
        if !ok {
            logger.Error("该节点没有'type'字段")
            continue
        }

        // 从proxyMap中获取'name'字段, 如果获取失败或类型不正确, 记录错误并跳过当前配置
        name, ok := proxyMap["name"].(string)
        if !ok {
            logger.Error("该节点没有'name'字段")
            continue
        }

        // 根据'protocol'字段的值, 选择相应的代理类型进行处理
        switch protocol {
        case "ss":
            // 对于'SS'协议, 创建一个ShadowSocks对象, 尝试解析配置并添加到outbounds中
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
            // 对于'VMess'协议, 创建一个VMess对象, 尝试解析配置并添加到outbounds中
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
            // 对于'Trojan'协议, 创建一个Trojan对象, 尝试解析配置并添加到outbounds中
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

    // 返回解析后的出站代理配置切片
    return outbounds, nil
}