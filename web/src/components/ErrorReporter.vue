<template>
  <p class="p-3 bg-red-100 rounded-md mb-4" v-if="show">Error: {{ message }}</p>
</template>

<script>
export default {
  name: 'error-reporter',
  props: {message: String},
  data() {
    return {
      timeout: null,
      show: false
    }
  },
  watch: {
    message: {
      handler() {
        if (this.timeout) {
          clearTimeout(this.timeout)
        }
        if (!this.message) {
          this.show = false
          return
        }
        this.show = true
        this.timeout = setTimeout(() => {
          this.show = false
        }, 5000)
      }
    }
  },
}
</script>