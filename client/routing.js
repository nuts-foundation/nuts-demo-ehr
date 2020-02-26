import Thimbleful from 'thimbleful'

import patientOverview from './components/patient-overview'
import inbox from './components/inbox'
import transactions from './components/transactions'
import patient from './components/patient/patient'
import irmaLogin from './components/irma-login'
import remoteOrganisation from './components/patient/remote/organisation'
const router = new Thimbleful.Router()

export default {
  load: () => {
    // Root redirects to patient overview
    if (!window.location.hash) window.location.hash = 'dashboard'

    router.addRoute('dashboard', async link => {
      await patientOverview.render()
      openPage(link)

      // These may come in later, that's ok
      inbox.render()
      transactions.render()
    })

    router.addRoute('irma-login', async link => {
      irmaLogin.render()
      openPage('irma-login')
    })

    router.addRoute(/patient-details\/([\da-z\-]+)(\/.*)?/, async (link, matches) => {
      await patient.render(matches[1])
      openPage('patient')
    })

    router.addRoute(/patient-network\/([\da-z\-]+)\/(.*)?/, async (link, matches) => {
      if (!patient.rendered()) { await patient.render(matches[1]) }
      await remoteOrganisation.render(matches[1], matches[2])
      openPage('patient')
    })
  }
}

// Show the given page, hide others
function openPage (page) {
  document.querySelector('.page.active').classList.remove('active')
  document.querySelector(`#${page}`).classList.add('active')
  window.scrollTo(0, 0)
}
