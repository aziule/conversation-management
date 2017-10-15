package nlu

import (
	"github.com/aziule/conversation-management/core"
	"errors"
)

var (
	ErrInvalidDataType = errors.New("Invalid data type provided")
	ErrInvalidEntityDataType = errors.New("Invalid entity data type provided")

	parserFactories = make(map[ParsingService]ParserFactory) // The list of available factories
)

// ParsingService represents any available ParsingService to use for NLU
type ParsingService string

// Parser is the interface responsible for parsing any given text using NLU
type Parser interface {
	ParseData(data interface{}) (*ParsedData, error)
}

// What is returned by our NLU services after parsing a given text
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

// ParserFactory is the main factory func used to instantiate new Parser implementations
type ParserFactory func(*core.Config) Parser

// RegisterFactory allows us to register factory methods for creating new parsers
func RegisterFactory(ParsingService ParsingService, factory ParserFactory) {
	if factory == nil {
		return
	}

	_, registered := parserFactories[ParsingService]

	if registered {
		return
	}

	parserFactories[ParsingService] = factory
}

// NewParserFromConfig instantiates the correct Parser service given the configuration
func NewParserFromConfig(config *core.Config) Parser {
	return parserFactories[ParsingService(config.NluService)](config)
}
