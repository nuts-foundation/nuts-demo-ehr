const config = require('../util/config');
const router = require('express').Router();

const {
  patient,
  observation,
  accessLog
} = require('../resources/database');

const {
  consentStore,
  auth,
  registry
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
    req.patient = await patient.byBSN(req.params.patient_id);
    next();
  } catch(e) {
    // Give 500 in all cases so we don't leak information about patients in care
    return console.error(`Can't find patient with bsn ${req.params.patient_id}`) && res.status(500).end();
  }
}

async function handleNutsAuth(req, res, next) {
  // Who is making the request?
  let actor;
  try {
    // TODO: This should in future be determined based on client certificate:
    actor = await registry.organizationById(req.headers.urn);
  } catch(e) {
    // Give 500 in all cases so we don't leak information about patients in care
    return console.error(`Could not fetch organisation ${req.headers.urn}: ${e}`) && res.status(500).end();
  }

  // Is this organisation allowed to make this request?
  try{
    const consent = await consentStore.checkConsent({
      subject:   req.patient,
      actor:     actor,
      custodian: config.organisation
    });

    if ( !consent || consent.consentGiven != 'yes' )
      // Consent was not given for this organisation
      // Give 500 in all cases so we don't leak information about patients in care
      return console.error(`No consent found for actor ${actor}`) && res.status(500);
  } catch(e) {
    // Give 500 in all cases so we don't leak information about patients in care
    return console.error(`Could not fetch consent for triple ${req.patient}, ${actor}, ${config.organisation}: ${e}`) && res.status(500).end();
  }

  // Is this user validly authenticated with IRMA?
  let contract;
  try {
    contract = await auth.validateContract({
      contract_format: 'JWT',
      contract_string: req.headers.user,
      acting_party_cn: 'Demo EHR' // << This is a TODO on the part of the Nuts node
    });

    if ( !contract || contract.validation_result != 'VALID' )
      // IRMA contract is not valid according to Nuts node
      // Give 500 in all cases so we don't leak information about patients in care
      return console.error(`No valid IRMA contract`) && res.status(500);
  } catch(e) {
    // Give 500 in all cases so we don't leak information about patients in care
    return console.error(`Could not validate contract: ${e}`) && res.status(500).end();
  }

  try {
    // Log this request for our own audits
    accessLog.store({
      timestamp: Date.now(),
      patientId: req.patient.id,
      actor:     actor,
      user:      contract.signer_attributes
    });
  } catch(e) {
    // Give 500 in all cases so we don't leak information about patients in care
    return console.error(`Could not save request audit: ${e}`) && res.status(500).end();
  }

  next();
}

module.exports = router;
