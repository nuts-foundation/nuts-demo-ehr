package fhir

import (
	"context"
	"fmt"
	"strings"

	"github.com/monarko/fhirgo/STU3/resources"
)

type Repository interface {
	GetTask(ctx context.Context, fhirTaskID string) (resources.Task, error)
	UpdateTask(ctx context.Context, task resources.Task) error
	ResolveComposition(ctx context.Context, compositionID string) (*Composition, map[string][]interface{}, *resources.Patient, error)

	CreateOrUpdateResource(ctx context.Context, resource interface{}) error
}

type fhirRepository struct {
	client Client
}

func NewFHIRRepository(client Client) Repository {
	return &fhirRepository{client: client}
}

func (s fhirRepository) GetTask(ctx context.Context, fhirTaskID string) (resources.Task, error) {
	task := resources.Task{}
	err := s.client.ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while fetching task (task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}

func (s fhirRepository) UpdateTask(ctx context.Context, task resources.Task) error {
	return s.client.CreateOrUpdate(ctx, task)
}

func (s fhirRepository) ResolveComposition(ctx context.Context, compositionPath string) (*Composition, map[string][]interface{}, *resources.Patient, error) {
	composition := Composition{}
	patient := resources.Patient{}
	sections := map[string][]interface{}{}

	err := s.client.ReadOne(ctx, "/"+compositionPath, &composition)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while fetching the advance notice composition(composition-id=%s): %w", compositionPath, err)
	}

	// Fetch the Patient
	err = s.client.ReadOne(ctx, "/"+FromStringPtr(composition.Subject.Reference), &patient)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while fetching the transfer subject (patient): %w", err)
	}

	// Fill a map with sections by code
	for _, l1 := range composition.Section {
		l1code := FromCodePtr(l1.Code.Coding[0].Code)
		if l1code == "" {
			continue
		}
		for _, l2 := range l1.Section {
			l2code := FromCodePtr(l2.Code.Coding[0].Code)
			if l2code == "" {
				continue
			}
			var entryResources []interface{}
			for _, l2entry := range l2.Entry {
				entryPath := FromStringPtr(l2entry.Reference)
				resourceTypeStr := strings.Split(entryPath, "/")[0]
				switch resourceTypeStr {
				case "Condition":
					resource := resources.Condition{}
					err := s.client.ReadOne(ctx, entryPath, &resource)
					if err != nil {
						return nil, nil, nil, err
					}
					entryResources = append(entryResources, resource)
				case "Procedure":
					resource := Procedure{}
					err := s.client.ReadOne(ctx, entryPath, &resource)
					if err != nil {
						return nil, nil, nil, err
					}
					entryResources = append(entryResources, resource)
				}

			}
			sections[l1code+"/"+l2code] = entryResources
		}
	}

	return &composition, sections, &patient, nil
}

func (s fhirRepository) CreateOrUpdateResource(ctx context.Context, resource interface{}) error {
	return s.client.CreateOrUpdate(ctx, resource)
}
