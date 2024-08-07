package singbox

import (
	"encoding/base64"
	"fmt"
	"sifu-box/utils"
	"strings"

	"github.com/gocolly/colly/v2"
	"gopkg.in/yaml.v3"
)

// FetchProxies 根据给定的URL和名称获取代理配置
// 它首先尝试从URL中解析YAML格式的代理配置,如果失败,则尝试解析Base64编码的配置
// 如果两种方法都失败,则返回错误
// 参数:
//   url - 代理配置的URL
//   name - 代理配置的名称,用于日志记录和错误信息
// 返回值:
//   []map[string]interface{} - 解析到的代理配置,每个代理作为一个map[string]interface{}元素
//   error - 如果在解析过程中发生错误,返回该错误
func FetchProxies(url,name string) ([]map[string]interface{},error) {
    var proxies []map[string]interface{}
    var err error
    c := colly.NewCollector()

    // 当收到响应时,尝试解析并处理响应内容
    c.OnResponse(func(r *colly.Response) {
        var results []map[string]interface{}
        
        // 尝试将响应内容解析为YAML格式
        content := map[string]interface{}{}
        if err = yaml.Unmarshal(r.Body, &content); err != nil {
            // 如果YAML解析失败,则记录错误并尝试以Base64格式解码响应内容
            utils.LoggerCaller(fmt.Sprintf("解析'%s'yaml配置文件失败",name), err, 1)
            var base64msg []byte
            base64msg, err = base64.StdEncoding.DecodeString(string(r.Body))
            if err != nil {
                // 如果Base64解码也失败,则记录错误
                utils.LoggerCaller(fmt.Sprintf("'%s'base64解码失败",name), err, 1)
                return
            }
            // 使用Base64解码后的内容作为URL,并尝试解析出代理配置
            results, err = ParseUrl(strings.Split(string(base64msg), "\n"), name)
            if err != nil {
                // 如果解析URL失败,则记录错误
                utils.LoggerCaller(fmt.Sprintf("生成'%s'配置文件失败",name), err, 1)
            }
        } else {
            // 如果响应内容成功解析为YAML格式,则尝试从其中提取代理配置
            if proxiesMsg,ok := content["proxies"].([]interface{}); ok {
                results, err = ParseYaml(proxiesMsg, name)
            }else{
                // 如果YAML内容中没有"proxies"字段,则记录错误
                err = fmt.Errorf("'%s'配置没有proxies字段",name)
            }
            if err != nil {
                // 如果从YAML解析过程中发生错误,则记录错误
                utils.LoggerCaller(fmt.Sprintf("生成'%s'配置文件失败",name), err, 1)
            }
        }
        // 将解析得到的代理配置保存到proxies变量中
        proxies = results
    })

    // 当发生错误时,记录错误信息并尝试修正
    c.OnError(func(r *colly.Response, e error) {
        // 记录连接错误
        utils.LoggerCaller(fmt.Sprintf("连接'%s'失败", name), e, 1)
        err = e
        // 尝试移除URL查询参数中的"flag",并重新访问URL
        request_url := r.Request.URL
        params := request_url.Query()
        for k, v := range params {
            if k == "flag" && v[0] == "clash" {
                params.Del("flag")
                request_url.RawQuery = params.Encode()
                c.Visit(request_url.String())
            }
        }
    })
    // 访问指定的URL,开始抓取过程
    c.Visit(url)
    if err != nil {
        // 如果在访问过程中发生错误,则返回错误
        return nil, err
    }
    // 返回成功解析到的代理配置
    return proxies, nil
}