const router = require('express').Router()
const config = require('../util/config')
const axios = require('axios')

const {
  patient,
  observation
} = require('../resources/database')

const {
  consentStore,
  registry,
  auth
} = require('../resources/nuts-node')

router.get('/byPatientId/:id', async (req, res) => {
  try {
    const observations = await observation.byPatientId(req.params.id)
    if (observations) res.status(200).send(observations).end()
    else res.status(404).send('Observation not found').end()
  } catch (e) {
    res.status(500).send(`Error in database query for observations by patient id: ${e}`)
  }
})

router.get('/remoteByPatientId/:patient_id/:urn', findPatient, async (req, res) => {
  try {
    // Is our application allowed to make this request?
    const consents = await consentStore.consentsFor({
      subject: req.patient,
      actor: config.organisation,
      custodian: { urn: req.params.urn }
    })
    if (!consents || consents.totalResults === 0) { return res.status(403).send('You don\'t have consent for this request').end() }

    // Is this user authenticated with IRMA?
    if (!req.session.nuts_auth_token) { return res.status(401).send('You\'re not authenticated').end() }

    // Can we find remote endpoints?
    const DEMO_ENDPOINT_TYPE = 'urn:oid:1.3.6.1.4.1.54851.2:demo-ehr'
    const endpoints = await registry.endpointsByOrganisationId(req.params.urn, DEMO_ENDPOINT_TYPE)
    if (!endpoints || endpoints.length === 0) { return res.status(500).send('Can\'t find remote endpoints').end() }

    // create a context
    const context = {
      subject: `urn:oid:2.16.840.1.113883.2.4.6.3:${req.patient.bsn}`,
      actor: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`,
      custodian: req.params.urn,
      identity: req.session.nuts_auth_token,
      scope: 'demo' // note that this is not an official scope
    }

    const observations = await Promise.all(
      endpoints.map(endpoint => {
        const url = `${endpoint.URL}/observations`

        // Access tokens not supported in v0.12
        if (config.nuts.version == "0.12") {
          return axios.get(url, {
            headers: {
              sid: `urn:oid:2.16.840.1.113883.2.4.6.3:${req.patient.bsn}`,
              sub: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`,
              name: req.session.user,
              Authorization: req.session.nuts_auth_token
            }
          })
          .then(response => response.data);
        } else {
          return auth.obtainAccessToken(context, endpoint)
            .then(accessToken =>
              // Fetch available observations from all available endpoints
              axios.get(url, { headers: { Authorization: accessToken } })
                .then(response => response.data)
                .catch((e) => {
                  throw Error(`Could not get observations from ${url}: ${e}`)
                })
            )
        }
      })
    )

    res.status(200).send(observations.flat()).end()
  } catch (e) {
    res.status(500).send(`Error while getting remote observations: ${e}`)
  }
})

router.put('/', async (req, res) => {
  try {
    // Store a new observation in the "database"
    const o = await observation.store({
      patientId: req.body.patientId,
      content: req.body.content,
      timestamp: new Date().toLocaleString('nl')
    })
    res.status(201).send(o).end()
  } catch (e) {
    res.status(500).send(`Error in database query for storing a new observation: ${e}`)
  }
})

async function findPatient (req, res, next) {
  try {
    req.patient = await patient.byId(req.params.patient_id)
    next()
  } catch (e) {
    res.status(404).send(`Could not find a patient with id ${req.params.patient_id}: ${e}`)
  }
}

module.exports = router
