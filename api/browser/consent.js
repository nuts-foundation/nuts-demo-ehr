const router = require('express').Router()
const config = require('../../util/config')
const { consentLogic } = require('../../resources/nuts-node')

const {
  patient,
  consentInTransit
} = require('../../resources/database')

router.put('/:patient_id', findPatient, async (req, res) => {
  try {
    // Create the "real" consent at the Nuts node
    const result = await consentLogic.createConsent({
      subject: req.patient,
      actor: {
        urn: req.body.organisationURN
      },
      custodian: config.organisation
    }, req.body.reason)

    if (result.resultCode !== 'OK') { return res.status(500).send(`Error in Nuts node query for storing a new consent, result is not OK: ${result}`) }

    // Create a "consent in transit" record locally to track its status
    await consentInTransit.store({
      jobId: result.jobId,
      subject: req.patient.bsn,
      actor: req.body.organisationURN,
      custodian: config.organisation.agb,
      proofTitle: req.body.reason
    })

    res.status(201).send(result).end()
  } catch (e) {
    res.status(500).send(`Error in Nuts node query for storing a new consent: ${e}`)
  }
})

async function findPatient (req, res, next) {
  try {
    req.patient = await patient.byId(req.params.patient_id)
    next()
  } catch (e) {
    res.status(404).send(`Could not find a patient with id ${req.params.patient_id}: ${e}`)
  }
}

module.exports = () => router
