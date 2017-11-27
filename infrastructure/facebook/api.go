// Package facebook provides a Facebook API to be used by bots running on Facebook.
package facebook

import (
	"net/http"
	"net/url"

	"github.com/aziule/conversation-management/core/api"
	"github.com/aziule/conversation-management/core/utils"
)

// facebookApi is the real-world implementation of the API
type facebookApi struct {
	pageAccessToken string
	client          *http.Client
	baseUrl         *url.URL
}

// newFacebookApi is the constructor that creates a new Facebook API, using
// user-defined variables such as the FB API version or the pageAccessToken.
func newFacebookApi(conf utils.BuilderConf) (interface{}, error) {
	pageAccessToken, ok := utils.GetParam(conf, "page_access_token").(string)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("page_access_token")
	}

	version, ok := utils.GetParam(conf, "version").(string)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("version")
	}

	client, ok := utils.GetParam(conf, "client").(*http.Client)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("client")
	}

	rawBaseUrl := "https://graph.facebook.com/v" + version
	baseUrl, _ := url.Parse(rawBaseUrl)

	return &facebookApi{
		pageAccessToken: pageAccessToken,
		client:          client,
		baseUrl:         baseUrl,
	}, nil
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

func init() {
	api.RegisterFacebookApiBuilder("facebook", newFacebookApi)
}
