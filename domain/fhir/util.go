package fhir

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/tidwall/gjson"
)

func Filter(resources []gjson.Result, predicate func(resource gjson.Result) bool) []gjson.Result {
	var result []gjson.Result
	for _, resource := range resources {
		if predicate(resource) {
			result = append(result, resource)
		}
	}
	return result
}

func ToIntegerPtr(input int) *datatypes.Integer {
	result := datatypes.Integer(int32(input))
	return &result
}

func ToStringPtr(str string) *datatypes.String {
	result := datatypes.String(str)
	return &result
}

func FromStringPtr(str *datatypes.String) string {
	if str == nil {
		return ""
	}
	return string(*str)
}

func ToUriPtr(str string) *datatypes.URI {
	result := datatypes.URI(str)
	return &result
}

func ToCodePtr(str string) *datatypes.Code {
	result := datatypes.Code(str)
	return &result
}

func FromCodePtr(str *datatypes.Code) string {
	if str == nil {
		return ""
	}
	return string(*str)
}

func FromIDPtr(str *datatypes.ID) string {
	if str == nil {
		return ""
	}
	return string(*str)
}

func ToIDPtr(str string) *datatypes.ID {
	result := datatypes.ID(str)
	return &result
}

func toCodePtr(str string) *datatypes.Code {
	result := datatypes.Code(str)
	return &result
}
