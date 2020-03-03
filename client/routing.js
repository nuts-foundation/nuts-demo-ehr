import Thimbleful from 'thimbleful'

import patientOverview from './components/patient-overview'
import inbox from './components/inbox'
import transactions from './components/transactions'
import patient from './components/patient/patient'
import irmaLogin from './components/irma-login'
import remoteOrganisation from './components/patient/remote/organisation'
import sso from './components/sso'
import header from './components/header'

const router = new Thimbleful.Router()

export default {
  load: () => {
    // Root redirects to patient overview
    if (!window.location.hash) window.location.hash = 'dashboard'

    router.addRoute('dashboard', async link => {
      openPage('private', link)

      // Render organisation name, colour and user
      header.render()

      await patientOverview.render()

      // These may come in later, that's ok
      inbox.render()
      transactions.render()
    })

    router.addRoute('irma-login', async link => {
      openPage('public', 'irma-login')
      irmaLogin.render()
    })

    router.addRoute(/patient-details\/([\da-z\-]+)(\/.*)?/, async (link, matches) => {
      openPage('private', 'patient')
      // Render organisation name, colour and user
      header.render()
      await patient.render(matches[1])
    })

    router.addRoute(/patient-network\/([\da-z\-]+)\/(.*)?/, async (link, matches) => {
      openPage('private', patient)
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
