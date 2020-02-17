const router = require('express').Router();
const config = require('../util/config');
const {auth} = require('../resources/nuts-node');

router.get('/new-session', async (req, res) => {
  try {
    const contract = await auth.createLoginSession();
    res.status(200).send(contract).end();
  } catch(error) {
    return res.status(500).send(`Error in creating contract: ${error}`);
  }
});

module.exports = router;
