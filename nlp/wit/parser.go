package wit

import (
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/labstack/gommon/log"
	"time"
)

// WitParser is the NLP parser for Wit.
// It implements the nlp.Parser interface
type WitParser struct {
	dataTypeMap nlp.DataTypeMap
}

// NewParser is the constructor method for WitParser
func NewParser(dataTypeMap nlp.DataTypeMap) nlp.Parser {
	return &WitParser{
		dataTypeMap: dataTypeMap,
	}
}

// ParseNlpData parses raw data and returns parsed data
func (parser *WitParser) ParseNlpData(rawData []byte) (*nlp.ParsedData, error) {
	var intent *nlp.Intent
	var entities []nlp.Entity

	data, err := jason.NewObjectFromBytes(rawData)

	if err != nil {
		// @todo: better stuff
		return nil, errors.New("Error when getting json from data")
	}

	for key, value := range data.Map() {
		dataType, ok := parser.dataTypeMap[key]

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
			_, err := e.GetString("value")

			if err != nil {
				// If there's an error, then look for interval datetimes, parsed as "from" & "to"
				from, err := e.GetObject("from")

				if err != nil {
					log.Info("Unable to find the from value")
					// @todo: log, handle error better and use a predefined error
					return nil, err
				}

				fromTime, fromGranularity, err := extractDateTimeInformation(from)

				if err != nil {
					log.Info("Unable to parse the from information")
					// @todo: log, handle error better and use a predefined error
					return nil, err
				}

				to, err := e.GetObject("to")

				if err != nil {
					log.Info("Unable to find the to value")
					// @todo: log, handle error better and use a predefined error
					return nil, err
				}

				toTime, toGranularity, err := extractDateTimeInformation(to)

				if err != nil {
					log.Info("Unable to parse the to information")
					// @todo: log, handle error better and use a predefined error
					return nil, err
				}

				return nlp.NewIntervalDateTimeEntity(
					name,
					float32(confidence),
					fromTime,
					toTime,
					fromGranularity,
					toGranularity,
				), nil
			}

			t, granularity, err := extractDateTimeInformation(e)

			if err != nil {
				log.Info("Unable to parse the single datetime information")
				// @todo: log, handle error better and use a predefined error
				return nil, err
			}

			return nlp.NewSingleDateTimeEntity(name, float32(confidence), t, granularity), nil
		}
	}

	return nil, errors.New("Data type not handled")
}

// extractDateTimeInformation extracts useful date time information from JSON
// This is heavily coupled with what Wit returns us
func extractDateTimeInformation(object *jason.Object) (time.Time, nlp.DateTimeGranularity, error) {
	value, err := object.GetString("value")

	if err != nil {
		log.Info("Unable to find from's value")
		// @todo: log, handle error better and use a predefined error
		return time.Time{}, "", err
	}

	t, err := stringToTime(value)

	if err != nil {
		log.Info("Unable to case the from time to time.Time")
		// @todo: log, handle error better and use a predefined error
		return time.Time{}, "", err
	}

	grain, err := object.GetString("grain")

	if err != nil {
		log.Info("Unable to find from's grain")
		// @todo: log, handle error better and use a predefined error
		return time.Time{}, "", err
	}

	// @todo: use a converter that will switch between values
	return t, nlp.DateTimeGranularity(grain), nil
}

// stringToTime converts a given string to a tine.Time
// @TODO: move somewhere else
func stringToTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)

	if err != nil {
		// @todo: use custom error and handle better
		return time.Time{}, errors.New("Invalid value")
	}

	return t, nil
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
