package fhir

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/monarko/fhirgo/STU3/resources"
	"net/http"
	"strconv"
)

// InitializeTenant sets up the FHIR server for the given tenant, if the FHIR server if of a supported type.
func InitializeTenant(fhirServerType string, fhirServerURL string, tenant string) error {
	if fhirServerType != "hapi" {
		return nil
	}

	restClient := resty.New()

	// Check if tenant already exists
	response, err := restClient.R().SetQueryParam("id", tenant).Get(buildRequestURI(fhirServerURL, "DEFAULT", "$partition-management-read-partition"))
	if err != nil {
		return err
	}
	if response.IsSuccess() {
		// Tenant exists
		return nil
	}
	if response.IsError() && response.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("error while checking for HAPI FHIR Server tenant (status-code=%d,tenant=%s): %s", response.StatusCode(), tenant, string(response.Body()))
	}

	// Tenant doesn't exist (yet), create it
	idAsInt, err := strconv.Atoi(tenant)
	if err != nil {
		return fmt.Errorf("tenant is not an integer (tenant=%s): %w", tenant, err)
	}
	parameters := resources.Parameters{
		Base: resources.Base{
			ResourceType: "Parameters",
		},
		Parameter: []resources.ParametersParameter{
			{Name: ToStringPtr("id"), ValueInteger: ToIntegerPtr(idAsInt)},
			{Name: ToStringPtr("name"), ValueCode: ToCodePtr(tenant)},
		},
	}
	response, err = restClient.R().SetHeader("Content-Type", "application/json").SetBody(parameters).Post(buildRequestURI(fhirServerURL, "DEFAULT", "$partition-management-create-partition"))
	if err != nil {
		return fmt.Errorf("unable create new HAPI FHIR Server partition: %w", err)
	}
	if !response.IsSuccess() {
		return fmt.Errorf("unable create new HAPI FHIR Server partition (status-code=%d)", response.StatusCode())
	}
	return nil
}
