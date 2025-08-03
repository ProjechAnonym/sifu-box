package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Provider holds the schema definition for the Provider entity.
type Provider struct {
	ent.Schema
}

// Fields of the Provider.
func (Provider) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(30).Unique(),
		field.String("path").NotEmpty().MaxLen(1000),
		field.JSON("nodes", []map[string]any{}).Optional(), field.Bool("remote"),
	}
}

// Edges of the Provider.
func (Provider) Edges() []ent.Edge {
	return nil
}
