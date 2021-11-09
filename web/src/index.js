import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import EHRApp from './ehr/EHRApp.vue'
import Login from './Login.vue'
import PasswordAuthentication from './components/auth/PasswordAuthentication.vue'
import IRMALogin from './components/auth/IRMALogin.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Api from './plugins/api'
import StatusReporter from './plugins/StatusReporter.js'
import Patients from './ehr/patient/Patients.vue'
import Patient from './ehr/patient/Patient.vue'
import PatientOverview from './ehr/patient/PatientOverview.vue'
import NewCollaboration from './ehr/collaboration/New.vue'
import EditCollaboration from './ehr/collaboration/Edit.vue'
import NewPatient from './ehr/patient/NewPatient.vue'
import EditPatient from "./ehr/patient/EditPatient.vue"
import NewDossier from "./ehr/patient/dossier/New.vue"
import NewTransfer from "./ehr/transfer/NewTransfer.vue"
import EditTransfer from "./ehr/transfer/EditTransfer.vue"
import TransferRequest from "./ehr/inbox/TransferRequest.vue"
import Inbox from "./ehr/inbox/Inbox.vue"
import Settings from "./ehr/Settings.vue"
import Components from "./Components.vue"
import Elevation from "./components/auth/SessionElevation.vue"
import NewReport from "./ehr/patient/dossier/NewReport.vue";

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
    name: 'auth.passwd',
    path: '/auth/passwd/',
    component: PasswordAuthentication,
    props: route => ({redirectPath: route.query.redirect})
  },
  {
    name: 'auth.irma',
    path: '/auth/irma/',
    component: IRMALogin,
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
            path: 'dossier/newReport',
            name: 'ehr.patient.dossier.newReport',
            component: NewReport
          },
          {
            path: 'collaboration/new',
            name: 'ehr.patient.collaboration.new',
            component: NewCollaboration
          },
          {
            path: 'collaboration/edit/:collaborationID',
            name: 'ehr.patient.collaboration.edit',
            component: EditCollaboration
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
  // Check before rendering the target route that we're authenticated, if it's required by the particular route.
  if (to.meta.requiresAuth === true) {
    if (!localStorage.getItem("session")) {
      return next({name: 'login', props: true, query: {redirect: to.path}})
    }
  }
  if (to.meta.requiresElevation === true) {
    let sessionStr = localStorage.getItem("session")
    let rawToken = atob(sessionStr.split('.')[1])
    let token = JSON.parse(rawToken)
    if (!token["elv"]) {
      return next({name: 'ehr.elevate', props: true, query: {redirect: to.path}})
    }
  }
  next()
})

app.use(router)
app.use(StatusReporter)
app.use(Api, {forbiddenRoute: {name: 'logout'}})
app.mount('#app')
