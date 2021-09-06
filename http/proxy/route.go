package proxy

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type fhirRoute struct {
	url        url.URL
	operation  string
	resourceID string
}

// parseRoute maps the HTTP request to a FHIR route. It can differentiate between read and vread only by looking at the route.
// And there isn't an HTTP method which could be mapped to the history operation for example.
// But as it's a proxy for demo purpose this should be fine for now.
func parseRoute(request *http.Request) *fhirRoute {
	var operation string

	switch request.Method {
	case http.MethodGet:
		operation = "read"
	case http.MethodPost:
		operation = "create"
	case http.MethodPut:
		operation = "update"
	case http.MethodPatch:
		operation = "patch"
	case http.MethodDelete:
		operation = "delete"
	}

	var resourceID string
	split := strings.Split(request.URL.Path, "/")
	last := split[len(split)-1]
	if _, err := uuid.Parse(last); err == nil {
		resourceID = last
	}

	return &fhirRoute{url: *request.URL, operation: operation, resourceID: resourceID}
}

func (fr fhirRoute) path() string {
	return fr.url.Path
}
