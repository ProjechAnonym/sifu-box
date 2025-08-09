package schema

import (
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
		field.JSON("dns", []map[string]any{}),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return nil
}
