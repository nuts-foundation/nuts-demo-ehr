// Package discovery provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	externalRef0 "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common"
	"github.com/oapi-codegen/runtime"
)

const (
	JwtBearerAuthScopes = "jwtBearerAuth.Scopes"
)

// SearchResult defines model for SearchResult.
type SearchResult struct {
	// Fields Input descriptor IDs and their mapped values that from the Verifiable Credential.
	Fields map[string]interface{} `json:"fields"`

	// Id The ID of the Verifiable Presentation.
	Id string `json:"id"`

	// SubjectId The ID of the Verifiable Credential subject (holder), typically a DID.
	SubjectId string `json:"subject_id"`

	// Vp Verifiable Presentation
	Vp VerifiablePresentation `json:"vp"`
}

// ServiceDefinition defines model for ServiceDefinition.
type ServiceDefinition struct {
	// Endpoint The endpoint of the Discovery Service.
	Endpoint string `json:"endpoint"`

	// Id The ID of the Discovery Service.
	Id string `json:"id"`

	// PresentationDefinition The Presentation Definition of the Discovery Service.
	PresentationDefinition map[string]interface{} `json:"presentation_definition"`

	// PresentationMaxValidity The maximum validity (in seconds) of a Verifiable Presentation of the Discovery Service.
	PresentationMaxValidity int `json:"presentation_max_validity"`
}

// VerifiablePresentation Verifiable Presentation
type VerifiablePresentation = externalRef0.VerifiablePresentation

