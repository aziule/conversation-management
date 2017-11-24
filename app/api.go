package app

import (
	"net/http"
	"strings"

	"encoding/json"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// api is the struct containing API endpoints and responsible for
// setting up their routes.
//
// The API represents both:
// - The public-facing API: Facebook webhooks, etc.
// - The private-facing API: bots management, metrics, etc.
type api struct {
	router       *chi.Mux
	app          *app
	ApiEndpoints []*bot.ApiEndpoint
	Webhooks     []*bot.Webhook
}

// NewApi is the constructor method for Api. The API endpoints
// are not mounted by default. They can be mounted using Mount()
func NewApi(router *chi.Mux, app *app) *api {
	return &api{
		app:          app,
		router:       router,
		ApiEndpoints: nil,
		Webhooks:     nil,
	}
}

// Mount mounts all of the Api endpoints (management API, bots APIs, bots webhooks)
func (api *api) Mount() {
	api.bindDefaultEndpoints()
	api.RegisterBotsEndpoints(api.app.Bots...)
}

// bindDefaultEndpoints binds the default API endpoints to the api object
func (api *api) bindDefaultEndpoints() {
	api.RegisterApiEndpoints(
		bot.NewApiEndpoint(
			"GET",
			"/bots",
			api.handleListBots,
		),
		bot.NewApiEndpoint(
			"POST",
			"/bots",
			api.handleCreateBot,
		),
	)
}

// RegisterApiEndpoints registers a set of new API endpoints
// and binds them to the router.
func (api *api) RegisterApiEndpoints(endpoints ...*bot.ApiEndpoint) {
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
func (api *api) RegisterBotsEndpoints(bots ...bot.Bot) {
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
func (api *api) bind(method, path string, handler http.HandlerFunc) {
	log.Debugf("%s %s", string(method), path)

	switch method {
	case "GET":
		api.router.Get(path, handler)
	case "POST":
		api.router.Post(path, handler)
	}
}

// handleListBots is the handler func that lists the available bots
func (api *api) handleListBots(w http.ResponseWriter, r *http.Request) {
	var metadatas []*bot.Metadata

	for _, b := range api.app.Bots {
		metadatas = append(metadatas, b.Metadata())
	}

	j, _ := json.Marshal(metadatas)

	w.Write(j)
}

// handleCreateBot creates a new bot
func (api *api) handleCreateBot(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var metadata bot.Metadata
	err := decoder.Decode(&metadata)

	if err != nil {
		log.Errorf("Could not decode the request body: %s", err)
		// @todo: return an error to the user
		return
	}

	err = api.app.botRepository.Save(&metadata)

	if err != nil {
		log.Errorf("Could not save the bot: %s", err)
		// @todo: return an error to the user
		return
	}

	j, _ := json.Marshal(metadata)

	// @todo: return an proper response
	w.Write(j)
}
