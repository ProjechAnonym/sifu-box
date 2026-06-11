package nodes

type TLS struct {
	Enabled    bool     `json:"enabled" yaml:"enabled"`
	Insecure   bool     `json:"insecure" yaml:"insecure"`
	ServerName string   `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	Alpn       []string `json:"alpn,omitempty" yaml:"alpn,omitempty"`
	Reality    Reality  `json:"reality" yaml:"reality"`
	Utls       Utls     `json:"utls" yaml:"utls"`
}

type Transport struct {
	Type    string         `json:"type" yaml:"type"`
	Headers map[string]any `json:"headers" yaml:"headers"`
	Path    string         `json:"path,omitempty" yaml:"path,omitempty"`
}

type Reality struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	PublicKey string `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	ShortId   string `json:"short_id,omitempty" yaml:"short_id,omitempty"`
}

type Utls struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Fingerprint string `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
}
