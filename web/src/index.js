import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'

import './style.css'
import store from "./ehr/store"
import App from './App.vue'
import EHRApp from './ehr/EHRApp.vue'
import Login from './Login.vue'
import PasswordAuthentication from './components/auth/PasswordAuthentication.vue'
import OpenID4VPLogin from './components/auth/OpenID4VPAuthentication.vue'
import Logout from './Logout.vue'
import Close from './Close.vue'
import NotFound from './NotFound.vue'
import Api from './plugins/api'
import StatusReporter from './plugins/StatusReporter.js'
import Patients from './ehr/patient/Patients.vue'
import Patient from './ehr/patient/Patient.vue'
import PatientOverview from './ehr/patient/PatientOverview.vue'
import NewEpisode from './ehr/episode/New.vue'
import EditEpisode from './ehr/episode/Edit.vue'
import NewCarePlan from './ehr/careplan/New.vue'
import EditCarePlan from './ehr/careplan/Edit.vue'
import NewPatient from './ehr/patient/NewPatient.vue'
import EditPatient from "./ehr/patient/EditPatient.vue"
import ViewRemotePatient from "./ehr/patient/ViewRemotePatient.vue"
import NewDossier from "./ehr/patient/dossier/New.vue"
import NewTransfer from "./ehr/transfer/NewTransfer.vue"
import EditTransfer from "./ehr/transfer/EditTransfer.vue"
import TransferRequest from "./ehr/inbox/TransferRequest.vue"
import Inbox from "./ehr/inbox/Inbox.vue"
import Settings from "./ehr/Settings.vue"
import Components from "./Components.vue"
import NewReport from "./ehr/patient/dossier/NewReport.vue"

const routes = [
  {path: '/', component: Login},
  {
    name: 'login',
    path: '/login',
    component: Login,
    props: route => ({redirectPath: route.query.redirect})
  },
  {
    name: 'logout',
    path: '/logout',
    component: Logout
  },
  {
    name: 'close',
    path: '/close',
    component: Close
  },
  {
    name: 'auth.passwd',
    path: '/auth/passwd/',
    component: PasswordAuthentication,
    props: route => ({redirectPath: route.query.redirect})
  },
  {
    name: 'auth.openid4vp',
    path: '/auth/openid4vp/',
    component: OpenID4VPLogin,
    props: route => ({redirectPath: route.query.redirect})
  },
  {
    path: '/ehr',
    components: {
      default: EHRApp,
    },
    children: [
      {
        path: '',
        name: 'ehr.home',
        redirect: '/ehr/patients'
      },
      {
        path: 'patients',
        name: 'ehr.patients',
        component: Patients
      },
      {
        path: 'patients/new',
        name: 'ehr.patients.new',
        component: NewPatient
      },
      {
        path: 'patients/remote',
        name: 'ehr.patients.remote',
        component: ViewRemotePatient
      },
      {
        path: 'patient/:id',
        name: 'ehr.patient',
        component: Patient,
        redirect: {name: 'ehr.patient.overview'},
        children: [
          {
            path: 'overview',
            name: 'ehr.patient.overview',
            component: PatientOverview
          },
          {
            path: 'edit',
            name: 'ehr.patient.edit',
            component: EditPatient
          },
          {
            path: 'dossier/new',
            name: 'ehr.patient.dossier.new',
            component: NewDossier
          },
          {
            path: 'episode/new',
            name: 'ehr.patient.episode.new',
            component: NewEpisode
          },
          {
            path: 'episode/edit/:episodeID',
            name: 'ehr.patient.episode.edit',
            component: EditEpisode,
            children: [{
              path: 'newReport',
              name: 'ehr.patient.episode.newReport',
              component: NewReport
            }]
          },
          {
            path: 'careplan/new',
            name: 'ehr.patient.careplan.new',
            component: NewCarePlan,
          },
          {
            path: 'careplan/edit/:dossierID',
            name: 'ehr.patient.careplan.edit',
            component: EditCarePlan,
          },
          {
            path: 'transfer',
            name: 'ehr.patient.transfer.new',
            component: NewTransfer
          },
          {
            path: 'transfer/:transferID/edit',
            name: 'ehr.patient.transfer.edit',
            component: EditTransfer
          },
        ],
      },
      {
        path: 'transfer-request/:requestorDID/:fhirTaskID',
        name: 'ehr.transferRequest.show',
        component: TransferRequest
      },
      {
        path: 'inbox',
        name: 'ehr.inbox',
        component: Inbox,
      },
      {
        path: 'settings',
        name: 'ehr.settings',
        component: Settings
      }
    ],
    meta: {requiresAuth: true}
  },
  {path: '/test/components', component: Components},
  {path: '/:pathMatch*', name: 'NotFound', component: NotFound}
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

const app = createApp(App)

router.beforeEach((to, from, next) => {
  // Check before rendering the target route that we're authenticated, if it's required by the particular route.
  if (to.meta.requiresAuth === true) {
    if (!localStorage.getItem("session")) {
      return next({name: 'login', props: true, query: {redirect: to.path}})
    }
  }
  next()
})

app.use(router)
app.use(StatusReporter)
app.use(Api, {forbiddenRoute: {name: 'logout'}})
app.use(store)
app.mount('#app')
