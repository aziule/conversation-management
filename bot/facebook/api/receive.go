package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/core/nlu"
	"io/ioutil"
	"net/http"
)

// Message is the base struct for messages
type Message struct {
	mid               string
	senderId          string
	recipientId       string
	sentAt            time.Time
	text              string
	quickReplyPayload string
	parsedData        *nlu.ParsedData
}

// ParseJsonBody creates a Message from json bytes and returns an error if a parsing issue occurred
func (api *FacebookApi) ParseRequestMessageReceived(r *http.Request) (*Message, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		// @todo: log and handle error types
		return nil, err
	}

	json, err := jason.NewObjectFromBytes(body)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not parse JSON") // @todo: error types
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

	//nlp, err := messageData.GetObject("message", "nlp", "entities")
	//nlpData := make(map[string]interface{})
	//
	//if err == nil {
	//	for key, data := range nlp.Map() {
	//		nlpDataStruct, err := makeNlpData(key, data)
	//
	//		if err != nil {
	//			// @todo: log
	//			fmt.Println("Error", err)
	//			continue
	//		}
	//
	//		nlpData[key] = nlpDataStruct
	//	}
	//}

	return &Message{
		mid:               mid,
		senderId:          senderId,
		recipientId:       recipientId,
		sentAt:            time.Unix(sentAt, 0),
		text:              text,
		quickReplyPayload: quickReplyPayload,
		//nlpData: nlpData,
	}, nil
}

// Getters
func (m *Message) SenderId() string            { return m.senderId }
func (m *Message) RecipientId() string         { return m.recipientId }
func (m *Message) SentAt() time.Time           { return m.sentAt }
func (m *Message) Mid() string                 { return m.mid }
func (m *Message) Text() string                { return m.text }
func (m *Message) QuickReplyPayload() string   { return m.quickReplyPayload }
func (m *Message) ParsedData() *nlu.ParsedData { return m.parsedData }
