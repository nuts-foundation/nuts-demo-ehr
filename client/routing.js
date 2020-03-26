import Thimbleful from 'thimbleful'
import patientOverview from './components/patient-overview'
import inbox from './components/inbox'
import transactions from './components/transactions'
import patient from './components/patient/patient'
import login from './components/login'
import escalate from './components/escalate'
import logout from './components/logout'
import remoteOrganisation from './components/patient/remote/organisation'
import header from './components/header'

const Router = new Thimbleful.Router()
login.load()

export default {
  load: () => {
    // Root redirects to patient overview
    if (!window.location.hash) window.location.hash = 'dashboard'

    Router.addRoute('dashboard', async link => {
      openPage('private', link)

      // Render organisation name, colour and user
      header.render()

      await patientOverview.render()

      // These may come in later, that's ok
      inbox.render()
      transactions.render()
    })

    Router.addRoute(/login\/?([\da-z-]+)?/, async (link, matches, evnt) => {
      openPage('public', 'login')
      login.render(matches[1], evnt)
    })

    Router.addRoute('escalate', async link => {
      openPage('private', 'escalate')
      escalate.render()
    })

    Router.addRoute('logout', async link => {
      logout.render()
    })

    Router.addRoute(/patient-details\/([\da-z\-]+)(\/.*)?/, async (link, matches) => {
      openPage('private', 'patient')
      // Render organisation name, colour and user
      header.render()
      await patient.render(matches[1])
    })

    Router.addRoute(/patient-network\/([\da-z\-]+)\/(.*)?/, async (link, matches) => {
      openPage('private', 'patient')
      // Render organisation name, colour and user
      header.render()
      if (!patient.rendered()) { await patient.render(matches[1]) }
      await remoteOrganisation.render(matches[1], matches[2])
    })
  }
}

// Show the given page, hide others
function openPage (layout, page) {
  document.querySelector('.layout.active').classList.remove('active')
  document.querySelector('.page.active').classList.remove('active')
  document.querySelector(`#${page}`).classList.add('active')
  document.querySelector(`#${layout}`).classList.add('active')
  window.scrollTo(0, 0)
}
