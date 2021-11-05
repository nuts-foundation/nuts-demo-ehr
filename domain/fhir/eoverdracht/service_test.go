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

var conditionFHIR = `
{
    "resourceType": "Condition",
    "id": "ae889298-d09c-477a-bd80-227da1868b85",
    "meta": {
        "extension": [
            {
                "url": "http://hapifhir.io/fhir/StructureDefinition/resource-meta-source",
                "valueUri": "#7ei18dobJKzqoHON"
            }
        ],
        "versionId": "1",
        "lastUpdated": "2021-09-22T14:58:12.496+00:00"
    },
    "note": [
        {
            "text": "we have  problem"
        }
    ]
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
	mockID := idGenerator.GenerateID()
	expected["id"] = mockID
	path := "Task/" + mockID

	mockClient := fhir.NewMockClientWithReadMock(t, map[string]map[string]interface{}{path: expected})
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
	json.Unmarshal([]byte(conditionFHIR), &expected)
	path := "Condition/ae889298-d09c-477a-bd80-227da1868b85"
	mockClient := fhir.NewMockClientWithReadMock(t, map[string]map[string]interface{}{path: expected})
	service := transferService{fhirClient: mockClient}
	section := fhir.CompositionSection{
		Entry: []datatypes.Reference{{
			Reference: fhir.ToStringPtr(path),
		}, {
			Reference: fhir.ToStringPtr("Procedure/7295c1db-b0b3-4e2f-a5ca-1ff95cecea70"),
		},
		},
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
	assert.Equal(t, "ae889298-d09c-477a-bd80-227da1868b85", fhir.FromIDPtr(condition.(*resources.Condition).ID))
}
