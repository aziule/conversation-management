package bot

import (
	"net/http"
)

// Endpoint is the struct that is reponsible for opening the bot to the world (API endpoints, Webhooks, etc.)
type Endpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Webhook represents an endpoint called by the platform
type Webhook Endpoint

// ApiEndpoint represents an endpoint accessible on the microservice
type ApiEndpoint Endpoint

// NewWebhook is the constructor method for Webhook
func NewWebHook(method, path string, handler http.HandlerFunc) *Webhook {
	return &Webhook{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

// NewApiEndpoint is the constructor method for ApiEndpoint
func NewApiEndpoint(method, path string, handler http.HandlerFunc) *ApiEndpoint {
	return &ApiEndpoint{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
