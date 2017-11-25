package utils

import (
	"errors"
)

var ErrInvalidOrMissingParam = func(param string) error { return errors.New("Missing param or invalid param type: " + param) }

// BuilderConf represents the base conf variable that is passed to any builder
type BuilderConf map[string]interface{}

// GetParam tries to get a param from a given BuilderConf and returns it
func GetParam(conf BuilderConf, paramName string) interface{} {
	param, ok := conf[paramName]

	if !ok {
		return nil
	}

	return param
}
