package facebook

import (
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
)

// processStepGetIntent processes the "book_table_entrypoint" step
func (b *facebookBot) processStepBookTable(step *conversation.Step, data *nlp.ParsedData) error {
	log.Info("BOOK TABLE")
	return nil
}

// processStepBookTableGetNbPersons processes the "book_table_get_nb_persons" step
func (b *facebookBot) processStepBookTableGetNbPersons(step *conversation.Step, data *nlp.ParsedData) error {
	log.Info("BOOK TABLE - GET NB PERSONS")
	return nil
}

// processStepGetIntent processes the "book_table_get_time" step
func (b *facebookBot) processStepBookTableGetTime(step *conversation.Step, data *nlp.ParsedData) error {
	log.Info("BOOK TABLE - GET TIME")
	return nil
}
