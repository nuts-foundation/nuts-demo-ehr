const config = require('../util/config');
const router = require('express').Router();

const {
  patient,
  observation
} = require('../resources/database');

const {
  consentStore
} = require('../resources/nuts-node');

router.get('/:patient_id', findPatient, handleNutsAuth, async (req, res) => {
  // Work already done by middleware
  res.status(200).send(req.patient).end();
});

router.get('/:patient_id/observations', findPatient, handleNutsAuth, async (req, res) => {
  try {
    // Query the "database" for the requested observations
    const o = await observation.byPatientId(req.patient.id);
    if ( o )
      res.status(200).send(o).end();
    else
      res.status(404).send('Observation not found').end();
  } catch(e) {
    res.status(500).send(`Error in database query for observations by patient id: ${e}`);
  }
});

async function findPatient(req, res, next) {
  try {
    req.patient = await patient.byId(req.params.patient_id);
    next();
  } catch(e) {
    res.status(404).send(`Could not find a patient with id ${req.params.patient_id}: ${e}`);
  }
}

async function handleNutsAuth(req, res, next) {
  try {
    if ( !req.headers.token )
      return res.status(403).send(`Request not allowed without access token`).end();

    const allowed = await consentStore.checkConsent({
      subject:   req.patient,
      custodian: config.organisation,
      actor:     {
        agb: req.headers.token
      }
    });

    if ( allowed.consentGiven !== 'true' )
      return res.status(403).send(`Request not allowed for actor with AGB ${req.headers.token}`).end();

    next();
  } catch(e) {
    res.status(500).send(`Error in Nuts query for consent: ${e}`);
  }
}

module.exports = router;
