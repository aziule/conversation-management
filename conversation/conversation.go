package conversation

import (
	"github.com/aziule/conversation-management/conversation/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Status string

const (
	StatusOngoing           Status = "ongoing"
	StatusHumanIntervention Status = "human"
	StatusOver              Status = "over"
)

// Conversation is the struct that will handle our conversations between
// the bot and the various users.
type Conversation struct {
	Status       Status
	MessagesFlow *MessagesFlow
}

type Repository interface {
	FindLatestConversation(user *User) (*Conversation, error)
	SaveConversation(conversation *Conversation) error
	FindUser(userId string) (*User, error)
	InsertUser(user *User) error
}

type repository struct {
	db *mongo.Db
}

func NewMongodbRepository(db *mongo.Db) Repository {
	return &repository{
		db: db,
	}
}

func (repository *repository) FindLatestConversation(user *User) (*Conversation, error) {
	session := repository.db.Session.Clone()
	defer session.Close()

	return nil, nil
}

// FindUser tries to find a user based on its UserId
func (repository *repository) FindUser(userId string) (*User, error) {
	session := repository.db.NewSession()
	defer session.Close()

	user := &User{}

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

func (repository *repository) InsertUser(user *User) error {
	session := repository.db.NewSession()
	defer session.Close()

	err := session.DB(repository.db.Params.DbName).C("user").Insert(user)

	if err != nil {
		// @todo: handle and log
		return err
	}

	return nil
}

func (repository *repository) SaveConversation(conversation *Conversation) error {
	session := repository.db.Session.Clone()
	defer session.Close()

	return nil
}
