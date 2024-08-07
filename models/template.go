package models

type FakeIP struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Inet4_range string `json:"inet4_range,omitempty" yaml:"inet4_range,omitempty"`
	Inet6_range string `json:"inet6_range,omitempty" yaml:"inet6_range,omitempty"`
}
type DnsServer struct {
	Tag              string `json:"tag" yaml:"tag"`
	Address          string `json:"address" yaml:"address"`
	Address_resolver string `json:"address_resolver,omitempty" yaml:"address_resolver,omitempty"`
	Address_strategy string `json:"address_strategy,omitempty" yaml:"address_strategy,omitempty"`
	Strategy         string `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Detour           string `json:"detour,omitempty" yaml:"detour,omitempty"`
	Client_subnet    string `json:"client_subnet,omitempty" yaml:"client_subnet,omitempty"`
}
type Dns struct {
	Final             string                   `json:"final,omitempty" yaml:"final,omitempty"`
	Strategy          string                   `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Disable_cache     bool                     `json:"disable_cache,omitempty" yaml:"disable_cache,omitempty"`
	Disable_expire    bool                     `json:"disable_expire,omitempty" yaml:"disable_expire,omitempty"`
	Independent_cache bool                     `json:"independent_cache,omitempty" yaml:"independent_cache,omitempty"`
	Reverse_mapping   bool                     `json:"reverse_mapping,omitempty" yaml:"reverse_mapping,omitempty"`
	Client_subnet     string                   `json:"client_subnet,omitempty" yaml:"client_subnet,omitempty"`
	Fakeip            *FakeIP                  `json:"fakeip,omitempty" yaml:"fakeip,omitempty"`
	Servers           []DnsServer              `json:"servers" yaml:"servers"`
	Rules             []map[string]interface{} `json:"rules" yaml:"rules"`
}
type Route struct {
	Rule_set              []Ruleset                `json:"rule_set,omitempty" yaml:"rule_set,omitempty"`
	Rules                 []map[string]interface{} `json:"rules" yaml:"rules"`
	Final                 string                   `json:"final,omitempty" yaml:"final,omitempty"`
	Auto_detect_interface bool                     `json:"auto_detect_interface" yaml:"auto_detect_interface"`
	Override_android_vpn  bool                     `json:"override_android_vpn,omitempty" yaml:"override_android_vpn,omitempty"`
	Default_interface     string                   `json:"default_interface,omitempty" yaml:"default_interface,omitempty"`
	Default_mark          uint                     `json:"default_mark,omitempty" yaml:"default_mark,omitempty"`
}
type Template struct {
	Name            string                   `json:"-" yaml:"-"`
	Log             map[string]interface{}   `json:"log,omitempty" yaml:"log,omitempty"`
	Ntp             map[string]interface{}   `json:"ntp,omitempty" yaml:"ntp,omitempty"`
	Experimental    map[string]interface{}   `json:"experimental,omitempty" yaml:"experimental,omitempty"`
	Inbounds        []map[string]interface{} `json:"inbounds" yaml:"inbounds"`
	Dns             Dns                      `json:"dns" yaml:"dns"`
	Route           Route                    `json:"route" yaml:"route"`
	Outbounds       []map[string]interface{} `json:"outbounds" yaml:"outbounds"`
	CustomOutbounds []map[string]interface{} `json:"-" yaml:"customOutbounds,omitempty"`
}