// Package api provides a Facebook API to be used by bots running on Facebook.
package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

// FacebookApi is the interface representing a Facebook API
type FacebookApi interface {
	ParseRequestMessageReceived(r *http.Request) (*receivedMessage, error)
	SendTextToUser(recipientId, text string) error
}

// facebookApi is the real-world implementation of the API
type facebookApi struct {
	pageAccessToken string
	client          *http.Client
	baseUrl         *url.URL
}

// NewfacebookApi is the constructor that creates a new Facebook API, using
// user-defined variables such as the FB API version or the pageAccessToken.
func NewfacebookApi(version, pageAccessToken string, client *http.Client) *facebookApi {
	rawBaseUrl := "https://graph.facebook.com/v" + version
	baseUrl, _ := url.Parse(rawBaseUrl)

	return &facebookApi{
		pageAccessToken: pageAccessToken,
		client:          client,
		baseUrl:         baseUrl,
	}
}

// @todo: store it and avoid recreating it every time
// getSendTextUrl returns the url to ping to send text messages to a user
func (api *facebookApi) getSendTextUrl() *url.URL {
	baseUrl := api.baseUrl

	u, _ := url.Parse(baseUrl.String() + "/me/messages")

	q := u.Query()
	q.Set("access_token", api.pageAccessToken)

	u.RawQuery = q.Encode()

	return u
}

// prettyPrint prints JSON the pretty way
func prettyPrint(data []byte) {
	var v interface{}
	json.Unmarshal(data, &v)
	prettyJson, _ := json.MarshalIndent(v, "", "    ")
	os.Stdout.Write(prettyJson)
}
