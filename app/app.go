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
		app.Api.RegisterApiEndpoints(curr.ApiEndpoints()...)
	}

	r := chi.NewRouter()

	// Automatically set the bot's webhooks routes
	for _, webhook := range b.Webhooks() {
		log.Debugf("%s %s", string(webhook.Method), webhook.Path)

		switch webhook.Method {
		case "GET":
			r.Get(webhook.Path, webhook.Handler)
		case "POST":
			r.Post(webhook.Path, webhook.Handler)
		}
	}

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), r)
}
