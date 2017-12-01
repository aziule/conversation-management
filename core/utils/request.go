package utils

import (
	"errors"
	"net/http"
	"net/url"
)

var ErrCouldNotCreateRequest = errors.New("Could not create the request")

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
