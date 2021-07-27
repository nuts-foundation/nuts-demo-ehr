package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func Test_fhirTask_MarshalToTask(t *testing.T) {
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
	"for": {
		"reference": "Patient/999"
	},
	"owner": {
		"reference": "Organization/1832473e-2fe0-452d-abe9-3cdb9879522f",
	},
}`)
	fTask := fhirTask{data: gsonData}
	task, err := fTask.MarshalToTask()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "1832473e-2fe0-452d-abe9-3cdb9879522f", task.OwnerID)
	assert.Equal(t, "123-22", task.ID)
	assert.Equal(t, "requested", task.Status)
	assert.Equal(t, "999", task.PatientID)
}
