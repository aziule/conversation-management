package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/aziule/conversation-management/bot"
	//"github.com/aziule/conversation-management/facebook"
	"github.com/aziule/conversation-management/facebook"
)

func main() {
	r := chi.NewRouter()

	b := bot.NewBot()

	r.Get("/", b.HandleValidateWebhook)
	r.Post("/", b.HandleMessageReceived)
api := facebook.NewFacebookApi(
	"2.6",
	"EAALQ00uMgf8BACxsTu75poGIqpkEtepvAZAZCzbuZC8TfQDv6wbZAFjWAjhaV0XwS7lNFcyZC8OrOe1AQrGeAFiCCIU683DekRRmDEy3B6EFsRshNsx8tP9SPusNcJ0Cty3Qt2HedwCUihFShFPbHXP5qZAuZBXCPAorZCNLGR8tAgZDZD",
	http.DefaultClient,
)
	err := api.SendTextToUser("1429733950458154", "HEY")

	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":3000", r)
}