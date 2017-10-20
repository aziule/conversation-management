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
	method  HttpMethod
	path    string
	handler http.HandlerFunc
}

// NewWebhook is the constructor method for Webhook
func NewWebHook(method HttpMethod, path string, handler http.HandlerFunc) *Webhook {
	return &Webhook{
		method:  method,
		path:    path,
		handler: handler,
	}
}

// Method returns the webhook's method.
// This method is required in order to implement the Webhook interface
func (webhook *Webhook) Method() HttpMethod {
	return webhook.method
}

// Path returns the webhook's path.
// This method is required in order to implement the Webhook interface
func (webhook *Webhook) Path() string {
	return webhook.path
}

// Handler returns the webhook's handler func.
// This method is required in order to implement the Webhook interface
func (webhook *Webhook) Handler() http.HandlerFunc {
	return webhook.handler
}
