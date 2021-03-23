const config = require('../../util/config')
const apiHelper = require('./open-api-helper')

const call = apiHelper.call({
    baseURL: `http://${config.nuts.node}`,
    definition: 'https://raw.githubusercontent.com/nuts-foundation/nuts-node/master/docs/_static/vcr/v1.yaml'
});

const vcr = {
    create: async (type, issuer, subject) => call('create', null, {
        type: [type],
        issuer: issuer,
        credentialSubject: subject
    }),
    trust: async (type, issuer) => call('trustIssuer', null, {issuer: issuer, credentialType: type}),
    search: async (type, params) => call('search', type, {params: params}),
};

module.exports = vcr