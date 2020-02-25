import io from '../socketio';
const socket = io.consent();

socket.on('inbox', m => {
  const json = m.sort((a,b) => a.bsn.localeCompare(b.bsn))
  document.getElementById('inbox').innerHTML = template(json);
});

export default {
  render: () => {
    socket.emit('get', 'inbox');
    return Promise.resolve();
  }
}

const template = (inbox) => `
  <h2>Inbox</h2>

  <p><em>Patients not (yet) in care that you have been given permission for.</em></p>

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
          <td colspan="2" style="text-align: center"><em>None</em></td>
        </tr>
      `}
    </tbody>

  </table>
`;
