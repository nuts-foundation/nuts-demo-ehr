import irmaFlow from './irma-flow'

export default {
  render: async (subroute, evnt) => {
    document.getElementById('escalate').innerHTML = template()
    irmaFlow.render(
      document.getElementById('escalate-irma'),
      {
        header:    'Identify yourself with <i class="irma-web-logo">IRMA</i>',
        cancelled: 'We have not received the signed contract. We\'re sorry, but because of this we can\'t identify you and you can\'t view data'
      }
    )
  }
}

const template = () => `
  <section class='irma-web-center-child' style='height: 80vh; flex-direction: column;'>
    <p style="max-width: 450px; text-align: center;">
      You are about to view data from an <b>external organisation</b>.
      You will need to identify yourself for this using IRMA.
    </p>
    <section id='escalate-irma' style="width: 100%"></section>
    <p><a href="javascript:window.history.go(-2);">&laquo; Back</a></p>
  </section>
`
