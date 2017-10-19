package nlp

import (
	"errors"
)

var (
	ErrInvalidDataType       = errors.New("Invalid data type provided")
	ErrInvalidEntityDataType = errors.New("Invalid entity data type provided")
)

// Parser is the main interface for parsing raw data and returning parsed data
type Parser interface {
	ParseNlpData([]byte) (*ParsedData, error)
}

// ParsedData represents intents and entities as understood after using NLP services
type ParsedData struct {
	intent   *Intent
	entities []Entity
}

// NewParsedData is the constructor method for ParsedData
func NewParsedData(intent *Intent, entities []Entity) *ParsedData {
	return &ParsedData{
		intent:   intent,
		entities: entities,
	}
}

// Getters
func (parsedData *ParsedData) Intent() *Intent    { return parsedData.intent }
func (parsedData *ParsedData) Entities() []Entity { return parsedData.entities }
