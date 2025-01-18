package models

type RouteLogic struct {
	Type  string      `json:"type,omitempty" yaml:"type,omitempty"`
	Mode  string      `json:"mode,omitempty" yaml:"mode,omitempty"`
	Rules []RouteRule `json:"rules,omitempty" yaml:"rules,omitempty"`
}

type RouteAction struct {
	Action                    string   `json:"action" yaml:"action"`
	Outbound                  string   `json:"outbound,omitempty" yaml:"outbound,omitempty"`
	OverRideAddress           string   `json:"override_address,omitempty" yaml:"override_address,omitempty"`
	OverRidePort              int      `json:"override_port,omitempty" yaml:"override_port,omitempty"`
	NetworkStrategy           string   `json:"network_strategy,omitempty" yaml:"network_strategy,omitempty"`
	FallbackDelay             string   `json:"fallback_delay,omitempty" yaml:"fallback_delay,omitempty"`
	UDPDisableDomainUnmapping bool     `json:"udp_disable_domain_unmapping,omitempty" yaml:"udp_disable_domain_unmapping,omitempty"`
	UDPConnect                bool     `json:"udp_connect,omitempty" yaml:"udp_connect,omitempty"`
	UDPTimeout                string   `json:"udp_timeout,omitempty" yaml:"udp_timeout,omitempty"`
	Method                    string   `json:"method,omitempty" yaml:"method,omitempty"`
	NoDrop                    bool     `json:"no_drop,omitempty" yaml:"no_drop,omitempty"`
	TimeOut                   string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Sniffer                   []string `json:"sniffer,omitempty" yaml:"sniffer,omitempty"`
	Strategy                  string   `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Server                    string   `json:"server,omitempty" yaml:"server,omitempty"`
}

type RouteRule struct {
	Inbound                  []string `json:"inbound,omitempty" yaml:"inbound,omitempty"`
	IPVersion                int      `json:"ip_version,omitempty" yaml:"ip_version,omitempty"`
	Network                  []string `json:"network,omitempty" yaml:"network,omitempty"`
	AuthUser                 []string `json:"auth_user,omitempty" yaml:"auth_user,omitempty"`
	Protocol                 []string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Client                   []string `json:"client,omitempty" yaml:"client,omitempty"`
	Domain                   []string `json:"domain,omitempty" yaml:"domain,omitempty"`
	DomainSuffix             []string `json:"domain_suffix,omitempty" yaml:"domain_suffix,omitempty"`
	DomainKeyword            []string `json:"domain_keyword,omitempty" yaml:"domain_keyword,omitempty"`
	DomainRegex              []string `json:"domain_regex,omitempty" yaml:"domain_regex,omitempty"`
	SourceIPCIDR             []string `json:"source_ip_cidr,omitempty" yaml:"source_ip_cidr,omitempty"`
	SourceIPIsPrivate        bool     `json:"source_ip_is_private,omitempty" yaml:"source_ip_is_private,omitempty"`
	IPCIDR                   []string `json:"ip_cidr,omitempty" yaml:"ip_cidr,omitempty"`
	IPIsPrivate              bool     `json:"ip_is_private,omitempty" yaml:"ip_is_private,omitempty"`
	SourcePort               []int    `json:"source_port,omitempty" yaml:"source_port,omitempty"`
	SourcePortRange          []string `json:"source_port_range,omitempty" yaml:"source_port_range,omitempty"`
	Port                     []int    `json:"port,omitempty" yaml:"port,omitempty"`
	PortRange                []string `json:"port_range,omitempty" yaml:"port_range,omitempty"`
	ProcessName              []string `json:"process_name,omitempty" yaml:"process_name,omitempty"`
	ProcessPath              []string `json:"process_path,omitempty" yaml:"process_path,omitempty"`
	ProcessPathRegex         []string `json:"process_path_regex,omitempty" yaml:"process_path_regex,omitempty"`
	PackageName              []string `json:"package_name,omitempty" yaml:"package_name,omitempty"`
	User                     []string `json:"user,omitempty" yaml:"user,omitempty"`
	UserID                   []int    `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	ClashMode                string   `json:"clash_mode,omitempty" yaml:"clash_mode,omitempty"`
	NetworkType              []string `json:"network_type,omitempty" yaml:"network_type,omitempty"`
	NetworkIsExpensive       bool     `json:"network_is_expensive,omitempty" yaml:"network_is_expensive,omitempty"`
	NetworkIsConstrained     bool     `json:"network_is_constrained,omitempty" yaml:"network_is_constrained,omitempty"`
	WifiSSID                 []string `json:"wifi_ssid,omitempty" yaml:"wifi_ssid,omitempty"`
	WifiBSSID                []string `json:"wifi_bssid,omitempty" yaml:"wifi_bssid,omitempty"`
	RuleSet                  []string `json:"rule_set,omitempty" yaml:"rule_set,omitempty"`
	RuleSetIPCIDRMatchSource bool     `json:"rule_set_ip_cidr_match_source,omitempty" yaml:"rule_set_ip_cidr_match_source,omitempty"`
	Invert                   bool     `json:"invert,omitempty" yaml:"invert,omitempty"`
	RouteAction              `json:",inline" yaml:",inline"`
	RouteLogic               `json:",inline" yaml:",inline"`
}

type RouteRuleSet struct {
	Type           string `json:"type" yaml:"type"`
	Tag            string `json:"tag" yaml:"tag"`
	Format         string `json:"format" yaml:"format"`
	URL            string `json:"url,omitempty" yaml:"url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty" yaml:"download_detour,omitempty"`
	UpdateInterval string `json:"update_interval,omitempty" yaml:"update_interval,omitempty"`
	Path           string `json:"path,omitempty" yaml:"path,omitempty"`
}
type Route struct {
	RuleSet             []RouteRuleSet `json:"rule_set,omitempty" yaml:"rule_set,omitempty"`
	Rules               []RouteRule    `json:"rules" yaml:"rules"`
	Final               string         `json:"final,omitempty" yaml:"final,omitempty"`
	AutoDetectInterface bool           `json:"auto_detect_interface" yaml:"auto_detect_interface"`
	OverrideAndroidVpn  bool           `json:"override_android_vpn,omitempty" yaml:"override_android_vpn,omitempty"`
	DefaultInterface    string         `json:"default_interface,omitempty" yaml:"default_interface,omitempty"`
	DefaultMark         uint           `json:"default_mark,omitempty" yaml:"default_mark,omitempty"`
}