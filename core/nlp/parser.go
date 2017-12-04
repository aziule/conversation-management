package nlp

import (
	"github.com/aziule/conversation-management/core/utils"
	"time"
)

const nlpParserBuilderPrefix = "nlp_parser_"

// DateTimeGranularity represents the available date time granularity
// for a given time value.
type DateTimeGranularity string

// ParsedIntent represents an intent parsed from any given sentence
type ParsedIntent struct {
	Intent     *Intent `json:",inline"`
	Confidence float32 `json:"confidence"`
}

// ParsedEntity represents an entity parsed from any given sentence.
// As there may be any kind of entities, we use an interface to store the
// entities' formatted data.
type ParsedEntity struct {
	Entity     *Entity     `json:",inline"`
	Role       string      `json:"role"`
	Confidence float32     `json:"confidence"`
	Data       interface{} `json:"data"`
}

type parsedSingleDateTime struct {
	Date        time.Time
	Granularity DateTimeGranularity
}

type parsedDateTimeInterval struct {
	From *parsedSingleDateTime
	To   *parsedSingleDateTime
}

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
	Intent   *ParsedIntent   `bson:"intent"`
	Entities []*ParsedEntity `bson:"entities"`
}

func NewParsedIntent(name string) *ParsedIntent {
	return &ParsedIntent{
		Intent: NewIntent(name),
	}
}

func NewParsedIntEntity(name string, confidence float32, value int, role string) *ParsedEntity {
	return &ParsedEntity{
		Entity:     NewIntEntity(name),
		Confidence: confidence,
		Data:       value,
		Role:       role,
	}
}

func NewParsedSingleDateTimeEntity(name string, confidence float32, date time.Time, granularity DateTimeGranularity, role string) *ParsedEntity {
	return &ParsedEntity{
		Entity:     NewSingleDateTimeEntity(name),
		Confidence: confidence,
		Data: &parsedSingleDateTime{
			Date:        date,
			Granularity: granularity,
		},
		Role: role,
	}
}

func NewParsedDateTimeIntervalEntity(name string, confidence float32, fromDate, toDate time.Time, fromGran, toGran DateTimeGranularity, role string) *ParsedEntity {
	return &ParsedEntity{
		Entity:     NewDateTimeIntervalEntity(name),
		Confidence: confidence,
		Data: &parsedDateTimeInterval{
			From: &parsedSingleDateTime{
				Date:        fromDate,
				Granularity: fromGran,
			},
			To: &parsedSingleDateTime{
				Date:        toDate,
				Granularity: toGran,
			},
		},
		Role: role,
	}
}

// NewParsedData is the constructor method for ParsedData
func NewParsedData(intent *ParsedIntent, entities []*ParsedEntity) *ParsedData {
	return nil
	//var entitiesWithType []*EntityWithType
	//
	//for _, e := range entities {
	//	entitiesWithType = append(entitiesWithType, &EntityWithType{
	//		Type:   e.Type(),
	//		Entity: e,
	//	})
	//}
	//
	//return &ParsedData{
	//	Intent:   intent,
	//	Entities: entitiesWithType,
	//}
}
