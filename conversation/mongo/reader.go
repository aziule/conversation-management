package mongo

import (
	"github.com/aziule/conversation-management/conversation"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

type MongoConversationReader struct {
	session *mgo.Session
}

func NewMongoConversationReader() *MongoConversationReader {
	return &MongoConversationReader{
		session: session.Copy(),
	}
}

func (reader *MongoConversationReader) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	return nil, nil
}
