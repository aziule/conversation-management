// mongo package provides all of the required methods to interact with
// a mongodb database, using mgo as the driver.
package mongo

import (
	"time"

	"github.com/aziule/conversation-management/app/conversation"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongoDbRepository is the unexported struct that implements the Repository interface
type mongoDbRepository struct {
	db *Db
}

// @todo: give it a variable for the mapping between messages <=> facebook messages (implementation)
// NewMongodbRepository creates a new conversation repository using MongoDb as the data source
func NewMongodbRepository(db *Db) conversation.Repository {
	return &mongoDbRepository{
		db: db,
	}
}

// SaveConversation saves a conversation to the database.
// The conversation can be an existing one or a new one
func (repository *mongoDbRepository) SaveConversation(c *conversation.Conversation) error {
	session := repository.db.Session.Clone()
	defer session.Close()

	var err error
	collection := session.DB(repository.db.Params.DbName).C("conversation")

	// Convert the conversation to our own mongo object
	converted := toMongoConversation(c)
	converted.UpdatedAt = time.Now()

	if c.Id == "" {
		converted.Id = bson.NewObjectId()
		converted.CreatedAt = time.Now()

		log.WithField("conversation", c).Debugf("Inserting conversation: %s", converted.Id)
		err = collection.Insert(converted)
	} else {
		log.WithField("conversation", converted).Debugf("Updating conversation", converted.Id)
		log.WithField("id", converted.Id).Debug("Updating conv with id")
		err = collection.UpdateId(converted.Id, converted)
	}

	if err != nil {
		// @todo: handle and log
		return err
	}

	// Update the conversation pointer with the new values - only if the transaction succeeded
	c = toConversation(converted)

	return nil
}

// FindLatestConversation tries to find the latest conversation that happened with a user.
// In case this is a new user, then no conversation is returned. Otherwise the latest one,
// which can be the current one, is returned.
// Returns a conversation.ErrNotFound error when the user is not found.
func (repository *mongoDbRepository) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := repository.db.Session.Clone()
	defer session.Close()

	// Store the result of the query in our own mongo struct
	var converted *mongoConversation

	log.WithField("fbid", user.FbId).Debug("Finding latest conversation")

	err := session.DB(repository.db.Params.DbName).C("conversation").Find(bson.M{
		"messages": bson.M{
			"$elemMatch": bson.M{
				"message.sender.fbid": user.FbId,
			},
		},
	}).Sort("-created_at").One(&converted)

	if err != nil {
		if err == mgo.ErrNotFound {
			log.Debug("Latest conversation not found")
			return nil, conversation.ErrNotFound
		}

		// @todo: handle and log
		log.Debugf("Error: %s", err)
		return nil, err
	}

	log.WithField("conversation", converted).Debug("Found latest conversation")

	// We return the domain Conversation object
	return toConversation(converted), nil
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

// InsertUser creates a new user in the DB
func (repository *mongoDbRepository) InsertUser(user *conversation.User) error {
	// @todo: if the user has an ID, return an error

	session := repository.db.NewSession()
	defer session.Close()

	err := session.DB(repository.db.Params.DbName).C("user").Insert(user)

	if err != nil {
		// @todo: handle and log
		return err
	}

	return nil
}
