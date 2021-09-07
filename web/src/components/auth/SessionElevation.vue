<template>

  <div v-if="!means">
    Choose a means to elevate your session:
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


  <iframe v-if="means === 'bono'" width="560" height="315" src="https://www.youtube-nocookie.com/embed/19KstSgU-c0"
          title="YouTube video player" frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen></iframe>

  <irma-login v-if="means === 'irma'" @success="onIRMASuccess"></irma-login>
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
    elevateWithIRMA() {
      this.means = "irma"
    },
    onIRMASuccess(token) {
      console.log("irma success!", token)
      localStorage.setItem("session", token)
      this.$router.push(this.redirectPath)
    },
    elevateWithBONO() {
      this.means = "bono"
    },
    elevateWithDummy() {

    }
  }
}
</script>