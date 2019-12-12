const config        = require('../../util/config');
const openAPIHelper = require('./open-api-helper');

const call = openAPIHelper.call({
  baseURL:    `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-consent-store/raw/master/docs/_static/nuts-consent-store.yaml'
});

module.exports = {
  consentsFor: async (parties) => {
    const res = await call('queryConsent', null, query(parties));
    console.log(res);
    return res;
  }
};

function query(parties) {
  return {
    custodian: openAPIHelper.urn(parties.custodian),
    actor:     openAPIHelper.urn(parties.actor),
    subject:   openAPIHelper.urn(parties.subject),
    page: {
      offset: 0,
      limit: 0
    },
    validAt: null
  };
}
