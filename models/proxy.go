package models

type Ruleset struct {
	Id              uint64 `json:"id,omitempty" gorm:"primaryKey" yaml:"-"`
	Type            string `json:"type" gorm:"not null;type:varchar(255)" yaml:"type"`
	Path            string `json:"path,omitempty" gorm:"path" yaml:"path,omitempty"`
	Url             string `json:"url,omitempty" gorm:"url" yaml:"url,omitempty"`
	Format          string `json:"format" gorm:"not null;type:varchar(255)" yaml:"format"`
	Tag             string `json:"tag" gorm:"not null;type:varchar(255);unique" yaml:"tag"`
	Download_detour string `json:"download_detour,omitempty" gorm:"not null;type:varchar(255)" yaml:"download_detour,omitempty"`
	Update_interval string `json:"update_interval,omitempty" gorm:"not null;type:varchar(255)" yaml:"update_interval,omitempty"`
	Label           string `json:"label,omitempty" gorm:"type:varchar(255);" yaml:"label,omitempty"`
	China           bool   `json:"china,omitempty" gorm:"not null;bool" yaml:"china"`
	DnsRule         string `json:"dnsRule,omitempty" gorm:"not null;bool" yaml:"dnsRule,omitempty"`
}
type Provider struct {
	Id     uint64 `json:"id" gorm:"primaryKey" yaml:"-"`
	Name   string `json:"name" gorm:"not null;type:varchar(255);unique" yaml:"name"`
	Proxy  bool   `json:"proxy" gorm:"not null;bool" yaml:"proxy"`
	Path   string `json:"path" gorm:"not null;type:varchar(255)" yaml:"path"`
	Remote bool   `json:"remote" gorm:"not null;bool" yaml:"remote"`
}
type Proxy struct {
	Providers []Provider `yaml:"providers"`
	Rulesets  []Ruleset  `yaml:"rulesets"`
}