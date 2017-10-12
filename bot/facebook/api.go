package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aziule/conversation-management/core/bot"
	"io/ioutil"
	"net/http"
	"net/url"
)

// facebookApi is the real-world implementation of the API
type FacebookApi struct {
	version         string
	pageAccessToken string
	client          *http.Client
	webhooks        []*bot.Webhook
}

// NewFacebookApi is the factory method to create a new facebook api implementation
func NewFacebookApi(version, pageAccessToken string, client *http.Client) *FacebookApi {
	return &FacebookApi{
		version:         version,
		pageAccessToken: pageAccessToken,
		client:          client,
	}
}

// Getters
func (api *FacebookApi) Version() string         { return api.version }
func (api *FacebookApi) PageAccessToken() string { return api.pageAccessToken }

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

// SendTextToUser is the FacebookApi's interface method responsible for sending a 1-to-1 message to a user
func (api *FacebookApi) SendTextToUser(recipientId, text string) error {
	url := api.getSendTextUrl()

	object := newTextToUserEnvelope(recipientId, text)
	jsonObject, err := json.Marshal(object)

	if err != nil {
		// @todo: fix
		return err
	}

	request, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonObject))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		// @todo: fix
		return err
	}

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		// @todo: fix
		return err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))

	return nil
}

// getSendTextUrl returns the url to ping to send text messages to a user
func (api *FacebookApi) getSendTextUrl() *url.URL {
	baseUrl := api.getBaseUrl()

	u, _ := url.Parse(baseUrl.String() + "/me/messages")

	q := u.Query()
	q.Set("access_token", api.pageAccessToken)

	u.RawQuery = q.Encode()

	return u
}
