<template>
  <p>Done, closing window...</p>
</template>
<script>
import {useRouter} from 'vue-router'

export default {
  mounted() {
    // close window after 2 seconds
    window.opener = self;
    setTimeout(window.close, 20)
  }
}
</script>
