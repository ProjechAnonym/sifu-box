package models

type V2rayStats struct {
	Enabled   bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Inbounds  []string `json:"inbounds,omitempty" yaml:"inbounds,omitempty"`
	Outbounds []string `json:"outbounds,omitempty" yaml:"outbounds,omitempty"`
	Users     []string `json:"users,omitempty" yaml:"users,omitempty"`
}
type V2rayAPI struct {
	Listen string      `json:"listen,omitempty" yaml:"listen,omitempty"`
	Stats  *V2rayStats `json:"stats,omitempty" yaml:"stats,omitempty"`
}
type ClashAPI struct {
	ExternalController               string   `json:"external_controller,omitempty" yaml:"external_controller,omitempty"`
	ExternalUI                       string   `json:"external_ui,omitempty" yaml:"external_ui,omitempty"`
	ExternalUIDownloadURL            string   `json:"external_ui_download_url,omitempty" yaml:"external_ui_download_url,omitempty"`
	ExternalUIDownloadDetour         string   `json:"external_ui_download_detour,omitempty" yaml:"external_ui_download_detour,omitempty"`
	Secret                           string   `json:"Secret,omitempty" yaml:"secret,omitempty"`
	DefaultMode                      string   `json:"default_mode,omitempty" yaml:"default_mode,omitempty"`
	AccessControlAllowOrigin         []string `json:"access_control_allow_origin,omitempty" yaml:"access_control_allow_origin,omitempty"`
	AccessControlAllowPrivateNetwork bool     `json:"access_control_allow_private_network,omitempty" yaml:"access_control_allow_private_network,omitempty"`
}
type CacheFile struct {
	Enabled      bool   `json:"enabled" yaml:"enabled"`
	Path         string `json:"path,omitempty" yaml:"path,omitempty"`
	CacheID      string `json:"cache_id,omitempty" yaml:"cache_id,omitempty"`
	StoreFakeIP  bool   `json:"store_fakeip,omitempty" yaml:"store_fakeip,omitempty"`
	StoreRdrc    bool   `json:"store_rdrc,omitempty" yaml:"store_rdrc,omitempty"`
	Rdrc_timeout string `json:"rdrc_timeout,omitempty" yaml:"rdrc_timeout,omitempty"`
}
type Experimental struct {
	CacheFile *CacheFile `json:"cache_file,omitempty" yaml:"cache_file,omitempty"`
	ClashAPI  *ClashAPI  `json:"clash_api,omitempty" yaml:"clash_api,omitempty"`
	V2rayAPI  *V2rayAPI  `json:"v2ray_api,omitempty" yaml:"v2ray_api,omitempty"`
}