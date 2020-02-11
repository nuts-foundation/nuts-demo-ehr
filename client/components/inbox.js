export default {
  // Fetch all new consents and render to the inbox
  render: () => {
    return fetch('/api/consent/inbox')
    .then(response => response.json())
    .then(json => {
      if ( json.length > 0 )
        document.getElementById('inbox').innerHTML = template(json);
    });
  }
}

const template = (inbox) => `
  <h2>Inbox</h2>

  <table class="table table-borderless table-bordered table-hover">

    <thead class="thead-dark">
      <tr>
        <th>BSN</th>
        <th>Organisation</th>
        <!-- <th>Actions</th> -->
      </tr>
    </thead>

    <tbody>
      ${inbox.length > 0 ? inbox.map(incoming => `
        <tr>
          <td>${incoming.bsn}</td>
          <td>${incoming.organisation.name}</td>
          <!-- <td><a href="#stuff">Accept</a> / <a href="#stuff">Reject</a></td> -->
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
