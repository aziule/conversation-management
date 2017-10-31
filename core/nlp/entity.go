package nlp

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// DateTimeGranularity represents the available date time granularity
// for a given time value.
type DateTimeGranularity string

// DataTypeMap maps a data type from the entity's name to its data type
type DataTypeMap map[string]DataType

// @todo: use only an EntityType?
// DataType represents the possible data types contained in the entities
type DataType string

const (
	// Data types that we handle
	DataTypeInt      DataType = "int"
	DataTypeDateTime DataType = "datetime"
	DataTypeIntent   DataType = "intent"

	GranularityHour DateTimeGranularity = "hour"
	GranularityDay  DateTimeGranularity = "day"
)

// Entity is the interface grouping any entity's behaviour and must-have fields
type Entity interface {
	Name() string
	Confidence() float32
	Type() DataType
}

// EntityWithType is the struct grouping an entity along with its type.
// Its purpose is to be used with mongodb so that we can unmarshal any Message interface easily.
type EntityWithType struct {
	Type   DataType `bson:"type"`
	Entity Entity   `bson:"entity"`
}

type entity struct {
	Name       string  `bson:"name"`
	Confidence float32 `bson:"confidence"`
}

type IntEntity struct {
	entity `bson:",inline"`
	Value  int `bson:"value"`
}

func (e *IntEntity) Name() string        { return e.entity.Name }
func (e *IntEntity) Confidence() float32 { return e.entity.Confidence }
func (e *IntEntity) Type() DataType      { return DataTypeInt }

type SingleDateTimeEntity struct {
	entity
	Value       time.Time           `bson:"value"`
	Granularity DateTimeGranularity `bson:"granularity"`
}

func (e *SingleDateTimeEntity) Name() string        { return e.entity.Name }
func (e *SingleDateTimeEntity) Confidence() float32 { return e.entity.Confidence }
func (e *SingleDateTimeEntity) Type() DataType      { return DataTypeDateTime }

type IntervalDateTimeEntity struct {
	entity
	FromValue       time.Time           `bson:"from"`
	FromGranularity DateTimeGranularity `bson:"from_granularity"`
	ToValue         time.Time           `bson:"to"`
	ToGranularity   DateTimeGranularity `bson:"to_granularity"`
}

func (e *IntervalDateTimeEntity) Name() string        { return e.entity.Name }
func (e *IntervalDateTimeEntity) Confidence() float32 { return e.entity.Confidence }
func (e *IntervalDateTimeEntity) Type() DataType      { return DataTypeDateTime }

// newEntity creates a new base entity
func newEntity(name string, confidence float32) entity {
	return entity{
		Name:       name,
		Confidence: confidence,
	}
}

// NewIntEntity is the factory method for StringEntity
func NewIntEntity(name string, confidence float32, value int) *IntEntity {
	return &IntEntity{newEntity(name, confidence), value}
}

// NewSingleDateTimeEntity is the factory method for SingleDateTimeEntity
func NewSingleDateTimeEntity(name string, confidence float32, value time.Time, granularity DateTimeGranularity) *SingleDateTimeEntity {
	return &SingleDateTimeEntity{newEntity(name, confidence), value, granularity}
}

// NewIntervalDateTimeEntity is the factory method for IntervalDateTimeEntity
func NewIntervalDateTimeEntity(name string, confidence float32, fromValue, toValue time.Time, fromGranularity, toGranularity DateTimeGranularity) *IntervalDateTimeEntity {
	return &IntervalDateTimeEntity{newEntity(name, confidence), fromValue, fromGranularity, toValue, toGranularity}
}

// SetBSON converts the EntityWithType's entity, of type interface, to the corresponding Entity struct.
// We have moved the SetBSON method here for now, as it was faster to develop. In the future,
// we may think about extracting it completely to the "conversation/mongo" package.
func (e *EntityWithType) SetBSON(raw bson.Raw) error {
	decodedType := struct {
		Type DataType `bson:"type"`
	}{}

	raw.Unmarshal(&decodedType)

	// @todo: handle datetimes
	switch decodedType.Type {
	case DataTypeInt:
		decodedEntity := struct {
			Entity *IntEntity `bson:"entity"`
		}{}
		raw.Unmarshal(&decodedEntity)
		e.Entity = decodedEntity.Entity
		break
	default:
		// @todo: handle and log
		log.WithField("type", decodedType.Type).Error("Unhandled entity type to decode")
		return errors.New("Unhandled case")
	}

	e.Type = decodedType.Type

	return nil
}
