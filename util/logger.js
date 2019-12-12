const config = require('../util/config');

module.exports = {

  log: (message, sender=false) => {
    if ( config.server.verbose )
      console.log(`[${new Date()}]${sender ? '['+sender+']' : ''} ${message}`);
  },

  error: (message, sender=false) => {
    if ( config.server.verbose )
      console.error(`[${new Date()}]${sender ? '['+sender+']' : ''} ${message}`);
  }

}
