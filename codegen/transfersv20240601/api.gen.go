// Package transfersv20240601 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package transfersv20240601

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

// Defines values for AssignmentType.
const (
	DEFAULT AssignmentType = "DEFAULT"
)

// Defines values for PaymentInstrumentType.
const (
	BANKACCOUNT  PaymentInstrumentType = "BANK_ACCOUNT"
	CARD         PaymentInstrumentType = "CARD"
	SELLERWALLET PaymentInstrumentType = "SELLER_WALLET"
)

// AssignmentFilter The list of default payment instruments that are returned when you use the `assignmentFilter`.
type AssignmentFilter struct {
	// AssignmentTypes The list of assignment types.
	AssignmentTypes *ListOfAssignmentType `json:"assignmentTypes,omitempty"`
}

// AssignmentType The filter that returns payment methods based on the type of assignment.
type AssignmentType string

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
type ErrorList struct {
	// Errors array of errors
	Errors []Error `json:"errors"`
}

// ExpiryDate The expiration date of a card used for payment. If the payment instrument is not a card, the expiration date is null.
type ExpiryDate struct {
	// Month A string with digits
	Month *StringWithDigits `json:"month,omitempty"`

	// Year A string with digits
	Year *StringWithDigits `json:"year,omitempty"`
}

// GetPaymentMethodsRequest The request schema for the `getPaymentMethods` operation.
type GetPaymentMethodsRequest struct {
	// MarketplaceId A string with no white spaces.
	MarketplaceId *StringWithNoWhiteSpace `json:"marketplaceId,omitempty"`

	// OwningCustomerId A string with no white spaces.
	OwningCustomerId *StringWithNoWhiteSpace `json:"owningCustomerId,omitempty"`

	// PaymentMethodFilter The object used to filter payment methods based on different factors.
	PaymentMethodFilter *PaymentMethodFilter `json:"paymentMethodFilter,omitempty"`

	// RequestId A string with no white spaces.
	RequestId *StringWithNoWhiteSpace `json:"requestId,omitempty"`
}

// GetPaymentMethodsResponse The response schema for the `getPaymentMethodsResponse` operation.
type GetPaymentMethodsResponse struct {
	// PiDetails The list of payment instruments.
	PiDetails *ListOfPIDetails `json:"piDetails,omitempty"`
}

// InitiatePayoutRequest The request schema for the `initiatePayout` operation.
type InitiatePayoutRequest struct {
	// AccountType The account type in the selected marketplace for which a payout has to be initiated. For the supported EU marketplaces, the only account type is `Standard Orders`.
	AccountType string `json:"accountType"`

	// MarketplaceId A marketplace identifier. For the list of possible marketplace identifiers, refer to [Marketplace IDs](https://developer-docs.amazon.com/sp-api/docs/marketplace-ids).
	MarketplaceId string `json:"marketplaceId"`
}

// InitiatePayoutResponse The response schema for the `initiatePayout` operation.
type InitiatePayoutResponse struct {
	// PayoutReferenceId The financial event group ID for a successfully initiated payout. You can use this ID to track payout information.
	PayoutReferenceId string `json:"payoutReferenceId"`
}

// ListOfAssignmentType The list of assignment types.
type ListOfAssignmentType = []AssignmentType

// ListOfPIDetails The list of payment instruments.
type ListOfPIDetails = []PIDetails

// ListOfPaymentInstrumentType The list of payment instrument types that are present.
type ListOfPaymentInstrumentType = []PaymentInstrumentType

// ListOfStrings A list of strings.
type ListOfStrings = []StringWithNoWhiteSpace

// PIDetails The details of a payment instrument.
type PIDetails struct {
	// CountryCode A string with no white spaces.
	CountryCode *StringWithNoWhiteSpace `json:"countryCode,omitempty"`

	// DefaultMarketplaces A list of strings.
	DefaultMarketplaces *ListOfStrings `json:"defaultMarketplaces,omitempty"`

	// EncryptedAccountHolderName A string with no white spaces.
	EncryptedAccountHolderName *StringWithNoWhiteSpace `json:"encryptedAccountHolderName,omitempty"`

	// ExpiryDate The expiration date of a card used for payment. If the payment instrument is not a card, the expiration date is null.
	ExpiryDate *ExpiryDate `json:"expiryDate,omitempty"`

	// PaymentInstrumentId A string with no white spaces.
	PaymentInstrumentId *StringWithNoWhiteSpace `json:"paymentInstrumentId,omitempty"`

	// PaymentInstrumentType The payment instrument type.
	PaymentInstrumentType *PaymentInstrumentType `json:"paymentInstrumentType,omitempty"`

	// Tail A string with digits
	Tail *StringWithDigits `json:"tail,omitempty"`
}

// PaymentInstrumentType The payment instrument type.
type PaymentInstrumentType string

// PaymentMethodFilter The object used to filter payment methods based on different factors.
type PaymentMethodFilter struct {
	// AssignmentFilter The list of default payment instruments that are returned when you use the `assignmentFilter`.
	AssignmentFilter *AssignmentFilter `json:"assignmentFilter,omitempty"`

	// PaymentInstrumentId A string with no white spaces.
	PaymentInstrumentId *StringWithNoWhiteSpace `json:"paymentInstrumentId,omitempty"`

	// PaymentInstrumentTypes The list of payment instrument types that are present.
	PaymentInstrumentTypes *ListOfPaymentInstrumentType `json:"paymentInstrumentTypes,omitempty"`

	// SubscribedMarketplaces A list of strings.
	SubscribedMarketplaces *ListOfStrings `json:"subscribedMarketplaces,omitempty"`
}

// StringWithDigits A string with digits
type StringWithDigits = string

