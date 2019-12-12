export default {
  // Fetch all new consents and render to the inbox
  render: () => {
    return fetch('/api/consent/inbox')
    .then(response => response.json())
    .then(json => {
      document.getElementById('inbox').innerHTML = template(json);
    });
  }
}

const template = (events) => `
  <h2>Inbox</h2>

  <table class="table table-borderless table-bordered table-hover">

    <thead class="thead-dark">
      <tr>
        <th>BSN</th>
        <th>External organisation</th>
        <th>Actions</th>
      </tr>
    </thead>

    <tbody>
      ${events.length > 0 ? events.map(evnt => `
        <tr>
          <td>Unknown</td>
          <td>${evnt.organisation.name}</td>
          <td><a href="#stuff">Accept</a> / <a href="#stuff">Reject</a></td>
        </tr>
      `).join('') : `
        <tr>
          <td></td>
          <td>None</td>
          <td></td>
        </tr>
      `}
    </tbody>

  </table>
`;
