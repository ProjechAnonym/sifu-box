package singbox

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/huandu/go-clone"
)

// AddClashTag 函数用于检查一组Provider中是否有带有“clash”标签的远程Provider
// 如果没有,则自动添加该标签
// 参数providers: 一个Provider对象的切片,代表一组Provider
// 返回值:
// - 第一个返回值为更新后的Provider切片,其中未带有“clash”标签的远程Provider将被添加该标签
// - 第二个返回值为处理过程中可能遇到的错误切片,例如解析URL失败
func AddClashTag(providers []models.Provider) ([]models.Provider, []error) {
    // 初始化一个错误切片,用于存储处理过程中遇到的错误
    var errors []error

    // 遍历Provider切片
    for i, provider := range providers {
        // 检查Provider是否为远程类型
        if provider.Remote {
            // 尝试解析Provider的路径URL
            parsedUrl, err := url.Parse(provider.Path)
            if err != nil {
                // 如果解析失败,记录错误日志并添加错误到错误切片
                utils.LoggerCaller(fmt.Sprintf("解析'%s'链接失败", provider.Name), err, 1)
                errors = append(errors, fmt.Errorf("解析'%s'链接失败", provider.Name))
                continue
            }

            // 获取URL查询参数
            params := parsedUrl.Query()
            // 初始化一个标志变量,用于标记是否已存在“clash”标签
            clashTag := false

            // 遍历查询参数,检查是否存在“flag”为“clash”的参数
            for key, values := range params {
                if key == "flag" && values[0] == "clash" {
                    clashTag = true
                    break
                }
            }

            // 如果不存在“clash”标签,则添加该标签到查询参数中
            if !clashTag {
                params.Add("flag", "clash")
                // 更新URL的查询参数部分
                parsedUrl.RawQuery = params.Encode()
                // 更新Provider的路径为带有“clash”标签的URL
                providers[i].Path = parsedUrl.String()
            }
        }
    }
    // 返回更新后的Provider切片和错误切片
    return providers, errors
}

// SortRulesets sorts a list of rulesets by their labels.
// Receives a list of rulesets as a parameter and returns a map where the key is the label of the ruleset and the value is a list of rulesets with the same label.
// If the input list is empty, it returns nil.
func SortRulesets(rulesets []models.Ruleset) map[string][]models.Ruleset{
    // Return nil if the length of the rulesets array is 0 to avoid unnecessary operations
    if len(rulesets) == 0{
        return nil
    }
    // Initialize a map to store the sorted rulesets, where the key is the label of the ruleset, and the value is a list of rulesets
    serviceMap := make(map[string][]models.Ruleset)
    // Iterate through the list of rulesets to group them by label
    for _,ruleset := range rulesets{
        // If the ruleset with the current label has not yet been added to the map, initialize it as a list containing only the current ruleset
        if serviceMap[ruleset.Label] == nil{
            serviceMap[ruleset.Label] = []models.Ruleset{ruleset}
        }else{
            // If the ruleset with the current label has already been added to the map, then append the current ruleset to its list
            serviceMap[ruleset.Label] = append(serviceMap[ruleset.Label], ruleset)
        }
    }
    // Return the map of sorted rulesets
    return serviceMap
}

// GetMapValue 从嵌套的map中获取值
// dstMap: 目标map,可能包含嵌套的map
// keys: 需要获取的值对应的键名,支持多级嵌套
// 返回值: 获取到的值,如果过程中遇到任何问题,返回错误信息
func GetMapValue(dstMap map[string]interface{}, keys ...string) (interface{}, error) {
    // 初始化tempMap为dstMap,用于后续的逐级检索
    tempMap := dstMap

    // 遍历keys,尝试从map中获取值
    for i, key := range keys {
        // 检查当前key对应的值是否存在
        if tempMap[key] != nil {
            // 如果已经是最后一级key,则尝试返回该值的副本
            if i == len(keys)-1 {
                return clone.Clone(tempMap[key]), nil
            }
            // 如果当前值是一个map,将其作为新的tempMap继续检索
            if subMap, ok := tempMap[key].(map[string]interface{}); ok {
                tempMap = subMap
            } else {
                // 如果当前值不是map,返回错误信息
                return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
            }
        } else {
            // 如果当前key不存在,返回错误信息
            return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
        }
    }

    // 如果遍历完keys没有返回值,说明参数不足,返回错误信息
    return nil, fmt.Errorf("参数不足,缺少键值参数")
}
// Struct2map 将结构体转换为map这个函数支持三种类型的结构体：Vmess、ShadowSocks和Trojan
// 参数P是待转换的结构体实例,class是结构体的类别字符串,用于错误信息中标识结构体类型
// 返回值是一个map,其中包含了结构体转换后的键值对,或者在转换过程中发生错误时返回错误
func Struct2map[P models.Vmess | models.ShadowSocks | models.Trojan](proxy P,class string) (map[string]interface{},error){
	
	// 将结构体实例转换为JSON字节切片,以便下一步将其解析到map中
	proxyBytes, err := json.Marshal(proxy)
	if err != nil{
		// 如果JSON序列化失败,记录错误并返回
		utils.LoggerCaller(fmt.Sprintf("json序列化'%s'失败",class),err,1)
		return nil,err
	}
	
	// 初始化一个空的map,用于接收解析后的结构体数据
	var proxyMap map[string]interface{}
	
	// 将JSON字节切片解析到map中
	err = json.Unmarshal(proxyBytes, &proxyMap)
	if err != nil {
		// 如果解析失败,记录错误并返回
		utils.LoggerCaller(fmt.Sprintf("'%s'json序列转换为字典失败",class),err,1)
		return nil,err
	}
	// 返回解析后的map,如果没有错误发生
	return proxyMap,nil
}