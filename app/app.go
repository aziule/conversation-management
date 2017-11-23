// Package app provides the main methods and config to initialise and run bots
// on any available platform.
package app

import (
	"encoding/json"
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

	router := chi.NewRouter()

	api := NewApi(router)

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
			continue
		}

		app.Bots = append(app.Bots, b)
	}

	api.RegisterApiEndpoints(bot.NewApiEndpoint(
		"GET",
		"/bots",
		app.handleListBots,
	))

	// Listen to each of the bot's webhooks and API endpoints
	api.RegisterBotsEndpoints(app.Bots...)

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), router)
}

// handleListBots is the handler func that lists the available bots
func (app *App) handleListBots(w http.ResponseWriter, r *http.Request) {
	var metadatas []*bot.Metadata

	for _, b := range app.Bots {
		metadatas = append(metadatas, b.Metadata())
	}

	j, _ := json.Marshal(metadatas)

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
