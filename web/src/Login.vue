<template>
  <div class="flex justify-center">

    <div class="mt-12 bg-white shadow-sm border rounded-md w-96 p-8 flex flex-col">
      <h1>Nuts Demo EHR</h1>
      <h2>Login</h2>

      <form class="w-full mt-4" @submit.stop.prevent="">
        <div class="space-y-4">
          <div>
            <label for="customer_select">Choose your organization</label>
            <div class="custom-select">
              <select id="customer_select" v-model="selectedCustomer">
                <option v-for="c in customers" v-bind:value="c">
                  {{ c.name }}
                </option>
              </select>

              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#444">
                <path d="M24 24H0V0h24v24z" fill="none" opacity=".87"/>
                <path d="M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6-1.41-1.41z"/>
              </svg>
            </div>
          </div>

          <p v-if="!!loginError" class="p-2 text-center bg-red-100 rounded-md">{{ loginError }}</p>

          <div class="pt-6">
            <h2 class="mb-4">Pick a method:</h2>

            <div class="grid grid-cols-2 gap-2">
              <button class="btn btn-primary flex justify-center items-center w-full" @click="loginWithIRMA" v-bind:disabled="selectedCustomer === null">
              <span class="w-8 mr-2">
                <img class="max-w-full" v-bind:src="irmaLogo">
              </span>

                <span>
                IRMA
              </span>
              </button>

              <button class="btn btn-primary block w-full" @click="loginWithPassword" v-bind:disabled="selectedCustomer === null">
                Password
              </button>
            </div>
          </div>
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
  props: ['redirectPath'],
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
        this.$router.push({name: 'auth.passwd', params: {customer: JSON.stringify(this.selectedCustomer)}, query: {redirect: this.redirectPath}})
      }
    },
    loginWithIRMA() {
      if (this.selectedCustomer) {
        this.$router.push({name: 'auth.irma', params: {customer: JSON.stringify(this.selectedCustomer)}, query: {redirect: this.redirectPath}})
      }
    }
  },
}
</script>
