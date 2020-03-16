import Thimbleful from 'thimbleful'
import routing from './routing'

// Load the routes
routing.load()

// Enable data attributes for interface components
new Thimbleful.Energize('#app')
