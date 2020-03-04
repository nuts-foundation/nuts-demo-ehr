const IrmaCore = require('irma-core')
const Server = require('irma-server')
const Web = require('irma-web')

export default {
  render: async () => {
    document.getElementById('irma-login').innerHTML = template()

    const irma = new IrmaCore({
      debugging: true,
      element: '#irma-web-form',

      language: 'en',
      translations: {
        header: 'Identify yourself with <i class="irma-web-logo">IRMA</i>',
        cancelled: 'We have not received the signed contract. We\'re sorry, but because of this we can\'t identify you and you can\'t request data'
      },

      session: {
        url: '/api/authentication',

        start: {
          url: o => `${o.url}/new-session`,
          method: 'GET',
          qrFromResult: r => r
        },
        result: {
          url: o => `${o.url}/session-done`,
          method: 'GET'
        }
      },

      state: {
        serverSentEvents: false
      }
    })

    irma.use(Server)
    irma.use(Web)

    try {
      const result = await irma.start()
      window.setTimeout(() => {
        let callbackUrl = window.localStorage.getItem('afterLoginReturnUrl')
        window.localStorage.removeItem('afterLoginReturnUrl')
        if (!callbackUrl) {
          callbackUrl = '#dashboard'
        }
        window.location = callbackUrl
      }, 1200)
    } catch (e) {
      console.error('Trouble running IRMA flow: ', e)
    }
  }
}

const template = () => `
  <section class='irma-web-center-child' style='height: 80vh; flex-direction: column;'>
    <p style="max-width: 450px; text-align: center;">
      You can now login to the DEMO EHR using IRMA.
    </p>
    <section id='irma-web-form' style="margin: 2em 0;"></section>
  </section>
`
