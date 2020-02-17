import call from '../component-loader';
let interval = false;

export default {
  render: () => {
    const element = document.getElementById('transactions');
    if ( !interval )
       interval = window.setInterval(() => update(element), 3000);
    update(element);
    return Promise.resolve();
  }
}

function update(element) {
  call('/api/consent/transactions', element)
  .then(json => {
    const transactions = json.map(o => ({
      status: o.name,
      organisations: o.payload.consentRecords.map(r => r.organisations).flat()
    }));
    document.getElementById('transactions').innerHTML = template(transactions);
  })
  .catch(error => {
    element.innerHTML = `<h2>Transactions</h2><p>Could not load transactions: ${error}</p>`;
  });
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
