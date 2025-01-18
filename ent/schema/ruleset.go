package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// RuleSet holds the schema definition for the RuleSet entity.
type RuleSet struct {
	ent.Schema
}

// Fields of the RuleSet.
func (RuleSet) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag").NotEmpty().MaxLen(30).Unique(),
		field.String("type").NotEmpty().MaxLen(10),
		field.String("path").NotEmpty().MaxLen(100).Unique(),
		field.String("format").NotEmpty().MaxLen(10),
		field.String("label").NotEmpty().MaxLen(30),
		field.String("download_detour").NotEmpty().MaxLen(30),
		field.String("update_interval").NotEmpty().MaxLen(10),
		field.String("outbound").NotEmpty().MaxLen(30),
		field.Bool("china"),
	}
}

// Edges of the RuleSet.
func (RuleSet) Edges() []ent.Edge {
	return nil
}
