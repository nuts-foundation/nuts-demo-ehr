// Package vdr provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package vdr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// DIDCreateRequest defines model for DIDCreateRequest.
type DIDCreateRequest struct {
	// indicates if the generated key pair can be used for assertions.
	AssertionMethod *bool `json:"assertionMethod,omitempty"`

	// indicates if the generated key pair can be used for authentication.
	Authentication *bool `json:"authentication,omitempty"`

	// indicates if the generated key pair can be used for capability delegations.
	CapabilityDelegation *bool `json:"capabilityDelegation,omitempty"`

	// indicates if the generated key pair can be used for altering DID Documents.
	// In combination with selfControl = true, the key can be used to alter the new DID Document.
	// Defaults to true when not given.
	// default: true
	CapabilityInvocation *bool `json:"capabilityInvocation,omitempty"`

	// List of DIDs that can control the new DID Document. If selfControl = true and controllers is not empty,
	// the newly generated DID will be added to the list of controllers.
	Controllers *[]string `json:"controllers,omitempty"`

	// indicates if the generated key pair can be used for Key agreements.
	KeyAgreement *bool `json:"keyAgreement,omitempty"`

	// whether the generated DID Document can be altered with its own capabilityInvocation key.
	SelfControl *bool `json:"selfControl,omitempty"`
}

// A DID document according to the W3C spec following the Nuts Method rules as defined in [Nuts RFC006]
type DIDDocument struct {
	// List of KIDs that may sign JWTs, JWSs and VCs
	AssertionMethod *[]string `json:"assertionMethod,omitempty"`

	// List of KIDs that may alter DID documents that they control
	Authentication *[]string `json:"authentication,omitempty"`

	// List of URIs
	Context *[]string `json:"context,omitempty"`

	// Single DID (as string) or List of DIDs that have control over the DID document
	Controller *interface{} `json:"controller,omitempty"`

	// DID according to Nuts specification
	Id string `json:"id"`

	// List of supported services by the DID subject
	Service *[]Service `json:"service,omitempty"`

	// list of keys
	VerificationMethod *[]VerificationMethod `json:"verificationMethod,omitempty"`
}

// The DID document metadata.
type DIDDocumentMetadata struct {
	// Time when DID document was created in rfc3339 form.
	Created string `json:"created"`

	// Whether the DID document has been deactivated.
	Deactivated bool `json:"deactivated"`

	// Sha256 in hex form of the DID document contents.
	Hash string `json:"hash"`

	// Sha256 in hex form of the previous version of this DID document.
	PreviousHash *string `json:"previousHash,omitempty"`

	// txs lists the transaction(s) that created the current version of this DID Document.
	// If multiple transactions are listed, the DID Document is conflicted
	Txs []string `json:"txs"`

	// Time when DID document was updated in rfc3339 form.
	Updated *string `json:"updated,omitempty"`
}

// DIDResolutionResult defines model for DIDResolutionResult.
type DIDResolutionResult struct {
	// A DID document according to the W3C spec following the Nuts Method rules as defined in [Nuts RFC006]
	Document DIDDocument `json:"document"`

	// The DID document metadata.
	DocumentMetadata DIDDocumentMetadata `json:"documentMetadata"`
}

// DIDUpdateRequest defines model for DIDUpdateRequest.
type DIDUpdateRequest struct {
	// The hash of the document in hex format.
	CurrentHash string `json:"currentHash"`

	// A DID document according to the W3C spec following the Nuts Method rules as defined in [Nuts RFC006]
	Document DIDDocument `json:"document"`
}

// A service supported by a DID subject.
type Service struct {
	// ID of the service.
	Id string `json:"id"`

	// Either a URI or a complex object.
	ServiceEndpoint map[string]interface{} `json:"serviceEndpoint"`

	// The type of the endpoint.
	Type string `json:"type"`
}

// A public key in JWK form.
type VerificationMethod struct {
	// The DID subject this key belongs to.
	Controller string `json:"controller"`

	// The ID of the key, used as KID in various JWX technologies.
	Id string `json:"id"`

	// The public key formatted according rfc7517.
	PublicKeyJwk map[string]interface{} `json:"publicKeyJwk"`

	// The type of the key.
	Type string `json:"type"`
}

// CreateDIDJSONBody defines parameters for CreateDID.
type CreateDIDJSONBody DIDCreateRequest

// GetDIDParams defines parameters for GetDID.
type GetDIDParams struct {
	// If a versionId DID parameter is provided, the DID resolution algorithm returns a specific version of the DID document.
	// The version is the Sha256 hash of the document.
	// See [the did resolution spec about versioning](https://w3c-ccg.github.io/did-resolution/#versioning)
	VersionId *string `json:"versionId,omitempty"`
}

// UpdateDIDJSONBody defines parameters for UpdateDID.
type UpdateDIDJSONBody DIDUpdateRequest

