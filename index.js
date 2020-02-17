const express    = require('express');
const bodyParser = require('body-parser');
const config     = require('./util/config');
const Logger     = require('./util/logger');
const api        = require('./client-api');
const external   = require('./external-api');
const {crypto}   = require('./resources/nuts-node');

// Run server

Logger.log(`Starting server at port ${config.server.port}`);

const app = express();

app.use('/', (req, res, next) => {
  Logger.log(`Received request for ${req.url}`, req.headers['x-forwarded-for'] || req.connection.remoteAddress);
  next();
});

app.use(express.static('public'));
app.use(express.json());
app.use(bodyParser.json());
app.use('/api', api);
app.use('/external', external);

app.listen(config.server.port, () =>
  Logger.log(`Server is listening on port ${config.server.port}`));

// Register our organisation with the Nuts node on startup

Logger.log(`Registering our organisation ${config.organisation.name}`);

crypto.getPublicKey(config.organisation.agb)
.then(() => Logger.log("Organisation already registered"))
.catch(e => {
  if ( !e.response || e.response.status != 404 )
    return Logger.error("Error registering organisation, is your Nuts node up?",e);

  crypto.generateKeyPair(config.organisation.agb)
  .then(pubKey => Logger.log(`Registered! Public key for copy'n'pastin:\n${JSON.stringify({
    name: config.organisation.name,
    identifier: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`,
    publicKey: pubKey.replace(/\n/g, "\n")
  })}`))
  .catch(e => {
    Logger.error("Error registering organisation, is your Nuts node up?",e);
  });
});
