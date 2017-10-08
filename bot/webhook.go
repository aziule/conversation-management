package bot

import "net/http"

type Webhook struct {
	method string
	handler func(http.ResponseWriter, *http.Request)
}

// Getters
func (webhook *Webhook) Method() { return webhook.method }
func (webhook *Webhook) Handler() { return webhook.handler }