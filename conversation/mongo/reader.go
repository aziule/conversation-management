package mongo

import (
	"github.com/aziule/conversation-management/conversation"
)

type MongoConversationReader struct {
	db *Db
}

func NewMongoConversationReader(db *Db) *MongoConversationReader {
	return &MongoConversationReader{
		db: db,
	}
}

func (reader *MongoConversationReader) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := reader.db.Session.Clone()
	defer session.Close()

	return nil, nil
}

func (reader *MongoConversationReader) FindUser(userId conversation.UserId) (*conversation.User, error) {
	session := reader.db.Session.Clone()
	defer session.Close()

	return nil, nil
}
