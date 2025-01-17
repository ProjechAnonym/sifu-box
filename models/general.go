package models

type Dial struct {
	Detour              string   `json:"detour,omitempty" yaml:"detour,omitempty"`
	BindInterface       string   `json:"bind_interface,omitempty" yaml:"bind_interface,omitempty"`
	Inet4BindAddress    string   `json:"inet4_bind_address,omitempty" yaml:"inet4_bind_address,omitempty"`
	Inet6BindAddress    string   `json:"inet6_bind_address,omitempty" yaml:"inet6_bind_address,omitempty"`
	RoutingMark         int      `json:"routing_mark,omitempty" yaml:"routing_mark,omitempty"`
	ReuseAddr           bool     `json:"reuse_addr,omitempty" yaml:"reuse_addr,omitempty"`
	ConnectTimeout      string   `json:"connect_timeout,omitempty" yaml:"connect_timeout,omitempty"`
	TcpFastOpen         bool     `json:"tcp_fast_open,omitempty" yaml:"tcp_fast_open,omitempty"`
	TcpMultiPath        bool     `json:"tcp_multi_path,omitempty" yaml:"tcp_multi_path,omitempty"`
	UdpFragment         bool     `json:"udp_fragment,omitempty" yaml:"udp_fragment,omitempty"`
	DomainStrategy      string   `json:"domain_strategy,omitempty" yaml:"domain_strategy,omitempty"`
	NetworkStrategy     string   `json:"network_strategy,omitempty" yaml:"network_strategy,omitempty"`
	NetworkType         []string `json:"network_type,omitempty" yaml:"network_type,omitempty"`
	FallbackNetworkType []string `json:"fallback_network_type,omitempty" yaml:"fallback_network_type,omitempty"`
	FallbackDelay       string   `json:"fallback_delay,omitempty" yaml:"fallback_delay,omitempty"`
}
type Listen struct {
	Listen                    string `json:"listen,omitempty" yaml:"listen,omitempty"`
	ListenPort                int    `json:"listen_port,omitempty" yaml:"listen_port,omitempty"`
	TCPFastOpen               bool   `json:"tcp_fast_open,omitempty" yaml:"tcp_fast_open,omitempty"`
	TCPMultiPath              bool   `json:"tcp_multi_path,omitempty" yaml:"tcp_multi_path,omitempty"`
	UDPFragment               bool   `json:"udp_fragment,omitempty" yaml:"udp_fragment,omitempty"`
	UDPTimeout                string `json:"udp_timeout,omitempty" yaml:"udp_timeout,omitempty"`
	Detour                    string `json:"detour,omitempty" yaml:"detour,omitempty"`
	Sniff                     bool   `json:"sniff,omitempty" yaml:"sniff,omitempty"`
	SniffOverrideDestination  bool   `json:"sniff_override_destination,omitempty" yaml:"sniff_override_destination,omitempty"`
	SniffTimeout              string `json:"sniff_timeout,omitempty" yaml:"sniff_timeout,omitempty"`
	DomainStrategy            string `json:"domain_strategy,omitempty" yaml:"domain_strategy,omitempty"`
	UDPDisableDomainUnmapping bool   `json:"udp_disable_domain_unmapping,omitempty" yaml:"udp_disable_domain_unmapping,omitempty"`
}