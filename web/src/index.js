import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import EHRApp from './ehr/EHRApp.vue'
import Login from './Login.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Api from './plugins/api'
import Patients from './ehr/Patients.vue'

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

router.beforeEach((to, from) => {
  if (to.meta.requiresAuth) {
    if (localStorage.getItem("session")) {
      return true
    }
    return '/login'
  }
})

const app = createApp(App)

app.use(router)
app.use(Api, {forbiddenRoute: {name: 'logout'}})
app.mount('#app')
