package inbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const tasks = `
{
  "resourceType": "Bundle",
  "id": "900575ae-bc60-4627-9d78-14308b3b1afb",
  "meta": {
    "lastUpdated": "2021-07-29T08:11:47.116+00:00"
  },
  "type": "searchset",
  "total": 2,
  "link": [ {
    "relation": "self",
    "url": "http://localhost:4004/hapi-fhir-jpaserver/fhir/Task"
  } ],
  "entry": [ {
    "fullUrl": "http://localhost:4004/hapi-fhir-jpaserver/fhir/Task/219157",
    "resource": {
      "resourceType": "Task",
      "id": "219157",
      "meta": {
        "versionId": "1",
        "lastUpdated": "2021-07-28T05:26:27.570+00:00",
        "source": "#awDxdDgFlFh79UR9"
      },
      "text": {
        "status": "generated",
        "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\"><p><b>Generated Narrative with Details</b></p><p><b>id</b>: example1</p><p><b>contained</b>: </p><p><b>identifier</b>: 20170201-001 (OFFICIAL)</p><p><b>basedOn</b>: General Wellness Careplan</p><p><b>groupIdentifier</b>: G20170201-001 (OFFICIAL)</p><p><b>status</b>: in-progress</p><p><b>businessStatus</b>: waiting for specimen <span>(Details )</span></p><p><b>intent</b>: order</p><p><b>priority</b>: routine</p><p><b>code</b>: Lipid Panel <span>(Details )</span></p><p><b>description</b>: Create order for getting specimen, Set up inhouse testing,  generate order for any sendouts and submit with specimen</p><p><b>focus</b>: <a>Lipid Panel Request</a></p><p><b>for</b>: <a>Peter James Chalmers</a></p><p><b>encounter</b>: <a>Example In-Patient Encounter</a></p><p><b>executionPeriod</b>: 31/10/2016 8:25:05 AM --&gt; (ongoing)</p><p><b>authoredOn</b>: 31/10/2016 8:25:05 AM</p><p><b>lastModified</b>: 31/10/2016 9:45:05 AM</p><p><b>requester</b>: <a>Dr Adam Careful</a></p><p><b>performerType</b>: Performer <span>(Details : {http://terminology.hl7.org/CodeSystem/task-performer-type code 'performer' = 'performer', given as 'Performer'})</span></p><p><b>owner</b>: <a>Clinical Laboratory @ Acme Hospital</a></p><p><b>reasonCode</b>: The Task.reason should only be included if there is no Task.focus or if it differs from the reason indicated on the focus <span>(Details )</span></p><p><b>note</b>: This is an example to demonstrate using task for actioning a servicerequest and to illustrate how to populate many of the task elements - this is the parent task that will be broken into subtask to grab the specimen and a sendout lab test </p><p><b>relevantHistory</b>: Author's Signature. Generated Summary: id: signature; recorded: 31/10/2016 8:25:05 AM; </p><h3>Restrictions</h3><table><tr><td>-</td><td><b>Repetitions</b></td><td><b>Period</b></td></tr><tr><td>*</td><td>1</td><td>?? --&gt; 02/11/2016 9:45:05 AM</td></tr></table></div>"
      },
      "contained": [ {
        "resourceType": "Provenance",
        "id": "signature",
        "target": [ {
          "reference": "ServiceRequest/physiotherapy"
        } ],
        "recorded": "2016-10-31T08:25:05+10:00",
        "agent": [ {
          "role": [ {
            "coding": [ {
              "system": "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
              "code": "AUT"
            } ]
          } ],
          "who": {
            "reference": "Practitioner/f202",
            "display": "Luigi Maas"
          }
        } ],
        "signature": [ {
          "type": [ {
            "system": "urn:iso-astm:E1762-95:2013",
            "code": "1.2.840.10065.1.12.1.1",
            "display": "Author's Signature"
          } ],
          "when": "2016-10-31T08:25:05+10:00",
          "who": {
            "reference": "Practitioner/example",
            "display": "Dr Adam Careful"
          },
          "targetFormat": "application/fhir+xml",
          "sigFormat": "application/signature+xml",
          "data": "dGhpcyBibG9iIGlzIHNuaXBwZWQ="
        } ]
      } ],
      "identifier": [ {
        "use": "official",
        "system": "http:/goodhealth.org/identifiers",
        "value": "20170201-001"
      } ],
      "basedOn": [ {
        "display": "General Wellness Careplan"
      } ],
      "groupIdentifier": {
        "use": "official",
        "system": "http:/goodhealth.org/accession/identifiers",
        "value": "G20170201-001"
      },
      "status": "in-progress",
      "businessStatus": {
        "text": "waiting for specimen"
      },
      "intent": "order",
      "priority": "routine",
      "code": {
        "text": "Lipid Panel"
      },
      "description": "Create order for getting specimen, Set up inhouse testing,  generate order for any sendouts and submit with specimen",
      "focus": {
        "reference": "ServiceRequest/lipid",
        "display": "Lipid Panel Request"
      },
      "for": {
        "reference": "Patient/example",
        "display": "Peter James Chalmers"
      },
      "encounter": {
        "reference": "Encounter/example",
        "display": "Example In-Patient Encounter"
      },
      "executionPeriod": {
        "start": "2016-10-31T08:25:05+10:00"
      },
      "authoredOn": "2016-10-31T08:25:05+10:00",
      "lastModified": "2016-10-31T09:45:05+10:00",
      "requester": {
        "reference": "Practitioner/example",
        "display": "Dr Adam Careful"
      },
      "performerType": [ {
        "coding": [ {
          "system": "http://terminology.hl7.org/CodeSystem/task-performer-type",
          "code": "performer",
          "display": "Performer"
        } ],
        "text": "Performer"
      } ],
      "owner": {
        "reference": "Organization/1832473e-2fe0-452d-abe9-3cdb9879522f",
        "display": "Clinical Laboratory @ Acme Hospital"
      },
      "reasonCode": {
        "text": "The Task.reason should only be included if there is no Task.focus or if it differs from the reason indicated on the focus"
      },
      "note": [ {
        "text": "This is an example to demonstrate using task for actioning a servicerequest and to illustrate how to populate many of the task elements - this is the parent task that will be broken into subtask to grab the specimen and a sendout lab test "
      } ],
      "relevantHistory": [ {
        "reference": "#signature",
        "display": "Author's Signature"
      } ],
      "restriction": {
        "repetitions": 1,
        "period": {
          "end": "2016-11-02T09:45:05+10:00"
        }
      }
    },
    "search": {
      "mode": "match"
    }
  }, {
    "fullUrl": "http://localhost:4004/hapi-fhir-jpaserver/fhir/Task/220760",
    "resource": {
      "resourceType": "Task",
      "id": "220760",
      "meta": {
        "versionId": "1",
        "lastUpdated": "2021-07-29T08:11:43.777+00:00",
        "source": "#qRmF5jVO8xMHnVLw"
      },
      "status": "requested",
      "code": {
        "coding": [ {
          "system": "http://snomed.info/sct",
          "code": "308292007",
          "display": "Overdracht van zorg"
        } ]
      },
      "requester": {
        "identifier": {
          "system": "http://nuts.nl",
          "value": "test1"
        }
      },
      "owner": {
        "identifier": {
          "system": "http://nuts.nl",
          "value": "did:nuts:JCx4c3ufdKNgaZJ4h54AghY8ZgCznptNpjHUtzvVgcvW"
        }
      }
    },
    "search": {
      "mode": "match"
    }
  } ]
}`

func Test_getInboxEntries(t *testing.T) {
	entries := getInboxEntries(tasks)
	if !assert.Len(t, entries, 1) {
		return
	}
	assert.Equal(t, "Overdracht van zorg", entries[0].Title)
}