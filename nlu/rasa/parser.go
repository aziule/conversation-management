package rasa

import (
	"net/http"

	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/nlu"
	"fmt"
)

// RasaNluParser is the Rasa NLU implementation of the core/nlu/Parser
type RasaNluParser struct {
	client *http.Client
}

// NewRasaNluParser is the constructor method for RasaNluParser
func NewRasaNluParser(config *core.Config) nlu.Parser {
	return &RasaNluParser{
		client: http.DefaultClient,
	}
}

// ParseText is using the Rasa NLU service to try to understand some text
// and extract some relevant data
func (parser *RasaNluParser) ParseData(data interface{}) (*nlu.ParsedData, error) {
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
