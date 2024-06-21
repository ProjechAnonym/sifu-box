package singbox

type Grpc struct {
	Type                  string `json:"type"`
	Service_name          string `json:"service_name"`
	Idle_timeout          string `json:"idle_timeout"`
	Ping_timeout          string `json:"ping_timeout"`
	Permit_without_stream bool   `json:"permit_without_stream"`
}

type Websocket struct {
	Type                   string            `json:"type"`
	Path                   string            `json:"path"`
	Headers                map[string]string `json:"headers"`
	Early_data_header_name string            `json:"early_data_header_name"`
}