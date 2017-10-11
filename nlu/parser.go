package nlu

import (
	"github.com/aziule/dev-board/config"
	"github.com/aziule/conversation-management/core/nlu"
	"net/http"
)

// RasaNluParser is the Rasa NLU implementation of the core/nlu/Parser
type RasaNluParser struct {
	client          *http.Client
}

// NewRasaNluParser is the constructor method for RasaNluParser
func NewRasaNluParser(config *config.Config) nlu.Parser {
	return &RasaNluParser{
		client: http.DefaultClient,
	}
}

// ParseText is using the Rasa NLU service to try to understand some text
// and extract some relevant data
func (parser *RasaNluParser) ParseText(text string) (*nlu.ParsedText, error) {
	entities := []*nlu.Entity{}
	entities = append(entities, nlu.NewEntity("entity_1", "value 1"))
	entities = append(entities, nlu.NewEntity("entity_2", "value 2"))

	return nlu.NewParsedText(text, nlu.NewIntent("intent_name"), entities), nil
}
