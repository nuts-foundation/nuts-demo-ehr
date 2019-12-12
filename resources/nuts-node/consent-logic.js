const uuid          = require('uuid/v4');
const config        = require('../../util/config');
const openAPIHelper = require('./open-api-helper');

const call = openAPIHelper.call({
  baseURL:    `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-consent-logic/raw/master/docs/_static/nuts-consent-logic.yaml'
});

module.exports = {
  createConsent: async (parties, proofTitle) => {
    const res = await call('createOrUpdateConsent', null, consentRecord(parties, proofTitle));
    console.log(res);
    return res;
  },

  // This is "almost right". Needs to reference the previous proof though
  deleteConsent: async (parties) => {
    const record = consentRecord(parties);
    record.records.period.end = new Date();
    return await call('createOrUpdateConsent', null, record);
  }
};

function consentRecord(parties, proofTitle) {
  return {
    custodian: openAPIHelper.urn(parties.custodian),
    actor:     openAPIHelper.urn(parties.actor),
    subject:   openAPIHelper.urn(parties.subject),
    records: [
      {
        period: {
          start: new Date()
        },
        consentProof: {
          ID: uuid(),
          title: proofTitle
        },
        dataClass: [
          'urn:oid:1.3.6.1.4.1.54851.1:MEDICAL'
        ]
      }
    ]
  };
}
