package conversation

import (
	"github.com/aziule/conversation-management/nlu"
	"fmt"
)

// Validate a given step and return either the next step or an error
func Progress(user User, parsedText *nlu.ParsedText) {
	context, _ := getContext(user)

	fmt.Println(context)
}
