<template>
  <irma-auth
      v-if="customer"
      :header-message="'Authenticate for ' + customer.name"
      @success="successHandler"
  />
</template>
<script>
import irmaAuth from './IRMAAuthentication.vue'

export default {
  props: ['redirectPath'],
  components: {irmaAuth},
  data() {
    return {
      customer: null,
    }
  },
  mounted() {
    if ('customer' in this.$route.params) {
      this.customer = JSON.parse(this.$route.params.customer)
    } else {
      // Missing required params, redirect to landing page
      this.$router.push("/")
    }
  },
  methods: {
    successHandler(token) {
      console.log("success!", token)
      localStorage.setItem("session", token)
      if (this.redirectPath) {
        return this.$router.push(this.redirectPath)
      }
      this.$router.push("/ehr/")
    },
  }
}
</script>