// StringWithNoWhiteSpace A string with no white spaces.
type StringWithNoWhiteSpace = string

// GetPaymentMethodsJSONRequestBody defines body for GetPaymentMethods for application/json ContentType.
type GetPaymentMethodsJSONRequestBody = GetPaymentMethodsRequest

// InitiatePayoutJSONRequestBody defines body for InitiatePayout for application/json ContentType.
type InitiatePayoutJSONRequestBody = InitiatePayoutRequest

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
	// GetPaymentMethodsWithBody request with any body
	GetPaymentMethodsWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	GetPaymentMethods(ctx context.Context, body GetPaymentMethodsJSONRequestBody) (*http.Response, error)

	// InitiatePayoutWithBody request with any body
	InitiatePayoutWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	InitiatePayout(ctx context.Context, body InitiatePayoutJSONRequestBody) (*http.Response, error)
}

func (c *Client) GetPaymentMethodsWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewGetPaymentMethodsRequestWithBody(c.Server, contentType, body)
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

func (c *Client) GetPaymentMethods(ctx context.Context, body GetPaymentMethodsJSONRequestBody) (*http.Response, error) {
	req, err := NewGetPaymentMethodsRequest(c.Server, body)
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

func (c *Client) InitiatePayoutWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewInitiatePayoutRequestWithBody(c.Server, contentType, body)
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

func (c *Client) InitiatePayout(ctx context.Context, body InitiatePayoutJSONRequestBody) (*http.Response, error) {
	req, err := NewInitiatePayoutRequest(c.Server, body)
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

// NewGetPaymentMethodsRequest calls the generic GetPaymentMethods builder with application/json body
func NewGetPaymentMethodsRequest(server string, body GetPaymentMethodsJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGetPaymentMethodsRequestWithBody(server, "application/json", bodyReader)
}

// NewGetPaymentMethodsRequestWithBody generates requests for GetPaymentMethods with any type of body
func NewGetPaymentMethodsRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/finances/transfers/2024-06-01/paymentmethods")
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

// NewInitiatePayoutRequest calls the generic InitiatePayout builder with application/json body
func NewInitiatePayoutRequest(server string, body InitiatePayoutJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewInitiatePayoutRequestWithBody(server, "application/json", bodyReader)
}

// NewInitiatePayoutRequestWithBody generates requests for InitiatePayout with any type of body
func NewInitiatePayoutRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/finances/transfers/2024-06-01/payouts")
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
	// GetPaymentMethodsWithBodyWithResponse request with any body
	GetPaymentMethodsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*GetPaymentMethodsResp, error)

	GetPaymentMethodsWithResponse(ctx context.Context, body GetPaymentMethodsJSONRequestBody) (*GetPaymentMethodsResp, error)

	// InitiatePayoutWithBodyWithResponse request with any body
	InitiatePayoutWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*InitiatePayoutResp, error)

	InitiatePayoutWithResponse(ctx context.Context, body InitiatePayoutJSONRequestBody) (*InitiatePayoutResp, error)
}

type GetPaymentMethodsResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetPaymentMethodsResponse
	JSON400      *ErrorList
	JSON403      *ErrorList
	JSON404      *ErrorList
	JSON413      *ErrorList
	JSON415      *ErrorList
	JSON429      *ErrorList
	JSON500      *ErrorList
	JSON503      *ErrorList
}

// Status returns HTTPResponse.Status
func (r GetPaymentMethodsResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetPaymentMethodsResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type InitiatePayoutResp struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *InitiatePayoutResponse
	JSON400      *ErrorList
	JSON403      *ErrorList
	JSON404      *ErrorList
	JSON413      *ErrorList
	JSON415      *ErrorList
	JSON429      *ErrorList
	JSON500      *ErrorList
	JSON503      *ErrorList
}

// Status returns HTTPResponse.Status
func (r InitiatePayoutResp) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r InitiatePayoutResp) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetPaymentMethodsWithBodyWithResponse request with arbitrary body returning *GetPaymentMethodsResp
func (c *ClientWithResponses) GetPaymentMethodsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*GetPaymentMethodsResp, error) {
	rsp, err := c.GetPaymentMethodsWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseGetPaymentMethodsResp(rsp)
}

func (c *ClientWithResponses) GetPaymentMethodsWithResponse(ctx context.Context, body GetPaymentMethodsJSONRequestBody) (*GetPaymentMethodsResp, error) {
	rsp, err := c.GetPaymentMethods(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseGetPaymentMethodsResp(rsp)
}

// InitiatePayoutWithBodyWithResponse request with arbitrary body returning *InitiatePayoutResp
func (c *ClientWithResponses) InitiatePayoutWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*InitiatePayoutResp, error) {
	rsp, err := c.InitiatePayoutWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseInitiatePayoutResp(rsp)
}

func (c *ClientWithResponses) InitiatePayoutWithResponse(ctx context.Context, body InitiatePayoutJSONRequestBody) (*InitiatePayoutResp, error) {
	rsp, err := c.InitiatePayout(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseInitiatePayoutResp(rsp)
}

// ParseGetPaymentMethodsResp parses an HTTP response from a GetPaymentMethodsWithResponse call
func ParseGetPaymentMethodsResp(rsp *http.Response) (*GetPaymentMethodsResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetPaymentMethodsResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetPaymentMethodsResponse
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

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 413:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON413 = &dest

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

// ParseInitiatePayoutResp parses an HTTP response from a InitiatePayoutWithResponse call
func ParseInitiatePayoutResp(rsp *http.Response) (*InitiatePayoutResp, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &InitiatePayoutResp{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest InitiatePayoutResponse
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

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 413:
		var dest ErrorList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON413 = &dest

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
