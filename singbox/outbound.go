package singbox

import (
	"fmt"
	"os"
	"sifu-box/models"
	"sifu-box/utils"

	"gopkg.in/yaml.v3"
)

func outboundSelect(tags []string,label string) map[string]interface{} {
    selectMap := map[string]interface{}{"type":"selector","interrupt_exist_connections":false,"tag":label} 
	
    tags = append(tags, "auto")
    
    selectMap["outbounds"] = tags
    
    return selectMap
}



func outboundAuto(tags []string) map[string]interface{}{
    autoMap := map[string]interface{}{"type":"urltest","interrupt_exist_connections":false,"tag":"auto"}
    autoMap["outbounds"] = tags
    
    return autoMap
}
// MergeOutbound 整合出站配置
// 参数:
// provider - 提供者信息
// serviceMap - 服务映射表,用于区分不同服务的规则集
// outbounds - 初始出站配置列表
// 返回值:
// []map[string]interface{} - 整合后的出站配置列表
// error - 错误信息,如果操作过程中发生错误
func MergeOutbound(provider models.Provider, serviceMap map[string][]models.Ruleset, outbounds []map[string]interface{}) ([]map[string]interface{}, error) {
    // 初始化代理配置列表
    var proxies []map[string]interface{}
    // 定义content变量用于存储文件内容
    var content []byte
    // 定义err变量用于存储错误信息
    var err error

    // 根据provider的远程标志,决定是获取远程代理配置还是读取本地文件
    if provider.Remote {
        // 获取远程代理配置
        proxies, err = FetchProxies(provider.Path, provider.Name)
        if err != nil {
            return nil, err
        }
    } else {
        // 读取本地配置文件
        content, err = os.ReadFile(provider.Path)
        if err != nil {
            // 记录读取yaml失败的错误
            utils.LoggerCaller("读取yaml失败", err, 1)
            return nil, err
        }

        // 解析yaml文件为map结构
        var data map[string]interface{}
        if err = yaml.Unmarshal(content, &data); err != nil {
            // 记录解析yaml失败的错误
            utils.LoggerCaller("解析yaml失败", err, 1)
            return nil, err
        }

        // 从解析的数据中获取代理配置
        if proxiesMsg, ok := data["proxies"].([]interface{}); ok {
            proxies, err = ParseYaml(proxiesMsg, provider.Name)
        } else {
            // 如果proxies字段不存在,返回错误
            err = fmt.Errorf("proxies字段不存在")
        }
        if err != nil {
            return nil, err
        }
    }

    // 将获取到的代理配置添加到出站配置列表中
    outbounds = append(outbounds, proxies...)
    // 初始化tags列表,用于存储出站配置的标签
    tags := make([]string, len(outbounds))
    for i, outbound := range outbounds {
        // 提取并存储每个出站配置的标签
        tags[i] = outbound["tag"].(string)
    }

    // 添加一个选择器类型的出站配置
    proxies = append(proxies, outboundSelect(tags, "select"))

    // 遍历服务映射表,根据规则集添加相应的出站配置
    for key, rulesets := range serviceMap {
        if key == "" {
            // 对于没有服务名的规则集,遍历并添加对应的出站配置
            for _, ruleset := range rulesets {
                if !ruleset.China {
                    // 如果规则集不指向中国,添加选择器类型的出站配置
                    proxies = append(proxies, outboundSelect(tags, fmt.Sprintf("select-%s", ruleset.Tag)))
                }
            }
        } else {
            // 对于有服务名的规则集,添加相应的出站配置
            if !rulesets[0].China {
                proxies = append(proxies, outboundSelect(tags, fmt.Sprintf("select-%s", key)))
            }
        }
    }

    // 添加自动选择类型的出站配置
    proxies = append(proxies, outboundAuto(tags))

    // 返回整合后的出站配置列表和nil错误,表示操作成功
    return proxies, nil
}