// Package solicitationsv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package solicitationsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	runt "runtime"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// CreateProductReviewAndSellerFeedbackSolicitationResponse The response schema for the createProductReviewAndSellerFeedbackSolicitation operation.
type CreateProductReviewAndSellerFeedbackSolicitationResponse struct {
	// Errors A list of error responses returned when a request is unsuccessful.
	Errors *ErrorList `json:"errors,omitempty"`
}

// Error Error response returned when the request is unsuccessful.
type Error struct {
	// Code An error code that identifies the type of error that occurred.
	Code string `json:"code"`

	// Details Additional details that can help the caller understand or fix the issue.
	Details *string `json:"details,omitempty"`

	// Message A message that describes the error condition in a human-readable form.
	Message string `json:"message"`
}

// ErrorList A list of error responses returned when a request is unsuccessful.
type ErrorList = []Error

// GetSchemaResponse defines model for GetSchemaResponse.
type GetSchemaResponse struct {
	Links *struct {
		// Self A Link object.
		Self LinkObject `json:"self"`
	} `json:"_links,omitempty"`

	// Errors A list of error responses returned when a request is unsuccessful.
	Errors *ErrorList `json:"errors,omitempty"`

	// Payload A JSON schema document describing the expected payload of the action. This object can be validated against <a href=http://json-schema.org/draft-04/schema>http://json-schema.org/draft-04/schema</a>.
	Payload *Schema `json:"payload,omitempty"`
}

// GetSolicitationActionResponse Describes a solicitation action that can be taken for an order. Provides a JSON Hypertext Application Language (HAL) link to the JSON schema document that describes the expected input.
type GetSolicitationActionResponse struct {
	Embedded *struct {
		Schema *GetSchemaResponse `json:"schema,omitempty"`
	} `json:"_embedded,omitempty"`
	Links *struct {
		// Schema A Link object.
		Schema LinkObject `json:"schema"`

		// Self A Link object.
		Self LinkObject `json:"self"`
	} `json:"_links,omitempty"`

	// Errors A list of error responses returned when a request is unsuccessful.
	Errors *ErrorList `json:"errors,omitempty"`

	// Payload A simple object containing the name of the template.
	Payload *SolicitationsAction `json:"payload,omitempty"`
}

// GetSolicitationActionsForOrderResponse The response schema for the getSolicitationActionsForOrder operation.
type GetSolicitationActionsForOrderResponse struct {
	Embedded *struct {
		Actions []GetSolicitationActionResponse `json:"actions"`
	} `json:"_embedded,omitempty"`
	Links *struct {
		// Actions Eligible actions for the specified amazonOrderId.
		Actions []LinkObject `json:"actions"`

		// Self A Link object.
		Self LinkObject `json:"self"`
	} `json:"_links,omitempty"`

	// Errors A list of error responses returned when a request is unsuccessful.
	Errors *ErrorList `json:"errors,omitempty"`
}

// LinkObject A Link object.
type LinkObject struct {
	// Href A URI for this object.
	Href string `json:"href"`

	// Name An identifier for this object.
	Name *string `json:"name,omitempty"`
}

// Schema A JSON schema document describing the expected payload of the action. This object can be validated against <a href=http://json-schema.org/draft-04/schema>http://json-schema.org/draft-04/schema</a>.
type Schema map[string]interface{}

// SolicitationsAction A simple object containing the name of the template.
type SolicitationsAction struct {
	Name string `json:"name"`
}

// GetSolicitationActionsForOrderParams defines parameters for GetSolicitationActionsForOrder.
type GetSolicitationActionsForOrderParams struct {
	// MarketplaceIds A marketplace identifier. This specifies the marketplace in which the order was placed. Only one marketplace can be specified.
	MarketplaceIds []string `form:"marketplaceIds" json:"marketplaceIds"`
}

// CreateProductReviewAndSellerFeedbackSolicitationParams defines parameters for CreateProductReviewAndSellerFeedbackSolicitation.
type CreateProductReviewAndSellerFeedbackSolicitationParams struct {
	// MarketplaceIds A marketplace identifier. This specifies the marketplace in which the order was placed. Only one marketplace can be specified.
	MarketplaceIds []string `form:"marketplaceIds" json:"marketplaceIds"`
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
	// GetSolicitationActionsForOrder request
	GetSolicitationActionsForOrder(ctx context.Context, amazonOrderId string, params *GetSolicitationActionsForOrderParams) (*http.Response, error)

	// CreateProductReviewAndSellerFeedbackSolicitation request
	CreateProductReviewAndSellerFeedbackSolicitation(ctx context.Context, amazonOrderId string, params *CreateProductReviewAndSellerFeedbackSolicitationParams) (*http.Response, error)
}

