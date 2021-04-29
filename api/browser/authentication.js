const router = require('express').Router()
const config = require('../../util/config')
const { auth } = require('../../resources/nuts-node')

module.exports = function(organisation) {
  router.get('/new-session', async (req, res) => {
    try {
      const contract = await auth.drawUpContract(organisation.did)
      const session = await auth.createLoginSession(contract)
      req.session.nuts_session = session.sessionID
      res.status(200).send(session.sessionPtr.clientPtr).end()
    } catch (error) {
      return res.status(500).send(`Error in creating contract: ${error}`)
    }
  })

  router.get('/session-done', async (req, res) => {
    try {
      const status = await auth.sessionRequestStatus(req.session.nuts_session)
      if (status.status === 'DONE') {
        const verifyResponse = await auth.verifySignature(status.verifiablePresentation)
        req.session.nuts_auth_token = status.nuts_auth_token

        const attrs = verifyResponse.issuerAttributes;
        // Extract the user's full name from the disclosed attributes.
        // TODO: don't depend on the irma-demo schemeManager here
        req.session.user = [
          attrs["gemeente.personalData.initials"],
          attrs["gemeente.personalData.prefix"],
          attrs["gemeente.personalData.familyname"],
        ].join(' ')
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

  router.get('/logged-in', (req, res) => {
    if ( req.session.user )
      res.status(200).end();
    else
      res.status(403).end();
  });

  return router;
};
