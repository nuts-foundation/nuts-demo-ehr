const router        = require('express').Router();
const {observation} = require('../resources/database');
const config        = require('../util/config');
const axios         = require('axios');

const {
  patient,
} = require('../resources/database');

const {
  consentStore,
  registry
} = require('../resources/nuts-node');

router.get('/byPatientId/:id', async (req, res) => {
  try {
    // Query the "database" for the requested observations
    const o = await observation.byPatientId(req.params.id);
    if ( o )
      res.status(200).send(o).end();
    else
      res.status(404).send('Observation not found').end();
  } catch(e) {
    res.status(500).send(`Error in database query for observations by patient id: ${e}`);
  }
});

router.get('/remoteByPatientId/:patient_id', findPatient, async(req, res)=> {
  try {
    var patientId = req.params.patient_id
    var patientBSN = req.patient.bsn

    // get received consents
    const consents = await consentStore.consentsFor({
      subject: patientId,
      actor:   config.organisation
    });

    if ( consents.totalResults == 0 )
      return res.status(200).send([])

    // Endpoint type of this custom jston health data api endpoint
    const DEMO_ENDPOINT_TYPE = "urn:ietf:rfc:3986:urn:oid:1.3.6.1.4.1.54851.2:demo-ehr"
    // for each consent, get endpoints
    const endpoints = await Promise.all(
        consents.results.map( async (item, i) => {
          return registry.endpointsByOrganisationId(item.custodian, DEMO_ENDPOINT_TYPE)
        })
      )
    const urls = endpoints.flat().map((ep) => ep.URL)

    // for each endpoint, get observations
    const remoteObservations = await Promise.all(
        urls.map(async (baseURL) => {
          const url = baseURL + '/external/patient/' + patientBSN + '/observations';
          return axios.get(url)
        }))

    const observations = remoteObservations.flatMap((response)=> response.data)

    res.status(200).send(observations).end()

  } catch(e) {
    res.status(500).send(`Error while getting remote observations: ${e}`)
  }
});

router.put('/', async (req, res) => {
  try {
    // Store a new observation in the "database"
    const o = await observation.store({
      patientId: req.body.patientId,
      content:   req.body.content,
      timestamp: new Date().toLocaleString('nl')
    });
    res.status(201).send(o).end();
  } catch(e) {
    res.status(500).send(`Error in database query for storing a new observation: ${e}`);
  }
});

async function findPatient(req, res, next) {
  try {
    req.patient = await patient.byId(req.params.patient_id);
    next();
  } catch(e) {
    res.status(404).send(`Could not find a patient with id ${req.params.patient_id}: ${e}`);
  }
}

module.exports = router;
