package wit

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/core/utils"
)

var ErrApiErr = errors.New("API error")

// witApi is the struct used to make calls to Wit
type witApi struct {
	client      *http.Client
	baseUrl     *url.URL
	bearerToken string
}

// newWitApi creates a new witApi using the given conf
func newWitApi(conf utils.BuilderConf) (interface{}, error) {
	client, ok := utils.GetParam(conf, "client").(*http.Client)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("client")
	}

	token, ok := utils.GetParam(conf, "bearer_token").(string)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("bearer_token")
	}

	rawBaseUrl := "https://api.wit.ai"
	baseUrl, _ := url.Parse(rawBaseUrl)

	return &witApi{
		client:      client,
		baseUrl:     baseUrl,
		bearerToken: token,
	}, nil
}

// GetIntents gets the list of intents from Wit
// @todo: make the API return api-specific objects instead of domain objects,
// and put the logic inside the repository.
func (api *witApi) GetIntents() ([]*nlp.Intent, error) {
	type envelope struct {
		Values []struct {
			Value       string   `json:"value"`
			Expressions []string `json:"expression"`
		} `json:"values"`
	}

	witIntents := &envelope{}

	err := api.callApi("GET", api.getIntentsUrl(), &witIntents)

	if err != nil {
		// @todo: log this and return a proper error
		return nil, err
	}

	intents := []*nlp.Intent{}

	for _, witIntent := range witIntents.Values {
		intents = append(intents, &nlp.Intent{Name: witIntent.Value})
	}

	return intents, nil
}

// GetEntities gets the list of entities from Wit
// @todo: make the API return api-specific objects instead of domain objects,
// and put the logic inside the repository.
func (api *witApi) GetEntities() ([]*nlp.Entity, error) {
	var envelope []string
	err := api.callApi("GET", api.getEntitiesUrl(), &envelope)

	if err != nil {
		// @todo: log this and return a proper error
		return nil, err
	}

	entities := []*nlp.Entity{}

	for _, entityName := range envelope {
		// @todo: find a better way to exclude the "intent" entity
		if entityName == "intent" {
			continue
		}

		entities = append(entities, &nlp.Entity{Name: entityName})
	}

	return entities, nil
}

// callApi calls the API given a method, an URL and an envelope. If it is a success, then
// the data is parsed and stored inside the envelope (using JSON).
// Returns an error if anything happens or if the status code != 200.
func (api *witApi) callApi(method string, u *url.URL, envelope interface{}) error {
	specs := utils.NewRequestSpecifications()
	specs.WithMethod(method)
	specs.WithUrl(u)
	specs.WithAuthorisationHeader("Bearer " + api.bearerToken)

	return utils.ParseJsonFromRequest(specs, envelope)
}

// @todo: store it and avoid recreating it every time
// getSendTextUrl returns the url to ping to send text messages to a user
func (api *witApi) getEntitiesUrl() *url.URL {
	u, _ := url.Parse(api.baseUrl.String() + "/entities")

	return u
}

// @todo: store it and avoid recreating it every time
// getSendTextUrl returns the url to ping to send text messages to a user
func (api *witApi) getIntentsUrl() *url.URL {
	// @todo: find a better way to access the "intent" entity => use
	// the default data type map?
	u, _ := url.Parse(api.baseUrl.String() + "/entities/intent")

	return u
}

func init() {
	nlp.RegisterApiBuilder("wit", newWitApi)
}
