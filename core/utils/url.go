package utils

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrCouldNotFetchParam = func(key string) error { return errors.New(fmt.Sprintf("Could not fetch param: %s", key)) }

// GetSingleQueryParam fetches a single query param using the given url values.
// Returns an error in case the param could not be fetched.
func GetSingleQueryParam(values url.Values, key string) (string, error) {
	params, ok := values[key]

	if !ok || len(params) != 1 {
		return "", ErrCouldNotFetchParam(key)
	}

	return params[0], nil
}
