package conversation

import (
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

type Reader struct {
	db *Db
}

func NewReader(db *Db) *Reader {
	return &Reader{
		db: db,
	}
}

func (reader *Reader) FindLatestConversation(user *User) (*Conversation, error) {
	session := reader.db.Session.Clone()
	defer session.Close()

	return nil, nil
}

// FindUser tries to find a user based on its UserId
func (reader *Reader) FindUser(userId string) (*User, error) {
	session := reader.db.NewSession()
	defer session.Close()

	user := &User{}

	// Tied to Facebook at the moment...
	err := session.DB(reader.db.Params.DbName).C("user").Find(bson.M{
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

type Writer struct {
	db *Db
}

func NewWriter(db *Db) *Writer {
	return &Writer{
		db: db,
	}
}

func (writer *Writer) InsertUser(user *User) error {
	session := writer.db.NewSession()
	defer session.Close()

	err := session.DB(writer.db.Params.DbName).C("user").Insert(user)

	if err != nil {
		// @todo: handle and log
		return err
	}

	return nil
}

func (writer *Writer) Save(conversation *Conversation) error {
	session := writer.db.Session.Clone()
	defer session.Close()

	return nil
}
