package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Blog holds the schema definition for the Blog entity.
type Blog struct {
	ent.Schema
}

// Fields of the Blog.
func (Blog) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("content"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
	}
}

// Edges of the Blog.
func (Blog) Edges() []ent.Edge {
	return nil
}
