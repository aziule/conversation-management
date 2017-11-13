package facebook

import (
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
)

// processStepGetIntent processes the "get_intent" step
func (b *facebookBot) processStepGetIntent(step *conversation.Step, data *nlp.ParsedData) error {
	log.Info("Processing step")
	return nil
}
