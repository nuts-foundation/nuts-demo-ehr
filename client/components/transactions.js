import events from '../events';

export default {
  render: () => {
    events.subscribe({
      path: '/api/consent',
      topic: 'transactions',
      message: m => {
        const json = JSON.parse(m)
                         .map(o => ({
                           status: o.name,
                           organisations: o.payload.consentRecords.map(r => r.organisations).flat()
                         }))

        document.getElementById('transactions').innerHTML = template(json);
      }
    });

    return Promise.resolve();
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
          <td colspan="2" style="text-align: center"><em>None</em></td>
        </tr>
      `}
    </tbody>

  </table>
`;
