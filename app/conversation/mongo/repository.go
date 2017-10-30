package mongo

import (
	"fmt"
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

// FindLatestConversation tries to find the latest conversation that happened with a user.
// In case this is a new user, then no conversation is returned. Otherwise the latest one,
// which can be the current one, is returned.
// Returns a conversation.ErrNotFound error when the user is not found.
func (repository *mongoDbRepository) FindLatestConversation(user *conversation.User) (*conversation.Conversation, error) {
	session := repository.db.Session.Clone()
	defer session.Close()

	// Store the result of the query in our own struct
	var converted *mongoConversation

	err := session.DB(repository.db.Params.DbName).C("conversation").Find(bson.M{
		"messages": bson.M{
			"$elemMatch": bson.M{
				"sender.fbid": "1429733950458154",
			},
		},
	}).Sort("-created_at").One(&converted)

	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Println("Conversation not found")
			return nil, conversation.ErrNotFound
		}

		// @todo: handle and log
		fmt.Println("Another error: ", err)
		return nil, err
	}

	// Return a pointer to the core Conversation struct
	// @todo: check that we have no memory leak here, for example after
	// unmarshalling 10000 conversations
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

// SaveConversation saves a conversation to the database.
// The conversation can be an existing one or a new one
func (repository *mongoDbRepository) SaveConversation(c *conversation.Conversation) error {
	session := repository.db.Session.Clone()
	defer session.Close()

	var err error
	collection := session.DB(repository.db.Params.DbName).C("conversation")
	fmt.Println(c.Messages)
	c.UpdatedAt = time.Now()

	// @todo: convert to mongoConversation beforehand

	if c.Id == "" {
		c.Id = bson.NewObjectId().String()
		log.Infof("Inserting conversation: %s", c.Id)
		err = collection.Insert(c)
	} else {
		log.Infof("Updating conversation: %s", c.Id)
		err = collection.Update(bson.M{"_id,omitempty": c.Id}, c)
	}

	if err != nil {
		// @todo: handle and log
		return err
	}

	return nil
}
