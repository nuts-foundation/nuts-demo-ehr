package fhir

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"github.com/tidwall/gjson"
	"strings"
)

func NewClient(baseURL string) Client {
	return &httpClient{
		restClient: resty.New().SetHeader("Content-Type", "application/json"),
		url:        baseURL,
	}
}

type Client interface {
	WriteResource(ctx context.Context, resourcePath string, resource interface{}) (gjson.Result, error)
	GetResources(ctx context.Context, path string, params map[string]string) ([]gjson.Result, error)
	GetResource(ctx context.Context, path string) (gjson.Result, error)
}

type httpClient struct {
	restClient *resty.Client

	url string
}

func (h httpClient) WriteResource(ctx context.Context, resourcePath string, resource interface{}) (gjson.Result, error) {
	resp, err := h.restClient.R().SetBody(resource).SetContext(ctx).Put(h.buildRequestURI(resourcePath))
	if err != nil {
		return gjson.Result{}, fmt.Errorf("unable to write FHIR resource (path=%s): %w", resourcePath, err)
	}
	if !resp.IsSuccess() {
		log.Warnf("FHIR server replied: %s", resp.String())
		return gjson.Result{}, fmt.Errorf("unable to write FHIR resource (path=%s,http-status=%d): %w", resourcePath, resp.StatusCode, err)
	}
	return gjson.ParseBytes(resp.Body()), nil
}

func (h httpClient) GetResources(ctx context.Context, path string, params map[string]string) ([]gjson.Result, error) {
	resource, err := h.getResource(ctx, path, params)
	if err == nil {
		return resource.Get("entry.#.resource").Array(), nil
	}
	return nil, err
}

func (h httpClient) GetResource(ctx context.Context, path string) (gjson.Result, error) {
	return h.getResource(ctx, path, nil)
}

func (h httpClient) getResource(ctx context.Context, path string, params map[string]string) (gjson.Result, error) {
	resp, err := h.restClient.R().SetQueryParams(params).SetContext(ctx).Get(h.buildRequestURI(path))
	if err != nil {
		return gjson.Result{}, err
	}
	if !resp.IsSuccess() {
		log.Warnf("FHIR server replied: %s", resp.String())
		return gjson.Result{}, fmt.Errorf("unable to read FHIR resource (path=%s,http-status=%d): %w", path, resp.StatusCode, err)
	}
	return gjson.ParseBytes(resp.Body()), nil
}

func (h httpClient) buildRequestURI(path string) string {
	if !strings.HasSuffix(path, "/") && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return h.url + path
}
