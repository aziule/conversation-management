package utils

import "errors"

var (
	ErrUndefinedParam = func(param string) error { return errors.New("Missing param: " + param) }
	ErrInvalidParam   = func(param string) error { return errors.New("Invalid param type: " + param) }
)

// BuilderConf represents the base conf variable that is passed to any builder
type BuilderConf map[string]interface{}

//func ParseParam(conf BuilderConf, paramName string, expectedType interface{}) (interface{}, error) {
//
//}
