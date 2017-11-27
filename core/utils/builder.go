package utils

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// BuilderConf represents the base conf variable that is passed to service builders
type BuilderConf map[string]interface{}

// ServiceBuilder is the main func used to build a service
type ServiceBuilder func(conf BuilderConf) (interface{}, error)

var (
	// serviceBuilders stores all of the available service builders
	serviceBuilders = make(map[string]ServiceBuilder)

	// ErrInvalidOrMissingParam is a generic error used when parsing a param isn't possible
	ErrInvalidOrMissingParam  = func(param string) error { return errors.New("Missing param or invalid param type: " + param) }
	ErrServiceBuilderNotFound = errors.New("Could not find the service builder")
)

// RegisterServiceBuilder registers a service builder
func RegisterServiceBuilder(name string, serviceBuilder ServiceBuilder) {
	_, registered := serviceBuilders[name]

	if registered {
		log.WithField("name", name).Warning("Service builder already registered, ignoring")
	}

	serviceBuilders[name] = serviceBuilder
}

// GetServiceBuilder returns a service from the available services builders.
// Returns an error if the service does not exist.
func GetServiceBuilder(name string) (ServiceBuilder, error) {
	serviceBuilder, ok := serviceBuilders[name]

	if !ok {
		return nil, ErrServiceBuilderNotFound
	}

	return serviceBuilder, nil
}

// GetParam tries to get a param from a given BuilderConf and returns it
func GetParam(conf BuilderConf, paramName string) interface{} {
	param, ok := conf[paramName]

	if !ok {
		return nil
	}

	return param
}
