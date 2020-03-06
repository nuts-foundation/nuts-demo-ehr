export default {
  render: async () => {
    fetch('/api/authentication/logout').then((response) => {
      if (response.ok && response.status === 204) {
        document.getElementById('logout').innerHTML = template()
      } else {
        throw('Error logging you out' + response)
      }
    })
  }
}

const template = () => `
  <section class='irma-web-center-child' style='height: 80vh; flex-direction: column;'>
    <p style="max-width: 450px; text-align: center;">
      You are logged out.
      <a href="#irma-login">Click here to login again</a>
    </p>
  </section>
`
