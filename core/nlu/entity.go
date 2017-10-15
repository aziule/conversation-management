package nlu

import "time"

type EntityDataType string

// The available data types
const (
	DataTypeString EntityDataType = "string"
	DataTypeNumber EntityDataType = "number"
	DataTypeDateTime EntityDataType = "datetime"
)

// Intent represents the data of a text, as understood by the NLU service
type Entity struct {
	name  string
	value interface{}
	dataType EntityDataType
}

// NewEntity is the factory method for Entity
// It tries to create a new Entity based on the expected data type, by converting
// the value to the required type
func NewEntity(name string, value interface{}, dataType EntityDataType) (*Entity, error) {
	convertedValue, err := convertValue(value, dataType)

	if err != nil {
		// @todo: Log the error
		return nil, err
	}

	return &Entity{
		name:  name,
		value: convertedValue,
		dataType: dataType,
	}, nil
}

// convertValue takes a value and a data type and tries to convert the given value
// to the correct type, given the data type.
// Returns an error if the data type is not supported or the value is malformed.
func convertValue(value interface{}, dataType EntityDataType) (interface{}, error) {
	var convertedValue interface{}
	var ok bool = false

	switch dataType {
	case DataTypeString:
		convertedValue, ok = value.(string)
		break
	case DataTypeNumber:
		convertedValue, ok = value.(int)
		break
	case DataTypeDateTime:
		timeAsString, canConvert := value.(string)

		if !canConvert {
			return nil, ErrInvalidEntityDataType
		}

		convertedTime, err := time.Parse(time.RFC3339, timeAsString)

		if err != nil {
			ok = true
		}

		convertedValue = convertedTime

		break
	}

	if ! ok {
		return nil, ErrInvalidEntityDataType
	}

	return convertedValue, nil
}

// Getters
func (entity *Entity) Name() string  { return entity.name }
func (entity *Entity) Value() interface{} { return entity.value }
func (entity *Entity) DataType() EntityDataType { return entity.dataType }
