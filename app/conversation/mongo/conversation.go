package mongo

import (
	"fmt"
	"time"

	"github.com/aziule/conversation-management/app/conversation"
	"gopkg.in/mgo.v2/bson"
)

type mongoConversation struct {
	Id        bson.ObjectId       `bson:"_id"`
	Status    conversation.Status `bson:"status"`
	Messages  mongoMessagesList   `bson:"messages"`
	CreatedAt time.Time           `bson:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at"`
}

type mongoMessagesList []*mongoMessageWithType

type mongoMessageWithType struct {
	Type    conversation.MessageType
	Message conversation.Message
}

func toMongoMessagesList(list conversation.MessagesList) mongoMessagesList {
	var mongoList mongoMessagesList

	for _, message := range list {
		mongoList = append(mongoList, newMessageWithType(message))
	}

	return mongoList
}

func toMessagesList(mongoList mongoMessagesList) conversation.MessagesList {
	var list conversation.MessagesList

	for _, message := range mongoList {
		list = append(list, message.Message)
	}

	return list
}

func newMessageWithType(message conversation.Message) *mongoMessageWithType {
	return &mongoMessageWithType{
		Type:    message.Type(),
		Message: message,
	}
}

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

// @todo: absolutely test this method to check that the correct Message type
// is returned in every case.
// SetBSON is called when unmarshalling a conversation. Here, we need to explicitely
// recreate the messages that we want, as a regular Message is simply an interface.
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
