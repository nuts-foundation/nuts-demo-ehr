const router  = require('express').Router();
const config  = require('../util/config');

const {
  patient,
  consentInTransit
} = require('../resources/database');

const {
  consentLogic,
  consentStore,
  eventStore,
  registry
} = require('../resources/nuts-node');

router.get('/:patient_id/received', findPatient, async (req, res) => {
  try {
    // Fetch received consents
    const consents = await consentStore.consentsFor({
      subject: req.patient,
      actor:   config.organisation
    });

    // Map URNs to sane organisations
    let organisations = [];
    if ( consents.totalResults > 0 && consents.results )
      organisations = await urnsToOrgs(consents.results.map(o => o.custodian));

    res.status(200).send(organisations).end();
  } catch(e) {
    res.status(500).send(`Error in Nuts node query for finding consents: ${e}`);
  }
});

router.get('/:patient_id/given', findPatient, async (req, res) => {
  try {
    // Fetch confirmed consents
    const consents = await consentStore.consentsFor({
      subject:   req.patient,
      custodian: config.organisation
    });

    // Fetch consents in transit
    const pending = await consentInTransit.get({
      subject:   req.patient.bsn,
      custodian: config.organisation.agb
    });

    // Are these consents in transit still pending?
    const stillPending = [];
    for ( consent of pending ) {
      const currentStatus = await eventStore.getEvent(consent.jobId);
      if (currentStatus.name != 'completed')
        stillPending.push(consent);
    }

    // Map URNs to sane organisations
    let organisations = [];
    if ( consents.totalResults > 0 && consents.results )
      organisations = await urnsToOrgs(consents.results.map(o => o.actor));
    if ( stillPending.length > 0 ) {
      const orgs = await urnsToOrgs(stillPending.map(o => o.actor))
      organisations = organisations.concat(orgs.map(o => { o.name += ' <em>(pending acceptance)</em>'; return o }));
    }

    res.status(200).send(organisations).end();
  } catch(e) {
    res.status(500).send(`Error in Nuts node query for finding consents: ${e}`);
  }
});

router.get('/inbox', async (req, res) => {
  try {
    // Get all consents I have been granted
    const consents = await consentStore.consentsFor({
      actor: config.organisation
    });

    // Add those that have unknown patients to the inbox
    const inbox = [];
    for ( let consent of consents.results ) {
      const bsn = consent.subject.split(':').pop();
      if ( await patient.byBSN(bsn) === undefined )
        inbox.push({
          bsn:          bsn,
          organisation: await registry.organizationById(consent.custodian)
        });
    }

    res.status(200).send(inbox).end();
  } catch(e) {
    res.status(500).send(`Error in Nuts node query for consent events: ${e}`);
  }
});

router.get('/transactions', async (req, res) => {
  try {
    const events = await eventStore.allEvents();

    // Map URNs in events to sane organisations
    for ( let event of events.events || [] ) {
      // Event can have multiple consent records
      for ( let record of event.payload.consentRecords ) {
        const organisations = [];
        // Consent record has multiple organisations
        for ( let org of record.metadata.organisationSecureKeys ) {
          organisations.push(await registry.organizationById(org.legalEntity));
        }
        record.organisations = organisations;
      }
    }

    res.status(200).send(events.events || []).end();
  } catch(e) {
    res.status(500).send(`Error in Nuts node query for consent events: ${e}`);
  }
});

router.put('/:patient_id', findPatient, async (req, res) => {
  try {
    // Create the "real" consent at the Nuts node
    const result = await consentLogic.createConsent({
      subject: req.patient,
      actor: {
        urn: req.body.organisationURN
      },
      custodian: config.organisation
    }, req.body.reason);

    if ( result.resultCode !== 'OK' )
      return res.status(500).send(`Error in Nuts node query for storing a new consent, result is not OK: ${result}`);

    // Create a "consent in transit" record locally to track its status
    await consentInTransit.store({
      jobId:      result.jobId,
      subject:    req.patient.bsn,
      actor:      req.body.organisationURN,
      custodian:  config.organisation.agb,
      proofTitle: req.body.reason
    });

    res.status(201).send(result).end();
  } catch(e) {
    res.status(500).send(`Error in Nuts node query for storing a new consent: ${e}`);
  }
});

async function urnsToOrgs(urns) {
  const orgs = [];
  for ( urn of urns ) {
    try {
      orgs.push(await registry.organizationById(urn));
    } catch (e) {
      orgs.push({
        identifier: urn,
        name: 'Could not find organisation: ' + e
      });
    }
  }
  return orgs;
}

async function findPatient(req, res, next) {
  try {
    req.patient = await patient.byId(req.params.patient_id);
    next();
  } catch(e) {
    res.status(404).send(`Could not find a patient with id ${req.params.patient_id}: ${e}`);
  }
}

module.exports = router;
