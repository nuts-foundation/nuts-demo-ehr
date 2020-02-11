export default {
  // Fetch all new consents and render to the inbox
  render: () => {
    return fetch('/api/consent/transactions')
    .then(response => response.json())
    .then(json => {
      const transactions = json.map(o => ({
        status: o.name,
        organisations: o.payload.consentRecords.map(r => r.organisations).flat()
      }));
      if ( transactions.length > 0 )
        document.getElementById('transactions').innerHTML = template(transactions);
    });
  }
}

const template = (transactions) => `
  <h2>Transactions</h2>

  <table class="table table-borderless table-bordered table-hover">

    <thead class="thead-dark">
      <tr>
        <th>Status</th>
        <th>Organisations involved</th>
      </tr>
    </thead>

    <tbody>
      ${transactions.length > 0 ? transactions.map(transaction => `
        <tr>
          <td>${transaction.status}</td>
          <td>${transaction.organisations.map(o => o.name)}</td>
        </tr>
      `).join('') : `
        <tr>
          <td></td>
          <td></td>
        </tr>
      `}
    </tbody>

  </table>
`;
