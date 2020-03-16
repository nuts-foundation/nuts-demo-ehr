const config = require('../util/config')
const NATS = require('nats')
const nc = NATS.connect(config.nuts.nats)
nc.on('error', (err) => {
  console.log(err)
})

const {
  patient,
  consentInTransit
} = require('../resources/database')

const {
  consentStore,
  eventStore,
  registry
} = require('../resources/nuts-node')

module.exports = async io => {
  let inbox = await getInbox()
  let transactions = await getTransactions()
  let watchedPatients = []

  io.on('connection', socket => {
    socket.on('disconnect', () => {
      // Unsubscribe from any patients
      watchedPatients = watchedPatients.filter(p => p.socket !== socket)
    })

    // Get me the latest of these, now
    socket.on('get', type => {
      switch (type) {
        case 'inbox':
          return socket.emit('inbox', inbox)
        case 'transactions':
          return socket.emit('transactions', transactions)
      }
    })

    // Subscribe to events concerning this patient
    let myWatchedPatient = null
    socket.on('subscribe', async patientId => {
      // Unsubscribe from previous patient
      watchedPatients = watchedPatients.filter(p =>
        p.patientId != myWatchedPatient ||
        p.socket !== socket
      )

      // Subscribe to new patient
      myWatchedPatient = patientId
      watchedPatients.push({ socket, patientId })

      // Send current status of this patient
      const myPatient = await patient.byId(patientId)
      socket.emit('givenConsents', await getGivenConsents(myPatient))
      socket.emit('receivedConsents', await getReceivedConsents(myPatient))
    })
  })

  // Subscribe to NATS events
  nc.subscribe('*.*.*.consentRequest', async (msg, reply, subject) => {
    // Only parse the bit that smells like JSON 😉
    const json = JSON.parse(msg.match('\{.*\}').pop())
    json.payload = JSON.parse(Buffer.from(json.payload, 'base64').toString())

    console.log(`Received consentRequest event on '${subject}':`, json)

    // For now just republish everything we've got 😅

    inbox = await getInbox()
    io.emit('inbox', inbox)

    transactions = await getTransactions()
    io.emit('transactions', transactions)

    for (const p of watchedPatients) {
      const patientObj = await patient.byId(p.patientId)
      p.socket.emit('givenConsents', await getGivenConsents(patientObj))
      p.socket.emit('receivedConsents', await getReceivedConsents(patientObj))
    }
  })
}

async function getInbox () {
  try {
    // Get all consents I have been granted
    const consents = await consentStore.consentsFor({
      actor: config.organisation
    })

    // If no consent found
    if (!consents || consents.totalResults === 0) { return [] }

    // Add those that have unknown patients to the inbox
    const inbox = []
    for (const consent of consents.results) {
      const bsn = consent.subject.split(':').pop()
      if (await patient.byBSN(bsn) === undefined) {
        inbox.push({
          bsn: bsn,
          organisation: await registry.organizationById(consent.custodian)
        })
      }
    }

    return inbox
  } catch (e) {
    console.error(`Error in Nuts node query for consent events: ${e}`)
  }
}

async function getTransactions () {
  try {
    const evnts = await eventStore.allEvents()

    const transactions = []
    // Map URNs in events to sane organisations
    for (const event of evnts.events || []) {
      // Event can have multiple consent records
      const organisations = []
      if (event.payload.legalEntity) { organisations.push(await registry.organizationById(event.payload.legalEntity)) }
      if (event.payload.legalEntities) {
        for (const org of event.payload.legalEntities) { organisations.push(await registry.organizationById(org)) }
      }
      if (event.payload.consentRecords) {
        for (const record of event.payload.consentRecords) {
        // Consent record has multiple organisations
          for (const org of record.metadata.organisationSecureKeys) {
            organisations.push(await registry.organizationById(org.legalEntity))
          }
          record.organisations = organisations
        }
      }

      transactions.push({
        status: event.name,

        // Only unique organisations
        organisations: organisations.filter((obj, pos, arr) =>
          arr.map(mapObj =>
            mapObj.name
          ).indexOf(obj.name) === pos
        )
      })
    }

    return transactions
  } catch (e) {
    console.error(`Error in Nuts node query for consent events: ${e}`)
  }
}

async function getReceivedConsents (patient) {
  try {
    // Fetch received consents for this patient
    const consents = await consentStore.consentsFor({
      subject: patient,
      actor: config.organisation
    })

    // Map URNs to sane organisations
    let organisations = []
    if (consents.totalResults > 0 && consents.results) { organisations = await urnsToOrgs(consents.results.map(o => o.custodian)) }

    return organisations
  } catch (e) {
    console.error(`Error in Nuts node query for finding consents: ${e}`)
  }
}

async function getGivenConsents (patient) {
  try {
    // Fetch confirmed consents
    const consents = await consentStore.consentsFor({
      subject: patient,
      custodian: config.organisation
    })

    // Fetch consents in transit
    const pending = await consentInTransit.get({
      subject: patient.bsn,
      custodian: config.organisation.agb
    })

    // Are these consents in transit still pending?
    const stillPending = []
    for (const consent of pending) {
      const currentStatus = await eventStore.getEvent(consent.jobId)
      if (currentStatus.name !== 'completed') { stillPending.push(consent) }
    }

    // Map URNs to sane organisations
    let organisations = []
    if (consents.totalResults > 0 && consents.results) { organisations = await urnsToOrgs(consents.results.map(o => o.actor)) }
    if (stillPending.length > 0) {
      const orgs = await urnsToOrgs(stillPending.map(o => o.actor))
      organisations = organisations.concat(orgs.map(o => { o.name += ' <em>(pending acceptance)</em>'; return o }))
    }

    return organisations
  } catch (e) {
    console.error(`Error in Nuts node query for finding consents: ${e}`)
  }
};

async function urnsToOrgs (urns) {
  const orgs = []
  for (const urn of urns) {
    try {
      orgs.push(await registry.organizationById(urn))
    } catch (e) {
      orgs.push({
        identifier: urn,
        name: 'Could not find organisation: ' + e
      })
    }
  }
  return orgs
}
