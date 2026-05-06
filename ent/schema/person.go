package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Person holds the schema definition for the Person entity.
type Person struct {
	ent.Schema
}

// Fields of the Person.
func (Person) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").
			Default("?"),
		field.String("last_name").
			Default("?"),
		field.Int("age").
			Default(0),
		field.String("email").
			Default("?"),
	}
}

// Edges of the Person.
func (Person) Edges() []ent.Edge {

	return []ent.Edge{
		edge.To("phones", Phone.Type),
	}

}
