const router = require('express').Router();
let clients  = [];
const lastMessages = {};

router.get('/:topic', (req, res) => {
  // Headers to keep connection open
  res.writeHead(200, {
    'Content-Type':  'text/event-stream',
    'Connection':    'keep-alive',
    'Cache-Control': 'no-cache'
  });

  // Add client and requested topic to the client list
  const client = {
    topic: req.params.topic,
    res
  };
  clients.push(client);

  // Send client last known message for this topic
  const lastMessage = lastMessages[req.params.topic];
  if ( lastMessage )
    res.write(`data: ${JSON.stringify(lastMessage)}\n\n`);

  // Clean up client when connection closes
  req.on('close', () => clients = clients.filter(c => c !== client));
});

router.publish = (options) => {
  lastMessages[options.topic] = options.message;
  clients.filter(c => c.topic == options.topic)
         .forEach(c => c.res.write(`data: ${JSON.stringify(options.message)}\n\n`));
};

module.exports = router;
