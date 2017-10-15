package nlu

import (
	"errors"
)

var (
	ErrInvalidDataType = errors.New("Invalid data type provided")
	ErrInvalidEntityDataType = errors.New("Invalid entity data type provided")
)

// ParsedData represents intents and entities as understood after using NLU services
type ParsedData struct {
	intent       *Intent
	entities     []*Entity
}

// NewParsedData is the constructor method for ParsedData
func NewParsedData(intent *Intent, entities []*Entity) *ParsedData {
	return &ParsedData{
		intent:       intent,
		entities:     entities,
	}
}

// Getters
func (parsedData *ParsedData) Intent() *Intent      { return parsedData.intent }
func (parsedData *ParsedData) Entities() []*Entity  { return parsedData.entities }
