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
api := facebook.NewFacebookApi("2.6", "EAALQ00uMgf8BAHnz4y711RUBOQyiQrUpy4ZAIXeXvL4L0mIZAZCq6WXKZBnwwhT8Xfw2So5DZABaRSfxjuO97mdQklTxZCdZATKFH7xvJ5VEwqsCyQRTXh9yTq9ZBSGATaSZCSsS7xhv3TeHvFyx5s0xcQ88BxiZBqmv8zPFRcTX9iJAZDZD")
	err := api.SendTextToUser("1429733950458154", "HEY")

	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":3000", r)
}
