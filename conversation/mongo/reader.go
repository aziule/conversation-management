package mongo

import (
	"github.com/aziule/conversation-management/conversation"
	"gopkg.in/mgo.v2"
)

type MongoConversationReader struct {
	session *mgo.Session
}

func NewMongoConversationReader(session *mgo.Session) *MongoConversationReader {
	return &MongoConversationReader{
		session: session,
	}
}

func (reader *MongoConversationReader) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := reader.session.Clone()
	defer session.Close()

	return nil, nil
}

func (reader *MongoConversationReader) FindUser(userId conversation.UserId) (*conversation.User, error) {
	session := reader.session.Clone()
	defer session.Close()

	return nil, nil
}
