package singbox

import (
	"fmt"
	"sifu-box/models"

	"go.uber.org/zap"
)
func filterRulesetList(rulesetsList []models.RuleSet) []string {
	targets := []string{}
	targetsMap := map[string]bool{}
	for _, ruleset := range rulesetsList {
		if targetsMap[ruleset.Label] {continue}
		targetsMap[ruleset.Label] = true
		targets = append(targets, ruleset.Label)
	}
	return targets
}
func addSelectorOutbound(provider string, outbounds []models.Outbound, rulesetsList []models.RuleSet, tags []string, logger *zap.Logger) ([]models.Outbound, error){
	targets := filterRulesetList(rulesetsList)
	var selector models.Selector
	selectorMap := map[string]interface{}{"type": "selector", "interrupt_exist_connections": false, "outbounds": tags, "tag": "select"}
	var outbound models.Outbound = &selector
	outbound, err := outbound.Transform(selectorMap, logger)
	if err != nil {
		return nil, err
	}
	outbounds = append(outbounds, outbound)
	for _, target := range targets {
		selectorMap["tag"] = target
		outbound = &selector
		outbound, err = outbound.Transform(selectorMap, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("'%s'生成%s出站失败: [%s]", provider, target, err.Error()))
		}
		outbounds = append(outbounds, outbound)
	}
	return outbounds, nil
}