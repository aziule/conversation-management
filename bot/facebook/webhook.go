package facebook

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// HandleMessageReceived is called when a new message is sent by the user to the page
// We parse the message, extract relevant NLU data, check the context, validate
// the data, and return a response.
func (bot *facebookBot) HandleMessageReceived(w http.ResponseWriter, r *http.Request) {
	receivedMessage, err := bot.fbApi.ParseRequestMessageReceived(r)

	if err != nil {
		panic(err)
	}

	if receivedMessage.NlpData() != nil {
		parsedData, err := ParseNlpData(receivedMessage.Nlp())

		if err != nil {
			panic(err)
		}

		fmt.Println(parsedData)
	}

	bot.fbApi.SendTextToUser(receivedMessage.SenderId(), receivedMessage.Text())
	//parsedData, err := bot.nluParser.ParseData(body)

	//if err != nil {
	//	panic(err)
	//}

	fmt.Println(receivedMessage)
	//fmt.Println(parsedData)
	fmt.Println("")
	fmt.Println("")
	//fmt.Println(string(body))
	//fmt.Println(parsedText)

	//// @todo: handle error here
	//step.HandleStep("step1")
	//
	//fmt.Println(message.SenderId())
	//fmt.Println(message.RecipientId())
}

// HandleValidateWebhook tries to validate the Facebook webhook
// More information here: https://developers.facebook.com/docs/messenger-platform/getting-started/quick-start
func (bot *facebookBot) HandleValidateWebhook(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	hubMode, err := getSingleQueryParam(queryParams, "hub.mode")

	if err != nil || hubMode != "subscribe" {
		return
	}

	verifyToken, err := getSingleQueryParam(queryParams, "hub.verify_token")

	if err != nil || verifyToken != bot.VerifyToken() {
		return
	}

	challenge, err := getSingleQueryParam(queryParams, "hub.challenge")

	if err != nil {
		return
	}

	// Validate the webhook by writing back the "hub.challenge" query param
	w.Write([]byte(challenge))
}

// getSingleQueryParam fetches a single query param using the given url values
func getSingleQueryParam(values url.Values, key string) (string, error) {
	params, ok := values[key]

	if !ok || len(params) != 1 {
		return "", errors.New("Could not fetch param")
	}

	return params[0], nil
}
