package conversation

import "time"

type message struct {
	senderId string
	recipientId string
	sentAt time.Time
	mid string
}

type textMessage struct {
	message
	text string
}

type quickReplyMessage struct {
	textMessage
	payload string
}

// Getters
func (m *message) SenderId() string { return m.senderId }
func (m *message) RecipientId() string { return m.recipientId }
func (m *message) SentAt() time.Time { return m.sentAt }
func (m *message) Mid() string { return m.mid }

func (m *textMessage) Text() string { return m.text }

func (m *quickReplyMessage) Payload() string { return m.payload }
