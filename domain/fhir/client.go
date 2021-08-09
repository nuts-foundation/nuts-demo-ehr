package fhir

import (
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
)

func NewClient(baseURL string) Client {
	return &httpClient{url: baseURL}
}

type Client interface {
	GetResources(path string, params map[string]string) ([]gjson.Result, error)
}

type httpClient struct {
	url string
}

func (h httpClient) GetResources(path string, params map[string]string) ([]gjson.Result, error) {
	requestURI, err := h.buildRequestURI(path, params)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	res, err := client.Get(requestURI)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return gjson.ParseBytes(data).Get("entry.#.resource").Array(), nil
}

func (h httpClient) buildRequestURI(path string, queryParams map[string]string) (string, error) {
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
