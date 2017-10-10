package facebook

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// FacebookApi is the interface that we will use for any API call needed to do on Facebook
type FacebookApi interface {
	SendTextToUser(recipientId, text string) error
}

// facebookApi is the real-world implementation of the API
type facebookApi struct {
	version string
	pageAccessToken string
}

// Getters
func (api *facebookApi) Version() string { return api.version }
func (api *facebookApi) PageAccessToken() string { return api.pageAccessToken }

// getBaseUrl returns the base url for a Facebook graph API call
func (api *facebookApi) getBaseUrl() string {
	return fmt.Sprintf("https://graph.facebook.com/v%s", api.Version())
}

// NewFacebookApi is the factory method to create a new facebook api implementation
func NewFacebookApi(version, pageAccessToken string) FacebookApi {
	return &facebookApi{
		version: version,
		pageAccessToken: pageAccessToken,
	}
}

// SendTextToUser is the FacebookApi's interface method responsible for sending a 1-to-1 message to a user
func (api *facebookApi) SendTextToUser(recipientId, text string) error {
	url := api.getSendTextUrl()

	object := newTextToUserEnvelope(recipientId, text)
	jsonObject, err := json.Marshal(object)

	if err != nil {
		// @todo: fix
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonObject))
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
func (api *facebookApi) getSendTextUrl() string {
	var buffer bytes.Buffer

	buffer.WriteString(api.getBaseUrl())
	buffer.WriteString("/me/messages?access_token=")
	buffer.WriteString(api.pageAccessToken)

	return buffer.String()
}
