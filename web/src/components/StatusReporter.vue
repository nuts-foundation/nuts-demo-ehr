<template>
  <p class="p-3 bg-red-100 rounded-md mb-4"
     :class="{ 'bg-red-100': type === 'error', 'bg-blue-100': type === 'info' }"
     v-if="show">{{ message }}</p>
</template>

<script>
export default {
  name: 'status-reporter',
  props: {
    message: String,
    type: {
      type: String,
      default: "info"
    }
  },
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