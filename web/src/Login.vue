<template>
  <div class="flex justify-center">

    <div class="mt-12 border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="text-3xl py-2">Nuts Demo EHR</h1>
      <form class="my-4 flex justify-center" @submit.stop.prevent="login">
        <div class="space-y-4">

          <div>
            <label for="customer_select" class="block text-sm font-medium text-gray-700">Organization</label>
            <select id="customer_select" v-model="customer">
              <option v-for="c in customers" v-bind:value="c">
                {{ c.name }}
              </option>
            </select>
          </div>

          <div>
            <span class="block text-sm font-medium text-gray-700">
              Selected: {{ customer? customer.name : 'none' }}
            </span>
          </div>
          <p v-if="!!loginError" class="p-2 text-center bg-red-100 rounded-md">{{ loginError }}</p>
          <button
              class="w-full btn-submit"
          >Login
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import irma from "@privacybydesign/irma-frontend";

export default {
  data() {
    return {
      loginError: "",
      customers: [],
      customer: null
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    redirectAfterLogin() {
      this.$router.push("/ehr/")
    },
    fetchData() {
      this.$api.get('web/customers')
          .then(data => this.customers = data)
          .catch(response => {
            console.error("failure", response)
            if (response.status === 403) {
              this.loginError = "Invalid credentials"
              return
            }
            this.loginError = response
          })
          .finally(() => this.loading = false)
    },
    login() {
      if (!this.customer) {
        return
      }
      let options = {
        // Developer options
        debugging: true,

        // Front-end options
        language: 'en',
        translations: {
          header: this.customer.name
        },

        // Back-end options
        session: {
          // Point to demo-ehr backend which forwards requests
          url: '/web/auth',

          // Define your disclosure request:
          start: {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(this.customer)
          },
          mapping: {
            sessionPtr:      r => r.sessionPtr.clientPtr,
            sessionToken:    r => r.sessionID
          }
        }
      };
      let irmaPopup = irma.newPopup(options);
      irmaPopup.start()
          .then(result => {
            console.log("success!")
            localStorage.setItem("session", result.token)
            this.redirectAfterLogin()
          })
          .catch(error => {
            if (error === 'Aborted') {
              console.log('Aborted');
              return;
            }
            console.error("error", error);
          })
          .finally(() => irmaPopup = irma.newPopup(options));
    }
  },
  mounted() {
    // Check if session still valid, if so just redirect to application
    this.$api.get('web/private')
        .then(() => this.redirectAfterLogin())
        .catch(() => {
          // session is invalid, need to authenticate
        })
  }
}
</script>
