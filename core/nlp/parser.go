package nlp

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var ErrParserNotFound = errors.New("Parser not found")

// parserBuilders stores the available Parser builders
var parserBuilders = make(map[string]ParserBuilder)

// ParserBuilder is the interface describing a builder for Parser
type ParserBuilder func() Parser

// RegisterParserBuilder adds a new ParserBuilder to the list of available builders
func RegisterParserBuilder(name string, builder ParserBuilder) {
	_, registered := parserBuilders[name]

	if registered {
		log.WithField("name", name).Warning("ParserBuilder already registered, ignoring")
	}

	parserBuilders[name] = builder
}

// NewParser tries to create a Parser using the available builders.
// Returns ErrParserNotFound if the parser builder isn't found.
func NewParser(name string) (Parser, error) {
	parserBuilder, ok := parserBuilders[name]

	if !ok {
		return nil, ErrParserNotFound
	}

	return parserBuilder(), nil
}

// Parser is the main interface for parsing raw data and returning parsed data
type Parser interface {
	ParseNlpData([]byte) (*ParsedData, error)
}

// ParsedData represents intents and entities as understood after using NLP services
type ParsedData struct {
	Intent   *Intent           `bson:"intent"`
	Entities []*EntityWithType `bson:"entities"`
}

// NewParsedData is the constructor method for ParsedData
func NewParsedData(intent *Intent, entities []Entity) *ParsedData {
	var entitiesWithType []*EntityWithType

	for _, e := range entities {
		entitiesWithType = append(entitiesWithType, &EntityWithType{
			Type:   e.Type(),
			Entity: e,
		})
	}

	return &ParsedData{
		Intent:   intent,
		Entities: entitiesWithType,
	}
}
