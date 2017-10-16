package facebook

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/core/nlu"
)

type NluEntry string
type NluDataType string

const (
	NluDataTypeNumber     NluDataType = "number"
	NluDataTypeDateTime   NluDataType = "datetime"
	NluDataTypeDateIntent NluDataType = "intent"
)

type NluEntryDataTypeMap map[NluEntry]NluDataType

// @todo: use an interface on top of that rather than a jason Object
// ParseNlpData returns an object of type ParsedData after parsing a jason object
func ParseNlpData(data *jason.Object) (*nlu.ParsedData, error) {
	for key, value := range data.Map() {
		v, _ := value.MarshalJSON()
		fmt.Println(key, string(v))
	}

	return nil, nil
}

func parseSingleNlpData() {

}
