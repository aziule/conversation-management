package bot

import "net/http"

type Webhook struct {
	method string
	handler func(http.ResponseWriter, *http.Request)
}

// Getters
func (webhook *Webhook) Method() string { return webhook.method }
func (webhook *Webhook) Handler() func(http.ResponseWriter, *http.Request) { return webhook.handler }