package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Appointment struct {
	ent.Schema
}

func (Appointment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("email").
			NotEmpty(),
		field.String("phone").
			NotEmpty(),
		field.Enum("type").NamedValues(
			"Goldankauf", "goldankauf",
			"Trauringe", "trauringe",
			"Verlobungsringe", "verlobungsringe",
			"Ohrlochstechen", "ohrlochstechen",
			"Sonstiges", "sonstiges",
		),
		field.String("delkey").
			NotEmpty(),
		field.Time("start_time"),
		field.Time("end_time"),
		field.String("description"),
	}
}

func (Appointment) Edges() []ent.Edge {
	return nil
}
