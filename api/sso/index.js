const router = require('express').Router()
const config = require('../../util/config')
const patientResource = require('../../resources/database/patient')
const { auth, registry } = require('../../resources/nuts-node')

router.get('/jump', findPatient, async (req, res) => {
  // Check if there is an existing nuts-auth-token (IRMA token)
  if (!req.session.nuts_auth_token) {
    res.redirect('/#escalate')
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

  // Get the correct endpoint
  let endpointResponse
  try {
    const type = 'urn:oid:1.3.6.1.4.1.54851.1:nuts-sso'
    endpointResponse = await registry.endpointsByOrganisationId(req.query.custodian, type)
  } catch (e) {
    return res.status(500).send(`something went wrong while fetching nuts-sso endpoint for ${req.query.custodian}: ${e}`)
  }
  if (!endpointResponse) {
    return res.status(500).send(`unable to find a nuts-sso endpoint for ${req.query.custodian}`)
  }
  console.log('endpointResponse', endpointResponse)
  const endpoint = endpointResponse.pop()

  // obtain access token based on the context and the endpoint.
  // The endpoint should contain the authorization server endpoint in the properties bag
  let accessToken
  try {
    accessToken = await auth.obtainAccessToken(context, endpoint)
  } catch (e) {
    return res.status(403).send(`unable to obtain an accessToken: ${e}`).end()
  }

  // Make the jump!
  const jumpEndpoint = endpoint.URL
  res.redirect(`${jumpEndpoint}/sso/land?token=${accessToken}`)
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
