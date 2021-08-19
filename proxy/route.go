package proxy

import "net/http"

type fhirRoute struct {
	path      string
	operation string
}

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

	return &fhirRoute{path: request.URL.Path, operation: operation}
}
