// mongo package provides all of the required methods to interact with
// a mongodb database, using mgo as the driver.
package mongo

import (
	"time"

	"github.com/aziule/conversation-management/core/conversation"
	db "github.com/aziule/conversation-management/infrastructure/mongo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongoDbRepository is the unexported struct that implements the Repository interface
type mongoDbRepository struct {
	db *db.Db
}

// @todo: give it a variable for the mapping between messages <=> facebook messages (implementation)
// NewMongodbRepository creates a new conversation repository using MongoDb as the data source
func NewMongodbRepository(db *db.Db) conversation.Repository {
	return &mongoDbRepository{
		db: db,
	}
}

// SaveConversation saves a conversation to the database.
// The conversation can be an existing one or a new one
func (repository *mongoDbRepository) SaveConversation(c *conversation.Conversation) error {
	session := repository.db.NewSession()
	defer session.Close()

	var err error
	collection := session.DB(repository.db.Params.DbName).C("conversation")

	c.UpdatedAt = time.Now()

	if c.Id == "" {
		c.Id = bson.NewObjectId()
		c.CreatedAt = time.Now()

		log.WithField("conversation", c).Debugf("Inserting conversation")
		err = collection.Insert(c)
	} else {
		log.WithField("conversation", c).Debugf("Updating conversation")
		err = collection.UpdateId(c.Id, c)
	}

	if err != nil {
		log.WithField("conversation", c).Infof("Could not save the conversation: %s", err)
		return err
	}

	return nil
}

// FindLatestConversation tries to find the latest conversation that happened with a user.
// In case this is a new user, then no conversation is returned. Otherwise the latest one,
// which can be the current one, is returned.
// Returns a conversation.ErrNotFound error when the user is not found.
func (repository *mongoDbRepository) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := repository.db.NewSession()
	defer session.Close()

	// Store the result of the query in our own mongo struct
	var c *conversation.Conversation

	log.WithField("fbid", user.FbId).Debug("Finding latest conversation for user")

	err := session.DB(repository.db.Params.DbName).C("conversation").Find(bson.M{
		"messages": bson.M{
			"$elemMatch": bson.M{
				"message.sender_id": user.Id,
			},
		},
	}).Sort("-created_at").One(&c)

	if err != nil {
		if err == mgo.ErrNotFound {
			log.Debug("Latest conversation not found")
			return nil, conversation.ErrNotFound
		}

		log.Infof("Could not find the latest conversation: %s", err)
		return nil, err
	}

	log.WithField("conversation", c.Id).Debug("Found latest conversation")

	return c, nil
}

// FindUserByFbId tries to find a user based on its fbId
// Returns a conversation.ErrNotFound error when the user is not found
// @todo: we should use a specification pattern
func (repository *mongoDbRepository) FindUserByFbId(fbId string) (*conversation.User, error) {
	session := repository.db.NewSession()
	defer session.Close()

	user := &conversation.User{}

	err := session.DB(repository.db.Params.DbName).C("user").Find(bson.M{
		"fbid": fbId,
	}).One(user)

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, conversation.ErrNotFound
		}

		log.WithField("fbId", fbId).Infof("Could not find the user: %s", err)
		return nil, err
	}

	return user, nil
}

// InsertUser creates a new user in the DB
func (repository *mongoDbRepository) InsertUser(user *conversation.User) error {
	// @todo: if the user has an ID, return an error
	// @todo: check that the user does not exist yet

	session := repository.db.NewSession()
	defer session.Close()

	err := session.DB(repository.db.Params.DbName).C("user").Insert(user)

	if err != nil {
		log.WithField("user", user).Infof("Could not insert the user: %s", err)
		return err
	}

	return nil
}
