import Thimbleful      from 'thimbleful';
import patientOverview from '../components/patient-overview';
import inbox           from '../components/inbox';
import transactions    from '../components/transactions';
import patient         from '../components/patient/patient';
import escalate        from '../components/escalate';
import header          from '../components/header';

const router = new Thimbleful.Router();

router.addRoute('dashboard', async link => {
  await patientOverview.render(); // Render patient list
  openPage('dashboard');

  // These may come in later, that's ok
  inbox.render();
  transactions.render();
});

router.addRoute('escalate', async link => {
  await escalate.render();
  openPage('escalate');
});

router.addRoute(/patient\/([\da-z\-]+)(\/(.*)?)?/, async (link, matches) => {
  await patient.render(matches[1], matches[3]);
  openPage('patient');
});

// Show the given page, hide others
function openPage (page) {
  document.querySelector('.page.active').classList.remove('active');
  document.getElementById(page).classList.add('active');
  window.scrollTo(0, 0);
}

export default {
  route: (route, evnt) => {
    // TODO: Redirect if not logged in here
    header.render(); // Render organisation name, colour and user
    router.route(route, evnt);
  }
}
