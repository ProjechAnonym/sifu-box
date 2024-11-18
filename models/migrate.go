package models

type Migrate struct {
	Hosts     []Host              `yaml:"hosts"`
	Providers []Provider          `yaml:"providers"`
	Rulesets  []Ruleset           `yaml:"rulesets"`
	Templates map[string]Template `yaml:"templates"`
}