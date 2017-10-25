package mongo

import (
	"github.com/aziule/conversation-management/conversation"
	"gopkg.in/mgo.v2"
)

type MongoConversationWriter struct {
	session *mgo.Session
}

func NewMongoConversationWriter(session *mgo.Session) *MongoConversationWriter {
	return &MongoConversationWriter{
		session: session,
	}
}

func (writer *MongoConversationWriter) InsertUser(user *conversation.User) error {
	session := writer.session.Clone()
	defer session.Close()

	return nil
}

func (writer *MongoConversationWriter) Save(conversation *conversation.Conversation) error {
	session := writer.session.Clone()
	defer session.Close()

	return nil
}
