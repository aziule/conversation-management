package facebook

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HandleMessageReceived is called when a new message is sent by the user to the page
func HandleMessageReceived(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Could not parse the request body", 500)
		return
	}

	message, err := NewMessageFromJson(body)
	// @todo: handle error here

	fmt.Println(message)
	fmt.Println(message.SenderId())
	fmt.Println(message.RecipientId())
}
