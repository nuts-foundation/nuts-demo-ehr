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

    if (status.status == 'DONE' && status.proofStatus == 'VALID') { req.session.nuts_auth_token = status.nuts_auth_token }

    res.status(200).send(status).end()
  } catch (error) {
    return res.status(500).send(`Error in fetching session status: ${error}`)
  }
})

module.exports = router
