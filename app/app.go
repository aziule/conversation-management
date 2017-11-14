// Package app provides the main methods and config to initialise and run bots
// on any available platform.
package app

import (
	"net/http"
	"strconv"

	"github.com/aziule/conversation-management/app/facebook"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/infrastructure/conversation/memory"
	"github.com/aziule/conversation-management/infrastructure/conversation/mongo"
	db "github.com/aziule/conversation-management/infrastructure/mongo"
	"github.com/aziule/conversation-management/infrastructure/nlp/wit"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// App defines the main structure, holding information about
// what bot is running, public-facing API endpoints, etc.
type App struct {
	Api  *Api
	Bots []bot.Bot
}

// Run starts the server and waits for interactions
func Run(configFilePath string) {
	config, err := LoadConfig(configFilePath)

	if err != nil {
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	api := &Api{}
	app := &App{
		Api: api,
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	db, err := db.CreateSession(db.DbParams{
		DbHost: config.DbHost,
		DbName: config.DbName,
		DbUser: config.DbUser,
		DbPass: config.DbPass,
	})

	if err != nil {
		log.Fatalf("An error occurred when connecting to the db: %s", err)
	}

	defer db.Close()

	// @todo: register all available implementations using a factory
	// pattern, and fetch them directly from the config passed
	b := facebook.NewBot(
		&facebook.Config{
			VerifyToken:            config.FbVerifyToken,
			ApiVersion:             config.FbApiVersion,
			PageAccessToken:        config.FbPageAccessToken,
			NlpParser:              wit.NewParser(),
			ConversationRepository: mongo.NewMongodbRepository(db),
			StoryRepository:        memory.NewStoryRepository(),
		},
	)

	app.Bots = append(app.Bots, b)

	for _, curr := range app.Bots {
		app.Api.RegisterEndpoints(curr.ApiEndpoints()...)
	}

	app.Api.RegisterEndpoint(bot.NewApiEndpoint(
		"GET",
		"/bots",
		app.handleListBots,
	))

	r := chi.NewRouter()

	// Automatically listen to the bot's webhooks routes
	for _, webhook := range b.Webhooks() {
		bindRoute(r, webhook.Method, webhook.Path, webhook.Handler)
	}

	// Same for the API
	for _, endpoint := range app.Api.Endpoints {
		bindRoute(r, endpoint.Method, endpoint.Path, endpoint.Handler)
	}

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), r)
}

// bindRoute binds a new route to the router given its method, path and handler func
func bindRoute(r *chi.Mux, method, path string, handler http.HandlerFunc) {
	log.Debugf("%s %s", string(method), path)

	switch method {
	case "GET":
		r.Get(path, handler)
	case "POST":
		r.Post(path, handler)
	}
}

// handleListBots is the handler func for listing the available bots
func (app *App) handleListBots(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
