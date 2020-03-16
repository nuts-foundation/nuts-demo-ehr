const config = require('../util/config')
const router = require('express').Router()

const {
  patient,
  observation,
  accessLog
} = require('../resources/database')

const {
  consentStore,
  auth,
  registry
} = require('../resources/nuts-node')

router.get('/patient', handleNutsAuth, findPatient, logRequest, async (req, res) => {
  // Work already done by middleware
  res.status(200).send(req.patient).end()
})

router.get('/observations', handleNutsAuth, findPatient, logRequest, async (req, res) => {
  try {
    // Query the "database" for the requested observations
    const o = await observation.byPatientId(req.patient.id)
    if (o) { res.status(200).send(o).end() } else { res.status(404).send('Observation not found').end() }
  } catch (e) {
    res.status(500).send(`Error in database query for observations by patient id: ${e}`)
  }
})

async function findPatient (req, res, next) {
  console.log('findPatient')
  const requestContext = req.requestContext
  const patientBsn = requestContext.sid.match(/urn:oid:2.16.840.1.113883.2.4.6.3:([0-9]{8,9})/).pop()
  try {
    req.patient = await patient.byBSN(patientBsn)
    next()
  } catch (e) {
    // Give 500 in all cases so we don't leak information about patients in care
    console.error(`Can't find patient with bsn ${patientBsn}`)
    return res.status(500).end()
  }
}

async function handleNutsAuth (req, res, next) {
  console.log('handleNutsAuh')
  console.log('header:', req.headers)
  const accessToken = req.headers.authorization
  if (!accessToken) {
    res.status(403).send('no authorization header provided')
  }

  // introspect token
  let introspectionResponse
  try {
    // Introspect the token at the local Nuts node
    introspectionResponse = await auth.introspectAccessToken(accessToken)
    if (!introspectionResponse.active) {
      res.status(401).send('invalid token')
    }
  } catch (e) {
    res.status(500).send(`error while introspecting access token: ${e}`)
  }

  // Is this organisation allowed to make this request?
  const consentQuery = {
    subject: introspectionResponse.sid,
    actor: introspectionResponse.sub,
    custodian: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`
  }
  try {
    console.log('consentQuery:', consentQuery)
    const consent = await consentStore.checkConsent(consentQuery)
    console.log('consentResponse', consent)

    if (!consent || consent.consentGiven != 'yes') {
      // Consent was not given for this organisation
      // Give 500 in all cases so we don't leak information about patients in care
      console.error(`No consent found for actor ${introspectionResponse.sub}`)
      return res.status(500).end()
    }
  } catch (e) {
    // Give 500 in all cases so we don't leak information about patients in care
    console.error(`Could not fetch consent for triple ${consentQuery.subject}, ${consentQuery.actor}, ${consentQuery.custodian}: ${e}`)
    return res.status(500).end()
  }

  req.requestContext = introspectionResponse

  next()
}

function logRequest (req, res, next) {
  console.log('logRequest')
  try {
    // Log this request for our own audits
    accessLog.store({
      timestamp: Date.now(),
      patientId: req.patient.id,
      actor: req.requestContext.sub,
      user: req.requestContext.name
    })
  } catch (e) {
    // Give 500 in all cases so we don't leak information about patients in care
    console.error(`Could not save request audit: ${e}`)
    return res.status(500).end()
  }
  next()
}

module.exports = router
