package facebook

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/core/nlp"
)

type NlpEntry string
type NlpDataType string

const (
	NlpDataTypeNumber     NlpDataType = "number"
	NlpDataTypeDateTime   NlpDataType = "datetime"
	NlpDataTypeDateIntent NlpDataType = "intent"
)

type NlpEntryDataTypeMap map[NlpEntry]NlpDataType

// @todo: use an interface on top of that rather than a jason Object
// ParseNlpData returns an object of type ParsedData after parsing a jason object
func ParseNlpData(data *jason.Object) (*nlp.ParsedData, error) {
	for key, value := range data.Map() {
		v, _ := value.MarshalJSON()
		fmt.Println(key, string(v))
	}

	return nil, nil
}

func parseSingleNlpData() {

}
