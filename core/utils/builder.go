package utils

import (
	"errors"

	"github.com/aziule/conversation-management/core/utils"
)

var (
	ErrUndefinedParam = func(param string) error { return errors.New("Missing param: " + param) }
	ErrInvalidParam   = func(param string) error { return errors.New("Invalid param type: " + param) }
)

// BuilderConf represents the base conf variable that is passed to any builder
type BuilderConf map[string]interface{}

func ValidateParam(conf BuilderConf, paramName string, expectedType interface{}) (interface{}, error) {
	param, ok := conf[paramName]

	if !ok {
		return nil, utils.ErrUndefinedParam(paramName)
	}

	parsed, ok := param.(expectedType.(T))

	if !ok {
		return nil, utils.ErrInvalidParam(paramName)
	}

	return parsed, nil
}
