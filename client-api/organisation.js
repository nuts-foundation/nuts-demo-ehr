const router     = require('express').Router();
const config     = require('../util/config');
const {registry} = require('../resources/nuts-node');

router.get('/me', async (req, res) => {
  res.status(200).send(config.organisation).end();
});

router.get('/search/:query', async (req, res) => {
  try {
    const results = await registry.searchOrganizations(req.params.query);
    res.status(200).send(results).end();
  } catch(error) {
    return res.status(500).send(`Error in search: ${error}`);
  }
});

module.exports = router;
