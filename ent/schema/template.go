package schema

import (
	config1 "sifu-box/config"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(30).Unique(),
		field.JSON("dns", config1.DNS{}).Optional(), field.JSON("log", config1.Log{}).Optional(),
		field.JSON("route", config1.Route{}).Optional(), field.JSON("inbounds", []config1.Inbound{}).Optional(),
		field.JSON("outbound_groups", []config1.OutboundGroup{}).Optional(), field.JSON("ntp", config1.Ntp{}).Optional(),
		field.JSON("experiment", config1.Experiment{}).Optional(), field.Strings("providers").Optional(),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return nil
}
