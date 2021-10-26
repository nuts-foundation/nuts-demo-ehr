package eoverdracht

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

func (m *mockIDGenerator) String() string {
	if m.uuid == nil {
		id := uuid.New()
		m.uuid = &id
	}
	return m.uuid.String()
}

func Test_transferService_CreateTask(t *testing.T) {
	idGenerator := &mockIDGenerator{}
	expected["id"] = idGenerator.String()
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
	expected["id"] = idGenerator.String()

	mockClient := fhir.NewMockClientWithReadMock(t, []map[string]interface{}{expected})
	service := transferService{fhirClient: mockClient}
	taskID := idGenerator.String()
	resolvedTask, err := service.GetTask(context.Background(), taskID)
	assert.NoError(t, err)
	assert.NotNil(t, resolvedTask)

	assert.Equal(t, idGenerator.String(), resolvedTask.ID)
	assert.Equal(t, "123", *resolvedTask.AdvanceNoticeID)
	assert.Equal(t, "456", *resolvedTask.NursingHandoffID)
	assert.Equal(t, "did:nuts:123", resolvedTask.SenderDID)
	assert.Equal(t, "did:nuts:456", resolvedTask.ReceiverDID)
}
