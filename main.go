package main

import (
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/aziule/conversation-management/bot"
	"github.com/aziule/conversation-management/bot/facebook"
	"github.com/aziule/conversation-management/conversation/mongo"
	"github.com/aziule/conversation-management/nlp/wit"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var configFlagPath = flag.String("config", "config.json", "Config file path")

func main() {
	flag.Parse()

	config, err := LoadConfig(*configFlagPath)

	if err != nil {
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	dbParams := mongo.DbParams{
		DbHost: config.DbHost,
		DbName: config.DbName,
		DbUser: config.DbUser,
		DbPass: config.DbPass,
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{dbParams.DbHost},
		Database: dbParams.DbName,
		Timeout:  2 * time.Second,
	})

	if err != nil {
		log.Fatalf("An error occurred when connecting to the db: %s", err)
	}

	defer session.Close()

	b := facebook.NewBot(
		&facebook.Config{
			VerifyToken:     config.FbVerifyToken,
			ApiVersion:      config.FbApiVersion,
			PageAccessToken: config.FbPageAccessToken,
			NlpParser:       wit.NewParser(facebook.DefaultDataTypeMap),
			ConversationRepository: mongo.NewMongodbRepository(&mongo.Db{
				Session: session,
				Params:  dbParams,
			}),
		},
	)

	r := chi.NewRouter()

	// Automatically set the bot's webhooks routes
	for _, webhook := range b.Webhooks() {
		log.Debugf("%s %s", string(webhook.Method), webhook.Path)

		switch webhook.Method {
		case bot.HttpMethodGet:
			r.Get(webhook.Path, webhook.Handler)
		case bot.HttpMethodPost:
			r.Post(webhook.Path, webhook.Handler)
		}
	}

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), r)
}
