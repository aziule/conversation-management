package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/aziule/conversation-management/bot/facebook"
	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

var configFlagPath = flag.String("config", "config.json", "Config file path")

func main() {
	flag.Parse()

	config, err := core.LoadConfig(*configFlagPath)

	if err != nil {
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	b := facebook.NewFacebookBot(config)

	r := chi.NewRouter()

	// Automatically set the bot's webhooks routes
	for _, webhook := range b.Webhooks() {
		log.Debugf("%s %s", string(webhook.Method()), webhook.Path())

		switch webhook.Method() {
		case bot.HttpMethodGet:
			r.Get(webhook.Path(), webhook.Handler())
		case bot.HttpMethodPost:
			r.Post(webhook.Path(), webhook.Handler())
		}
	}

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), r)
}
