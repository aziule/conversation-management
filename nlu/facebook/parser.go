package facebook

import (
	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/nlu"
	"fmt"
)

// FacebookNlp is the structure returned by Facebook that contains NLP information
type FacebookNlp struct {

}

// FacebookParser is the Facebook's built-in NLP implementation of the core/nlu/Parser
type FacebookParser struct {}

// NewFacebookParser is the constructor method for FacebookParser
func NewFacebookParser(config *core.Config) nlu.Parser {
	return &FacebookParser{
	}
}

// ParseText is using what Facebook's built-in NLP service returned
func (parser *FacebookParser) ParseData(data interface{}) (*nlu.ParsedData, error) {
	text, ok := data.(string)

	if !ok {
		// @todo: log
		return nil, nlu.ErrInvalidDataType
	}

	// @todo: parse it
	fmt.Println("Should now parse", text)

	entities := []*nlu.Entity{}

	e1, err := nlu.NewEntity("entity_1", "value 1", nlu.DataTypeString)
	if err != nil {return nil, err}

	e2, err := nlu.NewEntity("entity_2", 12, nlu.DataTypeNumber)
	if err != nil {return nil, err}

	e3, err := nlu.NewEntity("entity_3", "2017-10-17T17:00:00.000-07:00", nlu.DataTypeDateTime)
	if err != nil {return nil, err}

	entities = append(entities, e1, e2, e3)

	return nlu.NewParsedData(nlu.NewIntent("intent_name"), entities), nil
}
