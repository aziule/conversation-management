package app

import (
	"github.com/aziule/conversation-management/core/bot"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Api is the struct containing API endpoints and responsible for
// setting up their routes.
//
// The API represents both:
// - The public-facing API: Facebook webhooks, etc.
// - The private-facing API: bots management, metrics, etc.
type Api struct {
	router       *chi.Mux
	ApiEndpoints []*bot.ApiEndpoint
	Webhooks     []*bot.Webhook
}

// NewApi is the constructor method for Api. The endpoints are not set
// by default and must be registered using RegisterEndpoints
func NewApi(router *chi.Mux) *Api {
	return &Api{
		router:       router,
		ApiEndpoints: nil,
		Webhooks:     nil,
	}
}

// RegisterApiEndpoints registers a set of new API endpoints
// and binds them to the router.
func (api *Api) RegisterApiEndpoints(endpoints ...*bot.ApiEndpoint) {
	for _, endpoint := range endpoints {
		api.ApiEndpoints = append(api.ApiEndpoints, endpoint)
		api.bind(endpoint.Method, "/api"+endpoint.Path, endpoint.Handler)
	}
}

// RegisterWebhooks registers a set of new Webhooks
// and binds them to the router.
func (api *Api) RegisterWebhooks(endpoints ...*bot.Webhook) {
	for _, endpoint := range endpoints {
		api.Webhooks = append(api.Webhooks, endpoint)
		api.bind(endpoint.Method, "/webhooks"+endpoint.Path, endpoint.Handler)
	}
}

// bindRoute binds a new route to the router given its method, path and handler func
func (api *Api) bind(method, path string, handler http.HandlerFunc) {
	log.Debugf("%s %s", string(method), path)

	switch method {
	case "GET":
		api.router.Get(path, handler)
	case "POST":
		api.router.Post(path, handler)
	}
}
