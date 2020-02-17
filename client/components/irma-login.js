const IrmaCore = require('irma-core');
const Server   = require('irma-server');
const Web      = require('irma-web');

export default {
  render: () => {
    const irma = new IrmaCore({
      element: '#irma-login',
      debugging: true,

      session: {
        url: '/api/authentication',

        start: {
          url: o => `${o.url}/new-session`,
          method: 'GET',
          qrFromResult: r => r.qr_code_info
        },
        result: false
      }
    });

    irma.use(Server);
    irma.use(Web);

    try {
      const result = irma.start();
      window.irmaLogin = true;
      window.history.back();
    } catch (e) {
      console.error(`Trouble running IRMA flow: `, e);
    }
  }
}
