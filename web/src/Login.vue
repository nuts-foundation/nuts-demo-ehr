<template>
  <div class="flex justify-center">

    <div class="mt-12 border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="text-3xl py-2">Nuts Demo EHR</h1>
      <form class="my-4 flex justify-center" @submit.stop.prevent="">
        <div class="space-y-4">

          <div>
            <label for="customer_select" class="block text-sm font-medium text-gray-700">Organization</label>
            <select id="customer_select" v-model="selectedCustomer">
              <option v-for="c in customers" v-bind:value="c">
                {{ c.name }}
              </option>
            </select>
          </div>
          <p v-if="!!loginError" class="p-2 text-center bg-red-100 rounded-md">{{ loginError }}</p>
          <button class="w-full btn-submit grid justify-items-center" @click="loginWithIRMA" v-bind:disabled="selectedCustomer === null">
            <div>Login with IRMA</div>
            <img class="block my-3" v-bind:src="irmaLogo">
          </button>
          <button class="w-full btn-submit grid justify-items-center" @click="loginWithPassword" v-bind:disabled="selectedCustomer === null">
            Login with password
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style>
  .btn-login:disabled {
    filter: grayscale(1);
  }
</style>

<script>
import irma from "@privacybydesign/irma-frontend";
import irmaLogo from './img/irma-logo.png';

export default {
  data() {
    return {
      loginError: "",
      customers: [],
      selectedCustomer: null,
      irmaLogo: irmaLogo,
    }
  },
  created() {
    this.fetchData()
  },
  watch: {
    // changes to selected customer
    'selectedCustomer'() {
      localStorage.removeItem("session")
      this.$api.setCustomer({body: this.selectedCustomer})
          .then(responseData => {
            localStorage.setItem("session", responseData.token)
            this.loginError = ''
          })
          .catch(reason => {
            console.error("failure", reason)
            this.loginError = reason
          })
    }
  },
  methods: {
    redirectAfterLogin() {
      this.$router.push("/ehr/")
    },
    fetchData() {
      this.$api.listCustomers()
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
    loginWithPassword() {
      if (this.selectedCustomer) {
        this.$router.push({name: 'auth.passwd', params: {customer: JSON.stringify(this.selectedCustomer)}})
      }
    },
    loginWithIRMA() {
      if (this.selectedCustomer) {
        this.$router.push({name: 'auth.irma', params: {customer: JSON.stringify(this.selectedCustomer)}})
      }
    }
  },
  mounted() {
    // Disabled since it causes infinite loops since the $api service is redirecting to this page when a 401 is returned.
    // Can be enabled when switching back to a JWT. Only perform when it is present.
    // Check if session still valid, if so just redirect to application
    // this.$api.get('web/private')
    //     .then(() => this.redirectAfterLogin())
    //     .catch(() => {
    //       // session is invalid, need to authenticate
    //     })
  }
}
</script>
