package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("uuid").MaxLen(32).Optional(), field.Bool("updated").Optional(),
		field.Strings("templates").Optional(),
	}
}

// Edges of the Provider.
func (Provider) Edges() []ent.Edge {
	return nil
}
func (Provider) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "path").Unique(),
	}
}
