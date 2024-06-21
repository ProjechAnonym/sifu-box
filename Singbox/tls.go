package singbox

type Tls_config struct {
	Enabled     bool   `json:"enabled"`
	Insecure    bool   `json:"insecure"`
	Server_name string `json:"server_name"`
	Disable_sni bool   `json:"disable_sni"`
}