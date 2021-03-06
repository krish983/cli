package ccv2

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Request contains all the elements of a Cloud Controller request
type Request struct {
	// URIParams are the list URI route parameters
	URIParams map[string]string

	// Query is a list of HTTP query parameters
	Query url.Values

	// RequestName is the name of the request (see routes)
	RequestName string

	URI    string
	Method string

	// Body is the request body
	Body io.Reader
}

// newHTTPRequest returns a constructed HTTP.Request with some defaults.
// Defaults are applied when Request fields are not filled in.
func (client CloudControllerClient) newHTTPRequest(passedRequest Request) (*http.Request, error) {
	var body io.Reader
	if passedRequest.Body != nil {
		body = passedRequest.Body
	}

	var request *http.Request
	var err error
	if passedRequest.URI != "" {
		request, err = http.NewRequest(
			passedRequest.Method,
			fmt.Sprintf("%s%s", client.API(), passedRequest.URI),
			body,
		)
	} else {
		request, err = client.router.CreateRequest(
			passedRequest.RequestName,
			passedRequest.URIParams,
			body,
		)
		if err == nil {
			request.URL.RawQuery = passedRequest.Query.Encode()
		}
	}
	if err != nil {
		return nil, err
	}

	request.Header = http.Header{}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}
