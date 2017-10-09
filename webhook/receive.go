package webhook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/aziule/conversation-management/conversation/message"
)

// When a new message is received from the user
func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	//bot := bot.NewFacebookBot()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Could not parse the request body", 500)
		return
	}

	message, err := message.NewMessageFromJson(body)
	//fmt.Println(json)
	fmt.Println(message)
	//
	//m := textMessage{
	//	message{
	//		"sid.123456",
	//		"rid.123456",
	//		time.Now(),
	//		"mid.123456",
	//	},
	//	"This is the text",
	//}
	//
	//fmt.Println(m)

	//parser := &nlu.Parser{}
	//parsed, _ := parser.ParseText(m.Text)
	//
	//user := &FacebookUser{
	//	uuid: "uuid",
	//	fbid: "fbid",
	//	name: "Raoul",
	//}

	//entrypoint := data.GetDummyEntrypoint()
	//
	//for _, startingStep := range entrypoint.Stories()[0].StartingSteps() {
	//	fmt.Println(startingStep.Name())
	//}

	//conversation.Progress(user, parsed)
}
