const IrmaCore = require('irma-core');
const Dummy    = require('irma-dummy');
const Web      = require('irma-web');


export default {
  render: () => {
    console.log("irma login");
    document.getElementById('irma-login').innerHTML = template;

    const irma = new IrmaCore({
      debugging: true,
      dummy: 'happy path',
      element: '#irma',
      timing: {
        start: 1000,
        scan: 6000,
        app: 2000
      }
    });

    irma.use(Dummy);
    irma.use(Web);

    irma.start('localhost:21323', {request: 'content'}).then(() => {
        window.irmaLogin = true;
        window.history.back();
    })

  }
}

const template = `
<div id="irma"></div>
`
