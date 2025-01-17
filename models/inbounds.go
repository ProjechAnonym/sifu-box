package models

type Tun struct {
	InterfaceName          string   `json:"interface_name,omitempty" yaml:"interface_name,omitempty"`
	Address                []string `json:"address,omitempty" yaml:"address,omitempty"`
	MTU                    int      `json:"mtu,omitempty" yaml:"mtu,omitempty"`
	AutoRoute              bool     `json:"auto_route,omitempty" yaml:"auto_route,omitempty"`
	IPRoute2TableIndex     int      `json:"iproute2_table_index,omitempty" yaml:"iproute2_table_index,omitempty"`
	IPRoute2RuleIndex      int      `json:"iproute2_rule_index,omitempty" yaml:"iproute2_rule_index,omitempty"`
	AutoRedirect           bool     `json:"auto_redirect,omitempty" yaml:"auto_redirect,omitempty"`
	AutoRedirectInputMark  string   `json:"auto_redirect_input_mark,omitempty" yaml:"auto_redirect_input_mark,omitempty"`
	AutoRedirectOutputMark string   `json:"auto_redirect_output_mark,omitempty" yaml:"auto_redirect_output_mark,omitempty"`
	StrictRoute            bool     `json:"strict_route,omitempty" yaml:"strict_route,omitempty"`
	RouteAddress           []string `json:"route_address,omitempty" yaml:"route_address,omitempty"`
	RouteExcludeAddress    []string `json:"route_exclude_address,omitempty" yaml:"route_exclude_address,omitempty"`
	RouteAddressSet        []string `json:"route_address_set,omitempty" yaml:"route_address_set,omitempty"`
	RouteExcludeAddressSet []string `json:"route_exclude_address_set,omitempty" yaml:"route_exclude_address_set,omitempty"`
	EndpointIndependentNAT bool     `json:"endpoint_independent_nat,omitempty" yaml:"endpoint_independent_nat,omitempty"`
	UDPTimeout             string   `json:"udp_timeout,omitempty" yaml:"udp_timeout,omitempty"`
	Stack                  string   `json:"stack,omitempty" yaml:"stack,omitempty"`
	IncludeInterface       []string `json:"include_interface,omitempty" yaml:"include_interface,omitempty"`
	ExcludeInterface       []string `json:"exclude_interface,omitempty" yaml:"exclude_interface,omitempty"`
	IncludeUID             []int    `json:"include_uid,omitempty" yaml:"include_uid,omitempty"`
	IncludeUIDRange        []string `json:"include_uid_range,omitempty" yaml:"include_uid_range,omitempty"`
	ExcludeUID             []int    `json:"exclude_uid,omitempty" yaml:"exclude_uid,omitempty"`
	ExcludeUIDRange        []string `json:"exclude_uid_range,omitempty" yaml:"exclude_uid_range,omitempty"`
	IncludeAndroidUser     []int    `json:"include_android_user,omitempty" yaml:"include_android_user,omitempty"`
	IncludePackage         []string `json:"include_package,omitempty" yaml:"include_package,omitempty"`
	ExcludePackage         []string `json:"exclude_package,omitempty" yaml:"exclude_package,omitempty"`
	Platform               Platform `json:"platform,omitempty" yaml:"platform,omitempty"`
	Listen
}

type Platform struct {
	HTTPProxy HTTPProxy `json:"http_proxy,omitempty" yaml:"http_proxy,omitempty"`
}

type HTTPProxy struct {
	Enabled      bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Server       string   `json:"server" yaml:"server"`
	ServerPort   int      `json:"server_port" yaml:"server_port"`
	BypassDomain []string `json:"bypass_domain,omitempty" yaml:"bypass_domain,omitempty"`
	MatchDomain  []string `json:"match_domain,omitempty" yaml:"match_domain,omitempty"`
}
type Inbounds struct {
	Type string `json:"type" yaml:"type"`
	Tag  string `json:"tag" yaml:"tag"`
}