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
              Selected: {{ customer.name }}
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

export default {
  data() {
    return {
      loginError: "",
      customers: [],
      customer: {}
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
