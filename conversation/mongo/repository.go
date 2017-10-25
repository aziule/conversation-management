package mongo

import (
	"github.com/aziule/conversation-management/conversation"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoDbRepository struct {
	db *Db
}

func NewMongodbRepository(db *Db) conversation.Repository {
	return &mongoDbRepository{
		db: db,
	}
}

func (repository *mongoDbRepository) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := repository.db.Session.Clone()
	defer session.Close()

	return nil, nil
}

// FindUser tries to find a user based on its UserId
func (repository *mongoDbRepository) FindUser(userId string) (*conversation.User, error) {
	session := repository.db.NewSession()
	defer session.Close()

	user := &conversation.User{}

	// Tied to Facebook at the moment... Should use a specification pattern
	err := session.DB(repository.db.Params.DbName).C("user").Find(bson.M{
		"fbid": userId,
	}).One(user)

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}

		// @todo: handle and log
		return nil, err
	}

	return user, nil
}

func (repository *mongoDbRepository) InsertUser(user *conversation.User) error {
	session := repository.db.NewSession()
	defer session.Close()

	err := session.DB(repository.db.Params.DbName).C("user").Insert(user)

	if err != nil {
		// @todo: handle and log
		return err
	}

	return nil
}

func (repository *mongoDbRepository) SaveConversation(conversation *conversation.Conversation) error {
	session := repository.db.Session.Clone()
	defer session.Close()

	return nil
}
