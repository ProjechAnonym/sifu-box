package singbox

type Experiment struct {
	Clash_api Clash_api `json:"clash_api" yaml:"clash_api"`
}
type Clash_api struct {
	External_controller                  string   `json:"external_controller" yaml:"external_controller"`
	External_ui                          string   `json:"external_ui,omitempty" yaml:"external_ui,omitempty"`
	External_ui_download_url             string   `json:"external_ui_download_url,omitempty" yaml:"external_ui_download_url,omitempty"`
	External_ui_download_detour          string   `json:"external_ui_download_detour,omitempty" yaml:"external_ui_download_detour,omitempty"`
	Secret                               string   `json:"secret,omitempty" yaml:"secret,omitempty"`
	Default_mode                         string   `json:"default_mode,omitempty" yaml:"default_mode,omitempty"`
	Access_control_allow_origin          []string `json:"access_control_allow_origin,omitempty" yaml:"access_control_allow_origin,omitempty"`
	Access_control_allow_private_network bool     `json:"access_control_allow_private_network,omitempty" yaml:"access_control_allow_private_network,omitempty"`
}
