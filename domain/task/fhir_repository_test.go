package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func Test_fhirTask_Marshalling(t *testing.T) {
	gsonData := gjson.Parse(`
{
	"resourceType": "Task",
	"id": "123-22",
	"status": "requested",
	"code": {
		"coding": [{
			"system": "http://snomed.info/sct",
			"code": "308292007"
		}],
	},
	"requester": {
		"agent": {
			"identifier": {
				"system": "http://nuts.nl",
				"value": "did:nuts:7d032f65-a638-44e7-a7eb-d69b6d7b8c81",
			},
			"display": "Princes Amalia Ziekenhuis"
		}
	},
	"owner": {
		"identifier": {
			"system": "http://nuts.nl",
			"value": "did:nuts:1832473e-2fe0-452d-abe9-3cdb9879522f"
		},
		"display": "VVT Instelling de Regenboog"
	},
}`)
	fTask := fhirTask{data: gsonData}
	task, err := fTask.MarshalToTask()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "did:nuts:1832473e-2fe0-452d-abe9-3cdb9879522f", task.OwnerID, "expected correct ownerID")
	assert.Equal(t, "did:nuts:7d032f65-a638-44e7-a7eb-d69b6d7b8c81", task.RequesterID, "expected correct requesterID")
	assert.Equal(t, "123-22", task.ID)
	assert.Equal(t, "requested", task.Status)
	//assert.Equal(t, "999", task.PatientID, "expected correct patientID")

	newFHIRTask := fhirTask{}
	err = newFHIRTask.UnmarshalFromDomainTask(*task)
	if !assert.NoError(t, err) {
		return
	}

	task2, err := newFHIRTask.MarshalToTask()
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, task, task2, "task should be the same after double marshalling")
}
