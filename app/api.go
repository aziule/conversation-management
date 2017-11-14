package app

import (
	"github.com/aziule/conversation-management/core/bot"
)

// Api is the struct containing endpoints about the API
type Api struct {
	Endpoints []*bot.ApiEndpoint
}

// RegisterEndpoints registers a set of new endpoints to the API
func (api *Api) RegisterEndpoints(endpoints ...*bot.ApiEndpoint) {
	for _, endpoint := range endpoints {
		api.RegisterEndpoint(endpoint)
	}
}

// RegisterEndpoint registers a new endpoint to the API
func (api *Api) RegisterEndpoint(endpoint *bot.ApiEndpoint) {
	api.Endpoints = append(api.Endpoints, endpoint)
}
