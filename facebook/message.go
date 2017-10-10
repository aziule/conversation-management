package facebook

import (
	"errors"
	"fmt"
	"time"

	"github.com/antonholmquist/jason"
)

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

// NewMessageFromJson creates a Message from json bytes and returns an error if a parsing issue occurred
func NewMessageFromJson(bytes []byte) (*Message, error) {
	json, err := jason.NewObjectFromBytes(bytes)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not parse JSON")
	}

	// The message content itself is embedded inside the "entry" array
	entries, err := json.GetObjectArray("entry")

	if err != nil {
		return nil, errors.New("Missing entry key")
	}

	if len(entries) == 0 {
		return nil, errors.New("No message to parse")
	}

	entry := entries[0]

	messaging, err := entry.GetObjectArray("messaging")

	if err != nil {
		return nil, errors.New("Missing messaging key")
	}

	if len(messaging) == 0 {
		return nil, errors.New("No message to parse")
	}

	messageData := messaging[0]

	mid, err := messageData.GetString("message", "mid")

	if err != nil {
		return nil, errors.New("Missing message.id")
	}

	senderId, err := messageData.GetString("sender", "id")

	if err != nil {
		return nil, errors.New("Missing sender.id")
	}

	recipientId, err := messageData.GetString("recipient", "id")

	if err != nil {
		return nil, errors.New("Missing recipient.id")
	}

	sentAt, err := messageData.GetInt64("timestamp")

	if err != nil {
		return nil, errors.New("Missing timestamp")
	}

	text, _ := messageData.GetString("message", "text")
	quickReplyPayload, _ := messageData.GetString("quick_reply", "payload")

	return &Message{
		mid: mid,
		senderId: senderId,
		recipientId: recipientId,
		sentAt: time.Unix(sentAt, 0),
		text: text,
		quickReplyPayload: quickReplyPayload,
	}, nil
}
