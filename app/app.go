// Package app provides the main methods and config to initialise and run bots
// on any available platform.
package app

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aziule/conversation-management/app/facebook"
	"github.com/aziule/conversation-management/core/api"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/infrastructure/mongo"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	// Required for initialisation
	_ "github.com/aziule/conversation-management/infrastructure/facebook"
	_ "github.com/aziule/conversation-management/infrastructure/memory"
	_ "github.com/aziule/conversation-management/infrastructure/wit"
)

// app defines the main structure, holding information about
// what bot is running, public-facing API endpoints, etc.
type app struct {
	Bots          []bot.Bot
	botRepository bot.Repository
	nlpRepository nlp.Repository
}

// Run starts the server and waits for interactions
func Run(configFilePath string) {
	config, err := LoadConfig(configFilePath)

	if err != nil {
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	rand.Seed(time.Now().UnixNano())

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

	botRepository, err := bot.NewRepository("mongo", map[string]interface{}{
		"db": db,
	})

	if err != nil {
		log.Fatalf("An error occurred when creating the bot repository: %s", err)
	}

	conversationRepository, err := conversation.NewRepository("mongo", map[string]interface{}{
		"db": db,
	})

	if err != nil {
		log.Fatalf("An error occurred when creating the conversation repository: %s", err)
	}

	storyRepository, err := conversation.NewStoryRepository("memory", nil)

	if err != nil {
		log.Fatalf("An error occurred when creating the story repository: %s", err)
	}

	fbApi, err := api.NewFacebookApi("facebook", map[string]interface{}{
		"page_access_token": config.FbPageAccessToken,
		"version":           config.FbApiVersion,
		"client":            http.DefaultClient,
	})

	if err != nil {
		log.Fatalf("An error occurred when creating the Facebook API: %s", err)
	}

	definitions, err := botRepository.FindAll()

	if err != nil {
		log.Fatalf("An error occurred when finding the bots list: %s", err)
	}

	nlpParser, err := nlp.NewParser("wit", nil)

	if err != nil {
		log.Fatalf("An error occurred when creating the parser: %s", err)
	}

	nlpApi, err := nlp.NewApi("wit", map[string]interface{}{
		"client":       http.DefaultClient,
		"bearer_token": config.WitBearerToken,
	})

	if err != nil {
		log.Fatalf("An error occurred when creating the NLP API: %s", err)
	}

	nlpRepository, err := nlp.NewRepository("wit", map[string]interface{}{
		"api": nlpApi,
	})

	if err != nil {
		log.Fatalf("An error occurred when creating the repository: %s", err)
	}

	app := &app{
		botRepository: botRepository,
		nlpRepository: nlpRepository,
	}

	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	for _, definition := range definitions {
		var b bot.Bot

		switch definition.Platform {
		case bot.PlatformFacebook:
			// @todo: register all available implementations using a factory
			// pattern, and fetch them directly from the config passed
			b = facebook.NewBot(
				&facebook.Config{
					Definition:             definition,
					FbApi:                  fbApi,
					NlpParser:              nlpParser,
					ConversationRepository: conversationRepository,
					StoryRepository:        storyRepository,
				},
			)
		default:
			log.Errorf("Unhandled platform: %s", definition.Platform)
			continue
		}

		app.Bots = append(app.Bots, b)
	}

	// Mount the API
	appApi := NewAppApi(router, app)
	appApi.Mount()

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), router)
}
