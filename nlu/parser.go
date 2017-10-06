package nlu

// The parser service, responsible for parsing any given text using NLU
type Parser struct {
}

// Parse a given text and try to understand it using NLU
func (parser *Parser) ParseText(text string) (*ParsedText, error) {
	return &ParsedText{
		originalText: text,
		intent: &Intent{
			uuid: "intent_name",
		},
		entity: &Entity{
			uuid: "entity_name",
			value: "entity_value",
		},
	}, nil
}

// What is returned by our NLU services after parsing a given text
type ParsedText struct {
	originalText string
	intent *Intent
	entity *Entity
}

// Getters
func (parsedText *ParsedText) OriginalText() string { return parsedText.originalText }
func (parsedText *ParsedText) Intent() *Intent { return parsedText.intent }
func (parsedText *ParsedText) Entity() *Entity { return parsedText.entity }
