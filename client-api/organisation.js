const router = require('express').Router()
const config = require('../util/config')
const { registry } = require('../resources/nuts-node')

router.get('/me', async (req, res) => {
  let info = {
    ...config.organisation,
    user: req.session.user
  }
  if (!info.user) {
    return res.status(403).send("no active session")
  }
  res.status(200).send(info)
})

router.get('/search/:query', async (req, res) => {
  try {
    const results = await registry.searchOrganizations(req.params.query)
    res.status(200).send(results).end()
  } catch (error) {
    return res.status(500).send(`Error in search: ${error}`)
  }
})

router.get('/byURN/:urn', async (req, res) => {
  try {
    const organisation = await registry.organizationById(req.params.urn)
    res.status(200).send(organisation).end()
  } catch (error) {
    return res.status(500).send(`Error in fetching organisation: ${error}`)
  }
})

module.exports = router
