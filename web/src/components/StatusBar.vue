<template>
  <div class="fixed p-4 bottom-0 flex" v-if="show">
    <div class="rounded-md bg-green-300 p-4 flex justify-between">
      {{ statusMessage }}
      <div class="cursor-pointer text-white">
        <svg @click="show = false" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
             stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'status-bar',
  props: {statusMessage: String},
  data() {
    return {
      timeout: null,
      show: false
    }
  },
  watch: {
    statusMessage: {
      handler() {
        this.show = true
        if (this.timeout) {
          clearTimeout(this.timeout)
        }
        this.timeout = setTimeout(() => {
          this.show = false
        }, 5000)
      }
    }
  },
}
</script>