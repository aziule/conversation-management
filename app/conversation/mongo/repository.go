package mongo

import (
	"github.com/aziule/conversation-management/app/conversation"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongoDbRepository is the unexported struct that implements the Repository interface
type mongoDbRepository struct {
	db *Db
}

// NewMongodbRepository creates a new conversation repository using MongoDb as the data source
func NewMongodbRepository(db *Db) conversation.Repository {
	return &mongoDbRepository{
		db: db,
	}
}

// FindLatestConversation tries to find the latest conversation that happened with a user.
// In case this is a new user, then no conversation is returned. Otherwise the latest one,
// which can be the current one, is returned.
// Returns a conversation.ErrNotFound error when the user is not found.
func (repository *mongoDbRepository) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := repository.db.Session.Clone()
	defer session.Close()

	return nil, conversation.ErrNotFound
}

// FindUser tries to find a user based on its UserId
// Returns a conversation.ErrNotFound error when the user is not found
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
			return nil, conversation.ErrNotFound
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

// SaveConversation saves a conversation to the database.
// The conversation can be an existing one or a new one
func (repository *mongoDbRepository) SaveConversation(c *conversation.Conversation) error {
	session := repository.db.Session.Clone()
	defer session.Close()

	var err error
	collection := session.DB(repository.db.Params.DbName).C("conversation")

	if c.Id == "" {
		log.Infof("Inserting conversation: %s", c.Id.String())
		c.Id = bson.NewObjectId()
		err = collection.Insert(c)
	} else {
		log.Infof("Updating conversation: %s", c.Id.String())
		err = collection.Update(bson.M{"_id,omitempty": c.Id}, c)
	}

	if err != nil {
		// @todo: handle and log
		return err
	}

	return nil
}
