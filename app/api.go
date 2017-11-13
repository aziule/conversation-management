package app

import (
	"net/http"

	"github.com/aziule/conversation-management/core/bot"
)

// ApiEndpoint represents an API endpoint
type ApiEndpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Api is the struct containing endpoints about the API
type Api struct {
	ApiEndpoints []*bot.ApiEndpoint
}

// RegisterApiEndpoints registers a new endpoint to the API
func (api *Api) RegisterApiEndpoints(endpoints ...*bot.ApiEndpoint) {
	for _, endpoint := range endpoints {
		api.ApiEndpoints = append(api.ApiEndpoints, endpoint)
	}
}
