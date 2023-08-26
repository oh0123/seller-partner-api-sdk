// Package vendorDirectFulfillmentInventoryv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package vendorDirectFulfillmentInventoryv1

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

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// Error Error response returned when the request is unsuccessful.
type Error struct {
	// Code An error code that identifies the type of error that occurred.
	Code string `json:"code"`

	// Details Additional details that can help the caller understand or fix the issue.
	Details *string `json:"details,omitempty"`

	// Message A message that describes the error condition.
	Message string `json:"message"`
}

// ErrorList A list of error responses returned when a request is unsuccessful.
type ErrorList = []Error

// InventoryUpdate defines model for InventoryUpdate.
type InventoryUpdate struct {
	// IsFullUpdate When true, this request contains a full feed. Otherwise, this request contains a partial feed. When sending a full feed, you must send information about all items in the warehouse. Any items not in the full feed are updated as not available. When sending a partial feed, only include the items that need an update to inventory. The status of other items will remain unchanged.
	IsFullUpdate bool `json:"isFullUpdate"`

	// Items A list of inventory items with updated details, including quantity available.
	Items        []ItemDetails       `json:"items"`
	SellingParty PartyIdentification `json:"sellingParty"`
}

// ItemDetails Updated inventory details for an item.
type ItemDetails struct {
	// AvailableQuantity Details of item quantity.
	AvailableQuantity ItemQuantity `json:"availableQuantity"`

	// BuyerProductIdentifier The buyer selected product identification of the item. Either buyerProductIdentifier or vendorProductIdentifier should be submitted.
	BuyerProductIdentifier *string `json:"buyerProductIdentifier,omitempty"`

	// IsObsolete When true, the item is permanently unavailable.
	IsObsolete *bool `json:"isObsolete,omitempty"`

	// VendorProductIdentifier The vendor selected product identification of the item. Either buyerProductIdentifier or vendorProductIdentifier should be submitted.
	VendorProductIdentifier *string `json:"vendorProductIdentifier,omitempty"`
}

// ItemQuantity Details of item quantity.
type ItemQuantity struct {
	// Amount Quantity of units available for a specific item.
	Amount *int `json:"amount,omitempty"`

	// UnitOfMeasure Unit of measure for the available quantity.
	UnitOfMeasure string `json:"unitOfMeasure"`
}

// PartyIdentification defines model for PartyIdentification.
type PartyIdentification struct {
	// PartyId Assigned identification for the party.
	PartyId string `json:"partyId"`
}

// SubmitInventoryUpdateRequest The request body for the submitInventoryUpdate operation.
type SubmitInventoryUpdateRequest struct {
	Inventory *InventoryUpdate `json:"inventory,omitempty"`
}

// SubmitInventoryUpdateResponse The response schema for the submitInventoryUpdate operation.
type SubmitInventoryUpdateResponse struct {
	// Errors A list of error responses returned when a request is unsuccessful.
	Errors  *ErrorList            `json:"errors,omitempty"`
	Payload *TransactionReference `json:"payload,omitempty"`
}

// TransactionReference defines model for TransactionReference.
type TransactionReference struct {
	// TransactionId GUID to identify this transaction. This value can be used with the Transaction Status API to return the status of this transaction.
	TransactionId *string `json:"transactionId,omitempty"`
}

// SubmitInventoryUpdateJSONRequestBody defines body for SubmitInventoryUpdate for application/json ContentType.
type SubmitInventoryUpdateJSONRequestBody = SubmitInventoryUpdateRequest

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
	// SubmitInventoryUpdateWithBody request with any body
	SubmitInventoryUpdateWithBody(ctx context.Context, warehouseId string, contentType string, body io.Reader) (*http.Response, error)

	SubmitInventoryUpdate(ctx context.Context, warehouseId string, body SubmitInventoryUpdateJSONRequestBody) (*http.Response, error)
}

func (c *Client) SubmitInventoryUpdateWithBody(ctx context.Context, warehouseId string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewSubmitInventoryUpdateRequestWithBody(c.Server, warehouseId, contentType, body)
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

func (c *Client) SubmitInventoryUpdate(ctx context.Context, warehouseId string, body SubmitInventoryUpdateJSONRequestBody) (*http.Response, error) {
	req, err := NewSubmitInventoryUpdateRequest(c.Server, warehouseId, body)
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

// NewSubmitInventoryUpdateRequest calls the generic SubmitInventoryUpdate builder with application/json body
func NewSubmitInventoryUpdateRequest(server string, warehouseId string, body SubmitInventoryUpdateJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSubmitInventoryUpdateRequestWithBody(server, warehouseId, "application/json", bodyReader)
}

// NewSubmitInventoryUpdateRequestWithBody generates requests for SubmitInventoryUpdate with any type of body
func NewSubmitInventoryUpdateRequestWithBody(server string, warehouseId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "warehouseId", runtime.ParamLocationPath, warehouseId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/vendor/directFulfillment/inventory/v1/warehouses/%s/items", pathParam0)
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
	// SubmitInventoryUpdateWithBodyWithResponse request with any body
	SubmitInventoryUpdateWithBodyWithResponse(ctx context.Context, warehouseId string, contentType string, body io.Reader) (*SubmitInventoryUpdateResp, error)

	SubmitInventoryUpdateWithResponse(ctx context.Context, warehouseId string, body SubmitInventoryUpdateJSONRequestBody) (*SubmitInventoryUpdateResp, error)
}

type SubmitInventoryUpdateResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON202      *SubmitInventoryUpdateResponse
	JSON400      *SubmitInventoryUpdateResponse
	JSON403      *SubmitInventoryUpdateResponse
	JSON404      *SubmitInventoryUpdateResponse
	JSON413      *SubmitInventoryUpdateResponse
	JSON415      *SubmitInventoryUpdateResponse
	JSON429      *SubmitInventoryUpdateResponse
	JSON500      *SubmitInventoryUpdateResponse
	JSON503      *SubmitInventoryUpdateResponse
}

// Status returns HTTPResponse.Status
func (r SubmitInventoryUpdateResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SubmitInventoryUpdateResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SubmitInventoryUpdateWithBodyWithResponse request with arbitrary body returning *SubmitInventoryUpdateResp
func (c *ClientWithResponses) SubmitInventoryUpdateWithBodyWithResponse(ctx context.Context, warehouseId string, contentType string, body io.Reader) (*SubmitInventoryUpdateResp, error) {
	rsp, err := c.SubmitInventoryUpdateWithBody(ctx, warehouseId, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseSubmitInventoryUpdateResp(rsp)
}

func (c *ClientWithResponses) SubmitInventoryUpdateWithResponse(ctx context.Context, warehouseId string, body SubmitInventoryUpdateJSONRequestBody) (*SubmitInventoryUpdateResp, error) {
	rsp, err := c.SubmitInventoryUpdate(ctx, warehouseId, body)
	if err != nil {
		return nil, err
	}
	return ParseSubmitInventoryUpdateResp(rsp)
}

// ParseSubmitInventoryUpdateResp parses an HTTP response from a SubmitInventoryUpdateWithResponse call
func ParseSubmitInventoryUpdateResp(rsp *http.Response) (*SubmitInventoryUpdateResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SubmitInventoryUpdateResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 202:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON202 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 413:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON413 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON415 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 503:
		var dest SubmitInventoryUpdateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON503 = &dest

	}

	return response, nil
}
