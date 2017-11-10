package facebook

import (
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
)

// ProcessStepGetIntent processes the "get_intent" step
func (b *facebookBot) ProcessStepGetIntent(step *conversation.Step, data *nlp.ParsedData) error {
	log.Info("Processing step")
	return nil
}
