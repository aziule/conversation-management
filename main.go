package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/aziule/conversation-management/webhook"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", webhook.Validate)
	r.Post("/", webhook.ReceiveMessage)

	http.ListenAndServe(":3000", r)
}
