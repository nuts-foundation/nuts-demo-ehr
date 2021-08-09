package fhir

import (
	"github.com/tidwall/gjson"
	"strings"
)

type CodingSystem string

const SnomedCodingSystem CodingSystem = "http://snomed.info/sct"
const SnomedNursingHandoffCode = "371535009"
const SnomedTransferCode = "308292007"
const NutsCodingSystem = "http://nuts.nl"

func FilterResources(resources []gjson.Result, codingSystem CodingSystem, code string) []gjson.Result {
	var result []gjson.Result
	for _, resource := range resources {
		for _, coding := range resource.Get("code.coding").Array() {
			if strings.TrimSpace(coding.Get("system").String()) == string(codingSystem) && strings.TrimSpace(coding.Get("code").String()) == code {
				result = append(result, resource)
				break
			}
		}
	}
	return result
}