// SearchPresentationsParams defines parameters for SearchPresentations.
type SearchPresentationsParams struct {
	Query *map[string]string `form:"query,omitempty" json:"query,omitempty"`
}

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
	// GetServices request
	GetServices(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SearchPresentations request
	SearchPresentations(ctx context.Context, serviceID string, params *SearchPresentationsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeactivateServiceForDID request
	DeactivateServiceForDID(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetServiceActivation request
	GetServiceActivation(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ActivateServiceForDID request
	ActivateServiceForDID(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetServices(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetServicesRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SearchPresentations(ctx context.Context, serviceID string, params *SearchPresentationsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchPresentationsRequest(c.Server, serviceID, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeactivateServiceForDID(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeactivateServiceForDIDRequest(c.Server, serviceID, did)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetServiceActivation(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetServiceActivationRequest(c.Server, serviceID, did)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ActivateServiceForDID(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewActivateServiceForDIDRequest(c.Server, serviceID, did)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetServicesRequest generates requests for GetServices
func NewGetServicesRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/discovery/v1")
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

// NewSearchPresentationsRequest generates requests for SearchPresentations
func NewSearchPresentationsRequest(server string, serviceID string, params *SearchPresentationsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "serviceID", runtime.ParamLocationPath, serviceID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/discovery/v1/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Query != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "query", runtime.ParamLocationQuery, *params.Query); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
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

// NewDeactivateServiceForDIDRequest generates requests for DeactivateServiceForDID
func NewDeactivateServiceForDIDRequest(server string, serviceID string, did string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "serviceID", runtime.ParamLocationPath, serviceID)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1 = did

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/discovery/v1/%s/%s", pathParam0, pathParam1)
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

// NewGetServiceActivationRequest generates requests for GetServiceActivation
func NewGetServiceActivationRequest(server string, serviceID string, did string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "serviceID", runtime.ParamLocationPath, serviceID)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1 = did

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/discovery/v1/%s/%s", pathParam0, pathParam1)
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

// NewActivateServiceForDIDRequest generates requests for ActivateServiceForDID
func NewActivateServiceForDIDRequest(server string, serviceID string, did string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "serviceID", runtime.ParamLocationPath, serviceID)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1 = did

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/discovery/v1/%s/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

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
	// GetServicesWithResponse request
	GetServicesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetServicesResponse, error)

	// SearchPresentationsWithResponse request
	SearchPresentationsWithResponse(ctx context.Context, serviceID string, params *SearchPresentationsParams, reqEditors ...RequestEditorFn) (*SearchPresentationsResponse, error)

	// DeactivateServiceForDIDWithResponse request
	DeactivateServiceForDIDWithResponse(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*DeactivateServiceForDIDResponse, error)

	// GetServiceActivationWithResponse request
	GetServiceActivationWithResponse(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*GetServiceActivationResponse, error)

	// ActivateServiceForDIDWithResponse request
	ActivateServiceForDIDWithResponse(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*ActivateServiceForDIDResponse, error)
}

type GetServicesResponse struct {
	Body                          []byte
	HTTPResponse                  *http.Response
	JSON200                       *[]ServiceDefinition
	ApplicationproblemJSONDefault *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r GetServicesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetServicesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SearchPresentationsResponse struct {
	Body                          []byte
	HTTPResponse                  *http.Response
	JSON200                       *[]SearchResult
	ApplicationproblemJSONDefault *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r SearchPresentationsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SearchPresentationsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeactivateServiceForDIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON202      *struct {
		// Reason Description of why removal of the registration failed.
		Reason string `json:"reason"`
	}
	ApplicationproblemJSON400 *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
	ApplicationproblemJSONDefault *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r DeactivateServiceForDIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeactivateServiceForDIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetServiceActivationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Activated Whether the DID is activated on the Discovery Service.
		Activated bool `json:"activated"`

		// Vp Verifiable Presentation
		Vp *VerifiablePresentation `json:"vp,omitempty"`
	}
	ApplicationproblemJSONDefault *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r GetServiceActivationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetServiceActivationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ActivateServiceForDIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON202      *struct {
		// Reason Description of why registration failed.
		Reason string `json:"reason"`
	}
	ApplicationproblemJSON400 *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
	ApplicationproblemJSONDefault *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r ActivateServiceForDIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ActivateServiceForDIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetServicesWithResponse request returning *GetServicesResponse
func (c *ClientWithResponses) GetServicesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetServicesResponse, error) {
	rsp, err := c.GetServices(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetServicesResponse(rsp)
}

// SearchPresentationsWithResponse request returning *SearchPresentationsResponse
func (c *ClientWithResponses) SearchPresentationsWithResponse(ctx context.Context, serviceID string, params *SearchPresentationsParams, reqEditors ...RequestEditorFn) (*SearchPresentationsResponse, error) {
	rsp, err := c.SearchPresentations(ctx, serviceID, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchPresentationsResponse(rsp)
}

// DeactivateServiceForDIDWithResponse request returning *DeactivateServiceForDIDResponse
func (c *ClientWithResponses) DeactivateServiceForDIDWithResponse(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*DeactivateServiceForDIDResponse, error) {
	rsp, err := c.DeactivateServiceForDID(ctx, serviceID, did, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeactivateServiceForDIDResponse(rsp)
}

// GetServiceActivationWithResponse request returning *GetServiceActivationResponse
func (c *ClientWithResponses) GetServiceActivationWithResponse(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*GetServiceActivationResponse, error) {
	rsp, err := c.GetServiceActivation(ctx, serviceID, did, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetServiceActivationResponse(rsp)
}

// ActivateServiceForDIDWithResponse request returning *ActivateServiceForDIDResponse
func (c *ClientWithResponses) ActivateServiceForDIDWithResponse(ctx context.Context, serviceID string, did string, reqEditors ...RequestEditorFn) (*ActivateServiceForDIDResponse, error) {
	rsp, err := c.ActivateServiceForDID(ctx, serviceID, did, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseActivateServiceForDIDResponse(rsp)
}

// ParseGetServicesResponse parses an HTTP response from a GetServicesWithResponse call
func ParseGetServicesResponse(rsp *http.Response) (*GetServicesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetServicesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ServiceDefinition
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

// ParseSearchPresentationsResponse parses an HTTP response from a SearchPresentationsWithResponse call
func ParseSearchPresentationsResponse(rsp *http.Response) (*SearchPresentationsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SearchPresentationsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []SearchResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

// ParseDeactivateServiceForDIDResponse parses an HTTP response from a DeactivateServiceForDIDWithResponse call
func ParseDeactivateServiceForDIDResponse(rsp *http.Response) (*DeactivateServiceForDIDResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeactivateServiceForDIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 202:
		var dest struct {
			// Reason Description of why removal of the registration failed.
			Reason string `json:"reason"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON202 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

// ParseGetServiceActivationResponse parses an HTTP response from a GetServiceActivationWithResponse call
func ParseGetServiceActivationResponse(rsp *http.Response) (*GetServiceActivationResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetServiceActivationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Activated Whether the DID is activated on the Discovery Service.
			Activated bool `json:"activated"`

			// Vp Verifiable Presentation
			Vp *VerifiablePresentation `json:"vp,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

// ParseActivateServiceForDIDResponse parses an HTTP response from a ActivateServiceForDIDWithResponse call
func ParseActivateServiceForDIDResponse(rsp *http.Response) (*ActivateServiceForDIDResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ActivateServiceForDIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 202:
		var dest struct {
			// Reason Description of why registration failed.
			Reason string `json:"reason"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON202 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}
