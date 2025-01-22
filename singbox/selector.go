package singbox

import (
	"fmt"
	"sifu-box/models"

	"go.uber.org/zap"
)

// filterRulesetList 过滤规则集列表, 返回唯一的标签列表
// 该函数的目的是从给定的规则集列表中提取所有唯一的标签, 并以字符串切片的形式返回
// 参数:
//   rulesetsList ([]models.RuleSet): 一个规则集的切片, 每个规则集包含一个Label字段
// 返回值:
//   []string: 包含所有唯一标签的字符串切片
func filterRulesetList(rulesetsList []models.RuleSet) []string {
    // 初始化一个空的字符串切片, 用于存储唯一的标签
    targets := []string{}
    // 初始化一个映射, 用于记录标签是否已经存在, 以避免重复
    targetsMap := map[string]bool{}
    // 遍历规则集列表
    for _, ruleset := range rulesetsList {
        // 如果当前规则集的标签已经在映射中存在, 则跳过, 避免重复添加
        if targetsMap[ruleset.Label] || ruleset.China {
            continue
        }
        // 将新的标签添加到映射中, 标记为已存在
        targetsMap[ruleset.Label] = true
        // 将新的标签添加到目标列表中
        targets = append(targets, ruleset.Label)
    }
    // 返回包含所有唯一标签的切片
    return targets
}
// addSelectorOutbound 为给定的出站标签添加选择器类型的出站配置
// 该函数首先过滤规则集列表, 然后为每个目标标签生成一个选择器出站配置, 
// 并将其添加到出站配置列表中
// 参数:
//   provider - 出站配置的机场名称, 用于错误日志中
//   outbounds - 原始的出站配置列表, 函数将在其基础上添加新的出站配置
//   rulesetsList - 规则集列表, 用于确定需要添加的出站配置的目标标签
//   tags - 选择器出站配置中包含的标签列表
// 返回值:
//   添加了新的选择器出站配置的出站配置列表
//   如果生成出站配置过程中发生错误, 则返回该错误
func addSelectorOutbound(provider string, outbounds []models.Outbound, rulesetsList []models.RuleSet, tags []string, logger *zap.Logger) ([]models.Outbound, error) {
    // 过滤规则集列表, 仅保留适用的规则集
    targets := filterRulesetList(rulesetsList)

    // 初始化一个选择器实例, 用于后续生成选择器类型的出站配置
    var selector models.Selector
    // 定义选择器出站配置的基础映射, 包含选择器的通用配置
    selectorMap := map[string]interface{}{"type": "selector", "interrupt_exist_connections": false, "outbounds": tags, "tag": "select"}

    // 创建一个选择器类型的出站配置实例
    var outbound models.Outbound = &selector
    // 根据selectorMap中的配置, 生成选择器出站配置
    outbound, err := outbound.Transform(selectorMap, logger)
    if err != nil {
        return nil, err
    }
    // 将生成的选择器出站配置添加到出站配置列表中
    outbounds = append(outbounds, outbound)

    // 遍历每个目标标签, 为每个标签生成一个选择器出站配置
    for _, target := range targets {
        // 更新选择器出站配置的标签为目标标签
        selectorMap["tag"] = target
        outbound = &selector
        // 根据更新后的配置, 生成选择器出站配置
        outbound, err = outbound.Transform(selectorMap, logger)
        if err != nil {
            logger.Error(fmt.Sprintf("'%s'生成%s出站失败: [%s]", provider, target, err.Error()))
        }
        // 将生成的选择器出站配置添加到出站配置列表中
        outbounds = append(outbounds, outbound)
    }
    // 返回添加了新的选择器出站配置的出站配置列表
    return outbounds, nil
}