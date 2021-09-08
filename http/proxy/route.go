package proxy

import (
	"net/http"
	"net/url"
)

type fhirRoute struct {
	url       url.URL
	operation string
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

	return &fhirRoute{url: *request.URL, operation: operation}
}

func (fr fhirRoute) path() string {
	return fr.url.Path
}
