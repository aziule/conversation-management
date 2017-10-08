package facebook

import (
	"net/http"
	"time"
	"fmt"
	"github.com/aziule/conversation-management/nlu"
	"github.com/aziule/conversation-management/conversation"
	"github.com/aziule/conversation-management/test/data"
	"net/url"
	"errors"
	"io/ioutil"
	"github.com/antonholmquist/jason"
)

// When a new message is received from the user
func MessageReceived(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Could not parse the request body", 500)
		return
	}

	json, err := jason.NewValueFromBytes(body)

	fmt.Println(json)

	m := textMessage{
		message{
			"sid.123456",
			"rid.123456",
			time.Now(),
			"mid.123456",
		},
		"This is the text",
	}

	parser := &nlu.Parser{}
	parsed, _ := parser.ParseText(m.Text)

	user := &FacebookUser{
		uuid: "uuid",
		fbid: "fbid",
		name: "Raoul",
	}

	entrypoint := data.GetDummyEntrypoint()

	for _, startingStep := range entrypoint.Stories()[0].StartingSteps() {
		fmt.Println(startingStep.Name())
	}

	conversation.Progress(user, parsed)
}

// Try to validate the Facebook webhook
// More information here: https://developers.facebook.com/docs/messenger-platform/getting-started/quick-start
func ValidateWebhook(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	hubMode, err := getSingleQueryParam(queryParams, "hub.mode");

	if err != nil || hubMode != "subscribe" {
		return
	}

	verifyToken, err := getSingleQueryParam(queryParams, "hub.verify_token");

	if err != nil || verifyToken != "app_verify_token" {
		return
	}

	challenge, err := getSingleQueryParam(queryParams, "hub.challenge");

	if err != nil {
		return
	}

	// Validate the webhook by writing back the "hub.challenge" query param
	w.Write([]byte(challenge))
}

// Get a single query param using given url values
func getSingleQueryParam(values url.Values, key string) (string, error) {
	params, ok := values[key]

	if (!ok || len(params) != 1) {
		return "", errors.New("Could not fetch param")
	}

	return params[0], nil
}