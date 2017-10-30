package mongo

import (
	"fmt"
	"time"

	"github.com/aziule/conversation-management/app/core/conversation"
	"gopkg.in/mgo.v2/bson"
)

// mongoConversation is the representation of a Conversation for a mongodb storage.
// We use this structure in order to:
//     - Avoid mixing mongodb logic with our domain logic
//     - Be able to unmarshal interfaces back to the domain object
type mongoConversation struct {
	Id        bson.ObjectId       `bson:"_id"`
	Status    conversation.Status `bson:"status"`
	Messages  mongoMessagesList   `bson:"messages"`
	CreatedAt time.Time           `bson:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at"`
}

// mongoMessagesList is the representation of a MessagesList for a mongodb storage.
// However we do not store a class slice of Message, but rather a slice of a new
// struct type: mongoMessageWithType
type mongoMessagesList []*mongoMessageWithType

// mongoMessageWithType is the representation of a Message for a mongodb storage.
// We use this intermediate structure as a layer between mongodb and the domain Message,
// so that we can easily unmarshal the proper Message back to the Conversation, using
// the Type, in the SetBSON method
type mongoMessageWithType struct {
	Type    conversation.MessageType
	Message conversation.Message
}

// toMongoMessagesList converts a MessagesList to a mongoMessagesList
func toMongoMessagesList(list conversation.MessagesList) mongoMessagesList {
	var mongoList mongoMessagesList

	for _, message := range list {
		mongoList = append(mongoList, toMessageWithType(message))
	}

	return mongoList
}

// toMessagesList converts a mongoMessagesList to a MessagesList
func toMessagesList(mongoList mongoMessagesList) conversation.MessagesList {
	var list conversation.MessagesList

	for _, message := range mongoList {
		list = append(list, message.Message)
	}

	return list
}

// toMessageWithType converts a Message to a mongoMessageWithType
func toMessageWithType(message conversation.Message) *mongoMessageWithType {
	return &mongoMessageWithType{
		Type:    message.Type(),
		Message: message,
	}
}

// toMongoConversation converts a Conversation to a mongoConversation
func toMongoConversation(c *conversation.Conversation) *mongoConversation {
	var id bson.ObjectId

	if c.Id != "" {
		id = bson.ObjectIdHex(c.Id)
	}

	return &mongoConversation{
		Id:        id,
		Status:    c.Status,
		Messages:  toMongoMessagesList(c.Messages),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// toMongoConversation converts a mongoConversation to a Conversation
func toConversation(mongoConversation *mongoConversation) *conversation.Conversation {
	// @todo: fix this constructor
	return &conversation.Conversation{
		Id:        mongoConversation.Id.Hex(),
		Status:    mongoConversation.Status,
		Messages:  toMessagesList(mongoConversation.Messages),
		CreatedAt: mongoConversation.CreatedAt,
		UpdatedAt: mongoConversation.UpdatedAt,
	}
}

// @todo: absolutely test this method to check that the correct Message type is returned in every case.
// SetBSON is called when unmarshalling a mongoMessagesList, which is called
// when unmarshalling a mongoConversation.
// Thanks to the mongoMessagesList, we are able to know what kind of Message we want to
// return, and so we can easily manage the Message interface storage / fetching.
func (m *mongoMessagesList) SetBSON(raw bson.Raw) error {
	fmt.Println("In SetBSON: MessagesList")

	//*m = append(*m, conversation.NewUserMessage(
	//	"Test",
	//	time.Now(),
	//	nil,
	//	nil,
	//))
	//*m = append(*m, conversation.NewUserMessage(
	//	"Test",
	//	time.Now(),
	//	nil,
	//	nil,
	//))

	fmt.Println(m)

	return nil
}
