package inbox

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
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
	parsedData := gjson.Parse(string(data))
	var result []domain.InboxEntry
	for _, entry := range parsedData.Get("$.entry[].resource.id").Array() {
		result = append(result, domain.InboxEntry{Title: entry.String()})
	}
	return result, nil
}
