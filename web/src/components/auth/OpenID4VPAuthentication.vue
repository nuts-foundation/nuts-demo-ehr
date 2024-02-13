<template>
  <div class="flex justify-center">

    <div class="mt-12 bg-white border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="text-3xl py-2">Wallet login</h1>
      <p class="text-gray-600" id="status">Follow directions on screen</p>
    </div>
  </div>
</template>

<script>
export default {
  props: ['redirectPath'],
  data() {
    return {
      loginError: "",
      customer: {
        name: "",
      },
    }
  },
  mounted() {
    if ('customer' in this.$route.params) {
      this.customer = JSON.parse(this.$route.params.customer)
    } else {
      // Missing required params, redirect to landing page
      console.log("missing customer in params. Back to login page.")
      this.$router.push("/")
    }
    this.login();
  },
  emits: ['success', 'aborted', 'error'],
  methods: {
    redirectAfterLogin() {
      if (this.redirectPath) {
        return this.$router.push(this.redirectPath)
      }
      this.$router.push("/ehr/")
    },
    login() {
      this.$api.createAuthorizationRequest()
          .then(responseData => {
            // show a popup using the redirect URL
            window.open(responseData.redirect_uri, "_blank", "width=400,height=600")
            console.log("session_id: " + responseData.session_id)

            // update the status in the UI
            document.getElementById("status").innerText = "In progress"
            // start polling the server for the result using responseData.session_id

            // let's poll the server for the result
            // loop and call every 1 second
            const session_id = responseData.session_id
            let interval = setInterval(() => {
              this.$api.getOpenID4VPAuthenticationResult({token: session_id})
                  .then(responseData => {
                    if (responseData.status === "active") {
                      clearInterval(interval)
                      console.log("OpenID4VP authentication successful")
                      console.log("AccessToken: " + responseData.access_token)
                      document.getElementById("status").innerText = "Done"
                    }
                    //this.redirectAfterLogin()
                  })
                  .catch(response => {
                    console.log("OpenID4VP authentication failed: " + response)
                    //this.loginError = response
                  })
            }, 100)

            //localStorage.setItem("session", responseData.token)
            //this.redirectAfterLogin()
          })
          .catch(response => {
            console.log("OpenID4VP authentication failed: " + response)
            this.loginError = response
          })
    },
  },
}
</script>
