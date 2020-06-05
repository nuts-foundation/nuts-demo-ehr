const express = require('express')
const app = express()
const server = require('http').Server(app)
const io = require('socket.io')(server)
const redis = require('redis')
const session = require('express-session')
const bodyParser = require('body-parser')

const config = require('./util/config')
const Logger = require('./util/logger')
const clientAPI = require('./api/browser')
const externalAPI = require('./api/external')
const ssoHandler = require('./api/sso')
const eventAPI = require('./api/events')
const stateAPI = require('./api/state')
const { crypto } = require('./resources/nuts-node')

// Run server
Logger.log(`Starting server at port ${config.server.port}`)

app.use((req, res, next) => {
  Logger.log(`Received request for ${config.organisation.name}:${config.server.port}${req.url}`, req.headers['x-forwarded-for'] || req.connection.remoteAddress)
  next()
})

// Use Redis for session persistence.
// Remove the following 2 lines and the 'store' property if you prefer an in-memory store.
const RedisStore = require('connect-redis')(session)
const redisServer = process.env.REDIS_SERVER_ADDRESS || 'localhost'
const redisClient = redis.createClient({ host: redisServer })
redisClient.on('error', console.error)
app.use(session({
  store: new RedisStore({ client: redisClient, prefix: `session-${config.organisation.agb}:` }),
  secret: config.server.sessionSecret,
  resave: false,
  name: `session-${config.organisation.agb}`,
  saveUninitialized: false
}))

app.use(express.static('public'))
app.use(express.json())
app.use(bodyParser.json())

app.use('/api', clientAPI)
app.use('/external', externalAPI)
app.use('/sso', ssoHandler)
app.use('/', stateAPI)

eventAPI(io) // Mount events API using socket.io

server.listen(config.server.port, () =>
  Logger.log(`Server is listening on port ${config.server.port}`))

// Register our organisation with the Nuts node on startup

Logger.log(`Registering our organisation ${config.organisation.name}`)

crypto.getPublicKey(config.organisation.agb)
  .then(() => Logger.log('Organisation already registered'))
  .catch(e => {
    if (!e.response || e.response.status !== 404) { return Logger.error('Error registering organisation, is your Nuts node up?', e) }

    crypto.generateKeyPair(config.organisation.agb)
      .then(pubKey => {
        const exampleVendorClaimEvent = {
          type: 'VendorClaimEvent',
          payload: {
            vendorIdentifier: '',
            orgName: config.organisation.name,
            orgIdentifier: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`,
            orgKeys: [pubKey]
          }
        }

        Logger.log(`Registered! Public key for copy'n'pastin:\n${JSON.stringify(exampleVendorClaimEvent)}`)
      })
      .catch(e => {
        Logger.error('Error registering organisation, is your Nuts node up?', e)
      })
  })
