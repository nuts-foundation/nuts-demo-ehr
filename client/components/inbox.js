export default {
  // Fetch all new consents and render to the inbox
  render: () => {
    return fetch('/api/consent/inbox')
    .then(response => response.json())
    .then(json => {
      const transactions = json.map(o => ({
        status: o.name,
        organisations: o.payload.consentRecords.map(r => r.organisations).flat()
      }));
      console.log(transactions);
      document.getElementById('inbox').innerHTML = template(transactions);
    });
  }
}

const template = (transactions) => `
  <h2>Inbox</h2>

  <table class="table table-borderless table-bordered table-hover">

    <thead class="thead-dark">
      <tr>
        <th>Status</th>
        <th>BSN</th>
        <th>Organisations involved</th>
        <!-- <th>Actions</th> -->
      </tr>
    </thead>

    <tbody>
      ${transactions.length > 0 ? transactions.map(transaction => `
        <tr>
          <td>${transaction.status}</td>
          <td>Unknown</td>
          <td>${transaction.organisations.map(o => o.name)}</td>
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