func (c *Client) GetSolicitationActionsForOrder(ctx context.Context, amazonOrderId string, params *GetSolicitationActionsForOrderParams) (*http.Response, error) {
	req, err := NewGetSolicitationActionsForOrderRequest(c.Server, amazonOrderId, params)
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

func (c *Client) CreateProductReviewAndSellerFeedbackSolicitation(ctx context.Context, amazonOrderId string, params *CreateProductReviewAndSellerFeedbackSolicitationParams) (*http.Response, error) {
	req, err := NewCreateProductReviewAndSellerFeedbackSolicitationRequest(c.Server, amazonOrderId, params)
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

// NewGetSolicitationActionsForOrderRequest generates requests for GetSolicitationActionsForOrder
func NewGetSolicitationActionsForOrderRequest(server string, amazonOrderId string, params *GetSolicitationActionsForOrderParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "amazonOrderId", runtime.ParamLocationPath, amazonOrderId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/solicitations/v1/orders/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "marketplaceIds", runtime.ParamLocationQuery, params.MarketplaceIds); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				values := make([]string, len(v))
				copy(values, v)
				queryValues.Add(k, strings.Join(values, ","))
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateProductReviewAndSellerFeedbackSolicitationRequest generates requests for CreateProductReviewAndSellerFeedbackSolicitation
func NewCreateProductReviewAndSellerFeedbackSolicitationRequest(server string, amazonOrderId string, params *CreateProductReviewAndSellerFeedbackSolicitationParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "amazonOrderId", runtime.ParamLocationPath, amazonOrderId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/solicitations/v1/orders/%s/solicitations/productReviewAndSellerFeedback", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "marketplaceIds", runtime.ParamLocationQuery, params.MarketplaceIds); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				values := make([]string, len(v))
				copy(values, v)
				queryValues.Add(k, strings.Join(values, ","))
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
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
	// GetSolicitationActionsForOrderWithResponse request
	GetSolicitationActionsForOrderWithResponse(ctx context.Context, amazonOrderId string, params *GetSolicitationActionsForOrderParams) (*GetSolicitationActionsForOrderResp, error)

	// CreateProductReviewAndSellerFeedbackSolicitationWithResponse request
	CreateProductReviewAndSellerFeedbackSolicitationWithResponse(ctx context.Context, amazonOrderId string, params *CreateProductReviewAndSellerFeedbackSolicitationParams) (*CreateProductReviewAndSellerFeedbackSolicitationResp, error)
}

type GetSolicitationActionsForOrderResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetSolicitationActionsForOrderResponse
	JSON400      *GetSolicitationActionsForOrderResponse
	JSON403      *GetSolicitationActionsForOrderResponse
	JSON404      *GetSolicitationActionsForOrderResponse
	JSON413      *GetSolicitationActionsForOrderResponse
	JSON415      *GetSolicitationActionsForOrderResponse
	JSON429      *GetSolicitationActionsForOrderResponse
	JSON500      *GetSolicitationActionsForOrderResponse
	JSON503      *GetSolicitationActionsForOrderResponse
}

// Status returns HTTPResponse.Status
func (r GetSolicitationActionsForOrderResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSolicitationActionsForOrderResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateProductReviewAndSellerFeedbackSolicitationResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON400      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON403      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON404      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON413      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON415      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON429      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON500      *CreateProductReviewAndSellerFeedbackSolicitationResponse
	JSON503      *CreateProductReviewAndSellerFeedbackSolicitationResponse
}

// Status returns HTTPResponse.Status
func (r CreateProductReviewAndSellerFeedbackSolicitationResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateProductReviewAndSellerFeedbackSolicitationResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetSolicitationActionsForOrderWithResponse request returning *GetSolicitationActionsForOrderResp
func (c *ClientWithResponses) GetSolicitationActionsForOrderWithResponse(ctx context.Context, amazonOrderId string, params *GetSolicitationActionsForOrderParams) (*GetSolicitationActionsForOrderResp, error) {
	rsp, err := c.GetSolicitationActionsForOrder(ctx, amazonOrderId, params)
	if err != nil {
		return nil, err
	}
	return ParseGetSolicitationActionsForOrderResp(rsp)
}

// CreateProductReviewAndSellerFeedbackSolicitationWithResponse request returning *CreateProductReviewAndSellerFeedbackSolicitationResp
func (c *ClientWithResponses) CreateProductReviewAndSellerFeedbackSolicitationWithResponse(ctx context.Context, amazonOrderId string, params *CreateProductReviewAndSellerFeedbackSolicitationParams) (*CreateProductReviewAndSellerFeedbackSolicitationResp, error) {
	rsp, err := c.CreateProductReviewAndSellerFeedbackSolicitation(ctx, amazonOrderId, params)
	if err != nil {
		return nil, err
	}
	return ParseCreateProductReviewAndSellerFeedbackSolicitationResp(rsp)
}

// ParseGetSolicitationActionsForOrderResp parses an HTTP response from a GetSolicitationActionsForOrderWithResponse call
func ParseGetSolicitationActionsForOrderResp(rsp *http.Response) (*GetSolicitationActionsForOrderResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetSolicitationActionsForOrderResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 413:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON413 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON415 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 503:
		var dest GetSolicitationActionsForOrderResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON503 = &dest

	}

	return response, nil
}

// ParseCreateProductReviewAndSellerFeedbackSolicitationResp parses an HTTP response from a CreateProductReviewAndSellerFeedbackSolicitationWithResponse call
func ParseCreateProductReviewAndSellerFeedbackSolicitationResp(rsp *http.Response) (*CreateProductReviewAndSellerFeedbackSolicitationResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateProductReviewAndSellerFeedbackSolicitationResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 413:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON413 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON415 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 503:
		var dest CreateProductReviewAndSellerFeedbackSolicitationResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON503 = &dest

	}

	return response, nil
}
