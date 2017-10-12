package bot

import (
	"net/http"
)

type HttpMethod string

const (
	HTTP_METHOD_GET  HttpMethod = "GET"
	HTTP_METHOD_POST HttpMethod = "POST"
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

// Getters
func (webhook *Webhook) Method() HttpMethod        { return webhook.method }
func (webhook *Webhook) Path() string              { return webhook.path }
func (webhook *Webhook) Handler() http.HandlerFunc { return webhook.handler }
