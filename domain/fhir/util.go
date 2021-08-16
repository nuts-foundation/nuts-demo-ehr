package fhir

import (
	"github.com/tidwall/gjson"
	"strings"
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

func FilterResources(resources []gjson.Result, codingSystem CodingSystem, code Code) []gjson.Result {
	var result []gjson.Result
	for _, resource := range resources {
		for _, coding := range resource.Get("code.coding").Array() {
			if strings.TrimSpace(coding.Get("system").String()) == string(codingSystem) && strings.TrimSpace(coding.Get("code").String()) == string(code) {
				result = append(result, resource)
				break
			}
		}
	}
	return result
}
