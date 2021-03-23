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

const vdrAPI = require('./resources/nuts-node/vdr')
const vcrAPI = require('./resources/nuts-node/vcr')

// Run server
Logger.log(`Starting server at port ${config.server.port}`)

app.use((req, res, next) => {
    Logger.log(`Received request for ${config.organisation.name}:${config.server.port}${req.url}`, req.headers['x-forwarded-for'] || req.connection.remoteAddress)
    next()
})

async function lookupOrganisationDID(name, city) {
    const results = await vcrAPI.search('organization', [
        {key: 'organization.name', value: name},
        {key: 'organization.city', value: city},
    ])
    return results.length > 0 ? results[0].subject : null
}

async function createOrganisation(name, city) {
    Logger.log(`Setting up organisation: ${name} at ${city}`)
    // Create DID
    Logger.log(`   Creating DID`)
    let didDocument = await vdrAPI.create()
    const orgDID = didDocument.id
    let resolved = await vdrAPI.resolve(orgDID)
    Logger.log(`   Registering assertionMethod key`)
    didDocument.assertionMethod = [didDocument.authentication[0]];
    await vdrAPI.update(didDocument, resolved.documentMetadata.hash)
    if (config.server.verbose === true) {
        Logger.log(`Updated DID document: ${JSON.stringify(didDocument, null, 2)}`)
    }
    // Now issue NutsOrganizationCredential VC
    Logger.log(`   Marking self-issued Verifiable Credentials as trusted`)
    await vcrAPI.trust('NutsOrganizationCredential', orgDID);
    Logger.log(`   Issuing Verifiable Credential`)
    let vc = await vcrAPI.create('NutsOrganizationCredential', orgDID, {
        id: orgDID,
        organization: {
            name: name,
            city: city,
        }
    })
    if (config.server.verbose === true) {
        Logger.log(`Issued Verifiable Credential: ${JSON.stringify(vc, null, 2)}`)
    }
    return orgDID
}

async function configureOrganisation(name, city) {
    let orgDID = await lookupOrganisationDID(name, city);
    if (!orgDID) {
        orgDID = await createOrganisation(name, city)
    } else {
        Logger.log(`Organization already exists as DID: ${orgDID}`)
    }
    return {
        name: name,
        did: orgDID,
    };
}

function startServer(organisation) {
    // Use Redis for session persistence.
    // Remove the following 2 lines and the 'store' property if you prefer an in-memory store.
    const RedisStore = require('connect-redis')(session)
    const redisServer = process.env.REDIS_SERVER_ADDRESS || 'localhost'
    const redisClient = redis.createClient({host: redisServer})
    redisClient.on('error', Logger.error)
    app.use(session({
        store: new RedisStore({client: redisClient, prefix: `session-${organisation.did}:`}),
        secret: config.server.sessionSecret,
        resave: false,
        name: `session-${organisation.did}`,
        saveUninitialized: false
    }))

    app.use(express.static('public'))
    app.use(express.json())
    app.use(bodyParser.json())

    app.use('/api', clientAPI(organisation))
    app.use('/external', externalAPI)
    app.use('/sso', ssoHandler)
    app.use('/', stateAPI)

    eventAPI(io) // Mount events API using socket.io

    server.listen(config.server.port, () =>
        Logger.log(`Server is listening on port ${config.server.port}`))
}

configureOrganisation(config.organisation.name, config.organisation.city)
    .then((organization) => startServer(organization))
