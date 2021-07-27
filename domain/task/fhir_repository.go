package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/tidwall/gjson"
)

type fhirTask struct {
	data gjson.Result
}

func (task *fhirTask) UnmarshalFromDomainTask(domainTask domain.Task) error {
	fhirData := map[string]interface{}{
		"status": domainTask.Status,
		"for": map[string]string{
			"reference": fmt.Sprintf("Patient/%s", domainTask.PatientID),
		},
		"code": map[string]interface{}{
			"coding": []map[string]interface{}{{
				"system":  SnomedCodingSystem,
				"code":    SnomedTransferCode,
				"display": TransferDisplay,
			}},
		},
		"owner": map[string]string{
			"reference": fmt.Sprintf("Organization/%s", domainTask.OwnerID),
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
		ID:                   task.data.Get("id").String(),
		Status:               task.data.Get("status").String(),
		PatientID:            strings.Split(task.data.Get("for.reference").String(), "/")[1],
		OwnerID:              strings.Split(task.data.Get("owner.reference").String(), "/")[1],
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

func (r fhirTaskRepository) Create(ctx context.Context, task domain.Task) (*domain.Task, error) {
	fTask := fhirTask{}
	if err := fTask.UnmarshalFromDomainTask(task); err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Post(r.url+"/Task", "application/json", bytes.NewBuffer([]byte(fTask.data.Raw)))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, ioErr := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to create new patient. Unable to read error response: ioerr: %s", ioErr)
		}
		return nil, fmt.Errorf("unable to create new patient: %s", body)
	}

	return &task, nil
}
