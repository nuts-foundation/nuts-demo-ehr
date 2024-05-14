package fhir

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
)

const DateTimeLayout = "2006-01-02T15:04:05-07:00"

func ToIntegerPtr(input int) *datatypes.Integer {
	result := datatypes.Integer(int32(input))
	return &result
}

func ToStringPtr(str string) *datatypes.String {
	result := datatypes.String(str)
	return &result
}

func ToDateTimePtr(str string) *datatypes.DateTime {
	result := datatypes.DateTime(str)
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
