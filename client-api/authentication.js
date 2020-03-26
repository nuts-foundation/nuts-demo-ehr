const router = require('express').Router()
const config = require('../util/config')
const { auth } = require('../resources/nuts-node')

router.get('/new-session', async (req, res) => {
  try {
    const contract = await auth.createLoginSession()
    req.session.nuts_session = contract.session_id
    res.status(200).send(contract.qr_code_info).end()
  } catch (error) {
    return res.status(500).send(`Error in creating contract: ${error}`)
  }
})

router.get('/session-done', async (req, res) => {
  try {
    const status = await auth.sessionRequestStatus(req.session.nuts_session)

    if (status.status == 'DONE' && status.proofStatus === 'VALID') {
      req.session.nuts_auth_token = status.nuts_auth_token

      // Extract the user's full name from the disclosed attributes.
      // TODO: don't depend on the irma-demo schemeManager here
      const fullNameEntry = status.disclosed.find(el => el.identifier === 'irma-demo.gemeente.personalData.fullname')

      // This is still experimental, and depends on the SSO feature to land in
      // the Nuts node. When attribute not available, fall back on default user:
      if (fullNameEntry)
        req.session.user = fullNameEntry.rawvalue
      else
        req.session.user = req.session.user || config.organisation.user
    }

    res.status(200).send(status).end()
  } catch (error) {
    return res.status(500).send(`Error in fetching session status: ${error}`)
  }
})

router.post('/login', (req, res) => {
  req.session.user = req.body.username
  res.status(200).end()
})

router.get('/logout', (req, res) => {
  req.session.destroy()
  res.status(204).end()
})

module.exports = router
