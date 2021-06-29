import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import VueCookies from 'vue3-cookies'
import App from './App.vue'
import EHRApp from './ehr/EHRApp.vue'
import Login from './Login.vue'
import PasswordAuthentication from './components/auth/PasswordAuthentication.vue'
import IRMAAuthentication from './components/auth/IRMAAuthentication.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Api from './plugins/api'
import Patients from './ehr/Patients.vue'
import Patient from './ehr/Patient.vue'

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
        path: 'patient/:id',
        name: 'ehr.patient',
        component: Patient
      }
    ],
    meta: {requiresAuth: true}
  },
  {path: '/:pathMatch*', name: 'NotFound', component: NotFound}
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

const app = createApp(App)

router.beforeEach((to, from) => {
  // Check before rendering the target route that we're authenticated, if it's required by the particular route.
  if (to.meta.requiresAuth === true) {
    if (app.config.globalProperties.$cookies.get("session")) {
      return true
    }
    return '/login'
  }
})

app.use(router)
app.use(VueCookies)
app.use(Api, {forbiddenRoute: {name: 'logout'}})
app.mount('#app')
