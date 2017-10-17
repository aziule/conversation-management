package nlp

import "time"

type DateTimeGranularity string

// The available data types
const (
	GranularityHour DateTimeGranularity = "hour"
	GranularityDay  DateTimeGranularity = "day"
)

type Entity struct {
	name       string
	confidence float32
}

type StringEntity struct {
	*Entity
	Value string
}

type NumberEntity struct {
	*Entity
	Value int
}

type SingleDateTimeEntity struct {
	*Entity
	Value       time.Time
	Granularity DateTimeGranularity
}

type IntervalDateTimeEntity struct {
	*Entity
	FromValue       time.Time
	FromGranularity DateTimeGranularity
	ToValue         time.Time
	ToGranularity   DateTimeGranularity
}

func newEntity(name string, confidence float32) *Entity {
	return &Entity{
		name:       name,
		confidence: confidence,
	}
}

// NewStringEntity is the factory method for StringEntity
func NewStringEntity(name string, confidence float32, value string) *StringEntity {
	return &StringEntity{newEntity(name, confidence), value}
}

// NewStringEntity is the factory method for StringEntity
func NewNumberEntity(name string, confidence float32, value int) *NumberEntity {
	return &NumberEntity{newEntity(name, confidence), value}
}

// NewSingleDateTimeEntity is the factory method for SingleDateTimeEntity
func NewSingleDateTimeEntity(name string, confidence float32, value time.Time, granularity DateTimeGranularity) *SingleDateTimeEntity {
	return &SingleDateTimeEntity{newEntity(name, confidence), value, granularity}
}

// NewIntervalDateTimeEntity is the factory method for IntervalDateTimeEntity
func NewIntervalDateTimeEntity(name string, confidence float32, fromValue, toValue time.Time, fromGranularity, toGranularity DateTimeGranularity) *IntervalDateTimeEntity {
	return &IntervalDateTimeEntity{newEntity(name, confidence), fromValue, fromGranularity, toValue, toGranularity}
}

// Getters
func (entity *Entity) Name() string        { return entity.name }
func (entity *Entity) Confidence() float32 { return entity.confidence }
