package mongo

import (
	"github.com/aziule/conversation-management/conversation"
)

type MongoConversationWriter struct {
	db *Db
}

func NewMongoConversationWriter(db *Db) *MongoConversationWriter {
	return &MongoConversationWriter{
		db: db,
	}
}

func (writer *MongoConversationWriter) InsertUser(user *conversation.User) error {
	session := writer.db.Session.Clone()
	defer session.Close()

	return nil
}

func (writer *MongoConversationWriter) Save(conversation *conversation.Conversation) error {
	session := writer.db.Session.Clone()
	defer session.Close()

	return nil
}
