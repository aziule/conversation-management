package nlp

// Parser is the main interface for parsing raw data and returning parsed data
type Parser interface {
	ParseNlpData([]byte) (*ParsedData, error)
}

// ParsedData represents intents and entities as understood after using NLP services
type ParsedData struct {
	Intent   *Intent
	Entities []Entity
}

// NewParsedData is the constructor method for ParsedData
func NewParsedData(intent *Intent, entities []Entity) *ParsedData {
	return &ParsedData{
		Intent:   intent,
		Entities: entities,
	}
}
