package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

// FacebookApi is the real-world implementation of the API
type FacebookApi struct {
	pageAccessToken string
	client          *http.Client
	baseUrl         *url.URL
}

// NewFacebookApi is the constructor that creates a new Facebook API, using
// user-defined variables such as the FB API version or the pageAccessToken.
func NewFacebookApi(version, pageAccessToken string, client *http.Client) *FacebookApi {
	rawBaseUrl := "https://graph.facebook.com/v" + version
	baseUrl, _ := url.Parse(rawBaseUrl)

	return &FacebookApi{
		pageAccessToken: pageAccessToken,
		client:          client,
		baseUrl:         baseUrl,
	}
}

// @todo: store it and avoid recreating it every time
// getSendTextUrl returns the url to ping to send text messages to a user
func (api *FacebookApi) getSendTextUrl() *url.URL {
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
