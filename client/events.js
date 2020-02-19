export default {
  subscribe: (o) => {
    const eventSource = new window.EventSource(`${o.path}/events/${o.topic}`);
    eventSource.addEventListener('message', e => { if ( o.message ) o.message(JSON.parse(e.data)) });
    eventSource.addEventListener('error',   e => { if ( o.error ) o.error(e) });
    eventSource.addEventListener('open',    e => { if ( o.open ) o.open(e) });
  }
};
