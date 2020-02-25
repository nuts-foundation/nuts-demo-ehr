import observations from './observations';

export default {
  render: (patient, organisation) => {
    return fetch(`/api/organisation/byURN/${organisation}`)
    .then(result => result.json())
    .then(organisation => {
      // Add tab and pane if they aren't in the dom yet
      if ( !document.getElementById(`org-${hash(organisation.identifier)}`) ) {
        addTab(organisation.name, `org-${hash(organisation.identifier)}`);
      }

      // Find our organisation's tab and pane
      let tab  = document.querySelector(`a[data-open='#org-${hash(organisation.identifier)}']`);
      let pane = document.getElementById(`org-${hash(organisation.identifier)}`);

      // Fill it with observations
      observations.render(pane, patient, organisation);

      // And open it
      tab.click();
      return;
    });
  }
}

// Hash function taken from https://stackoverflow.com/questions/6122571/simple-non-secure-hash-function-for-javascript
function hash(string) {
  if (string.length == 0) return 0;

  let hash = 0;
  for (var i = 0; i < string.length; i++) {
    const char = string.charCodeAt(i);
    hash = ((hash<<5)-hash)+char;
    hash = hash & hash; // Convert to 32bit integer
  }
  return hash;
}

function addTab(label, id) {
  const patient = document.getElementById('patient');

  const tab = document.createElement('li');
  tab.classList.add('nav-item');
  tab.innerHTML = `<a class="nav-link" data-open="#${id}">${label}</a>`;
  patient.querySelector('ul.nav').appendChild(tab);

  const pane = document.createElement('section');
  pane.classList.add('tab-pane');
  pane.id = `${id}`;
  pane.setAttribute('data-group', 'patient-tab-panes');
  pane.setAttribute('data-follower', `a[data-open='#${id}']`);
  patient.appendChild(pane);

  return {tab, pane};
}
