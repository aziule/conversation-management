package message

import "time"

// Message is the base struct for messages
type Message struct {
	mid string
	senderId string
	recipientId string
	sentAt time.Time
	text string
	quickReplyPayload string
}

// Getters
func (m *Message) SenderId() string { return m.senderId }
func (m *Message) RecipientId() string { return m.recipientId }
func (m *Message) SentAt() time.Time { return m.sentAt }
func (m *Message) Mid() string { return m.mid }
func (m *Message) Text() string { return m.text }
func (m *Message) QuickReplyPayload() string { return m.quickReplyPayload }
