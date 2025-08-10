package singbox

type DNS struct {
	Final             string           `json:"final" yaml:"final"`
	Strategy          string           `json:"strategy" yaml:"strategy"`
	Disable_cache     bool             `json:"disable_cache,omitempty" yaml:"disable_cache,omitempty"`
	Disable_expire    bool             `json:"disable_expire,omitempty" yaml:"disable_expire,omitempty"`
	Independent_cache bool             `json:"independent_cache,omitempty" yaml:"independent_cache,omitempty"`
	Cache_capacity    int              `json:"cache_capacity,omitempty" yaml:"cache_capacity,omitempty"`
	Reverse_mapping   bool             `json:"reverse_mapping,omitempty" yaml:"reverse_mapping,omitempty"`
	Client_subnet     string           `json:"client_subnet,omitempty" yaml:"client_subnet,omitempty"`
	Rules             []map[string]any `json:"rules,omitempty" yaml:"rules,omitempty"`
	Servers           []map[string]any `json:"servers,omitempty" yaml:"servers,omitempty"`
}
