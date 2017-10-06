package facebook

import (
	"net/http"
	"time"
	"fmt"
	"github.com/aziule/conversation-management/nlu"
	"github.com/aziule/conversation-management/conversation"
)

// When a new message is received from the user
func MessageReceived(w http.ResponseWriter, r *http.Request) {
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

	conversation.Progress(user, parsed)

	fmt.Println(m)
	fmt.Println(parsed)
	fmt.Println(parsed.Intent())
	fmt.Println(parsed.Entity())
}
