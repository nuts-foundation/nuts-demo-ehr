package fhir

import (
	"context"
	"fmt"
)

type Service struct {
	ClientFactory Factory
}

func (s Service) GetTask(ctx context.Context, customerID int, taskID string) (map[string]interface{}, error) {
	fhirClient := s.ClientFactory(WithTenant(customerID))
	result := make(map[string]interface{}, 0)
	err := fhirClient.ReadOne(ctx, "Task/"+taskID, &result)
	if err != nil {
		return nil, fmt.Errorf("error while fetching task (task-id=%s): %w", taskID, err)
	}
	return result, nil
}
