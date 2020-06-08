const config = require('../../util/config');
const router = require('express').Router();
const axios = require('axios');
const NATS = require('nats');

router.get('/', async (req, res) => {
  const node = await nodeStatus();
  const nats = await natsStatus();

  const state = {
    demo_ehr_is_up:   true,
    nuts_node_is_up:  node == 'OK',
    nuts_node_status: node,
    nats_is_up:       nats == 'OK',
    nats_status:      nats
  }
  res.status(200).send(state).end();
});

async function nodeStatus() {
  try {
    return await axios.get(`http://${config.nuts.node}/status`)
      .then(response => {
        if ( response.status !== 200 )
          return `ERROR: ${response.data}`;
        return response.data;
      });
  } catch(e) {
    return `Can't reach Nuts node: ${e}`;
  }
}

function natsStatus() {
  return new Promise((resolve, reject) => {
    let nc;
    try {
      nc = NATS.connect(config.nuts.nats);
    } catch(e) {
      return resolve(e);
    }

    nc.on('error', (err) => {
      nc.close();
      resolve(err.message);
    });

    nc.on('connect', () => {
      nc.close();
      resolve('OK');
    });
  });
}

module.exports = router;