// CreateDIDJSONRequestBody defines body for CreateDID for application/json ContentType.
type CreateDIDJSONRequestBody CreateDIDJSONBody

// UpdateDIDJSONRequestBody defines body for UpdateDID for application/json ContentType.
type UpdateDIDJSONRequestBody UpdateDIDJSONBody

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
	// CreateDID request with any body
	CreateDIDWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateDID(ctx context.Context, body CreateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ConflictedDIDs request
	ConflictedDIDs(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeactivateDID request
	DeactivateDID(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetDID request
	GetDID(ctx context.Context, did string, params *GetDIDParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateDID request with any body
	UpdateDIDWithBody(ctx context.Context, did string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateDID(ctx context.Context, did string, body UpdateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AddNewVerificationMethod request
	AddNewVerificationMethod(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteVerificationMethod request
	DeleteVerificationMethod(ctx context.Context, did string, kid string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) CreateDIDWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateDIDRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateDID(ctx context.Context, body CreateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateDIDRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ConflictedDIDs(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewConflictedDIDsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeactivateDID(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeactivateDIDRequest(c.Server, did)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetDID(ctx context.Context, did string, params *GetDIDParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDIDRequest(c.Server, did, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateDIDWithBody(ctx context.Context, did string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateDIDRequestWithBody(c.Server, did, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateDID(ctx context.Context, did string, body UpdateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateDIDRequest(c.Server, did, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddNewVerificationMethod(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddNewVerificationMethodRequest(c.Server, did)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteVerificationMethod(ctx context.Context, did string, kid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteVerificationMethodRequest(c.Server, did, kid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewCreateDIDRequest calls the generic CreateDID builder with application/json body
func NewCreateDIDRequest(server string, body CreateDIDJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateDIDRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateDIDRequestWithBody generates requests for CreateDID with any type of body
func NewCreateDIDRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did")
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

// NewConflictedDIDsRequest generates requests for ConflictedDIDs
func NewConflictedDIDsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did/conflicted")
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

// NewDeactivateDIDRequest generates requests for DeactivateDID
func NewDeactivateDIDRequest(server string, did string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "did", runtime.ParamLocationPath, did)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did/%s", pathParam0)
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

// NewGetDIDRequest generates requests for GetDID
func NewGetDIDRequest(server string, did string, params *GetDIDParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "did", runtime.ParamLocationPath, did)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if params.VersionId != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "versionId", runtime.ParamLocationQuery, *params.VersionId); err != nil {
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

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateDIDRequest calls the generic UpdateDID builder with application/json body
func NewUpdateDIDRequest(server string, did string, body UpdateDIDJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateDIDRequestWithBody(server, did, "application/json", bodyReader)
}

// NewUpdateDIDRequestWithBody generates requests for UpdateDID with any type of body
func NewUpdateDIDRequestWithBody(server string, did string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "did", runtime.ParamLocationPath, did)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewAddNewVerificationMethodRequest generates requests for AddNewVerificationMethod
func NewAddNewVerificationMethodRequest(server string, did string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "did", runtime.ParamLocationPath, did)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did/%s/verificationmethod", pathParam0)
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

// NewDeleteVerificationMethodRequest generates requests for DeleteVerificationMethod
func NewDeleteVerificationMethodRequest(server string, did string, kid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "did", runtime.ParamLocationPath, did)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "kid", runtime.ParamLocationPath, kid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/vdr/v1/did/%s/verificationmethod/%s", pathParam0, pathParam1)
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
	// CreateDID request with any body
	CreateDIDWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateDIDResponse, error)

	CreateDIDWithResponse(ctx context.Context, body CreateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateDIDResponse, error)

	// ConflictedDIDs request
	ConflictedDIDsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ConflictedDIDsResponse, error)

	// DeactivateDID request
	DeactivateDIDWithResponse(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*DeactivateDIDResponse, error)

	// GetDID request
	GetDIDWithResponse(ctx context.Context, did string, params *GetDIDParams, reqEditors ...RequestEditorFn) (*GetDIDResponse, error)

	// UpdateDID request with any body
	UpdateDIDWithBodyWithResponse(ctx context.Context, did string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateDIDResponse, error)

	UpdateDIDWithResponse(ctx context.Context, did string, body UpdateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateDIDResponse, error)

	// AddNewVerificationMethod request
	AddNewVerificationMethodWithResponse(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*AddNewVerificationMethodResponse, error)

	// DeleteVerificationMethod request
	DeleteVerificationMethodWithResponse(ctx context.Context, did string, kid string, reqEditors ...RequestEditorFn) (*DeleteVerificationMethodResponse, error)
}

type CreateDIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r CreateDIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateDIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ConflictedDIDsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]DIDResolutionResult
}

// Status returns HTTPResponse.Status
func (r ConflictedDIDsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ConflictedDIDsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeactivateDIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeactivateDIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeactivateDIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetDIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *DIDResolutionResult
}

// Status returns HTTPResponse.Status
func (r GetDIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateDIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r UpdateDIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateDIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AddNewVerificationMethodResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r AddNewVerificationMethodResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddNewVerificationMethodResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteVerificationMethodResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteVerificationMethodResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteVerificationMethodResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// CreateDIDWithBodyWithResponse request with arbitrary body returning *CreateDIDResponse
func (c *ClientWithResponses) CreateDIDWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateDIDResponse, error) {
	rsp, err := c.CreateDIDWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateDIDResponse(rsp)
}

func (c *ClientWithResponses) CreateDIDWithResponse(ctx context.Context, body CreateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateDIDResponse, error) {
	rsp, err := c.CreateDID(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateDIDResponse(rsp)
}

// ConflictedDIDsWithResponse request returning *ConflictedDIDsResponse
func (c *ClientWithResponses) ConflictedDIDsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ConflictedDIDsResponse, error) {
	rsp, err := c.ConflictedDIDs(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConflictedDIDsResponse(rsp)
}

// DeactivateDIDWithResponse request returning *DeactivateDIDResponse
func (c *ClientWithResponses) DeactivateDIDWithResponse(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*DeactivateDIDResponse, error) {
	rsp, err := c.DeactivateDID(ctx, did, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeactivateDIDResponse(rsp)
}

// GetDIDWithResponse request returning *GetDIDResponse
func (c *ClientWithResponses) GetDIDWithResponse(ctx context.Context, did string, params *GetDIDParams, reqEditors ...RequestEditorFn) (*GetDIDResponse, error) {
	rsp, err := c.GetDID(ctx, did, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDIDResponse(rsp)
}

// UpdateDIDWithBodyWithResponse request with arbitrary body returning *UpdateDIDResponse
func (c *ClientWithResponses) UpdateDIDWithBodyWithResponse(ctx context.Context, did string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateDIDResponse, error) {
	rsp, err := c.UpdateDIDWithBody(ctx, did, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateDIDResponse(rsp)
}

func (c *ClientWithResponses) UpdateDIDWithResponse(ctx context.Context, did string, body UpdateDIDJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateDIDResponse, error) {
	rsp, err := c.UpdateDID(ctx, did, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateDIDResponse(rsp)
}

// AddNewVerificationMethodWithResponse request returning *AddNewVerificationMethodResponse
func (c *ClientWithResponses) AddNewVerificationMethodWithResponse(ctx context.Context, did string, reqEditors ...RequestEditorFn) (*AddNewVerificationMethodResponse, error) {
	rsp, err := c.AddNewVerificationMethod(ctx, did, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddNewVerificationMethodResponse(rsp)
}

// DeleteVerificationMethodWithResponse request returning *DeleteVerificationMethodResponse
func (c *ClientWithResponses) DeleteVerificationMethodWithResponse(ctx context.Context, did string, kid string, reqEditors ...RequestEditorFn) (*DeleteVerificationMethodResponse, error) {
	rsp, err := c.DeleteVerificationMethod(ctx, did, kid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteVerificationMethodResponse(rsp)
}

// ParseCreateDIDResponse parses an HTTP response from a CreateDIDWithResponse call
func ParseCreateDIDResponse(rsp *http.Response) (*CreateDIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &CreateDIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseConflictedDIDsResponse parses an HTTP response from a ConflictedDIDsWithResponse call
func ParseConflictedDIDsResponse(rsp *http.Response) (*ConflictedDIDsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ConflictedDIDsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []DIDResolutionResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeactivateDIDResponse parses an HTTP response from a DeactivateDIDWithResponse call
func ParseDeactivateDIDResponse(rsp *http.Response) (*DeactivateDIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &DeactivateDIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetDIDResponse parses an HTTP response from a GetDIDWithResponse call
func ParseGetDIDResponse(rsp *http.Response) (*GetDIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetDIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest DIDResolutionResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateDIDResponse parses an HTTP response from a UpdateDIDWithResponse call
func ParseUpdateDIDResponse(rsp *http.Response) (*UpdateDIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &UpdateDIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseAddNewVerificationMethodResponse parses an HTTP response from a AddNewVerificationMethodWithResponse call
func ParseAddNewVerificationMethodResponse(rsp *http.Response) (*AddNewVerificationMethodResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &AddNewVerificationMethodResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseDeleteVerificationMethodResponse parses an HTTP response from a DeleteVerificationMethodWithResponse call
func ParseDeleteVerificationMethodResponse(rsp *http.Response) (*DeleteVerificationMethodResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &DeleteVerificationMethodResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}
