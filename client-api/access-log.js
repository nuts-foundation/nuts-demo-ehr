const router        = require('express').Router();
const { accessLog } = require('../resources/database');

router.get('/byPatientId/:id', async (req, res) => {
  try {
    const logs = await accessLog.byPatientId(req.params.id);
    if ( logs ) res.status(200).send(logs).end();
    else        res.status(404).send('Logs not found').end();
  } catch(e) {
    res.status(500).send(`Error in database query for logs by patient id: ${e}`);
  }
});

module.exports = router;
