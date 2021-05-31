<template>
  <p>You are being logged out...</p>
</template>
<script>
import {useRouter} from 'vue-router'

export default {
  setup() {
    localStorage.removeItem("session")
    useRouter().push("/login")
  }
}
</script>