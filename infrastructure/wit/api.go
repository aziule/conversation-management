package wit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/core/utils"
	log "github.com/sirupsen/logrus"
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
func (api *witApi) GetIntents() ([]*nlp.Intent, error) {
	specs := utils.NewRequestSpecifications()
	specs.WithMethod("GET")
	specs.WithUrl(api.getIntentsUrl())
	specs.WithAuthorisationHeader("Bearer " + api.bearerToken)
	request, err := utils.NewRequest(specs)

	if err != nil {
		log.WithFields(log.Fields{
			"url": request.URL.String(),
		}).Infof("Could not create a new request: %s", err)
		// @todo: return a proper error
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+api.bearerToken)

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		log.Infof("Failed to send the request: %s", err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Infof("Failed to read the response's body: %s", err)
		return nil, err
	}

	if response.StatusCode != 200 {
		log.WithField("code", response.StatusCode).Info("API returned a non-200 code")
		return nil, ErrApiErr
	}

	type intentsEnvelope struct {
		Values []struct {
			Value       string   `json:"value"`
			Expressions []string `json:"expression"`
		} `json:"values"`
	}

	witIntents := &intentsEnvelope{}
	err = json.Unmarshal(body, &witIntents)

	if err != nil {
		log.Infof("Failed to unmarshal the response's body: %s", err)
		return nil, err
	}

	intents := []*nlp.Intent{}

	for _, witIntent := range witIntents.Values {
		intents = append(intents, &nlp.Intent{Name: witIntent.Value})
	}

	return intents, nil
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
	u, _ := url.Parse(api.baseUrl.String() + "/entities/intent")

	return u
}

func init() {
	nlp.RegisterApiBuilder("wit", newWitApi)
}
