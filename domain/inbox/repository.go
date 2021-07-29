package inbox

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type Repository interface {
	List(ctx context.Context) ([]domain.InboxEntry, error)
}

func NewFHIRRepository(url string) Repository {
	return fhirRepository{url: url}
}

type fhirRepository struct {
	url string
}

func (f fhirRepository) List(ctx context.Context) ([]domain.InboxEntry, error) {
	client := http.Client{}
	res, err := client.Get(f.url + "/Task")
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := getInboxEntries(string(data))
	return result, nil
}

func getInboxEntries(taskJSON string) []domain.InboxEntry {
	var result []domain.InboxEntry
	transferResources := fhir.FilterResources(gjson.Parse(taskJSON).Get("entry.#.resource").Array(), fhir.SnomedCodingSystem, fhir.SnomedTransferCode)
	for _, resource := range transferResources {
		result = append(result, domain.InboxEntry{Title: resource.Get("code.coding.0.display").String()})
	}
	return result
}
