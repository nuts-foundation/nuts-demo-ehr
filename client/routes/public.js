import Thimbleful from 'thimbleful';
import login      from '../components/login';
import logout     from '../components/logout';

const router = new Thimbleful.Router();

router.addRoute(/login\/?([\da-z-]+)?/, async (link, matches, evnt) => {
  await login.render(matches[1], evnt);
  openPage('login');
});

router.addRoute('logout', async link => {
  await logout.render();
});

// Show the given page, hide others
function openPage (page) {
  document.querySelector('.page.active').classList.remove('active');
  document.getElementById(page).classList.add('active');
  window.scrollTo(0, 0);
}

export default router;
