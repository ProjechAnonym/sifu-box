package nodes

type TLS struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Insecure   bool   `json:"insecure" yaml:"insecure"`
	ServerName string `json:"server_name" yaml:"server_name"`
}

type Transport struct {
	Type    string            `json:"type" yaml:"type"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Path    string            `json:"path" yaml:"path"`
}

func (t *Transport) SetWsYaml() {
}
