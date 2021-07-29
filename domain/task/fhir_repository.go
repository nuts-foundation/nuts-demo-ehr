package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/tidwall/gjson"
)

// Coding systems:
const SnomedCodingSystem = "http://snomed.info/sct"
const LoincCodingSystem = "http://loinc.org"
const NutsCodingSystem = "http://nuts.nl"

// Codes:
const SnomedTransferCode = "308292007"
const TransferDisplay = "Overdracht van zorg"
const LoincAdvanceNoticeCode = "57830-2"
const SnomedAlternaticeDateCode = "146851000146105"
const SnomedNursingHandoffCode = "371535009"

type fhirTask struct {
	data gjson.Result
}

func (task *fhirTask) UnmarshalFromDomainTask(domainTask domain.Task) error {
	fhirData := map[string]interface{}{
		"resourceType": "Task",
		"id":           domainTask.ID,
		"status":       domainTask.Status,
		// TODO: patient seems mandatory in the spec, but can only be sent when placer already
		// has patient in care to protect the identity of the patient during the negotiation phase.
		//"for": map[string]string{
		//	"reference": fmt.Sprintf("Patient/%s", domainTask.PatientID),
		//},
		"code": map[string]interface{}{
			"coding": []map[string]interface{}{{
				"system":  SnomedCodingSystem,
				"code":    SnomedTransferCode,
				"display": TransferDisplay,
			}},
		},
		"requester": map[string]interface{}{
			"identifier": map[string]interface{}{
				"value":  fmt.Sprintf("%s", domainTask.RequesterID),
				"system": NutsCodingSystem,
			},
		},
		"owner": map[string]interface{}{
			"identifier": map[string]interface{}{
				"value":  fmt.Sprintf("%s", domainTask.OwnerID),
				"system": NutsCodingSystem,
			},
		},
	}
	jsonData, err := json.Marshal(fhirData)
	if err != nil {
		return fmt.Errorf("error unmarshalling fhirTask from domain.TransferNegotiation: %w", err)
	}

	*task = *newFHIRTaskFromJSON(jsonData)
	return nil
}

func (task fhirTask) MarshalToTask() (*domain.Task, error) {
	if rType := task.data.Get("resourceType").String(); rType != "Task" {
		return nil, fmt.Errorf("invalid resource type. got: %s, expected Task", rType)
	}
	codeQuery := fmt.Sprintf("code.coding.#(system==%s).code", SnomedCodingSystem)
	if codeValue := task.data.Get(codeQuery).String(); codeValue != SnomedTransferCode {
		return nil, fmt.Errorf("unexpecting coding: %s", codeValue)
	}
	return &domain.Task{
		ID: task.data.Get("id").String(),
		TaskProperties: domain.TaskProperties{
			Status:      task.data.Get("status").String(),
			OwnerID:     task.data.Get("owner.identifier.value").String(),
			RequesterID: task.data.Get("requester.identifier.value").String(),
			PatientID:   strings.Split(task.data.Get("for.reference").String(), "/")[1],
		},
		FHIRAdvanceNoticeID:  nil,
		FHIRNursingHandoffID: nil,
		AlternativeDate:      nil,
	}, nil
}

func newFHIRTaskFromJSON(data []byte) *fhirTask {
	return &fhirTask{data: gjson.ParseBytes(data)}
}

type fhirTaskRepository struct {
	url     string
	factory Factory
}

func NewFHIRTaskRepository(factory Factory, url string) *fhirTaskRepository {
	return &fhirTaskRepository{
		url:     url,
		factory: factory,
	}
}

func (r fhirTaskRepository) Create(ctx context.Context, taskProperties domain.TaskProperties) (*domain.Task, error) {
	fTask := fhirTask{}
	newTask := r.factory.New(taskProperties)
	if err := fTask.UnmarshalFromDomainTask(*newTask); err != nil {
		return nil, err
	}

	// Use a PUT method here so we can provide client generated resource IDs.
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, r.url+"/Task/"+newTask.ID, bytes.NewBuffer([]byte(fTask.data.Raw)))
	if err != nil {
		return nil, fmt.Errorf("unable to buid PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, ioErr := io.ReadAll(resp.Body)
		if ioErr != nil {
			return nil, fmt.Errorf("unable to create new patient. Unable to read error response: ioerr: %s", ioErr)
		}
		return nil, fmt.Errorf("unable to create new patient: %s", body)
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Debug(body)
	}

	return newTask, nil
}
