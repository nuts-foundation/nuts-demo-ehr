import io from '../socketio';
const socket = io.consent();

socket.on('transactions', m => {
  try {
    const transactions = m.map(o => ({
      status: o.name,
      organisations: o.payload.consentRecords.map(r => r.organisations).flat()
    }));

    document.getElementById('transactions').innerHTML = template(transactions);
  } catch(e) {
    // In some states, the payload has no consentRecords.
    // Just silently fail here for now.
    console.info('Failing silently when mapping ', m);
  }
});

export default {
  render: () => {
    socket.emit('get', 'transactions');
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
