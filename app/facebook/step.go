package facebook

import (
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
)

func (b *facebookBot) ProcessStep1(step *conversation.Step, data *nlp.ParsedData) error {
	log.Info("Processing step")
	return nil
}

func (b *facebookBot) ProcessStep2(step *conversation.Step, data *nlp.ParsedData) error {
	return nil
}
