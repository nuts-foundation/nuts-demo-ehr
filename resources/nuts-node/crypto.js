const config = require('../../util/config');
const call   = require('./open-api-helper').call({
  baseURL:    `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-crypto/raw/master/docs/_static/nuts-service-crypto.yaml'
});

module.exports = {
  generateKeyPair:  async (agb) => await call('generateKeyPair', `urn:oid:2.16.840.1.113883.2.4.6.1:${agb}`),
  getPublicKey:     async (agb) => await call('publicKey',       `urn:oid:2.16.840.1.113883.2.4.6.1:${agb}`)
};
