const router = require('express').Router()
const config = require('../util/config')
const { auth, registry } = require('../resources/nuts-node')

const patientResource = require('../resources/database').patient

router.get('/jump', findPatient, async (req, res) => {
  if (!req.session.nuts_auth_token) {
    res.redirect('/#irma-login')
  }

  // build context
  const context = {
    subject: `urn:oid:2.16.840.1.113883.2.4.6.3:${req.patient.bsn}`,
    actor: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`,
    custodian: req.query.custodian,
    identity: req.session.nuts_auth_token,
    scope: 'nuts-sso'
  }
  console.log(context)

  // Get the JWT Bearer token at the local Nuts node
  let jwtBearerTokenResponse
  try {
    jwtBearerTokenResponse = await auth.createJwtBearerToken(context)
    console.log(jwtBearerTokenResponse)
  } catch (e) {
    console.log(e.response.data)
    res.status(500).send(`error while creating jwt bearer token: ${e.response.data}`)
  }

  // Get the correct endpoint
  let endpointResponse
  try {
    const type = 'urn:oid:1.3.6.1.4.1.54851.1:nuts-sso'
    endpointResponse = await registry.endpointsByOrganisationId(req.query.custodian, type)
  } catch (e) {
    if (e.response) {
      console.log(e.response.data)
    } else {
      console.log(e)
    }
  }

  console.log('endpointResponse', endpointResponse)
  const endpointEntry = endpointResponse.pop()
  const accessTokenEndpoint = endpointEntry.properties.authenticationServerURL
  const jumpEndpoint = endpointEntry.URL

  // Get the access token at the custodians Nuts node
  let accessTokenResponse
  try {
    accessTokenResponse = await auth.createAccessToken(accessTokenEndpoint, jwtBearerTokenResponse.bearer_token)
    console.log(accessTokenResponse)
    // Make the jump!
    res.redirect(`${jumpEndpoint}/sso/land?token=${accessTokenResponse.access_token}`)
  } catch (e) {
    if (e.response) {
      console.log(e.response.data)
      res.status(500).send(`error while creating access token: ${JSON.stringify(e.response.data)}`)
    } else {
      console.log(e)
      res.status(500).send(`error while creating access token: ${e}`)
    }
  }
})

router.get('/land', async (req, res) => {
  const accessToken = req.query.token
  if (!accessToken) {
    res.status(401).send('missing access token')
  }

  let introspectionResponse
  try {
    // Introspect the token at the local Nuts node
    introspectionResponse = await auth.introspectAccessToken(accessToken)
    if (!introspectionResponse.active) {
      res.status(401).send('invalid token')
    }
  } catch (e) {
    res.status(500).send('error while introspecting access token:', e)
  }

  // set the user session
  req.session.user = introspectionResponse.name

  // Get bsn from urn
  const subjectId = introspectionResponse.sid
  const patientBsn = subjectId.match(/urn:oid:2.16.840.1.113883.2.4.6.3:([0-9]{8,9})/).pop()
  const patient = await patientResource.byBSN(patientBsn)

  if (!patient) {
    res.status(401).send('patient not found')
  }

  res.redirect(`/#patient-details/${patient.id}`)
})

async function findPatient (req, res, next) {
  try {
    req.patient = await patientResource.byId(req.query.patient)
    next()
  } catch (e) {
    res.status(404).send(`Could not find a patient with id ${req.query.patient}: ${e}`)
  }
}

module.exports = router
