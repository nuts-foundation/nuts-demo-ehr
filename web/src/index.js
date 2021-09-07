import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import EHRApp from './ehr/EHRApp.vue'
import Login from './Login.vue'
import PasswordAuthentication from './components/auth/PasswordAuthentication.vue'
import IRMAAuthentication from './components/auth/IRMAAuthentication.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Api from './plugins/api'
import StatusReporter from './plugins/StatusReporter.js'
import Patients from './ehr/patient/Patients.vue'
import Patient from './ehr/patient/Patient.vue'
import PatientOverview from './ehr/patient/PatientOverview.vue'
import NewPatient from './ehr/patient/NewPatient.vue'
import EditPatient from "./ehr/patient/EditPatient.vue"
import NewDossier from "./ehr/patient/dossier/New.vue"
import NewTransfer from "./ehr/patient/transfer/NewTransfer.vue"
import EditTransfer from "./ehr/patient/transfer/EditTransfer.vue"
import TransferRequest from "./ehr/transfer/TransferRequest.vue"
import Inbox from "./ehr/Inbox.vue"
import Settings from "./ehr/Settings.vue"
import Components from "./Components.vue"
import Elevation from "./components/auth/SessionElevation.vue"

const routes = [
  {path: '/', component: Login},
  {
    name: 'login',
    path: '/login',
    component: Login
  },
  {
    name: 'logout',
    path: '/logout',
    component: Logout
  },
  {
    name: 'auth.passwd',
    path: '/auth/passwd/',
    component: PasswordAuthentication,
  },
  {
    name: 'auth.irma',
    path: '/auth/irma/',
    component: IRMAAuthentication,
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
        path: 'elevate',
        name: 'ehr.elevate',
        component: Elevation,
        props: route => ({redirectPath: route.query.redirect})
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
        meta: {requiresElevation: true}
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
  console.log("from: ", from.path, from.name, "to: ", to.path, to.name)
  // Check before rendering the target route that we're authenticated, if it's required by the particular route.
  if (to.meta.requiresAuth === true) {
    if (!localStorage.getItem("session")) {
      console.log("no cookie found, redirect to login")
      return '/login'
    }
  }
  if (to.meta.requiresElevation === true) {
    let sessionStr = localStorage.getItem("session")
    let rawToken = atob(sessionStr.split('.')[1])
    let token = JSON.parse(rawToken)
    if (!token["elv"]) {
      console.log("route requires elevation, redirect to elevate")
      next({name: 'ehr.elevate', props: true, query: {redirect: to.path }})
    }
  }
  next()
})

app.use(router)
app.use(StatusReporter)
app.use(Api, {forbiddenRoute: {name: 'logout'}})
app.mount('#app')
