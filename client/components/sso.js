export default {
  render: async (patient, orgId) => {
    document.getElementById('sso').innerHTML = template()
  }
}
const template = () => `
  <section class='irma-web-center-child' style='height: 80vh; flex-direction: column;'>
    <h1>SSO</h1>
  </section>
`
