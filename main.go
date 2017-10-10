package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/facebook"
	"github.com/go-chi/chi"
)

var configFlagPath = flag.String("config", "config.json", "Config file path")

func main() {
	flag.Parse()

	config, err := core.LoadConfig(*configFlagPath)

	if err != nil {
		panic(err)
	}

	bot := facebook.NewFacebookBot(config)

	r := chi.NewRouter()

	for _, webhook := range bot.Webhooks() {
		switch webhook.Method() {
		case core.HTTP_METHOD_GET:
			r.Get(webhook.Path(), webhook.Handler())
		case core.HTTP_METHOD_POST:
			r.Get(webhook.Path(), webhook.Handler())
		}
	}

	http.ListenAndServe(":" + strconv.Itoa(config.ListeningPort), r)
}