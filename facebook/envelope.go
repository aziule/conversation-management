package facebook

// recipientEnvelope is the envelope for a recipient
type recipientEnvelope struct {
	Id string `json:"id"`
}

// messageEnvelope represents the envelope for a message with text
type messageEnvelope struct {
	Text string `json:"text"`
}

// textToUserEnvelope is the JSON envelope that needs to be sent
type textToUserEnvelope struct {
	Recipient *recipientEnvelope `json:"recipient"`
	Message *messageEnvelope `json:"message"`
}

// newRecipientEnvelope is the constructor for a recipientEnvelope
func newRecipientEnvelope(recipientId string) *recipientEnvelope {
	return &recipientEnvelope{
		Id: recipientId,
	}
}

// newMessageEnvelope is the constructor for a messageEnvelope
func newMessageEnvelope(text string) *messageEnvelope {
	return &messageEnvelope{
		Text: text,
	}
}


// newTextToUserEnvelope is the constructor for a textToUserEnvelope
func newTextToUserEnvelope(recipientId, text string) *textToUserEnvelope {
	return &textToUserEnvelope{
		Recipient: &recipientEnvelope{
			Id: recipientId,
		},
		Message: &messageEnvelope{
			Text: text,
		},
	}
}
