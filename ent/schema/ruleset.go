package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Rulesets holds the schema definition for the Rulesets entity.
type Ruleset struct {
	ent.Schema
}

// Fields of the Rulesets.
func (Ruleset) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(30).Unique(),
		field.String("path").NotEmpty().MaxLen(1000), field.Bool("remote"),
		field.Bool("binary"), field.String("download_detour").Optional(),
		field.String("update_interval").Optional().MaxLen(30), field.Strings("templates").Optional(),
	}
}

// Edges of the Rulesets.
func (Ruleset) Edges() []ent.Edge {
	return nil
}
func (Ruleset) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "path").Unique(),
	}
}
