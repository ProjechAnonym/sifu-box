package singbox

import (
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
func addSelectorOutbound(outbounds []models.Outbound, rulesetsList []models.RuleSet, logger *zap.Logger) {
	targets := filterRulesetList(rulesetsList)
	var selector models.Selector
	selectorMap := map[string]interface{}{"type": "selector", "interrupt_exist_connections": false, "outbounds": tags, "tag": "select"}
	outbound = &selector
	outbound, err = outbound.Transform(selectorMap, logger)
}