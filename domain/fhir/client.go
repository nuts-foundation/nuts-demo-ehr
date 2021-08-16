package fhir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewClient(baseURL string) Client {
	return &httpClient{url: baseURL}
}

type Client interface {
	WriteResource(ctx context.Context, resourcePath string, resource interface{}) (gjson.Result, error)
	GetResources(path string, params map[string]string) ([]gjson.Result, error)
	GetResource(path string) (gjson.Result, error)
}

type httpClient struct {
	url string
}

func (h httpClient) WriteResource(ctx context.Context, resourcePath string, resource interface{}) (gjson.Result, error) {
	resourceAsJSON, err := json.Marshal(resource)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("unable to marshal FHIR resource (path=%s): %w", resourcePath, err)
	}
	requestURI, err := h.buildRequestURI(resourcePath, nil)
	if err != nil {
		return gjson.Result{}, err
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, requestURI, bytes.NewBuffer(resourceAsJSON))
	if err != nil {
		return gjson.Result{}, fmt.Errorf("unable to build FHIR resource write request (path=%s): %w", resourcePath, err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("unable to write FHIR resource (path=%s): %w", resourcePath, err)
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			log.Warnf("FHIR server replied: %s", string(body))
		}
		return gjson.Result{}, fmt.Errorf("unable to write FHIR resource (path=%s,http-status=%d): %w", resourcePath, resp.StatusCode, err)
	}
	return h.readResource(resp)
}

func (h httpClient) GetResources(path string, params map[string]string) ([]gjson.Result, error) {
	resource, err := h.getResource(path, params)
	if err == nil {
		return resource.Get("entry.#.resource").Array(), nil
	}
	return nil, err
}

func (h httpClient) GetResource(path string) (gjson.Result, error) {
	return h.getResource(path, nil)
}

func (h httpClient) getResource(path string, params map[string]string) (gjson.Result, error) {
	requestURI, err := h.buildRequestURI(path, params)
	if err != nil {
		return gjson.Result{}, err
	}
	client := http.Client{}
	res, err := client.Get(requestURI)
	if err != nil {
		return gjson.Result{}, err
	}
	return h.readResource(res)
}

func (h httpClient) readResource(res *http.Response) (gjson.Result, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	parsed := gjson.ParseBytes(data)
	return parsed, nil
}

func (h httpClient) buildRequestURI(path string, queryParams map[string]string) (string, error) {
	if !strings.HasSuffix(path, "/") && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	result, err := url.Parse(h.url + path)
	if err != nil {
		return "", err
	}
	if queryParams != nil {
		for key, value := range queryParams {
			result.Query().Add(key, value)
		}
	}
	return result.String(), nil
}
