package singbox

import (
	"fmt"
	"net/url"
	"sifu-box/models"
	"sifu-box/utils"
)

// SetRulesets 根据给定的服务映射创建新的规则集列表
// 该函数接收一个map,其中key是服务名称,value是与该服务相关的规则集列表
// 返回值是所有规则集的扁平化列表
func SetRulesets(serviceMap map[string][]models.Ruleset) []models.Ruleset {
    // 初始化一个新的规则集列表,用于存储所有服务的规则集
	var newRulesets []models.Ruleset
	
    // 遍历服务映射,处理每个服务及其相关的规则集
	for _, rulesets := range serviceMap {
        // 遍历单个服务的所有规则集
		for _,ruleset := range(rulesets){
            // 将每个规则集添加到新的规则集列表中这包括复制规则集的各个字段
			newRulesets = append(newRulesets, models.Ruleset{Type: ruleset.Type, Tag: ruleset.Tag, Download_detour: ruleset.Download_detour, Format: ruleset.Format, Path: ruleset.Path, Update_interval: ruleset.Update_interval, Url: ruleset.Url})
		}
	}
    // 返回包含所有规则集的扁平化列表
	return newRulesets
}
// SetRules 根据服务映射和提供者信息设置规则
// serviceMap: 一个映射,其中包含服务名称和相应的规则集列表
// provider: 提供者对象,包含提供者的相关信息和配置
// 返回值: 一个切片,包含设置的规则,每个规则是一个映射,其中包含“域”、“规则集”和“出站”等键值对
func SetRules(serviceMap map[string][]models.Ruleset, provider models.Provider) []map[string]interface{}{
    // 初始化规则切片
    var rules []map[string]interface{}

    // 处理远程提供者
    if provider.Remote {
        // 如果提供者是远程的,则添加一条包含其路径和选择性出站的规则
        providerPath,err := url.Parse(provider.Path)
        if err != nil {
            utils.LoggerCaller("解析域名失败",err,1)
        }else{
            providerDomain := providerPath.Hostname()
            rules = append(rules, map[string]interface{}{"domain": providerDomain, "outbound": "select"})
        }
    }

    // 遍历服务映射,根据服务名称和规则集设置规则
    for key, rulesets := range serviceMap {
        // 处理未指定服务名称的情况
        if key == "" {
            // 遍历规则集,根据是否针对中国用户添加不同的出站规则
            for _, ruleset := range rulesets {
                if ruleset.China {
                    // 如果规则针对中国用户,则添加一条指向直接出站的规则
                    rules = append(rules, map[string]interface{}{"rule_set": ruleset.Tag, "outbound": "direct"})
                } else {
                    // 如果规则不针对中国用户,则添加一条指向选择性出站的规则,出站名称由规则标签生成
                    rules = append(rules, map[string]interface{}{"rule_set": ruleset.Tag, "outbound": fmt.Sprintf("select-%s", ruleset.Tag)})
                }
            }
        } else {
            // 处理指定了服务名称的情况
            var rulesetsList []string
            var china bool
            var label string
            // 遍历规则集,收集规则标签,并确定是否针对中国用户
            for _, ruleset := range rulesets {
                china = ruleset.China
                label = ruleset.Label
                rulesetsList = append(rulesetsList, ruleset.Tag)
            }
            // 根据是否针对中国用户,添加相应的出站规则
            if china {
                rules = append(rules, map[string]interface{}{"rule_set": rulesetsList, "outbound": "direct"})
            } else {
                rules = append(rules, map[string]interface{}{"rule_set": rulesetsList, "outbound": fmt.Sprintf("select-%s", label)})
            }
        }
    }
    // 返回设置的规则切片
    return rules
}
