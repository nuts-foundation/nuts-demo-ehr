// Package pip provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package pip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Data defines model for Data.
type Data struct {
	// AuthInput Data used in OPA script
	AuthInput map[string]interface{} `json:"auth_input"`

	// ClientId client DID (for now)
	ClientId string `json:"client_id"`

	// Scope The scope. Corresponds to the auth scopes
	Scope string `json:"scope"`

	// VerifierId verifier DID (for now)
	VerifierId string `json:"verifier_id"`
}

// CreateDataJSONRequestBody defines body for CreateData for application/json ContentType.
type CreateDataJSONRequestBody = Data

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

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

// The interface specification for the client above.
type ClientInterface interface {
	// DeleteData request
	DeleteData(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetData request
	GetData(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateDataWithBody request with any body
	CreateDataWithBody(ctx context.Context, id string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateData(ctx context.Context, id string, body CreateDataJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) DeleteData(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteDataRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetData(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDataRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateDataWithBody(ctx context.Context, id string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateDataRequestWithBody(c.Server, id, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateData(ctx context.Context, id string, body CreateDataJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateDataRequest(c.Server, id, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewDeleteDataRequest generates requests for DeleteData
func NewDeleteDataRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0 = id

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/pip/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetDataRequest generates requests for GetData
func NewGetDataRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0 = id

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/pip/%s", pathParam0)
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

// NewCreateDataRequest calls the generic CreateData builder with application/json body
func NewCreateDataRequest(server string, id string, body CreateDataJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateDataRequestWithBody(server, id, "application/json", bodyReader)
}

// NewCreateDataRequestWithBody generates requests for CreateData with any type of body
func NewCreateDataRequestWithBody(server string, id string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0 = id

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/pip/%s", pathParam0)
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

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
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
	// DeleteDataWithResponse request
	DeleteDataWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*DeleteDataResponse, error)

	// GetDataWithResponse request
	GetDataWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetDataResponse, error)

	// CreateDataWithBodyWithResponse request with any body
	CreateDataWithBodyWithResponse(ctx context.Context, id string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateDataResponse, error)

	CreateDataWithResponse(ctx context.Context, id string, body CreateDataJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateDataResponse, error)
}

type DeleteDataResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteDataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteDataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetDataResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Data
}

// Status returns HTTPResponse.Status
func (r GetDataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateDataResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r CreateDataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateDataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// DeleteDataWithResponse request returning *DeleteDataResponse
func (c *ClientWithResponses) DeleteDataWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*DeleteDataResponse, error) {
	rsp, err := c.DeleteData(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteDataResponse(rsp)
}

// GetDataWithResponse request returning *GetDataResponse
func (c *ClientWithResponses) GetDataWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetDataResponse, error) {
	rsp, err := c.GetData(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDataResponse(rsp)
}

// CreateDataWithBodyWithResponse request with arbitrary body returning *CreateDataResponse
func (c *ClientWithResponses) CreateDataWithBodyWithResponse(ctx context.Context, id string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateDataResponse, error) {
	rsp, err := c.CreateDataWithBody(ctx, id, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateDataResponse(rsp)
}

func (c *ClientWithResponses) CreateDataWithResponse(ctx context.Context, id string, body CreateDataJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateDataResponse, error) {
	rsp, err := c.CreateData(ctx, id, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateDataResponse(rsp)
}

// ParseDeleteDataResponse parses an HTTP response from a DeleteDataWithResponse call
func ParseDeleteDataResponse(rsp *http.Response) (*DeleteDataResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteDataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetDataResponse parses an HTTP response from a GetDataWithResponse call
func ParseGetDataResponse(rsp *http.Response) (*GetDataResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Data
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateDataResponse parses an HTTP response from a CreateDataWithResponse call
func ParseCreateDataResponse(rsp *http.Response) (*CreateDataResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateDataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}
