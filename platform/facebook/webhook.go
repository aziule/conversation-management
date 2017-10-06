package facebook

import (
	"net/http"
	"time"
	"fmt"
	"github.com/aziule/conversation-management/nlu"
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

	p, _ := parser.ParseText(m.Text)

	fmt.Println(m)
	fmt.Println(p)
	fmt.Println(p.Intent())
	fmt.Println(p.Entity())
}
