import Thimbleful from 'thimbleful';
const router = new Thimbleful.Router();

import patientOverview from './components/patient-overview';
import inbox           from './components/inbox';
import patient         from './components/patient/patient';

export default {
  load: () => {
    // Root redirects to patient overview
    if ( !window.location.hash ) window.location.hash = 'dashboard';

    router.addRoute('dashboard', async link => {
      await patientOverview.render();
      openPage(link);
      inbox.render(); // This may come in later, that's ok
    });

    router.addRoute(/patient-details\/(\d+)(\/.*)?/, async (link, matches) => {
      await patient.render(matches[1]);
      openPage('patient');
    });
  }
}

// Show the given page, hide others
function openPage(page) {
  document.querySelector('.page.active').classList.remove('active');
  document.querySelector(`#${page}`).classList.add('active');
  window.scrollTo(0,0);
}
