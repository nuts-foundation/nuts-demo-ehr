package fhir

import (
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func NewClient(baseURL string) Client {
	return &httpClient{url: baseURL}
}

type Client interface {
	GetResources(path string) ([]gjson.Result, error)
}

type httpClient struct {
	url string
}

func (h httpClient) GetResources(path string) ([]gjson.Result, error) {
	client := http.Client{}
	res, err := client.Get(h.url + path)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return gjson.ParseBytes(data).Get("entry.#.resource").Array(), nil
}
