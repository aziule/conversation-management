package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

var baseUrl string

// FacebookApi is the real-world implementation of the API
type FacebookApi struct {
	Version         string
	PageAccessToken string
	client          *http.Client
}

// NewFacebookApi is the constructor that creates a new Facebook API, using
// user-defined variables such as the FB API version or the pageAccessToken.
func NewFacebookApi(version, pageAccessToken string, client *http.Client) *FacebookApi {
	return &FacebookApi{
		Version:         version,
		PageAccessToken: pageAccessToken,
		client:          client,
	}
}

// @todo: store it and avoid recreating it every time
// getSendTextUrl returns the url to ping to send text messages to a user
func (api *FacebookApi) getSendTextUrl() *url.URL {
	baseUrl := api.getBaseUrl()

	u, _ := url.Parse(baseUrl.String() + "/me/messages")

	q := u.Query()
	q.Set("access_token", api.PageAccessToken)

	u.RawQuery = q.Encode()

	return u
}

// getBaseUrl returns the base url for a Facebook graph API call
func (api *FacebookApi) getBaseUrl() *url.URL {
	rawUrl := "https://graph.facebook.com/v" + api.Version
	u, _ := url.Parse(rawUrl)

	return u
}

// prettyPrint prints JSON the pretty way
func prettyPrint(data []byte) {
	var v interface{}
	json.Unmarshal(data, &v)
	prettyJson, _ := json.MarshalIndent(v, "", "    ")
	os.Stdout.Write(prettyJson)
}
