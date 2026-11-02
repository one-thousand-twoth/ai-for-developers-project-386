package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// MeetingType is independent of the generated HTTP contract.
type MeetingType struct{ ent.Schema }

func (MeetingType) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(uuid.NewString).Immutable(),
		field.String("name"),
		field.String("name_key").Unique(),
		field.Int("duration_minutes").Min(1).Max(480),
		field.String("description").Optional().Nillable(),
	}
}
