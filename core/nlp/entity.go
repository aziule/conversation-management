// Package nlp provides objects to store, handle and manipulate NLP data.
package nlp

// DateTimeGranularity represents the available date time granularity
// for a given time value.
type DateTimeGranularity string

// DataTypeMap maps a data type from the entity name to its data type
type DataTypeMap map[string]EntityType

// @todo: use only an EntityType?
// DataType represents the kind of entity we are dealing with
type EntityType string

const (
	GranularityHour DateTimeGranularity = "hour"
	GranularityDay  DateTimeGranularity = "day"

	IntEntity      EntityType = "int"
	DateTimeEntity EntityType = "datetime"
)

// Entity is the struct that represents a base entity
type Entity struct {
	Name string     `bson:"name"`
	Type EntityType `bson:"type"`
}

//// SetBSON converts the EntityWithType's entity, of type interface, to the corresponding Entity struct.
//// We have moved the SetBSON method here for now, as it was faster to develop. In the future,
//// we may think about extracting it completely to the "conversation/mongo" package.
//func (e *EntityWithType) SetBSON(raw bson.Raw) error {
//	decodedType := struct {
//		Type DataType `bson:"type"`
//	}{}
//
//	raw.Unmarshal(&decodedType)
//
//	// @todo: handle datetimes
//	switch decodedType.Type {
//	case DataTypeInt:
//		decodedEntity := struct {
//			Entity *IntEntity `bson:"entity"`
//		}{}
//		raw.Unmarshal(&decodedEntity)
//		e.Entity = decodedEntity.Entity
//		break
//	default:
//		// @todo: handle and log
//		log.WithField("type", decodedType.Type).Error("Unhandled entity type to decode")
//		return errors.New("Unhandled case")
//	}
//
//	e.Type = decodedType.Type
//
//	return nil
//}
