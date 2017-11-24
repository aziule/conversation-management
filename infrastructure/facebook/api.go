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

// @todo: move the conf to utils
// newFacebookApi is the constructor that creates a new Facebook API, using
// user-defined variables such as the FB API version or the pageAccessToken.
func newFacebookApi(conf utils.BuilderConf) (api.FacebookApi, error) {
	// @todo: move this conf parsing to utils
	pageAccessTokenParam, ok := conf["page_access_token"]

	if !ok {
		return nil, utils.ErrUndefinedParam("page_access_token")
	}

	pageAccessToken, ok := pageAccessTokenParam.(string)

	if !ok {
		return nil, utils.ErrInvalidParam("page_access_token")
	}

	clientParam, ok := conf["client"]

	if !ok {
		return nil, utils.ErrUndefinedParam("client")
	}

	client, ok := clientParam.(*http.Client)

	if !ok {
		return nil, utils.ErrInvalidParam("client")
	}

	versionParam, ok := conf["version"]

	if !ok {
		return nil, utils.ErrUndefinedParam("version")
	}

	version, ok := versionParam.(string)

	if !ok {
		return nil, utils.ErrInvalidParam("version")
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
