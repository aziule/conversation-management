package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

var (
	ErrCouldNotCreateRequest = errors.New("Could not create the request")
	ErrInvalidStatusCode     = errors.New("Invalid status code returned")
)

// requestSpecifications represents specifications that can be used to create an http.Request
type requestSpecifications struct {
	Method              string
	Url                 *url.URL
	AuthorisationHeader string
}

// NewRequestSpecifications creates a new requestSpecifications object
func NewRequestSpecifications() *requestSpecifications {
	return &requestSpecifications{}
}

// WithMethod defines the method to be used by the request
func (specs *requestSpecifications) WithMethod(method string) *requestSpecifications {
	specs.Method = method
	return specs
}

// WithUrl defines the URL the request will reach
func (specs *requestSpecifications) WithUrl(url *url.URL) *requestSpecifications {
	specs.Url = url
	return specs
}

// WithAuthorisationHeader defines what the Authorization header will contain
func (specs *requestSpecifications) WithAuthorisationHeader(authorisation string) *requestSpecifications {
	specs.AuthorisationHeader = authorisation
	return specs
}

// NewRequest creates a new request based on a set of specifications
func NewRequest(specs *requestSpecifications) (*http.Request, error) {
	request, err := http.NewRequest(specs.Method, specs.Url.String(), nil)

	if err != nil {
		return nil, ErrCouldNotCreateRequest
	}

	if specs.AuthorisationHeader != "" {
		request.Header.Set("Authorization", specs.AuthorisationHeader)
	}

	return request, nil
}

// ParseJsonFromRequest parses a request, crafted using specifications, and stores
// the result inside the provided envelope.
func ParseJsonFromRequest(specs *requestSpecifications, envelope interface{}) error {
	request, err := NewRequest(specs)

	if err != nil {
		log.WithFields(log.Fields{
			"url": request.URL.String(),
		}).Infof("Could not create a new request: %s", err)
		// @todo: return a proper error
		return err
	}

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		log.Infof("Failed to send the request: %s", err)
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Infof("Failed to read the response body: %s", err)
		return err
	}

	if response.StatusCode != 200 {
		log.WithField("code", response.StatusCode).Info("API returned a non-200 code")
		return ErrInvalidStatusCode
	}

	err = json.Unmarshal(body, envelope)

	if err != nil {
		log.Infof("Failed to unmarshal the response body: %s", err)
		return err
	}

	return nil
}
