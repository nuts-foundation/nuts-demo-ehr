package eoverdracht

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/stretchr/testify/assert"
)

type mockIDGenerator struct {
	uuid *uuid.UUID
}

var expected = map[string]interface{}{
	"code": map[string]interface{}{"coding": []interface{}{map[string]interface{}{"code": "308292007", "display": "Overdracht van zorg", "system": "http://snomed.info/sct"}}},
	"id":   "13",
	"input": []interface{}{map[string]interface{}{
		"type": map[string]interface{}{
			"coding": []interface{}{map[string]interface{}{
				"code":   "57830-2",
				"system": "http://loinc.org",
			}}, "text": "Aanmeldbericht"},
		"valueReference": map[string]interface{}{"reference": "/Composition/123"},
	}, map[string]interface{}{
		"type": map[string]interface{}{
			"coding": []interface{}{map[string]interface{}{
				"code":    "371535009",
				"display": "verslag van zorg",
				"system":  "http://snomed.info/sct"}}},
		"valueReference": map[string]interface{}{"reference": "/Composition/456"},
	}},
	"owner":        map[string]interface{}{"identifier": map[string]interface{}{"system": "http://nuts.nl", "value": "did:nuts:456"}},
	"requester":    map[string]interface{}{"agent": map[string]interface{}{"identifier": map[string]interface{}{"system": "http://nuts.nl", "value": "did:nuts:123"}}},
	"resourceType": "Task",
	"status":       "requested",
}

var compositionFHIR = `
{
  "resourceType": "Composition",
  "id": "c1561592-dbcd-440d-834c-08ce2ab14015",
  "meta": {
    "extension": [ {
      "url": "http://hapifhir.io/fhir/StructureDefinition/resource-meta-source",
      "valueUri": "#Xtge9s39Ke9H4XTP"
    } ],
    "versionId": "1",
    "lastUpdated": "2021-09-22T14:58:12.677+00:00"
  },
  "type": {
    "coding": [ {
      "system": "http://loinc.org",
      "code": "57830-2"
    } ]
  },
  "subject": {
    "reference": "Patient/58327c5e-ef34-4d37-a6f2-8cde4526c1b0"
  },
  "title": "Advance notice",
  "section": [ {
    "extension": [ {
      "url": "http://nictiz.nl/fhir/StructureDefinition/eOverdracht-TransferDate",
      "valueDateTime": "2021-09-23T00:00:00Z"
    } ],
    "title": "Administrative data",
    "code": {
      "coding": [ {
        "system": "http://snomed.info/sct",
        "code": "405624007",
        "display": "Administrative documentation (record artifact)"
      } ]
    }
  }, {
    "code": {
      "coding": [ {
        "system": "http://snomed.info/sct",
        "code": "773130005",
        "display": "Nursing care plan (record artifact)"
      } ]
    },
    "section": [ {
      "title": "Current patient problems",
      "code": {
        "coding": [ {
          "system": "http://snomed.info/sct",
          "code": "86644006",
          "display": "Nursing diagnosis"
        } ]
      },
      "entry": [ {
        "reference": "Condition/ae889298-d09c-477a-bd80-227da1868b85"
      }, {
        "reference": "Condition/7c863eb7-243c-4bf6-b953-88757b14dd93"
      }, {
        "reference": "Procedure/7295c1db-b0b3-4e2f-a5ca-1ff95cecea70"
      } ]
    } ]
  } ]
}`

func (m *mockIDGenerator) GenerateID() string {
	if m.uuid == nil {
		id := uuid.New()
		m.uuid = &id
	}
	return m.uuid.String()
}

func Test_transferService_CreateTask(t *testing.T) {
	idGenerator := &mockIDGenerator{}
	expected["id"] = idGenerator.GenerateID()
	fhirBuilder := FHIRBuilder{IDGenerator: idGenerator}

	advanceNoticeID := "123"
	nursingHandoffID := "456"

	requestorID := "did:nuts:123"
	ownerID := "did:nuts:456"

	mockClient := fhir.NewMockClientWithExpectedCreateOrUpdate(t, expected)
	service := &transferService{fhirClient: mockClient, resourceBuilder: fhirBuilder}
	transferTask, err := service.CreateTask(context.Background(), TransferTask{
		SenderDID:        requestorID,
		ReceiverDID:      ownerID,
		AdvanceNoticeID:  &advanceNoticeID,
		NursingHandoffID: &nursingHandoffID,
	})
	assert.NotNil(t, transferTask)
	assert.NoError(t, err)
}

func Test_transferService_GetTask(t *testing.T) {
	idGenerator := &mockIDGenerator{}
	expected["id"] = idGenerator.GenerateID()

	mockClient := fhir.NewMockClientWithReadMock(t, []map[string]interface{}{expected})
	service := transferService{fhirClient: mockClient}
	taskID := idGenerator.GenerateID()
	resolvedTask, err := service.GetTask(context.Background(), taskID)
	assert.NoError(t, err)
	assert.NotNil(t, resolvedTask)

	assert.Equal(t, idGenerator.GenerateID(), resolvedTask.ID)
	assert.Equal(t, "123", *resolvedTask.AdvanceNoticeID)
	assert.Equal(t, "456", *resolvedTask.NursingHandoffID)
	assert.Equal(t, "did:nuts:123", resolvedTask.SenderDID)
	assert.Equal(t, "did:nuts:456", resolvedTask.ReceiverDID)
}

func TestTransferService_ResolveComposition(t *testing.T) {

}

func TestTransferService_resolveCompositionEntry(t *testing.T) {
	expected := map[string]interface{}{}
	json.Unmarshal([]byte(compositionFHIR), &expected)
	mockClient := fhir.NewMockClientWithReadMock(t, []map[string]interface{}{expected})
	service := transferService{fhirClient: mockClient}
	section := fhir.CompositionSection{
		Entry: []datatypes.Reference{{
			Reference: fhir.ToStringPtr("Composition/c1561592-dbcd-440d-834c-08ce2ab14015"),
		}},
	}
	conditions, err := service.resolveCompositionEntry(context.Background(), section, resources.Condition{})
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Len(t, conditions, 1) {
		return
	}
	condition := conditions[0]
	assert.IsType(t, &resources.Condition{}, condition)
	assert.Equal(t, "c1561592-dbcd-440d-834c-08ce2ab14015", fhir.FromIDPtr(condition.(*resources.Condition).ID))
}
