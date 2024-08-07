package singbox

import "sifu-box/models"

func SetDnsRules(serviceMap map[string][]models.Ruleset) []map[string]interface{}{
	var rules []map[string]interface{}
	for _,rulesets := range serviceMap{
		for _,ruleset := range(rulesets){
			if ruleset.DnsRule != "" {
				rules = append(rules, map[string]interface{}{"rule_set":ruleset.Tag,"server":ruleset.DnsRule})
			}
		}
	}
	return rules
}