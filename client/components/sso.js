export default {
  render: async (patient, orgId) => {
    return fetch(`/api/jump?patient=${patient}&custodian=${orgId}`)
      .then(response => response.json())
      .then(json => { })
    document.getElementById('sso').innerHTML = template()
  }
}
const template = () => `
  <section class='irma-web-center-child' style='height: 80vh; flex-direction: column;'>
    <h1>SSO</h1>
  </section>
`
