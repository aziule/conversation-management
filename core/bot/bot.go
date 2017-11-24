// Package bot defines the generic structs and interfaces for creating a bot
// on any platform.
package bot

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	ApiEndpoints() []*ApiEndpoint
	Definition() *Definition
}

// ParamName represents a bot's parameter name. It is used to identify
// the parameters easily, when put within the Definition.Parameters map.
type ParamName string

// Platform represents a platform where a bot is operating on
type Platform string

const PlatformFacebook Platform = "facebook"

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
