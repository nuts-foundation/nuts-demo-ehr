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
import Patients from './ehr/patient/Patients.vue'
import Patient from './ehr/patient/Patient.vue'
import NewPatient from './ehr/patient/NewPatient.vue'
import EditPatient from "./ehr/patient/EditPatient.vue";
import Settings from "./ehr/Settings.vue";

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
                path: 'patient/:id',
                name: 'ehr.patient',
                component: Patient
            },
            {
                path: 'patient/:id/edit',
                name: 'ehr.patient.edit',
                component: EditPatient
            },
            {
                path: 'settings',
                name: 'ehr.settings',
                component: Settings
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
    console.log("from: ", from.path, from.name, "to: ", to.path, to.name)
    // Check before rendering the target route that we're authenticated, if it's required by the particular route.
    if (to.meta.requiresAuth === true) {
        if (localStorage.getItem("session")) {
            return true
        }
        console.log("no cookie found, redirect to login")
        return '/login'
    }
})

app.use(router)
app.use(Api, {forbiddenRoute: {name: 'logout'}})
app.mount('#app')
