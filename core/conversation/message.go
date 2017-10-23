package conversation

import (
	"time"

	"github.com/aziule/conversation-management/core/nlp"
)

type Message interface {
	Text() string
	Sender() User
	ParsedData() nlp.ParsedData
	SentAt() time.Time
}
