package nlp

import (
	"github.com/aziule/conversation-management/core/utils"
)

const nlpParserBuilderPrefix = "nlp_parser_"

// RegisterRepositoryBuilder registers a new service builder using a package-level prefix
func RegisterParserBuilder(name string, builder utils.ServiceBuilder) {
	utils.RegisterServiceBuilder(nlpParserBuilderPrefix+name, builder)
}

// NewParser tries to create a Parser using the available builders.
// Returns ErrParserNotFound if the parser builder isn't found.
func NewParser(name string, conf utils.BuilderConf) (Parser, error) {
	parserBuilder, err := utils.GetServiceBuilder(nlpParserBuilderPrefix + name)

	if err != nil {
		return nil, err
	}

	parser, err := parserBuilder(conf)

	if err != nil {
		return nil, err
	}

	return parser.(Parser), nil
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
