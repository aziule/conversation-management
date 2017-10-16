package api

import "github.com/aziule/conversation-management/core/nlu"

const (
	NluDataTypeNumber     = "number"
	NluDataTypeDateTime   = "datetime"
	NluDataTypeDateIntent = "intent"
)

// parseNluDataFromJson parses JSON and returns a ParsedData object.
// Returns an error if the data is malformed or it there is any unhandled data type.
func parseNluDataFromJson(bytes []byte) (*nlu.ParsedData, error) {
	return nil, nil
}
