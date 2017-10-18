package facebook

import (
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/core/nlp"
	"time"
)

type DataTypeMap map[string]nlp.DataType

// @todo: use an interface on top of that rather than a jason Object
// ParseNlpData returns an object of type ParsedData after parsing a jason object
func (bot *facebookBot) ParseNlpData(data *jason.Object) (*nlp.ParsedData, error) {
	var intent *nlp.Intent
	var entities []nlp.Entity

	for key, value := range data.Map() {
		dataType, ok := bot.DataTypeMap[key]

		if !ok {
			// @todo: log
			fmt.Println("NOT HANDLED: ", key)
			continue
		}

		switch dataType {
		case nlp.DataTypeIntent:
			i, err := toIntent(value)

			if err != nil {
				// @todo: handle it better
				fmt.Println("Error: ", err)
				continue
			}

			intent = i
			break
		default:
			entity, err := toEntity(value, key, dataType)

			if err != nil {
				// @todo: handle it better
				fmt.Println("Error: ", err)
				continue
			}

			entities = append(entities, entity)
			break
		}
	}

	return nlp.NewParsedData(intent, entities), nil
}

// toEntity converts a jason entity to a built-in NLP representation of an entity
// Returns an error if the JSON is malformed or if we do not handle the data type correctly
func toEntity(value *jason.Value, name string, dataType nlp.DataType) (nlp.Entity, error) {
	object, err := value.ObjectArray()

	if err != nil {
		// @todo: log, handle error better and use a predefined error
		return nil, errors.New("Could not cast to object")
	}

	for _, e := range object {
		confidence, err := e.GetFloat64("confidence")

		if err != nil {
			// @todo: log, handle error better and use a predefined error
			return nil, err
		}

		switch dataType {
		case nlp.DataTypeInt:
			value, err := e.GetInt64("value")

			if err != nil {
				// @todo: log, handle error better and use a predefined error
				return nil, err
			}

			return nlp.NewIntEntity(name, float32(confidence), int(value)), nil
			break
		case nlp.DataTypeDateTime:
			return nlp.NewSingleDateTimeEntity(name, float32(confidence), time.Now(), nlp.GranularityDay), nil
			break
		}
	}

	return nil, errors.New("Data type not handled")
}

// toIntent converts a jason intent to a built-in NLP representation of an intent
// Returns an error if the JSON is malformed
func toIntent(value *jason.Value) (*nlp.Intent, error) {
	object, err := value.ObjectArray()

	if err != nil {
		// @todo: log, handle error better and use a predefined error
		return nil, errors.New("Could not cast to object")
	}

	// Handle single intents only
	intentName, err := object[0].GetString("value")

	if err != nil {
		// @todo: log, handle error better and use a predefined error
		return nil, err
	}

	return nlp.NewIntent(intentName), nil
}
