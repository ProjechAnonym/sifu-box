package singbox

import (
	"net/url"
	utils "sifu-box/Utils"
)

// 根据模板获取规则集
// 该函数结合了默认规则集和自定义规则集,根据模板配置生成最终的规则集列表
// 参数 template 是配置模板的字符串表示
// 返回值是一个包含多个规则集的切片,每个规则集是一个包含不同字段的映射表
func get_ruleset(template string) ([]map[string]interface{}, error) {
    // 从模板中获取默认规则集
    default_rulesets, err := utils.Get_value(template, "route", "rule_set")
    if err != nil {
        // 如果获取默认规则集失败,记录错误并返回
        utils.Logger_caller("Get default rulesets failed!", err,1)
        return nil, err
    }

    // 从配置中获取自定义规则集
    custom_rulesets, err := utils.Get_value("Proxy", "rule_set")
    if err != nil {
        // 如果获取自定义规则集失败,记录错误并返回
        utils.Logger_caller("Get custom_rule_set failed!", err,1)
        return nil, err
    }

    // 初始化规则集切片,大小为默认规则集和自定义规则集长度之和
    rulesets := make([]map[string]interface{}, len(default_rulesets.([]interface{}))+len(custom_rulesets.([]interface{})))

    // 遍历自定义规则集,根据类型创建新的规则集映射表,并追加到默认规则集中
    for _, rule := range custom_rulesets.([]interface{}) {
        // 提取自定义规则的标签和路径信息
        tag := rule.(map[string]interface{})["label"].(string)
        path := rule.(map[string]interface{})["value"].(map[string]interface{})["path"].(string)
		format := rule.(map[string]interface{})["value"].(map[string]interface{})["format"].(string)
        // 根据类型创建规则集映射表
        ruleset := make(map[string]interface{})
        switch rule.(map[string]interface{})["value"].(map[string]interface{})["type"] {
        case "local":
            ruleset = map[string]interface{}{"tag": tag,"type": "local", "format": format, "path": path}
        case "remote":
			detour := rule.(map[string]interface{})["value"].(map[string]interface{})["download_detour"].(string)
			interval := rule.(map[string]interface{})["value"].(map[string]interface{})["update_interval"].(string)
            ruleset = map[string]interface{}{"tag": tag, "type": "remote",  "format": format, "url": path, "download_detour": detour, "update_interval": interval}
        }

        // 将新创建的规则集追加到默认规则集中
        default_rulesets = append(default_rulesets.([]interface{}), ruleset)
    }

    // 将默认规则集中的每个规则集映射表添加到最终的规则集切片中
    for i, ruleset := range default_rulesets.([]interface{}) {
        rulesets[i] = ruleset.(map[string]interface{})
    }

    // 返回合并后的规则集切片
    return rulesets, nil
}

// get_rules 根据模板和配置文件生成路由规则
// template: 路由规则模板字符串
// 返回值:
// - 一个包含多个规则的map切片,每个规则是一个包含"rule_set"和"outbound"键值对的map
// - 错误对象,如果在获取规则过程中发生错误
func get_rules(template string,link string,proxy bool) ([]map[string]interface{}, error) {
    // 从模板中获取默认路由规则
    base_rules, err := utils.Get_value(template, "route", "rules", "default")
    if err != nil {
        // 如果获取默认规则失败,记录错误并返回
        utils.Logger_caller("Get origin rules failed!", err,1)
        return nil, err
    }

    // 解析链接以获取域名
    domain,err := url.Parse(link)
    if err != nil {
        // 日志记录域名解析失败
        utils.Logger_caller("Extract the domain failed!", err,1)
        return nil, err
    }
    host := domain.Host

    // 如果使用代理，添加针对当前域名的代理规则。
    if proxy{
        base_rules = append(base_rules.([]interface{}),map[string]interface{}{"domain_keyword":[]string{host},"outbound":"select"})
    }
    
    // 从配置中获取自定义路由规则
    custom_rules, err := utils.Get_value("Proxy", "rule_set")
    if err != nil {
        // 如果获取自定义规则失败,记录错误并返回
        utils.Logger_caller("Get custom rules failed!", err,1)
        return nil, err
    }

    // 从模板中获取分流路由规则
    // 获取分流规则
    shunt_rules, err := utils.Get_value(template, "route", "rules", "shunt")
    if err != nil {
        // 如果获取分流规则失败,记录错误并返回
        utils.Logger_caller("Get shunt rules failed!", err,1)
        return nil, err
    }

    // 根据默认规则、自定义规则和分流规则的总数初始化规则切片
    rules := make([]map[string]interface{}, len(base_rules.([]interface{}))+len(custom_rules.([]interface{}))+len(shunt_rules.([]interface{})))

    // 遍历自定义规则,根据规则的"china"值决定是直接路由还是选择路由
    // 此处逻辑于上面相同
    for _, rule := range custom_rules.([]interface{}) {
        tag := rule.(map[string]interface{})["label"].(string)
        switch rule.(map[string]interface{})["value"].(map[string]interface{})["china"] {
        case true:
            // 如果"china"为true,添加直接路由规则
            base_rules = append(base_rules.([]interface{}), map[string]interface{}{"rule_set": tag, "outbound": "direct"})
        case false:
            // 如果"china"为false,添加选择路由规则
            base_rules = append(base_rules.([]interface{}), map[string]interface{}{"rule_set": tag, "outbound": tag + "-select"})
        }
    }

    // 将分流规则追加到默认规则之后
    base_rules = append(base_rules.([]interface{}), shunt_rules.([]interface{})...)

    // 将合并后的规则复制到最终的规则切片中
    for i, rule_set := range base_rules.([]interface{}) {
        rules[i] = rule_set.(map[string]interface{})
    }

    // 返回处理后的规则切片和nil错误
    return rules, nil
}

// Merge_route 根据模板字符串合并路由规则和规则集
// template: 模板字符串包含路由配置信息
// 返回值:
// - 一个map[string]interface{}类型的路由配置,包含规则集和规则
// - 错误信息,如果在处理过程中出现错误
func Merge_route(template string,url string,proxy bool) (map[string]interface{}, error) {
    // 从模板中提取规则集信息
    rulesets, err := get_ruleset(template)
    if err != nil {
        return nil, err
    }

    // 从模板中提取路由规则信息
    rules, err := get_rules(template,url,proxy)
    if err != nil {
        return nil, err
    }

    // 从模板中提取路由配置
    route, err := utils.Get_value(template, "route")
    if err != nil {
        // 记录路由解析失败的日志
        utils.Logger_caller("Marshal route failed!", err,1)
        return nil, err
    }

    // 将提取的规则集和规则合并到路由配置中
    route.(map[string]interface{})["rule_set"] = rulesets
    route.(map[string]interface{})["rules"] = rules

    // 返回合并后的路由配置
    return route.(map[string]interface{}), nil
}