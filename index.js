const express    = require('express');
const app        = express();
const server     = require('http').Server(app);
const io         = require('socket.io')(server);
const session    = require('express-session');
const bodyParser = require('body-parser');

const config      = require('./util/config');
const Logger      = require('./util/logger');
const clientAPI   = require('./client-api');
const externalAPI = require('./external-api');
const eventAPI    = require('./event-api');
const {crypto}    = require('./resources/nuts-node');

// Run server

Logger.log(`Starting server at port ${config.server.port}`);

app.use('/', (req, res, next) => {
  Logger.log(`Received request for ${req.url}`, req.headers['x-forwarded-for'] || req.connection.remoteAddress);
  next();
});

app.use(session({
  secret:            config.server.sessionSecret,
  resave:            false,
  saveUninitialized: false
}));

app.use(express.static('public'));
app.use(express.json());
app.use(bodyParser.json());

app.use('/api',      clientAPI);
app.use('/external', externalAPI);

eventAPI(io);  // Mount events API using socket.io

server.listen(config.server.port, () =>
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
