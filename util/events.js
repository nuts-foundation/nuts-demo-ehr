module.exports = subscriptionCallback => {

  const events = require('express').Router();
  let clients  = [];
  const lastMessages = {};

  events.get('/:topic', (req, res) => {
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

    // Clean up client when connection closes
    req.on('close', () => clients = clients.filter(c => c !== client));

    // Send client last known message for this topic
    const lastMessage = lastMessages[req.params.topic];
    if ( lastMessage )
      res.write(`data: ${JSON.stringify(lastMessage)}\n\n`);

    // Notify outside world of new client
    subscriptionCallback(client);
  });

  events.publish = (options) => {
    lastMessages[options.topic] = options.message;
    clients.filter(c => c.topic == options.topic)
           .forEach(c => c.res.write(`data: ${JSON.stringify(options.message)}\n\n`));
  };

  events.topics = () => {
    return [...new Set(clients.map(c => c.topic))]; // Unique topics
  }

  return events;

};
