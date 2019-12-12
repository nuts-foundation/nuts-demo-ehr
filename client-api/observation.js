const router        = require('express').Router();
const {observation} = require('../resources/database');

router.get('/byPatientId/:id', async (req, res) => {
  try {
    // Query the "database" for the requested observations
    const o = await observation.byPatientId(req.params.id);
    if ( o )
      res.status(200).send(o).end();
    else
      res.status(404).send('Observation not found').end();
  } catch(e) {
    res.status(500).send(`Error in database query for observations by patient id: ${e}`);
  }
});

router.put('/', async (req, res) => {
  try {
    // Store a new observation in the "database"
    const o = await observation.store({
      patientId: req.body.patientId,
      content:   req.body.content,
      timestamp: new Date().toLocaleString('nl')
    });
    res.status(201).send(o).end();
  } catch(e) {
    res.status(500).send(`Error in database query for storing a new observation: ${e}`);
  }
});

module.exports = router;
