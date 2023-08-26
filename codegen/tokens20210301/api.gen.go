// Package tokens20210301 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package tokens20210301

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	runt "runtime"
	"strings"
)

// Defines values for RestrictedResourceMethod.
const (
	DELETE RestrictedResourceMethod = "DELETE"
	GET    RestrictedResourceMethod = "GET"
	POST   RestrictedResourceMethod = "POST"
	PUT    RestrictedResourceMethod = "PUT"
)

// CreateRestrictedDataTokenRequest The request schema for the createRestrictedDataToken operation.
type CreateRestrictedDataTokenRequest struct {
	// RestrictedResources A list of restricted resources.
	// Maximum: 50
	RestrictedResources []RestrictedResource `json:"restrictedResources"`

	// TargetApplication The application ID for the target application to which access is being delegated.
	TargetApplication *string `json:"targetApplication,omitempty"`
}

// CreateRestrictedDataTokenResponse The response schema for the createRestrictedDataToken operation.
type CreateRestrictedDataTokenResponse struct {
	// ExpiresIn The lifetime of the Restricted Data Token, in seconds.
	ExpiresIn *int `json:"expiresIn,omitempty"`

	// RestrictedDataToken A Restricted Data Token (RDT). This is a short-lived access token that authorizes calls to restricted operations. Pass this value with the x-amz-access-token header when making subsequent calls to these restricted resources.
	RestrictedDataToken *string `json:"restrictedDataToken,omitempty"`
}

// Error An error response returned when the request is unsuccessful.
type Error struct {
	// Code An error code that identifies the type of error that occurred.
	Code string `json:"code"`

	// Details Additional details that can help the caller understand or fix the issue.
	Details *string `json:"details,omitempty"`

	// Message A message that describes the error condition.
	Message string `json:"message"`
}

// ErrorList A list of error responses returned when a request is unsuccessful.
type ErrorList struct {
	Errors *[]Error `json:"errors,omitempty"`
}

// RestrictedResource Model of a restricted resource.
type RestrictedResource struct {
	// DataElements Indicates the type of Personally Identifiable Information requested. This parameter is required only when getting an RDT for use with the getOrder, getOrders, or getOrderItems operation of the Orders API. For more information, see the [Tokens API Use Case Guide](doc:tokens-api-use-case-guide). Possible values include:
	// - **buyerInfo**. On the order level this includes general identifying information about the buyer and tax-related information. On the order item level this includes gift wrap information and custom order information, if available.
	// - **shippingAddress**. This includes information for fulfilling orders.
	// - **buyerTaxInformation**. This includes information for issuing tax invoices.
	DataElements *[]string `json:"dataElements,omitempty"`

	// Method The HTTP method in the restricted resource.
	Method RestrictedResourceMethod `json:"method"`

	// Path The path in the restricted resource. Here are some path examples:
	// - ```/orders/v0/orders```. For getting an RDT for the getOrders operation of the Orders API. For bulk orders.
	// - ```/orders/v0/orders/123-1234567-1234567```. For getting an RDT for the getOrder operation of the Orders API. For a specific order.
	// - ```/orders/v0/orders/123-1234567-1234567/orderItems```. For getting an RDT for the getOrderItems operation of the Orders API. For the order items in a specific order.
	// - ```/mfn/v0/shipments/FBA1234ABC5D```. For getting an RDT for the getShipment operation of the Shipping API. For a specific shipment.
	// - ```/mfn/v0/shipments/{shipmentId}```. For getting an RDT for the getShipment operation of the Shipping API. For any of a selling partner's shipments that you specify when you call the getShipment operation.
	Path string `json:"path"`
}

// RestrictedResourceMethod The HTTP method in the restricted resource.
type RestrictedResourceMethod string

// CreateRestrictedDataTokenJSONRequestBody defines body for CreateRestrictedDataToken for application/json ContentType.
type CreateRestrictedDataTokenJSONRequestBody = CreateRestrictedDataTokenRequest

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
	// CreateRestrictedDataTokenWithBody request with any body
	CreateRestrictedDataTokenWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	CreateRestrictedDataToken(ctx context.Context, body CreateRestrictedDataTokenJSONRequestBody) (*http.Response, error)
}

func (c *Client) CreateRestrictedDataTokenWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewCreateRestrictedDataTokenRequestWithBody(c.Server, contentType, body)
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

func (c *Client) CreateRestrictedDataToken(ctx context.Context, body CreateRestrictedDataTokenJSONRequestBody) (*http.Response, error) {
	req, err := NewCreateRestrictedDataTokenRequest(c.Server, body)
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

// NewCreateRestrictedDataTokenRequest calls the generic CreateRestrictedDataToken builder with application/json body
func NewCreateRestrictedDataTokenRequest(server string, body CreateRestrictedDataTokenJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateRestrictedDataTokenRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateRestrictedDataTokenRequestWithBody generates requests for CreateRestrictedDataToken with any type of body
func NewCreateRestrictedDataTokenRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/tokens/2021-03-01/restrictedDataToken")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

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
	// CreateRestrictedDataTokenWithBodyWithResponse request with any body
	CreateRestrictedDataTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*CreateRestrictedDataTokenResp, error)

	CreateRestrictedDataTokenWithResponse(ctx context.Context, body CreateRestrictedDataTokenJSONRequestBody) (*CreateRestrictedDataTokenResp, error)
}

type CreateRestrictedDataTokenResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CreateRestrictedDataTokenResponse
	JSON400      *ErrorList
	JSON401      *ErrorList
	JSON403      *ErrorList
	JSON404      *ErrorList
	JSON415      *ErrorList
	JSON429      *ErrorList
	JSON500      *ErrorList
	JSON503      *ErrorList
}

// Status returns HTTPResponse.Status
func (r CreateRestrictedDataTokenResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateRestrictedDataTokenResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// CreateRestrictedDataTokenWithBodyWithResponse request with arbitrary body returning *CreateRestrictedDataTokenResp
func (c *ClientWithResponses) CreateRestrictedDataTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*CreateRestrictedDataTokenResp, error) {
	rsp, err := c.CreateRestrictedDataTokenWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseCreateRestrictedDataTokenResp(rsp)
}

func (c *ClientWithResponses) CreateRestrictedDataTokenWithResponse(ctx context.Context, body CreateRestrictedDataTokenJSONRequestBody) (*CreateRestrictedDataTokenResp, error) {
	rsp, err := c.CreateRestrictedDataToken(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseCreateRestrictedDataTokenResp(rsp)
}

// ParseCreateRestrictedDataTokenResp parses an HTTP response from a CreateRestrictedDataTokenWithResponse call
func ParseCreateRestrictedDataTokenResp(rsp *http.Response) (*CreateRestrictedDataTokenResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateRestrictedDataTokenResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CreateRestrictedDataTokenResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON415 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 503:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON503 = &dest

	}

	return response, nil
}
