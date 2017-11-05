// Package utils provides useful methods for data manipulation and parsing.
package utils

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrMalformedDateTime = errors.New("Malformed datetime value")

// stringToTime converts a given string to a time.Time according to RFC3339 standards
// Returns an ErrMalformedDateTime in case the value given is malformed
func ParseTime(value string) (time.Time, error) {
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
