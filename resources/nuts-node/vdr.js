const config = require('../../util/config')
const apiHelper = require('./open-api-helper')

const call = apiHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: 'https://raw.githubusercontent.com/nuts-foundation/nuts-node/master/docs/_static/vdr/v1.yaml'
});

const vdr = {
  create: async () => call('createDID', null, null),
  resolve: async (did) => call('getDID', did, null),
  update: async(didDocument, currentHash) => call('updateDID', didDocument.id, {document: didDocument, currentHash: currentHash}),
}

module.exports = vdr