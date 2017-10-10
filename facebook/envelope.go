package facebook

type recipientEnvelope struct {
	id string
}

func newRecipientEnvelope(recipientId string) *recipientEnvelope {
	return &recipientEnvelope{
		id: recipientId,
	}
}

type messageEnvelope struct {
	text string
}

func newMessageEnvelope(text string) *messageEnvelope {
	return &messageEnvelope{
		text: text,
	}
}

// textToUserEnvelope is the JSON envelope that needs to be sent
type textToUserEnvelope struct {
	recipient *recipientEnvelope `json:"recipient"`
	message *messageEnvelope
}

func newTextToUserEnvelope(recipientId, text string) *textToUserEnvelope {
	return &textToUserEnvelope{
		recipient: &recipientEnvelope{
			id: recipientId,
		},
		message: &messageEnvelope{
			text: text,
		},
	}
}
