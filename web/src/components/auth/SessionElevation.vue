<template>

  <div class="flex justify-center">

    <div v-if="!means" class="mt-12 border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="mb-2">Session elevation</h1>

      <p class="mb-6">The page you want to access requires a higher lever of authentication.</p>

      <h2 class="mb-4">Choose a method</h2>

      <div class="grid grid-cols-3 gap-4">
        <button class="btn btn-secondary" @click="elevateWithIRMA">
        <span class="w-10 mr-4">
          <img class="max-w-full" alt="IRMA logo" v-bind:src="irmaLogo">
        </span>

          <span class="text-lg font-semibold">IRMA</span>
        </button>

        <button class="btn btn-secondary" @click="elevateWithBONO">
          <span class="text-lg font-semibold">Bono</span>
        </button>

        <button class="btn btn-secondary" @click="elevateWithDummy">
          <span class="text-lg font-semibold">Dummy <small>(for easy testing)</small></span>
        </button>
      </div>
    </div>
    <button v-else @click="means = ''">back</button>


    <iframe v-if="means === 'bono'" width="560" height="315" src="https://www.youtube-nocookie.com/embed/19KstSgU-c0"
            title="YouTube video player" frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen></iframe>

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
      irmaLogo: irmaLogo
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
    elevateWithBONO() {
      this.means = "bono"
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
  }
}
</script>
