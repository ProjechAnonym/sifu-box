package singbox

import (
	"encoding/base64"
	"errors"
	"fmt"
	utils "sifu-box/Utils"
	"strings"

	"github.com/gocolly/colly/v2"
	"gopkg.in/yaml.v3"
)

// handle_yaml 处理配置文件中的代理设置,根据指定的模板格式化代理信息
// 参数:
// config: 包含配置信息的映射,其中"proxies"子映射包含代理配置
// host: 当前处理的主机名或服务名,用于错误日志中标识配置来源
// template: 用于格式化代理信息的模板
// 返回值:
// []map[string]interface{}: 格式化后的代理信息列表,每个代理是一个映射
// error: 如果处理过程中出现错误,则返回错误信息
func handle_yaml(config map[string]interface{}, host string, template string) ([]map[string]interface{}, error) {
    // 检查配置中的代理列表是否为空,如果为空,则记录错误日志并返回错误
    if len(config["proxies"].([]interface{})) == 0 {
        utils.Logger_caller("no proxies", fmt.Errorf("%s config without proxies", host), 1)
        return nil, errors.New("no proxies")
    }
    var proxies []map[string]interface{}

    // 遍历配置中的每个代理,格式化代理信息,并将格式化后的信息添加到代理列表中
    for _, proxy := range config["proxies"].([]interface{}) {
        result, err := Format_yaml(proxy.(map[string]interface{}), template)
        // 如果格式化成功,则将格式化后的代理信息添加到列表中
        if err == nil {
            proxies = append(proxies, result)
        }
    }
    // 返回格式化后的代理信息列表和nil错误,表示处理成功
    return proxies, nil
}
// handle_url 根据给定的URL列表、主机名和模板,处理并返回有效的代理配置
// urls: URL列表,用于获取代理配置信息
// host: 主机名,用于错误日志中标识配置来源
// template: URL模板,用于格式化获取的代理配置URL
// 返回值: 一个包含多个代理配置的map切片,以及可能出现的错误
func handle_url(urls []string, host string, template string) ([]map[string]interface{}, error) {
    // 检查URL列表是否为空,如果为空,则记录错误日志并返回错误
    if len(urls) == 0 {
        utils.Logger_caller("no proxies", fmt.Errorf("%s config without proxies", host), 1)
        return nil, errors.New("no proxies")
    }

    // 初始化一个空的代理配置切片
    var proxies []map[string]interface{}

    // 遍历URL列表,尝试对每个URL进行格式化处理
    for _, url := range urls {
        // 格式化URL并获取结果,如果格式化成功,则将结果添加到代理配置切片中
        result, err := Format_url(url, template)
        if err == nil {
            proxies = append(proxies, result)
        }
    }

    // 返回处理后的代理配置切片和nil错误
    return proxies, nil
}
// fetch_proxies 从给定的URL获取代理服务器配置
// url: 获取代理配置的URL
// template: 代理配置的模板字符串,用于解析返回的内容
// 返回值: 一个包含多个代理配置的map切片,以及可能的错误
func fetch_proxies(url string,template string) ([]map[string]interface{},error){
    // 初始化代理切片和错误变量
    var proxies []map[string]interface{}
    var err error
    // 创建一个新的Colly收集器实例
    c := colly.NewCollector()
    
    // 当收到响应时,执行以下回调函数
    c.OnResponse(func(r *colly.Response) {
        // 初始化用于存储解析结果的切片和错误变量
        var results []map[string]interface{}
        var handle_err error
        // 创建一个空的map,用于存储解析后的配置
        content := map[string]interface{}{}
        // 尝试解析响应体为yaml格式
        if err := yaml.Unmarshal(r.Body,&content); err != nil{
            // 解析失败时,记录错误日志
            msg := fmt.Sprintf("Parse %s proxies yaml failed!",r.Request.URL.Host)
            utils.Logger_caller(msg,err,1)
            // 尝试将响应体解码为base64
            decodedBytes, err := base64.StdEncoding.DecodeString(string(r.Body))
            if err != nil{
                // 解码失败时,记录错误日志
                msg := fmt.Sprintf("decode base64 %s proxies failed!",r.Request.URL.Host)
                utils.Logger_caller(msg,err,1)
                return
            }
            // 使用解码后的数据处理代理配置
            results,handle_err = handle_url(strings.Split(string(decodedBytes), "\n"),r.Request.URL.Host,template)
        }else{
            // 解析成功时,直接处理yaml格式的代理配置
            results,handle_err = handle_yaml(content,r.Request.URL.Host,template)
        }
        // 更新代理切片和错误变量
        proxies = results
        err = handle_err
    })
    
    // 当请求发生错误时,执行以下回调函数
    c.OnError(func(r *colly.Response, e error) {
        // 记录错误日志
        utils.Logger_caller(fmt.Sprintf("Connect to %s failed!",r.Request.URL.Host),e,1)
        err = e
        // 检查请求参数中是否包含"flag"为"clash"的项,如果存在,则尝试重新访问URL
        request_url := r.Request.URL
        params := request_url.Query()
        for k, v := range params {
            if k=="flag" && v[0]=="clash"{
                params.Del("flag")
                request_url.RawQuery = params.Encode()
                c.Visit(request_url.String())
            }
        }
    })
    // 访问指定的URL
    c.Visit(url)
    // 如果存在错误,返回空切片和错误；否则返回解析出的代理配置切片
    if err != nil{
        return nil,err
    }
    return proxies,nil
}
// outbound_select 根据给定的标签和模板,选择出站策略。
// 参数tags是一个字符串切片,代表可用的标签选项。
// 参数template是一个字符串,代表模板信息。
// 返回一个map[string]interface{},包含选中的出站策略信息,以及一个error,用于表示可能发生的错误。
func outbound_select(tags []string, template string) (map[string]interface{}, error) {
    // 从模板中提取出"outbounds"部分的"select"配置
    select_map, err := utils.Get_value(template, "outbounds", "select")
    if err != nil {
        // 如果提取过程中发生错误,记录日志并返回错误
        utils.Logger_caller("get select outbound template failed!", err, 1)
        return nil, err
    }
    // 将"auto"标签添加到标签切片中,确保自动选择选项可用
    tags = append(tags, "auto")
    // 将更新后的标签切片赋值回select_map中的"outbounds"字段
    select_map.(map[string]interface{})["outbounds"] = tags
    // 返回更新后的select_map以及nil错误
    return select_map.(map[string]interface{}), nil
}
// outbound_auto 根据给定的标签和模板生成自动外发配置。
// 参数tags是用于外发的标签列表,template是配置模板的字符串。
// 返回一个包含外发配置的map以及可能的错误。
func outbound_auto(tags []string,template string) (map[string]interface{},error){
    // 从模板中提取自动外发配置。
    // 获取模板信息
    auto_map,err := utils.Get_value(template,"outbounds","auto")
    if err != nil{
        // 如果提取配置失败,记录错误并返回。
        utils.Logger_caller("get auto outbound template failed!",err,1)
        return nil,err
    }
    
    // 将标签列表赋值给自动外发配置的outbounds字段。
    // 设置选择节点的出站标签
    auto_map.(map[string]interface{})["outbounds"] = tags
    // 返回更新后的自动外发配置以及nil错误。
    return auto_map.(map[string]interface{}),nil
}
// Merge_outbounds 根据给定的URL和模板字符串,合并代理配置和自定义出站规则。
// url: 用于获取代理配置的URL。
// template: 包含出站规则模板的字符串。
// 返回值: 一个包含合并后出站规则的切片,以及可能的错误。
func Merge_outbounds(url string,template string) ([]map[string]interface{},error){
    // 从URL获取代理配置,并处理可能的错误。
	proxies,err := fetch_proxies(url,template)
	if err != nil || len(proxies) == 0 {
		utils.Logger_caller("fetch proies failed!",err,1)
		return nil,err
	}

    // 从模板中提取自定义和默认出站规则。
    template_outbounds,err := utils.Get_value(template,"outbounds","custom_outbound")
    custom_outbounds := make([]map[string]interface{},len(template_outbounds.([]interface{})[1].([]interface{})))
    default_outbounds := make([]map[string]interface{},len(template_outbounds.([]interface{})[0].([]interface{})))
    // 设置自定义规则列表。
    for i,custom_outbound := range template_outbounds.([]interface{})[1].([]interface{}){
        custom_outbounds[i] = custom_outbound.(map[string]interface{})
    }
    // 设置默认规则列表
    for i,default_outbound := range template_outbounds.([]interface{})[0].([]interface{}){
        default_outbounds[i] = default_outbound.(map[string]interface{})
    }
    if err != nil{
        utils.Logger_caller("get template template outbound failed!",err,1)
        return nil,err
    }

    // 生成标签列表
    tags := make([]string,len(custom_outbounds) + len(proxies))
    proxies = append(proxies, custom_outbounds...)
    for i,proxy := range proxies{
        tags[i] = proxy["tag"].(string)
    }

    // 生成自动选择的出站规则,并处理可能的错误。
    auto,err := outbound_auto(tags,template)
    if err != nil{
        utils.Logger_caller("get auto outbound failed!",err,1)
        return nil,err
    }
    // 添加自动出站设置
    proxies = append(proxies,auto)

    // 添加手动选择出站
    select_outbound,err := outbound_select(tags,template)
    if err != nil{
        utils.Logger_caller("Get select outbounds failed!",err,1)
        return proxies,err
    }
    proxies = append(proxies, select_outbound)

    // 根据规则集生成并添加相应的选择器出站规则。

    proxy_config,err := utils.Get_value("Proxy")
    rulesets := proxy_config.(utils.Box_config).Rule_set
	if err != nil{
		utils.Logger_caller("Get ruleset msg failed!",err,1)
		return proxies,err
	}

	for _,rule := range(rulesets){
		// 如果china标签为否,说明是连国外的网站,为其生成selector类型出站
		if !rule.Value.China{
			// 生成selector类型出站
			ruleset_select_outbound,err := outbound_select(tags,template)
			if err != nil{
				utils.Logger_caller("Get select outbounds failed!",err,1)
				return proxies,err
			}
			// 生成selector类型出站成功则添加进出站列表
			ruleset_select_outbound["tag"] = rule.Label + "-select"
			proxies = append(proxies, ruleset_select_outbound)
		}
	}
    // 添加默认出站规则
    proxies = append(proxies, default_outbounds...)
	return proxies,err
}