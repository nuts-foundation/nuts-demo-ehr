<template>

  <div class="flex justify-center">

    <div v-if="!means" class="mt-12 border rounded-md max-w-7xl p-8 flex flex-col">
      <p>The page you want to access requires a higher lever of authentication.</p>

      <p>Choose a means to elevate your session:</p>
      <button class="w-full btn btn-submit grid justify-items-center" @click="elevateWithIRMA">
        <div>IRMA</div>
        <img class="block my-3" v-bind:src="irmaLogo">
      </button>
      <button class="w-full btn btn-submit grid justify-items-center" @click="elevateWithBONO">
        <div>BONO</div>
      </button>
      <button class="w-full btn btn-submit grid justify-items-center" @click="elevateWithDummy">
        Dummy (for easy testing)
      </button>
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