import Thimbleful    from 'thimbleful';
import publicRoutes  from './routes/public';
import privateRoutes from './routes/private';
import header        from './components/header';

// Enable data attributes for interface components
new Thimbleful.Energize('#app');

/** Load routes so the right components get rendered **/

const router = new Thimbleful.Router().install();

router.addRoute(/public\/([\da-z-\/]+)/, async (link, matches, evnt) => {
  header.render(); // Render organisation name, colour and user
  publicRoutes.route(matches[1], evnt);
  openLayout('public');
});

router.addRoute(/private\/([\da-z-\/:\.]+)/, async (link, matches, evnt) => {
  header.render(); // Render organisation name, colour and user
  privateRoutes.route(matches[1], evnt);
  openLayout('private');
});

// Root leads to login page
if (!window.location.hash) {
  fetch('/api/authentication/logged-in')
  .then(response => {
    if ( !response.ok ) {
      window.location.hash = '#public/login';
    } else {
      window.location.hash = '#private/dashboard';
    }
  });
}

// Show this layout, hide others
function openLayout(layout) {
  document.querySelector('.layout.active').classList.remove('active');
  document.getElementById(layout).classList.add('active');
}
