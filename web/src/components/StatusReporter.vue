<template>
  <div role="status" :data-message="message" class="px-6 py-4 bg-red-100 rounded-md mb-4 fixed top-10 right-10 shadow-md"
     :class="{ 'bg-red-100': type === 'error', 'bg-blue-100': type === 'info' }"
     v-if="show">{{ message }}</div>
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
