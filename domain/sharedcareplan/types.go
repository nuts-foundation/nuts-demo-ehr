package sharedcareplan

type SharedCarePlan struct {
	DossierID  string `db:"dossier_id"`
	CustomerID string `db:"customer_id"`
	// ServiceURL is the FHIR base URL of the shared Care Plan Service on which this CarePlan is stored.
	ServiceURL string `db:"service_url"`
	// Reference is the FHIR Reference to the CarePlan on the shared Care Plan Service.
	// It can be used by FHIR clients to resolve the CarePlan.
	Reference string `db:"reference"`
}
