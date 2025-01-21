package models

type DNSLogic struct {
	Type  string    `json:"type,omitempty" yaml:"type,omitempty"`
	Rules []DNSRule `json:"rules,omitempty" yaml:"rules,omitempty"`
	Mode  string    `json:"mode,omitempty" yaml:"mode,omitempty"`
}
type DNSAction struct {
	Action       string `json:"action,omitempty" yaml:"action,omitempty"`
	Server       string `json:"server,omitempty" yaml:"server,omitempty"`
	DisableCache bool   `json:"disable_cache,omitempty" yaml:"disable_cache,omitempty"`
	RewriteTTL   int    `json:"rewrite_ttl,omitempty" yaml:"rewrite_ttl,omitempty"`
	ClientSubnet string `json:"client_subnet,omitempty" yaml:"client_subnet,omitempty"`
	Method       string `json:"method,omitempty" yaml:"method,omitempty"`
	NoDrop       bool   `json:"no_drop,omitempty" yaml:"no_drop,omitempty"`
}
type DNSRule struct {
	Inbound                  []string      `json:"inbound,omitempty" yaml:"inbound,omitempty"`
	IPVersion                int           `json:"ip_version,omitempty" yaml:"ip_version,omitempty"`
	QueryType                []interface{} `json:"query_type,omitempty" yaml:"query_type,omitempty"`
	Network                  string        `json:"network,omitempty" yaml:"network,omitempty"`
	AuthUser                 []string      `json:"auth_user,omitempty" yaml:"auth_user,omitempty"`
	Protocol                 []string      `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Domain                   []string      `json:"domain,omitempty" yaml:"domain,omitempty"`
	DomainSuffix             []string      `json:"domain_suffix,omitempty" yaml:"domain_suffix,omitempty"`
	DomainKeyword            []string      `json:"domain_keyword,omitempty" yaml:"domain_keyword,omitempty"`
	DomainRegex              []string      `json:"domain_regex,omitempty" yaml:"domain_regex,omitempty"`
	SourceGeoIP              []string      `json:"source_geoip,omitempty" yaml:"source_geoip,omitempty"`
	SourceIPCIDR             []string      `json:"source_ip_cidr,omitempty" yaml:"source_ip_cidr,omitempty"`
	SourceIPIsPrivate        bool          `json:"source_ip_is_private,omitempty" yaml:"source_ip_is_private,omitempty"`
	IPCIDR                   []string      `json:"ip_cidr,omitempty" yaml:"ip_cidr,omitempty"`
	IPIsPrivate              bool          `json:"ip_is_private,omitempty" yaml:"ip_is_private,omitempty"`
	SourcePort               []int         `json:"source_port,omitempty" yaml:"source_port,omitempty"`
	SourcePortRange          []string      `json:"source_port_range,omitempty" yaml:"source_port_range,omitempty"`
	Port                     []int         `json:"port,omitempty" yaml:"port,omitempty"`
	PortRange                []string      `json:"port_range,omitempty" yaml:"port_range,omitempty"`
	ProcessName              []string      `json:"process_name,omitempty" yaml:"process_name,omitempty"`
	ProcessPath              []string      `json:"process_path,omitempty" yaml:"process_path,omitempty"`
	ProcessPathRegex         []string      `json:"process_path_regex,omitempty" yaml:"process_path_regex,omitempty"`
	PackageName              []string      `json:"package_name,omitempty" yaml:"package_name,omitempty"`
	User                     []string      `json:"user,omitempty" yaml:"user,omitempty"`
	UserID                   []int         `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	ClashMode                string        `json:"clash_mode,omitempty" yaml:"clash_mode,omitempty"`
	NetworkType              []string      `json:"network_type,omitempty" yaml:"network_type,omitempty"`
	NetworkIsExpensive       bool          `json:"network_is_expensive,omitempty" yaml:"network_is_expensive,omitempty"`
	NetworkIsConstrained     bool          `json:"network_is_constrained,omitempty" yaml:"network_is_constrained,omitempty"`
	WifiSSID                 []string      `json:"wifi_ssid,omitempty" yaml:"wifi_ssid,omitempty"`
	WifiBSSID                []string      `json:"wifi_bssid,omitempty" yaml:"wifi_bssid,omitempty"`
	RuleSet                  []string      `json:"rule_set,omitempty" yaml:"rule_set,omitempty"`
	RuleSetIPCIDRMatchSource bool          `json:"rule_set_ipcidr_match_source,omitempty" yaml:"rule_set_ipcidr_match_source,omitempty"`
	RuleSetIPCIDRAcceptEmpty bool          `json:"rule_set_ip_cidr_accept_empty,omitempty" yaml:"rule_set_ip_cidr_accept_empty,omitempty"`
	Invert                   bool          `json:"invert,omitempty" yaml:"invert,omitempty"`
	Outbound                 []string      `json:"outbound,omitempty" yaml:"outbound,omitempty"`
	DNSLogic                 `json:",inline" yaml:",inline"`
	DNSAction                `json:",inline" yaml:",inline"`
}

type FakeIP struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Inet4Range string `json:"inet4_range,omitempty" yaml:"inet4_range,omitempty"`
	Inet6Range string `json:"inet6_range,omitempty" yaml:"inet6_range,omitempty"`
}
type NameServer struct {
	Tag             string `json:"tag" yaml:"tag"`
	Address         string `json:"address" yaml:"address"`
	AddressResolver string `json:"address_resolver,omitempty" yaml:"address_resolver,omitempty"`
	AddressStrategy string `json:"address_strategy,omitempty" yaml:"address_strategy,omitempty"`
	Strategy        string `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Detour          string `json:"detour,omitempty" yaml:"detour,omitempty"`
	ClientSubnet    string `json:"client_subnet,omitempty" yaml:"client_subnet,omitempty"`
}
type DNS struct {
	Final            string       `json:"final,omitempty" yaml:"final,omitempty"`
	Strategy         string       `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	DisableCache     bool         `json:"disable_cache,omitempty" yaml:"disable_cache,omitempty"`
	DisableExpire    bool         `json:"disable_expire,omitempty" yaml:"disable_expire,omitempty"`
	IndependentCache bool         `json:"independent_cache,omitempty" yaml:"independent_cache,omitempty"`
	ReverseMapping   bool         `json:"reverse_mapping,omitempty" yaml:"reverse_mapping,omitempty"`
	ClientSubnet     string       `json:"client_subnet,omitempty" yaml:"client_subnet,omitempty"`
	Fakeip           *FakeIP      `json:"fakeip,omitempty" yaml:"fakeip,omitempty"`
	Servers          []NameServer `json:"servers" yaml:"servers"`
	Rules            []DNSRule    `json:"rules" yaml:"rules"`
}

func (d *DNS) SetDNSRules(rulesetList []RuleSet) {
	var rules []DNSRule
	rules = append(rules, d.Rules...)
	for _, ruleset := range rulesetList {
		if ruleset.NameServer != "" {
			rule := DNSRule{RuleSet: []string{ruleset.Tag}, DNSAction: DNSAction{Server: ruleset.NameServer, Action: "route"}}
			rules = append(rules, rule)
		}
	}
	d.Rules = rules
}