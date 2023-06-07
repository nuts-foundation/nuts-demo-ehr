<template>
  <div class="flex justify-center">

    <div v-if="!means" class="mt-12 border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="mb-2">Session elevation</h1>

      <p class="mb-6">The page you want to access requires a higher lever of authentication.</p>

      <h2 class="mb-4">Choose a method</h2>

      <div class="grid grid-cols-3 gap-4">
        <button class="btn btn-secondary" @click="elevateWithIRMA" id="elevate-irma-button">
        <span class="w-10 mr-4">
          <img class="max-w-full" alt="IRMA logo" v-bind:src="irmaLogo">
        </span>

          <span class="text-lg font-semibold">IRMA</span>
        </button>

        <button class="btn btn-secondary" @click="elevateWithEmployeeID" id="elevate-employeeid-button">
          <span class="text-lg font-semibold">EmployeeID</span>
        </button>

        <button class="btn btn-secondary" @click="elevateWithDummy" id="elevate-dummy-button">
          <span class="text-lg font-semibold">Dummy <small>(for easy testing)</small></span>
        </button>
      </div>
    </div>
    <div v-if="employeeIDSessionURL !== null" class="px-12 py-8">
      <h1 class="mb-6 mt-12">EmployeeID</h1>
      <iframe :src="employeeIDSessionURL" title="EmployeeID Authentication"
              width="560" height="650"></iframe>
    </div>

    <irma-login v-if="means === 'irma'"
                @success="onElevationSuccess"
                @aborted="means = ''"
    ></irma-login>
  </div>
</template>

<script>
import IrmaLogin from './IRMAAuthentication.vue'
import irmaLogo from '../../img/irma-logo.png'

export default {
  props: ['redirectPath'],
  components: {IrmaLogin},
  data() {
    return {
      means: null,
      irmaLogo: irmaLogo,
      employeeIDSessionURL: null,
      employeeIDResultPoller: null,
    }
  },
  methods: {
    onElevationSuccess(token) {
      console.log("elevation success!", token)
      localStorage.setItem("session", token)
      this.$router.push(this.redirectPath)
    },
    elevateWithIRMA() {
      this.means = "irma"
    },
    elevateWithEmployeeID() {
      this.means = "EmployeeID"
      // Elevation with EmployeeID means starting a EmployeeID means authentication session (providing the current auth session token),
      // showing the user the returned URL (which is the consent page) in an IFrame,
      // then polling the server for the result of the authentication session.
      this.$api.authenticateWithEmployeeID()
          .then(session => {
            if (!session.sessionPtr.url) {
              throw "No URL returned by server";
            }
            this.employeeIDSessionURL = session.sessionPtr.url;
            this.employeeIDResultPoller = setInterval(() => {
              this.$api.getEmployeeIDAuthenticationResult({sessionToken: session.sessionID})
                  .then((sessionResult) => {
                    clearInterval(this.employeeIDResultPoller);
                    if (sessionResult.token) {
                      // Wait for 3 seconds, so the user can see the result of the authentication session.
                      setTimeout(() => this.onElevationSuccess(sessionResult.token), 3000);
                    }
                  })
                  .catch(err => {
                    switch(err) {
                      case "created":
                      case "in-progress":
                        // Normal operation, user still has to choose to accept/reject
                        break;
                      default:
                        // All other cases are errors, unless the user rejected
                        clearInterval(this.employeeIDResultPoller);
                        if (err === "cancelled") {
                          console.log('User has rejected EmployeeID auth');
                        } else {
                          console.error('EmployeeID authentication error: ' + err)
                        }
                    }
                  })
            }, 2000)
          })
          .catch(err => console.log(err))
    },
    elevateWithDummy() {
      console.log("elevate with Dummy")
      this.$api.authenticateWithDummy()
          .then((res) => {
            console.log(res)
            return this.$api.getDummyAuthenticationResult({sessionToken: res.sessionID})
                .then((responseData) => {
                  this.onElevationSuccess(responseData.token)
                })
                .catch(err => console.log(err))
          })
          .catch(err => console.log(err))

    }
  },
  beforeUnmount() {
    if (this.employeeIDResultPoller) {
      clearInterval(this.employeeIDResultPoller);
    }
  }
}
</script>
