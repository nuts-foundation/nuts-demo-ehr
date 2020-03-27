import Thimbleful    from 'thimbleful';
import publicRoutes  from './routes/public';
import privateRoutes from './routes/private';

// Enable data attributes for interface components
new Thimbleful.Energize('#app');

/** Load routes so the right components get rendered **/

const router = new Thimbleful.Router().install();

router.addRoute(/public\/([\da-z-\/]+)/, async (link, matches, evnt) => {
  publicRoutes.route(matches[1], evnt);
  openLayout('public');
});

router.addRoute(/private\/([\da-z-\/]+)/, async (link, matches, evnt) => {
  privateRoutes.route(matches[1], evnt);
  openLayout('private');
});

// Root leads to login page
if (!window.location.hash) window.location.hash = '#public/login';

// Show this layout, hide others
function openLayout(layout) {
  document.querySelector('.layout.active').classList.remove('active');
  document.getElementById(layout).classList.add('active');
}
