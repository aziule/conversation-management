package nlp

import (
	"time"
)

type DateTimeGranularity string

type DataTypeMap map[string]DataType

type DataType string

const (
	// Data types that we handle
	DataTypeInt      DataType = "int"
	DataTypeDateTime DataType = "datetime"
	DataTypeIntent   DataType = "intent"

	GranularityHour DateTimeGranularity = "hour"
	GranularityDay  DateTimeGranularity = "day"
)

type Entity interface {
	Name() string
	Confidence() float32
	Type() DataType
}

type entity struct {
	name       string
	confidence float32
}

type IntEntity struct {
	*entity
	Value int
}

func (e *IntEntity) Name() string        { return e.name }
func (e *IntEntity) Confidence() float32 { return e.confidence }
func (e *IntEntity) Type() DataType      { return DataTypeInt }

type SingleDateTimeEntity struct {
	*entity
	Value       time.Time
	Granularity DateTimeGranularity
}

func (e *SingleDateTimeEntity) Name() string        { return e.name }
func (e *SingleDateTimeEntity) Confidence() float32 { return e.confidence }
func (e *SingleDateTimeEntity) Type() DataType      { return DataTypeDateTime }

type IntervalDateTimeEntity struct {
	*entity
	FromValue       time.Time
	FromGranularity DateTimeGranularity
	ToValue         time.Time
	ToGranularity   DateTimeGranularity
}

func (e *IntervalDateTimeEntity) Name() string        { return e.name }
func (e *IntervalDateTimeEntity) Confidence() float32 { return e.confidence }
func (e *IntervalDateTimeEntity) Type() DataType      { return DataTypeDateTime }

// newEntity creates a new base entity
func newEntity(name string, confidence float32) *entity {
	return &entity{
		name:       name,
		confidence: confidence,
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
