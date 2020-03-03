const router = require('express').Router()
const config = require('../util/config')
const { auth } = require('../resources/nuts-node')

const {
  patient,
} = require('../resources/database')

router.get('/jump', findPatient, async (req, res) => {


  if (!req.session.nuts_auth_token) {
    res.redirect("/#irma-login")
  }

  // build context
  let context = {
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

  // Get the access token at the custodians Nuts node
  let accessTokenResponse
  try {
    accessTokenResponse = await auth.createAccessToken("http://localhost:11323", jwtBearerTokenResponse.bearer_token)
    console.log(accessTokenResponse)
    // Make the jump!
    res.redirect(`http://localhost:80/sso/land?token=${accessTokenResponse.access_token}`)
  } catch(e) {
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
  let accessToken = req.query.token
  if (!accessToken) {
    res.status(404).send('missing access token')
  }

  // Introspect the token at the local Nuts node
  let introspectionResponse = await auth.introspectAccessToken(accessToken)
  if (!introspectionResponse.active) {
    res.status(401).send("invalid token")
  }

  res.status(200).send(introspectionResponse)
})

async function findPatient (req, res, next) {
  try {
    req.patient = await patient.byId(req.query.patient)
    next()
  } catch (e) {
    res.status(404).send(`Could not find a patient with id ${req.query.patient}: ${e}`)
  }
}

module.exports = router
