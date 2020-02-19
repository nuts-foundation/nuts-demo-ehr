const router  = require('express').Router();
const config  = require('../util/config');
const events  = require('../util/events');
const NATS    = require('nats')
const nc      = NATS.connect(config.nuts.nats);

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


/** Server-Sent Events API **/

// Mount events endpoint
router.use('/events', events);

// Publish initial values for these topics:
publishInbox();
publishTransactions();

// Subscribe to NATS events and keep client up to date
nc.subscribe('*.*.*.consentRequest', (msg, reply, subject) => {
  // Only parse the bit that smells like JSON 😉
  const json = JSON.parse(msg.match('\{.*\}').pop());
  json.payload = JSON.parse(Buffer.from(json.payload, 'base64').toString());

  console.log(`Received consentRequest event on '${subject}':`, json);

  // For now just republish everything we've got 😅
  publishInbox();
  publishTransactions();
});

async function publishInbox() {
  try {
    // Get all consents I have been granted
    const consents = await consentStore.consentsFor({
      actor: config.organisation
    });

    // If no consent found
    if ( !consents || consents.totalResults === 0 )
      return events.publish({
        topic: 'inbox',
        message: '[]'
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

    events.publish({
      topic: 'inbox',
      message: JSON.stringify(inbox)
    });
  } catch(e) {
    console.errror(`Error in Nuts node query for consent events: ${e}`);
  }
}

async function publishTransactions() {
  try {
    const evnts = await eventStore.allEvents();

    // Map URNs in events to sane organisations
    for ( let event of evnts.events || [] ) {
      // Event can have multiple consent records
      if ( event.payload.consentRecords )
        for ( let record of event.payload.consentRecords ) {
          const organisations = [];
          // Consent record has multiple organisations
          for ( let org of record.metadata.organisationSecureKeys ) {
            organisations.push(await registry.organizationById(org.legalEntity));
          }
          record.organisations = organisations;
        }
    }

    events.publish({
      topic: 'transactions',
      message: JSON.stringify(evnts.events) || '[]'
    });
  } catch(e) {
    console.error(`Error in Nuts node query for consent events: ${e}`);
  }
}


/** Normal API requests **/

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
