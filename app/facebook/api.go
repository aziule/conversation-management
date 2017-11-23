package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/core/bot"
)

// bindDefaultApiEndpoints initialises the default API endpoints.
func (b *facebookBot) bindDefaultApiEndpoints() {
	b.apiEndpoints = append(b.apiEndpoints, bot.NewApiEndpoint(
		"GET",
		"/",
		b.handleViewBot,
	))
}

// handleViewBot shows details about the bot
func (b *facebookBot) handleViewBot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
