const IrmaCore = require('irma-core');
const Server   = require('irma-server');
const Web      = require('irma-web');

export default {
  render: async () => {
    document.getElementById('irma-login').innerHTML = template();

    const irma = new IrmaCore({
      element: '#irma-web-form',
      debugging: true,

      session: {
        url: '/api/authentication',

        start: {
          url: o => `${o.url}/new-session`,
          method: 'GET',
          qrFromResult: r => r.qr_code_info
        },
        result: false
      },

      state: {
        serverSentEvents: false
      }
    });

    irma.use(Server);
    irma.use(Web);

    try {
      const result = await irma.start();
      window.irmaLogin = true;
      window.setTimeout(() => window.history.back(), 1200);
    } catch (e) {
      console.error(`Trouble running IRMA flow: `, e);
    }
  }
}

const template = () => `
  <section class='irma-web-center-child' style='height: 80vh;'>
    <section id='irma-web-form'></section>
  </section>
`;
