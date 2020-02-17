const config = require('../../util/config');
const call   = require('./open-api-helper').call({
  baseURL:    `http://${config.nuts.node}`,
  definition: 'https://raw.githubusercontent.com/nuts-foundation/nuts-auth/master/docs/_static/nuts-auth.yaml'
});

module.exports = {
  createLoginSession:   async ()         => await call('createSession', null, loginContract()),
  createSession:        async (contract) => await call('createSession', null, contract),
  sessionRequestStatus: async (id)       => await call('sessionRequestStatus', id),
  validateContract:     async (contract) => await call('validateContract', null, contract)
};

function loginContract() {
  return {
    type:        "BehandelaarLogin",
    language:    "NL",
    version:     "v1",
    legalEntity: config.organisation.name
  }
}
