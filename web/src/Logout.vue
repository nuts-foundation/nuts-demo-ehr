<template>
  <p>You are being logged out...</p>
</template>
<script>
import {useRouter} from 'vue-router'

export default {
  mounted() {
    this.$cookies.remove("session")
    useRouter().push("/login")
  }
}
</script>