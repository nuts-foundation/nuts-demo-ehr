// Module that takes care of server-internal events

const EventEmitter = require('events')
class AccessLogEmitter extends EventEmitter {}

module.exports = {
  accessLog: new AccessLogEmitter()
}
