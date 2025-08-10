package schema

import (
	"sifu-box/singbox"

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
		field.JSON("dns", singbox.DNS{}).Optional(), field.JSON("log", singbox.Log{}).Optional(),
		field.JSON("route", singbox.Route{}).Optional(), field.JSON("inbounds", []map[string]any{}).Optional(),
		field.JSON("outbound_groups", []singbox.OutboundGroup{}).Optional(), field.JSON("ntp", singbox.Ntp{}).Optional(),
		field.JSON("experiment", singbox.Experiment{}).Optional(), field.Strings("providers").Optional(),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return nil
}
