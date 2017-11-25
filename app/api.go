package app

import (
	"net/http"
	"strings"

	"encoding/json"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// appApi is the struct containing API endpoints and responsible for
// setting up their routes.
//
// The API represents both:
// - The public-facing API: Facebook webhooks, etc.
// - The private-facing API: bots management, metrics, etc.
type appApi struct {
	router          *chi.Mux
	app             *app
	AppApiEndpoints []*bot.ApiEndpoint
	Webhooks        []*bot.Webhook
}

// NewAppApi is the constructor method for AppApi. The API endpoints
// are not mounted by default. They can be mounted using Mount()
func NewAppApi(router *chi.Mux, app *app) *appApi {
	return &appApi{
		app:             app,
		router:          router,
		AppApiEndpoints: nil,
		Webhooks:        nil,
	}
}

// Mount mounts all of the AppApi endpoints (management API, bots APIs, bots webhooks)
func (appApi *appApi) Mount() {
	appApi.bindDefaultEndpoints()
	appApi.RegisterBotsEndpoints(appApi.app.Bots...)
}

// bindDefaultEndpoints binds the default API endpoints to the appApi object
func (appApi *appApi) bindDefaultEndpoints() {
	appApi.RegisterAppApiEndpoints(
		bot.NewApiEndpoint(
			"GET",
			"/bots",
			appApi.handleListBots,
		),
		bot.NewApiEndpoint(
			"POST",
			"/bots",
			appApi.handleCreateBot,
		),
	)
}

// RegisterAppApiEndpoints registers a set of new API endpoints
// and binds them to the router.
func (appApi *appApi) RegisterAppApiEndpoints(endpoints ...*bot.ApiEndpoint) {
	for _, endpoint := range endpoints {
		path := "/api" + endpoint.BasePath
		path = strings.TrimRight(path, "/")
		endpoint.MountedPath = path
		appApi.bind(endpoint.Method, path, endpoint.Handler)

		appApi.AppApiEndpoints = append(appApi.AppApiEndpoints, endpoint)
	}
}

// RegisterWebhooks registers a set of new Webhooks
// and binds them to the router.
func (appApi *appApi) RegisterBotsEndpoints(bots ...bot.Bot) {
	for _, b := range bots {
		for _, endpoint := range b.ApiEndpoints() {
			path := "/api/bots/" + b.Definition().Slug + endpoint.BasePath
			path = strings.TrimRight(path, "/")
			endpoint.MountedPath = path
			appApi.bind(endpoint.Method, path, endpoint.Handler)

			appApi.AppApiEndpoints = append(appApi.AppApiEndpoints, endpoint)
		}
		for _, endpoint := range b.Webhooks() {
			path := "/bots/" + b.Definition().Slug + "/webhooks" + endpoint.BasePath
			path = strings.TrimRight(path, "/")
			endpoint.MountedPath = path
			appApi.bind(endpoint.Method, path, endpoint.Handler)

			appApi.Webhooks = append(appApi.Webhooks, endpoint)
		}
	}
}

// bindRoute binds a new route to the router given its method, path and handler func
func (appApi *appApi) bind(method, path string, handler http.HandlerFunc) {
	log.Debugf("%s %s", string(method), path)

	switch method {
	case "GET":
		appApi.router.Get(path, handler)
	case "POST":
		appApi.router.Post(path, handler)
	}
}

// handleListBots is the handler func that lists the available bots
func (appApi *appApi) handleListBots(w http.ResponseWriter, r *http.Request) {
	var definitions []*bot.Definition

	for _, b := range appApi.app.Bots {
		definitions = append(definitions, b.Definition())
	}

	j, _ := json.Marshal(definitions)

	w.Write(j)
}

// handleCreateBot creates a new bot
func (appApi *appApi) handleCreateBot(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var definition bot.Definition
	err := decoder.Decode(&definition)

	if err != nil {
		log.Errorf("Could not decode the request body: %s", err)
		// @todo: return an error to the user
		return
	}

	err = appApi.app.botRepository.Save(&definition)

	if err != nil {
		log.Errorf("Could not save the bot: %s", err)
		// @todo: return an error to the user
		return
	}

	j, _ := json.Marshal(definition)

	// @todo: return an proper response
	w.Write(j)
}
