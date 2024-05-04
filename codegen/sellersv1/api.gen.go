// Package sellersv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package sellersv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	runt "runtime"
	"strings"
)

// Error Error response returned when the request is unsuccessful.
type Error struct {
	// Code An error code that identifies the type of error that occured.
	Code string `json:"code"`

	// Details Additional details that can help the caller understand or fix the issue.
	Details *string `json:"details,omitempty"`

	// Message A message that describes the error condition in a human-readable form.
	Message string `json:"message"`
}

// ErrorList A list of error responses returned when a request is unsuccessful.
type ErrorList = []Error

// GetMarketplaceParticipationsResponse The response schema for the getMarketplaceParticipations operation.
type GetMarketplaceParticipationsResponse struct {
	// Errors A list of error responses returned when a request is unsuccessful.
	Errors *ErrorList `json:"errors,omitempty"`

	// Payload List of marketplace participations.
	Payload *MarketplaceParticipationList `json:"payload,omitempty"`
}

// Marketplace Detailed information about an Amazon market where a seller can list items for sale and customers can view and purchase items.
type Marketplace struct {
	// CountryCode The ISO 3166-1 alpha-2 format country code of the marketplace.
	CountryCode string `json:"countryCode"`

	// DefaultCurrencyCode The ISO 4217 format currency code of the marketplace.
	DefaultCurrencyCode string `json:"defaultCurrencyCode"`

	// DefaultLanguageCode The ISO 639-1 format language code of the marketplace.
	DefaultLanguageCode string `json:"defaultLanguageCode"`

	// DomainName The domain name of the marketplace.
	DomainName string `json:"domainName"`

	// Id The encrypted marketplace value.
	Id string `json:"id"`

	// Name Marketplace name.
	Name string `json:"name"`
}

// MarketplaceParticipation defines model for MarketplaceParticipation.
type MarketplaceParticipation struct {
	// Marketplace Detailed information about an Amazon market where a seller can list items for sale and customers can view and purchase items.
	Marketplace Marketplace `json:"marketplace"`

	// Participation Detailed information that is specific to a seller in a Marketplace.
	Participation Participation `json:"participation"`
}

// MarketplaceParticipationList List of marketplace participations.
type MarketplaceParticipationList = []MarketplaceParticipation

// Participation Detailed information that is specific to a seller in a Marketplace.
type Participation struct {
	// HasSuspendedListings Specifies if the seller has suspended listings. True if the seller Listing Status is set to Inactive, otherwise False.
	HasSuspendedListings bool `json:"hasSuspendedListings"`
	IsParticipating      bool `json:"isParticipating"`
}

// RequestEditorFn is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// ResponseEditorFn is the function signature for the ResponseEditor callback function
type ResponseEditorFn func(ctx context.Context, rsp *http.Response) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn

	// A callback for modifying response which are generated after receive from the network.
	ResponseEditors []ResponseEditorFn

	// The user agent header identifies your application, its version number, and the platform and programming language you are using.
	// You must include a user agent header in each request submitted to the sales partner API.
	UserAgent string
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	// setting the default useragent
	if client.UserAgent == "" {
		client.UserAgent = fmt.Sprintf("selling-partner-api-sdk/v2.0 (Language=%s; Platform=%s-%s)", strings.Replace(runt.Version(), "go", "go/", -1), runt.GOOS, runt.GOARCH)
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// WithResponseEditorFn allows setting up a callback function, which will be
// called right after receive the response.
func WithResponseEditorFn(fn ResponseEditorFn) ClientOption {
	return func(c *Client) error {
		c.ResponseEditors = append(c.ResponseEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetMarketplaceParticipations request
	GetMarketplaceParticipations(ctx context.Context) (*http.Response, error)
}

func (c *Client) GetMarketplaceParticipations(ctx context.Context) (*http.Response, error) {
	req, err := NewGetMarketplaceParticipationsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if err := c.applyReqEditors(ctx, req); err != nil {
		return nil, err
	}
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := c.applyRspEditor(ctx, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

// NewGetMarketplaceParticipationsRequest generates requests for GetMarketplaceParticipations
func NewGetMarketplaceParticipationsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/sellers/v1/marketplaceParticipations")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyReqEditors(ctx context.Context, req *http.Request) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) applyRspEditor(ctx context.Context, rsp *http.Response) error {
	for _, r := range c.ResponseEditors {
		if err := r(ctx, rsp); err != nil {
			return err
		}
	}
	return nil
} // ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetMarketplaceParticipationsWithResponse request
	GetMarketplaceParticipationsWithResponse(ctx context.Context) (*GetMarketplaceParticipationsResp, error)
}

type GetMarketplaceParticipationsResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetMarketplaceParticipationsResponse
	JSON400      *GetMarketplaceParticipationsResponse
	JSON403      *GetMarketplaceParticipationsResponse
	JSON404      *GetMarketplaceParticipationsResponse
	JSON413      *GetMarketplaceParticipationsResponse
	JSON415      *GetMarketplaceParticipationsResponse
	JSON429      *GetMarketplaceParticipationsResponse
	JSON500      *GetMarketplaceParticipationsResponse
	JSON503      *GetMarketplaceParticipationsResponse
}

// Status returns HTTPResponse.Status
func (r GetMarketplaceParticipationsResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetMarketplaceParticipationsResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetMarketplaceParticipationsWithResponse request returning *GetMarketplaceParticipationsResp
func (c *ClientWithResponses) GetMarketplaceParticipationsWithResponse(ctx context.Context) (*GetMarketplaceParticipationsResp, error) {
	rsp, err := c.GetMarketplaceParticipations(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetMarketplaceParticipationsResp(rsp)
}

// ParseGetMarketplaceParticipationsResp parses an HTTP response from a GetMarketplaceParticipationsWithResponse call
func ParseGetMarketplaceParticipationsResp(rsp *http.Response) (*GetMarketplaceParticipationsResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMarketplaceParticipationsResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 413:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON413 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON415 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 503:
		var dest GetMarketplaceParticipationsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON503 = &dest

	}

	return response, nil
}
