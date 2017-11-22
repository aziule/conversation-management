// Package mongo provides all of the required methods to interact with
// a mongodb database, using mgo as the driver.
package mongo

import (
	"time"

	"github.com/aziule/conversation-management/core/bot"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

const BotMetadataCollection = "bot"

// botRepository is the unexported struct that implements the Repository interface
type botRepository struct {
	db *Db
}

// NewBotRepository creates a new bot repository using MongoDb as the data source
func NewBotRepository(db *Db) bot.Repository {
	return &botRepository{
		db: db,
	}
}

// FindAll returns all of the metadata available for each bot
func (repository *botRepository) FindAll() ([]*bot.Metadata, error) {
	session := repository.db.NewSession()
	defer session.Close()

	var metadatas []*bot.Metadata

	err := session.DB(repository.db.Params.DbName).C(BotMetadataCollection).Find(nil).All(&metadatas)

	if err != nil {
		log.Infof("Could not find the bots metadatas: %s", err)
		return nil, err
	}

	return metadatas, nil
}

// Save inserts / updates a bot's metadata
func (repository *botRepository) Save(metadata *bot.Metadata) error {
	session := repository.db.NewSession()
	defer session.Close()

	var err error
	collection := session.DB(repository.db.Params.DbName).C(BotMetadataCollection)

	metadata.UpdatedAt = time.Now()

	if metadata.Id == "" {
		metadata.Id = bson.NewObjectId()
		metadata.CreatedAt = time.Now()

		log.WithField("metadata", metadata).Debugf("Inserting metadata")
		err = collection.Insert(metadata)
	} else {
		log.WithField("metadata", metadata).Debugf("Updating metadata")
		err = collection.UpdateId(metadata.Id, metadata)
	}

	if err != nil {
		log.WithField("metadata", metadata).Infof("Could not save the metadata: %s", err)
		return err
	}

	return nil
}
