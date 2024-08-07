package models

type Utls struct {
	Enabled     bool   `json:"enabled,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}
type Tls struct {
	Enabled     bool   `json:"enabled"`
	Insecure    bool   `json:"insecure"`
	Server_name string `json:"server_name"`
	Utls        *Utls  `json:"utls,omitempty"`
}