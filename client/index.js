import Thimbleful from 'thimbleful';
import header     from './components/header';
import routing    from './routing';

// Render organisation name, colour and user
header.render();

// Load the routes
routing.load();

// Enable data attributes for interface components
new Thimbleful.Energize("#app");
