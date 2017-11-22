// Package app provides the main methods and config to initialise and run bots
// on any available platform.
package app

import (
	"net/http"
	"strconv"

	"github.com/aziule/conversation-management/app/facebook"
	"github.com/aziule/conversation-management/core/bot"
	fbApi "github.com/aziule/conversation-management/infrastructure/facebook/api"
	"github.com/aziule/conversation-management/infrastructure/memory"
	"github.com/aziule/conversation-management/infrastructure/mongo"
	"github.com/aziule/conversation-management/infrastructure/wit"
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

	db, err := mongo.CreateSession(mongo.DbParams{
		DbHost: config.DbHost,
		DbName: config.DbName,
		DbUser: config.DbUser,
		DbPass: config.DbPass,
	})

	if err != nil {
		log.Fatalf("An error occurred when connecting to the db: %s", err)
	}

	defer db.Close()

	botRepository := mongo.NewBotRepository(db)

	metadatas, err := botRepository.FindAll()

	if err != nil {
		log.Fatalf("An error occurred when finding the bots list: %s", err)
	}

	for _, metadata := range metadatas {
		var b bot.Bot

		switch metadata.Platform {
		case bot.PlatformFacebook:
			// @todo: register all available implementations using a factory
			// pattern, and fetch them directly from the config passed
			b = facebook.NewBot(
				&facebook.Config{
					Metadata:               metadata,
					FbApi:                  fbApi.NewfacebookApi(config.FbApiVersion, config.FbPageAccessToken, http.DefaultClient),
					NlpParser:              wit.NewParser(),
					ConversationRepository: mongo.NewConversationRepository(db),
					StoryRepository:        memory.NewStoryRepository(),
				},
			)
		default:
			log.Errorf("Unhandled platform: %s", metadata.Platform)
		}

		app.Bots = append(app.Bots, b)
	}

	api.RegisterEndpoint(bot.NewApiEndpoint(
		"GET",
		"/bots",
		app.handleListBots,
	))

	r := chi.NewRouter()

	// Listen to each of the bot's webhooks and API endpoints
	for _, b := range app.Bots {
		api.RegisterEndpoints(b.ApiEndpoints()...)

		for _, webhook := range b.Webhooks() {
			bindRoute(r, webhook.Method, webhook.Path, webhook.Handler)
		}
	}

	// Listen to the app's API endpoints
	for _, endpoint := range api.Endpoints {
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
