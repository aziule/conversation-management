package bot

import (
	"net/http"
)

type HttpMethod string

const (
	HttpMethodGet  HttpMethod = "GET"
	HttpMethodPost HttpMethod = "POST"
)

// Webhook is the struct that is reponsible for for making the link between a bot and its platform
type Webhook struct {
	Method  HttpMethod
	Path    string
	Handler http.HandlerFunc
}

// NewWebhook is the constructor method for Webhook
func NewWebHook(method HttpMethod, path string, handler http.HandlerFunc) *Webhook {
	return &Webhook{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
