package nlp

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
