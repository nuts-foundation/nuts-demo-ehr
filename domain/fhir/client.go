package fhir

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"github.com/tidwall/gjson"
	"strings"
)

type ClientOpt func(client *httpClient)

type Factory func(opts ...ClientOpt) Client

func WithURL(serverURL string) ClientOpt {
	return func(client *httpClient) {
		client.url = serverURL
	}
}

func WithTenant(tenant string) ClientOpt {
	return func(client *httpClient) {
		client.tenant = tenant
	}
}

func WithAuthToken(authToken string) ClientOpt {
	return func(client *httpClient) {
		client.restClient.SetAuthToken(authToken)
	}
}

func NewFactory(defaultOpts ...ClientOpt) Factory {
	return func(callerOpts ...ClientOpt) Client {
		client := &httpClient{
			restClient: resty.New().SetHeader("Content-Type", "application/json"),
		}
		for _, opt := range append(defaultOpts, callerOpts...) {
			opt(client)
		}
		return client
	}
}

type Client interface {
	String() string
	CreateOrUpdate(ctx context.Context, resource interface{}) error
	ReadMultiple(ctx context.Context, path string, params map[string]string, results interface{}) error
	ReadOne(ctx context.Context, path string, result interface{}) error
}

type httpClient struct {
	restClient *resty.Client
	url        string
	tenant     string
}

func (h httpClient) String() string {
	return h.url
}

func (h httpClient) CreateOrUpdate(ctx context.Context, resource interface{}) error {
	resourcePath, err := resolveResourcePath(resource)
	if err != nil {
		return fmt.Errorf("unable to determine resource path: %w", err)
	}
	requestURI := h.buildRequestURI(resourcePath)
	resp, err := h.restClient.R().SetBody(resource).SetContext(ctx).Put(requestURI)
	if err != nil {
		return fmt.Errorf("unable to write FHIR resource (path=%s): %w", requestURI, err)
	}
	if !resp.IsSuccess() {
		log.Warnf("FHIR server replied: %s", resp.String())
		return fmt.Errorf("unable to write FHIR resource (path=%s,http-status=%d): %w", requestURI, resp.StatusCode(), err)
	}
	return nil
}

func (h httpClient) ReadMultiple(ctx context.Context, path string, params map[string]string, results interface{}) error {
	raw, err := h.getResource(ctx, path, params)
	if err != nil {
		return err
	}
	resourcesJSON := raw.Get("entry.#.resource").String()
	if resourcesJSON == "" {
		resourcesJSON = "[]"
	}
	err = json.Unmarshal([]byte(resourcesJSON), results)
	if err != nil {
		log.Warnf("FHIR server replied: %s", raw.String())
		return fmt.Errorf("unable to unmarshal FHIR result (path=%s,target-type=%T): %w", path, results, err)
	}
	return nil
}

func (h httpClient) ReadOne(ctx context.Context, path string, result interface{}) error {
	raw, err := h.getResource(ctx, path, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(raw.String()), &result)
	if err != nil {
		log.Warnf("FHIR server replied: %s", raw.String())
		return fmt.Errorf("unable to unmarshal FHIR result (path=%s,target-type=%T): %w", path, result, err)
	}
	return nil
}

func (h httpClient) getResource(ctx context.Context, path string, params map[string]string) (gjson.Result, error) {
	resp, err := h.restClient.R().SetQueryParams(params).SetContext(ctx).SetHeader("Cache-Control", "no-cache").Get(h.buildRequestURI(path))
	if err != nil {
		return gjson.Result{}, err
	}
	if !resp.IsSuccess() {
		log.Warnf("FHIR server replied: %s", resp.String())
		return gjson.Result{}, fmt.Errorf("unable to read FHIR resource (path=%s,http-status=%d): %w", path, resp.StatusCode(), err)
	}
	return gjson.ParseBytes(resp.Body()), nil
}

func (h httpClient) buildRequestURI(path string) string {
	parts := []string{h.url, h.tenant, path}
	result := ""
	for i := 0; i < len(parts); i++ {
		// Make sure parts are separated by exactly 1 slash
		if result != "" && !strings.HasSuffix(result, "/") && !strings.HasPrefix(parts[i], "/") {
			result += "/"
		}
		result += parts[i]
	}
	return result
}

func resolveResourcePath(resource interface{}) (string, error) {
	data, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	js := gjson.ParseBytes(data)
	resourceType := js.Get("resourceType").String()
	if resourceType == "" {
		return "", fmt.Errorf("unable to determine resource type of %T", resource)
	}
	resourceID := js.Get("id").String()
	if resourceType == "" {
		return "", fmt.Errorf("unable to determine resource ID of %T", resource)
	}
	return resourceType + "/" + resourceID, nil
}
