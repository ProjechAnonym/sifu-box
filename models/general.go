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
	DomainResolver      string   `json:"domain_resolver,omitempty" yaml:"domain_resolver,omitempty"`
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

type TCPBrutal struct {
	Enabled  bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	UpMbps   int  `json:"up_mbps,omitempty" yaml:"up_mbps,omitempty"`
	DownMbps int  `json:"down_mbps,omitempty" yaml:"down_mbps,omitempty"`
}

type Multiplex struct {
	Enabled        bool      `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Protocol       string    `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	MaxConnections int       `json:"max_connections,omitempty" yaml:"max_connections,omitempty"`
	MinStreams     int       `json:"min_streams,omitempty" yaml:"min_streams,omitempty"`
	MaxStreams     int       `json:"max_streams,omitempty" yaml:"max_streams,omitempty"`
	Padding        bool      `json:"padding,omitempty" yaml:"padding,omitempty"`
	Brutal         TCPBrutal `json:"brutal,omitempty" yaml:"brutal,omitempty"`
}

type TransportWS struct {
	Path                string            `json:"path,omitempty" yaml:"path,omitempty"`
	Headers             map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	MaxEarlyData        int               `json:"max_early_data,omitempty" yaml:"max_early_data,omitempty"`
	EarlyDataHeaderName string            `json:"early_data_header_name,omitempty" yaml:"early_data_header_name,omitempty"`
}
type Transport struct {
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	TransportWS `json:",inline" yaml:",inline"`
}

type TLS struct {
	Enabled         bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Insecure        bool     `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	DisableSni      bool     `json:"disable_sni,omitempty" yaml:"disable_sni,omitempty"`
	ServerName      string   `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	ALPN            []string `json:"alpn,omitempty" yaml:"alpn,omitempty"`
	MinVersion      string   `json:"min_version,omitempty" yaml:"min_version,omitempty"`
	MaxVersion      string   `json:"max_version,omitempty" yaml:"max_version,omitempty"`
	CipherSuites    []string `json:"cipher_suites,omitempty" yaml:"cipher_suites,omitempty"`
	Certificate     [][]byte `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	CertificatePath string   `json:"certificate_path,omitempty" yaml:"certificate_path,omitempty"`
	Key             [][]byte `json:"key,omitempty" yaml:"key,omitempty"`
	KeyPath         string   `json:"key_path,omitempty" yaml:"key_path,omitempty"`
	ACME            *ACME    `json:"acme,omitempty" yaml:"acme,omitempty"`
	ECH             *ECH     `json:"ech,omitempty" yaml:"ech,omitempty"`
	Reality         *Reality `json:"reality,omitempty" yaml:"reality,omitempty"`
}

type ACME struct {
	Domain                  []string         `json:"domain,omitempty" yaml:"domain,omitempty"`
	DataDirectory           string           `json:"data_directory,omitempty" yaml:"data_directory,omitempty"`
	DefaultServerName       string           `json:"default_server_name,omitempty" yaml:"default_server_name,omitempty"`
	Email                   string           `json:"email,omitempty" yaml:"email,omitempty"`
	Provider                string           `json:"provider,omitempty" yaml:"provider,omitempty"`
	DisableHTTPChallenge    bool             `json:"disable_http_challenge,omitempty" yaml:"disable_http_challenge,omitempty"`
	DisableTLSALPNChallenge bool             `json:"disable_tls_alpn_challenge,omitempty" yaml:"disable_tls_alpn_challenge,omitempty"`
	AlternativeHTTPPort     int              `json:"alternative_http_port,omitempty" yaml:"alternative_http_port,omitempty"`
	AlternativeTLSPort      int              `json:"alternative_tls_port,omitempty" yaml:"alternative_tls_port,omitempty"`
	ExternalAccount         *ExternalAccount `json:"external_account,omitempty" yaml:"external_account,omitempty"`
	DNS01Challenge          *DNS01Challenge  `json:"dns01_challenge,omitempty" yaml:"dns01_challenge,omitempty"`
}

type ExternalAccount struct {
	KeyID  string `json:"key_id,omitempty" yaml:"key_id,omitempty"`
	MacKey string `json:"mac_key,omitempty" yaml:"mac_key,omitempty"`
}

type DNS01Challenge struct {
	// 省略具体字段，根据实际需求补充
}

type ECH struct {
	Enabled                     bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	PQSignatureSchemesEnabled   bool     `json:"pq_signature_schemes_enabled,omitempty" yaml:"pq_signature_schemes_enabled,omitempty"`
	DynamicRecordSizingDisabled bool     `json:"dynamic_record_sizing_disabled,omitempty" yaml:"dynamic_record_sizing_disabled,omitempty"`
	Key                         [][]byte `json:"key,omitempty" yaml:"key,omitempty"`
	KeyPath                     string   `json:"key_path,omitempty" yaml:"key_path,omitempty"`
	Config                      string   `json:"config,omitempty" yaml:"config,omitempty"`
	ConfigPath                  string   `json:"config_path,omitempty" yaml:"config_path,omitempty"`
}

type Reality struct {
	Enabled           bool      `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Handshake         Handshake `json:"handshake" yaml:"handshake"`
	PublicKey         string    `json:"public_key" yaml:"public_key"`
	PrivateKey        string    `json:"private_key" yaml:"private_key"`
	ShortID           []string  `json:"short_id" yaml:"short_id"`
	MaxTimeDifference string    `json:"max_time_difference" yaml:"max_time_difference"`
}

type Handshake struct {
	Server     string
	ServerPort int
	Dial       `json:",inline" yaml:",inline"`
}