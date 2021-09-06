package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPClient struct {
	NutsNodeAddress string
}

func (c HTTPClient) getNodeURL() string {
	url := c.NutsNodeAddress
	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("http://%v", c.NutsNodeAddress)
	}
	return url
}

func testAndReadResponse(status int, response *http.Response) ([]byte, error) {
	if response.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}
	if err := testResponseCode(status, response); err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func testResponseCode(expectedStatusCode int, response *http.Response) error {
	if response.StatusCode != expectedStatusCode {
		responseData, _ := io.ReadAll(response.Body)
		return fmt.Errorf("server returned HTTP %d (expected: %d), response: %s",
			response.StatusCode, expectedStatusCode, string(responseData))
	}
	return nil
}
