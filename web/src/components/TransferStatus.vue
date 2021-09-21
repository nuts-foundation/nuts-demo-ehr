<template>
  <div :class="className">
    {{status.status}}
  </div>
</template>
<script>
export default {
  props: {
    status: {
      type: Object,
      default: {
        status: ""
      }
    }
  },
  computed: {
    className() {
      const colors = {
        'completed': {bg: 'green-300', text: 'green-800'},
        'in-progress': {bg: 'yellow-300', text: 'yellow-800'},
        'cancelled': {bg: 'red-300', text: 'red-800'},
        'requested': {bg: 'gray-300', text: 'gray-800'},
      }[this.status.status] || {bg: 'white', text: 'black'}

      return `inline-block text-sm rounded-xl p-1 pl-3 pr-3 bg-${colors.bg} text-${colors.text}`
    }
  }
}
</script>
