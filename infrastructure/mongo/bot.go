// Package mongo provides all of the required methods to interact with
// a mongodb database, using mgo as the driver.
package mongo

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"

	"github.com/aziule/conversation-management/core/bot"
	log "github.com/sirupsen/logrus"
)

const BotDefinitionCollection = "bot"

var (
	ErrUndefinedParam = func(param string) error { return errors.New("Undefined param: " + param) }
	ErrInvalidParam   = func(param string) error { return errors.New("Invalid param type: " + param) }
)

// botRepository is the unexported struct that implements the Repository interface
type botRepository struct {
	db *Db
}

// newBotRepository creates a new bot repository using MongoDb as the data source
func newBotRepository(conf map[string]interface{}) (bot.Repository, error) {
	dbParam, ok := conf["db"]

	if !ok {
		return nil, ErrUndefinedParam("db")
	}

	db, ok := dbParam.(*Db)

	if !ok {
		return nil, ErrInvalidParam("db")
	}

	return &botRepository{
		db: db,
	}, nil
}

// FindAll returns all of the definition available for each bot
func (repository *botRepository) FindAll() ([]*bot.Definition, error) {
	session := repository.db.NewSession()
	defer session.Close()

	var definitions []*bot.Definition

	err := session.DB(repository.db.Params.DbName).C(BotDefinitionCollection).Find(nil).All(&definitions)

	if err != nil {
		log.Infof("Could not find the bots definitions: %s", err)
		return nil, err
	}

	return definitions, nil
}

// Save inserts / updates a bot's definition
func (repository *botRepository) Save(definition *bot.Definition) error {
	session := repository.db.NewSession()
	defer session.Close()

	var err error
	collection := session.DB(repository.db.Params.DbName).C(BotDefinitionCollection)

	definition.UpdatedAt = time.Now()

	if definition.Id == "" {
		definition.Id = bson.NewObjectId()
		definition.CreatedAt = time.Now()

		log.WithField("definition", definition).Debugf("Inserting definition")
		err = collection.Insert(definition)
	} else {
		log.WithField("definition", definition).Debugf("Updating definition")
		err = collection.UpdateId(definition.Id, definition)
	}

	if err != nil {
		log.WithField("definition", definition).Infof("Could not save the definition: %s", err)
		return err
	}

	return nil
}

func init() {
	bot.RegisterRepositoryBuilder("mongo", newBotRepository)
}
