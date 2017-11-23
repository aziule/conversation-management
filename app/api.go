package app

import (
	"net/http"
	"strings"

	"github.com/aziule/conversation-management/core/bot"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
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
		path := "/api" + endpoint.BasePath
		path = strings.TrimRight(path, "/")
		endpoint.MountedPath = path
		api.bind(endpoint.Method, path, endpoint.Handler)

		api.ApiEndpoints = append(api.ApiEndpoints, endpoint)
	}
}

// RegisterWebhooks registers a set of new Webhooks
// and binds them to the router.
func (api *Api) RegisterBotsEndpoints(bots ...bot.Bot) {
	for _, b := range bots {
		for _, endpoint := range b.ApiEndpoints() {
			path := "/api/bots/" + b.Metadata().Slug + endpoint.BasePath
			path = strings.TrimRight(path, "/")
			endpoint.MountedPath = path
			api.bind(endpoint.Method, path, endpoint.Handler)

			api.ApiEndpoints = append(api.ApiEndpoints, endpoint)
		}
		for _, endpoint := range b.Webhooks() {
			path := "/bots/" + b.Metadata().Slug + "/webhooks" + endpoint.BasePath
			path = strings.TrimRight(path, "/")
			endpoint.MountedPath = path
			api.bind(endpoint.Method, path, endpoint.Handler)

			api.Webhooks = append(api.Webhooks, endpoint)
		}
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
