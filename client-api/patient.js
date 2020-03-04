const router = require('express').Router()
const { patient } = require('../resources/database')

router.get('/all', async (req, res) => {
  try {
    // Query the "database" for the requested patients
    const p = await patient.all()
    res.status(200).send(p).end()
  } catch (e) {
    res.status(500).send(`Error in database query for all patients: ${e}`)
  }
})

router.get('/:id', async (req, res) => {
  try {
    // Query the "database" for the requested patient
    const p = await patient.byId(req.params.id)
    if (p) { res.status(200).send(p).end() } else { res.status(404).send('Patient not found').end() }
  } catch (e) {
    res.status(500).send(`Error in database query for single patient: ${e}`)
  }
})

module.exports = router
