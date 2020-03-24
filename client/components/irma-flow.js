const IrmaCore = require('@nuts-foundation/irma-core')
const Server = require('@nuts-foundation/irma-server')
const Web = require('@nuts-foundation/irma-web')

export default {
  render: async (element, translations) => {
    element.innerHTML = template()

    const irma = new IrmaCore({
      debugging: true,
      element: '#irma-web-form',

      language: 'en',
      translations,

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
      await irma.start()
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
  <section class='irma-web-center-child'>
    <section id='irma-web-form' style="margin: 2em 0;"></section>
  </section>
`
