package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/aziule/conversation-management/platform/facebook"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", facebook.ValidateWebhook)
	r.Post("/", facebook.MessageReceived)

	http.ListenAndServe(":3000", r)
}
