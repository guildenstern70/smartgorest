package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Phone holds the schema definition for the Phone entity.
type Phone struct {
	ent.Schema
}

// Fields of the Phone.
func (Phone) Fields() []ent.Field {

	return []ent.Field{
		field.String("prefix").
			Default("+39"),
		field.String("number").
			Default(""),
		field.Enum("kind").
			Values("cellular", "home", "work"),
	}
}

// Edges of the Phone.
func (Phone) Edges() []ent.Edge {
	return nil
}
