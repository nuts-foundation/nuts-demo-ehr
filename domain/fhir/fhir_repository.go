package fhir

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type fhirTask struct {
	data gjson.Result
}

func (task fhirTask) MarshalToTask() (*Task, error) {
	if rType := task.data.Get("resourceType").String(); rType != "Task" {
		return nil, fmt.Errorf("invalid resource type. got: %s, expected Task", rType)
	}
	codeQuery := fmt.Sprintf("code.coding.#(system==%s).code", SnomedCodingSystem)
	if codeValue := task.data.Get(codeQuery).String(); codeValue != string(SnomedTransferCode) {
		return nil, fmt.Errorf("unexpecting coding: %s", codeValue)
	}
	patientID := ""
	if parts := strings.Split(task.data.Get("for.reference").String(), "/"); len(parts) > 1 {
		patientID = parts[1]
	}
	return &Task{
		ID: task.data.Get("id").String(),
		TaskProperties: TaskProperties{
			Status:      task.data.Get("status").String(),
			OwnerID:     task.data.Get("owner.identifier.value").String(),
			RequesterID: task.data.Get("requester.agent.identifier.value").String(),
			PatientID:   patientID,
		},
	}, nil
}

type fhirRepository struct {
	client      Client
	taskFactory TaskFactory
}

func NewFHIRRepository(fhirClient Client) *fhirRepository {
	return &fhirRepository{
		client:      fhirClient,
		taskFactory: TaskFactory{},
	}
}

func (r fhirRepository) CreateTask(ctx context.Context, taskProperties TaskProperties) (*Task, error) {
	task := r.taskFactory.New(taskProperties)

	fhirData := map[string]interface{}{
		"resourceType": "Task",
		"id":           task.ID,
		"status":       task.Status,
		// TODO: patient seems mandatory in the spec, but can only be sent when placer already
		// has patient in care to protect the identity of the patient during the negotiation phase.
		//"for": map[string]string{
		//	"reference": fmt.Sprintf("Patient/%s", domainTask.PatientID),
		//},
		"code": CodeableConcept{Coding: Coding{
			System:  SnomedCodingSystem,
			Code:    SnomedTransferCode,
			Display: TransferDisplay,
		}},
		"requester": Requester{Agent: Organization{Identifier: Identifier{
			System: NutsCodingSystem,
			Value:  fmt.Sprintf("%s", task.RequesterID),
		}}},
		"owner": Organization{Identifier: Identifier{
			System: NutsCodingSystem,
			Value:  fmt.Sprintf("%s", task.OwnerID),
		}},
		"input":  taskProperties.Input,
		"output": taskProperties.Output,
	}
	_, err := r.client.WriteResource(ctx, "Task/"+task.ID, fhirData)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r fhirRepository) CreateComposition(ctx context.Context, elements map[string]interface{}) (*Composition, error) {
	id := uuid.New().String()
	fhirData := map[string]interface{}{
		"resourceType": "Composition",
		"id":           id,
		"status":       "final",
		"date":         time.Now().Format(time.RFC3339),
	}
	for key, value := range elements {
		fhirData[key] = value
	}
	reference := "Composition/" + id
	_, err := r.client.WriteResource(ctx, reference, fhirData)
	if err != nil {
		return nil, err
	}
	return &Composition{ID: id, Reference: reference}, nil
}
