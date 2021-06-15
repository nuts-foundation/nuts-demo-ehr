<template>
  <div class="flex justify-center">

    <div class="mt-12 border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="text-3xl py-2">Password login</h1>
      <form class="my-4 flex justify-center" @submit.stop.prevent="login">
        <div class="space-y-4">

          <div class="text-sm font-medium text-gray-700">Organization: {{ customer.name }}</div>

          <div>
            <label for="password_input" class="block text-sm font-medium text-gray-700">Password</label>
            <input
                id="password_input"
                v-model="credentials.password"
                type="password"
                placeholder="Password"
                class="flex-1 py-2 px-4 block border border-gray-300 rounded-md"
            />
          </div>
          <p v-if="!!loginError" class="p-2 text-center bg-red-100 rounded-md">{{ loginError }}</p>
          <button class="w-full btn-submit">Login</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loginError: "",
      credentials: {
        password: ''
      },
      customer: {
        name: "",
      },
    }
  },
  watch: {
    // Remove error when typing
    'credentials.password'() {
      this.loginError = ""
    },
    // Fetch customer from route params
    "$route.params": {
      handler(toParams, previousParams) {
        if (toParams && 'id' in toParams) {
          this.fetchCustomer(toParams.id)
        }
      },
      immediate: true
    }
  },
  methods: {
    redirectAfterLogin() {
      console.log('logged in, redirecting!')
      this.$router.push("/ehr/")
    },
    login() {
      this.$api.post('web/auth-passwd', this.credentials)
          .then(responseData => {
            console.log("Password authentication successful")
            this.redirectAfterLogin()
          })
          .catch(response => {
            this.loginError = response.statusText
          })
    },
    fetchCustomer(id) {
      this.$api.get(`web/customers/${id}`)
          .then((customer) => {
            this.customer = customer
            this.loading = false
          })
          .catch((reason) => {
            console.log("Could not retrieve customer, redirecting to landing page: ", reason)
            this.$router.push("/")
          })
    },
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