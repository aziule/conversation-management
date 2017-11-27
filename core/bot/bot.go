// Package bot defines the generic structs and interfaces for creating a bot
// on any platform.
package bot

import (
	"time"

	"github.com/aziule/conversation-management/core/utils"
	"gopkg.in/mgo.v2/bson"
)

// ParamName represents a bot's parameter name. It is used to identify
// the parameters easily, when put within the Definition.Parameters map.
type ParamName string

// Platform represents a platform where a bot is operating on
type Platform string

const (
	PlatformFacebook Platform = "facebook"
	builderPrefix             = "bot_"
)

// RegisterRepositoryBuilder registers a new service builder using a package-level prefix
func RegisterRepositoryBuilder(name string, builder utils.ServiceBuilder) {
	utils.RegisterServiceBuilder(builderPrefix+name, builder)
}

// NewRepository tries to create a Repository using the available builders.
// Returns ErrRepositoryNotFound if the repository builder isn't found.
// Returns an error in case of any error during the build process.
func NewRepository(name string, conf utils.BuilderConf) (Repository, error) {
	repositoryBuilder, err := utils.GetServiceBuilder(builderPrefix + name)

	if err != nil {
		return nil, err
	}

	repository, err := repositoryBuilder(conf)

	if err != nil {
		return nil, err
	}

	return repository.(Repository), nil
}

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	ApiEndpoints() []*ApiEndpoint
	Definition() *Definition
}

// @todo: add slug
// Definition is the struct describing the bot: what is its Id, what platform
// is it using, and some platform-specific parameters
type Definition struct {
	Id         bson.ObjectId             `json:"-" bson:"_id"`
	Slug       string                    `json:"slug" bson:"slug"`
	Platform   Platform                  `json:"platform" bson:"platform"`
	Parameters map[ParamName]interface{} `json:"parameters" bson:"parameters"`
	CreatedAt  time.Time                 `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at" bson:"updated_at"`
}

// Repository is the interface responsible for fetching / saving bots
type Repository interface {
	FindAll() ([]*Definition, error)
	Save(definition *Definition) error
}
