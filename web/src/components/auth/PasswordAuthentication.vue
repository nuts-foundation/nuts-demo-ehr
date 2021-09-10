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
          <button class="w-full btn btn-submit">Login</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  props: ['redirectPath'],
  data() {
    return {
      loginError: "",
      credentials: {
        password: '',
        customerID: null
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
  },
  methods: {
    redirectAfterLogin() {
      if (this.redirectPath) {
        return this.$router.push(this.redirectPath)
      }
      this.$router.push("/ehr/")
    },
    login() {
      this.$api.authenticateWithPassword({body: this.credentials})
          .then(responseData => {
            localStorage.setItem("session", responseData.token)
            console.log("Password authentication successful")
            this.redirectAfterLogin()
          })
          .catch(response => {
            console.log("Password authentication failed: " + response)
            this.loginError = response
          })
    },
  },
  mounted() {
    // // Check if session still valid, if so just redirect to application
    // Currently disabled, make conditional on existence of cookie / JWT
    // this.$api.get('web/private')
    //     .then(() => this.redirectAfterLogin())
    //     .catch(() => {
    //       // session is invalid, need to authenticate
    //     }))
    if ('customer' in this.$route.params) {
      this.customer = JSON.parse(this.$route.params.customer)
      this.credentials.customerID = this.customer.id
    } else {
      // Missing required params, redirect to landing page
      console.log("missing customer in params. Back to login page.")
      this.$router.push("/")
    }
  }
}
</script>
