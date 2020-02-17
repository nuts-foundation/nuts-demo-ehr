import Thimbleful from 'thimbleful';
const router = new Thimbleful.Router();

import patientOverview    from './components/patient-overview';
import inbox              from './components/inbox';
import transactions       from './components/transactions';
import patient            from './components/patient/patient';
import irmaLogin          from './components/irma-login';
import remoteOrganisation from './components/patient/remote/organisation';

export default {
  load: () => {
    // Root redirects to patient overview
    if ( !window.location.hash ) window.location.hash = 'dashboard';

    router.addRoute('dashboard', async link => {
      await patientOverview.render();
      openPage(link);

      // These may come in later, that's ok
      inbox.render();
      transactions.render();
    });

    router.addRoute('irma-login', async link => {
      irmaLogin.render();
      openPage('irma-login');
    })

    router.addRoute(/patient-details\/([\da-z\-]+)(\/.*)?/, async (link, matches) => {
      await patient.render(matches[1]);
      openPage('patient');
    });

    router.addRoute(/patient-network\/([\da-z\-]+)\/(.*)?/, async (link, matches) => {
      await remoteOrganisation.render(matches[1], matches[2]);
    });
  }
}

// Show the given page, hide others
function openPage(page) {
  document.querySelector('.page.active').classList.remove('active');
  document.querySelector(`#${page}`).classList.add('active');
  window.scrollTo(0,0);
}
