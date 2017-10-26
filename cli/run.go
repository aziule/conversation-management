package cli

import (
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/aziule/conversation-management/app"
	"github.com/aziule/conversation-management/app/bot"
	"github.com/aziule/conversation-management/app/bot/facebook"
	"github.com/aziule/conversation-management/app/conversation/mongo"
	"github.com/aziule/conversation-management/app/nlp/wit"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// RunCommand is the command responsible for running our bot using the given configuration.
// This is the main command.
type RunCommand struct {
	configFilePath string
}

// NewRunCommand returns a new RunCommand
func NewRunCommand() *RunCommand {
	return &RunCommand{}
}

// Usage returns the usage text for the command
func (c *RunCommand) Usage() string {
	return `run [-config=./config.json]:
	Runs the server and listens to incoming messages`
}

// Execute runs the command
func (c *RunCommand) Execute(f *flag.FlagSet) error {
	// Shared flags between the commands
	config, err := app.LoadConfig(c.configFilePath)

	if err != nil {
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	dbParams := mongo.DbParams{
		DbHost: config.DbHost,
		DbName: config.DbName,
		DbUser: config.DbUser,
		DbPass: config.DbPass,
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{dbParams.DbHost},
		Database: dbParams.DbName,
		Timeout:  2 * time.Second,
	})

	if err != nil {
		log.Fatalf("An error occurred when connecting to the db: %s", err)
	}

	defer session.Close()

	b := facebook.NewBot(
		&facebook.Config{
			VerifyToken:     config.FbVerifyToken,
			ApiVersion:      config.FbApiVersion,
			PageAccessToken: config.FbPageAccessToken,
			NlpParser:       wit.NewParser(facebook.DefaultDataTypeMap),
			ConversationRepository: mongo.NewMongodbRepository(&mongo.Db{
				Session: session,
				Params:  dbParams,
			}),
		},
	)

	r := chi.NewRouter()

	// Automatically set the bot's webhooks routes
	for _, webhook := range b.Webhooks() {
		log.Debugf("%s %s", string(webhook.Method), webhook.Path)

		switch webhook.Method {
		case bot.HttpMethodGet:
			r.Get(webhook.Path, webhook.Handler)
		case bot.HttpMethodPost:
			r.Post(webhook.Path, webhook.Handler)
		}
	}

	log.Debugf("Listening on port %d", config.ListeningPort)
	http.ListenAndServe(":"+strconv.Itoa(config.ListeningPort), r)

	return nil
}

// FlagSet returns the command's flag set
func (c *RunCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.configFilePath, "config", "config.json", "Config file path")
}

// Name returns the command's name, to be used when invoking it from the cli
func (c *RunCommand) Name() string {
	return "run"
}
