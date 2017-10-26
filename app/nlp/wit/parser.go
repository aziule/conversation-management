package wit

import (
	"errors"
	"fmt"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/app/nlp"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCouldNotParseJson       = errors.New("Could not parse JSON")
	ErrMalformedDateTime       = errors.New("Malformed datetime value")
	ErrCouldNotParseJsonObject = errors.New("Could not parse object from JSON")
	ErrMissingKey              = func(key string) error { return errors.New(fmt.Sprintf("Missing key: %s", key)) }
	ErrCouldNotCastValue       = func(key, expectedType string) error {
		return errors.New(fmt.Sprintf("Could not cast %s to %s", expectedType))
	}
	ErrUnhandledDataType = func(dataType string) error { return errors.New(fmt.Sprintf("Unhandled data type %s", dataType)) }
)

// witParser is the NLP parser for Wit.
// It implements the nlp.Parser interface
type witParser struct {
	dataTypeMap nlp.DataTypeMap
}

// NewParser is the constructor method for witParser
func NewParser(dataTypeMap nlp.DataTypeMap) nlp.Parser {
	return &witParser{
		dataTypeMap: dataTypeMap,
	}
}

// ParseNlpData parses raw data and returns parsed data
func (parser *witParser) ParseNlpData(rawData []byte) (*nlp.ParsedData, error) {
	var intent *nlp.Intent
	var entities []nlp.Entity

	data, err := jason.NewObjectFromBytes(rawData)

	if err != nil {
		log.WithField("rawData", string(rawData)).Infof("Could not parse JSON: %s", err)
		return nil, ErrCouldNotParseJson
	}

	for key, value := range data.Map() {
		dataType, ok := parser.dataTypeMap[key]

		if !ok {
			log.WithField("key", key).Warnf("Data type is not handled: %s", key)
			continue
		}

		switch dataType {
		case nlp.DataTypeIntent:
			i, err := toIntent(value)

			if err != nil {
				log.WithField("dataType", dataType).Warnf("Could not convert value to DataTypeIntent: %s", err)
				continue
			}

			intent = i
			break
		default:
			entity, err := toEntity(value, key, dataType)

			if err != nil {
				log.WithFields(log.Fields{
					"dataType": dataType,
					"key":      key,
				}).Warnf("Could not convert value to entity: %s", err)
				continue
			}

			entities = append(entities, entity)
			break
		}
	}

	return nlp.NewParsedData(intent, entities), nil
}

// toIntent converts a jason intent to a built-in NLP representation of an intent
// Returns an error if the JSON is malformed
func toIntent(value *jason.Value) (*nlp.Intent, error) {
	object, err := value.ObjectArray()

	if err != nil {
		return nil, ErrCouldNotParseJsonObject
	}

	// Handle single intents only
	intentName, err := object[0].GetString("value")

	if err != nil {
		return nil, ErrMissingKey("value")
	}

	return nlp.NewIntent(intentName), nil
}

// toEntity converts a jason entity to a built-in NLP representation of an entity
// Returns an error if the JSON is malformed or if we do not handle the data type correctly
func toEntity(value *jason.Value, name string, dataType nlp.DataType) (nlp.Entity, error) {
	object, err := value.ObjectArray()

	if err != nil {
		return nil, ErrCouldNotParseJsonObject
	}

	for _, e := range object {
		confidence, err := e.GetFloat64("confidence")

		if err != nil {
			return nil, ErrCouldNotCastValue("confidence", "float64")
		}

		switch dataType {
		case nlp.DataTypeInt:
			value, err := e.GetInt64("value")

			if err != nil {
				return nil, ErrCouldNotCastValue("value", "int64")
			}

			return nlp.NewIntEntity(name, float32(confidence), int(value)), nil
		case nlp.DataTypeDateTime:
			_, err := e.GetString("value")

			if err != nil {
				// If there's an error, then look for interval datetimes, parsed as "from" & "to"
				from, err := e.GetObject("from")

				if err != nil {
					return nil, ErrMissingKey("from")
				}

				fromTime, fromGranularity, err := extractDateTimeInformation(from)

				if err != nil {
					return nil, err
				}

				to, err := e.GetObject("to")

				if err != nil {
					return nil, ErrMissingKey("to")
				}

				toTime, toGranularity, err := extractDateTimeInformation(to)

				if err != nil {
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
				return nil, err
			}

			return nlp.NewSingleDateTimeEntity(name, float32(confidence), t, granularity), nil
		}
	}

	return nil, ErrUnhandledDataType(string(dataType))
}

// extractDateTimeInformation extracts useful date time information from JSON
// This is heavily coupled with what Wit returns us
func extractDateTimeInformation(object *jason.Object) (time.Time, nlp.DateTimeGranularity, error) {
	value, err := object.GetString("value")

	if err != nil {
		return time.Time{}, "", ErrMissingKey("value")
	}

	t, err := stringToTime(value)

	if err != nil {
		return time.Time{}, "", err
	}

	grain, err := object.GetString("grain")

	if err != nil {
		return time.Time{}, "", ErrMissingKey("grain")
	}

	// @todo: use a converter that will check that the granularity exists
	return t, nlp.DateTimeGranularity(grain), nil
}

// stringToTime converts a given string to a tine.Time
// @TODO: move somewhere else
func stringToTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)

	if err != nil {
		log.WithFields(log.Fields{
			"value":          value,
			"expectedLayout": "RFC3339",
		}).Infof("Could not cast string to time.Time: %s", err)
		return time.Time{}, ErrMalformedDateTime
	}

	return t, nil
}
