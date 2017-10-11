package nlu

import "github.com/aziule/dev-board/config"

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

// ParserFactory is the main factory interface used to instanciate new Parser implementations
type ParserFactory func(*config.Config) Parser

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
