package api

import (
	"net/http"
	"net/url"
)

// FacebookApi is the real-world implementation of the API
type FacebookApi struct {
	version         string
	pageAccessToken string
	client          *http.Client
}

// NewFacebookApi is the constructor that creates a new Facebook API, using
// user-defined variables such as the FB API version or the pageAccessToken.
func NewFacebookApi(version, pageAccessToken string, client *http.Client) *FacebookApi {
	return &FacebookApi{
		version:         version,
		pageAccessToken: pageAccessToken,
		client:          client,
	}
}

// @todo: store it and avoid recreating it every time
// getSendTextUrl returns the url to ping to send text messages to a user
func (api *FacebookApi) getSendTextUrl() *url.URL {
	baseUrl := api.getBaseUrl()

	u, _ := url.Parse(baseUrl.String() + "/me/messages")

	q := u.Query()
	q.Set("access_token", api.pageAccessToken)

	u.RawQuery = q.Encode()

	return u
}

// getBaseUrl returns the base url for a Facebook graph API call
func (api *FacebookApi) getBaseUrl() *url.URL {
	rawUrl := "https://graph.facebook.com/v" + api.Version()
	u, err := url.Parse(rawUrl)

	if err != nil {
		// @todo improve
		panic(err)
	}

	return u
}

// Getters
func (api *FacebookApi) Version() string         { return api.version }
func (api *FacebookApi) PageAccessToken() string { return api.pageAccessToken }
