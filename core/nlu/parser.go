package nlu

import (
	"github.com/aziule/conversation-management/core"
)

var (
	parserFactories = make(map[string]ParserFactory) // The list of available factories
)

// Parser is the interface responsible for parsing any given text using NLU
type Parser interface {
	ParseText(text string) (*ParsedText, error)
}

// What is returned by our NLU services after parsing a given text
type ParsedText struct {
	originalText string
	intent       *Intent
	entities     []*Entity
}

// NewParsedText is the constructor method for ParsedText
func NewParsedText(originalText string, intent *Intent, entities []*Entity) *ParsedText {
	return &ParsedText{
		originalText: originalText,
		intent:       intent,
		entities:     entities,
	}
}

// Getters
func (parsedText *ParsedText) OriginalText() string { return parsedText.originalText }
func (parsedText *ParsedText) Intent() *Intent      { return parsedText.intent }
func (parsedText *ParsedText) Entities() []*Entity  { return parsedText.entities }

// ParserFactory is the main factory func used to instantiate new Parser implementations
type ParserFactory func(*core.Config) Parser

// RegisterFactory allows us to register factory methods for creating new parsers
func RegisterFactory(name string, factory ParserFactory) {
	if factory == nil {
		return
	}

	_, registered := parserFactories[name]

	if registered {
		return
	}

	parserFactories[name] = factory
}

// NewParserFromConfig instantiates the correct Parser service given the configuration
func NewParserFromConfig(config *core.Config) Parser {
	return parserFactories[config.NluService](config)
}
