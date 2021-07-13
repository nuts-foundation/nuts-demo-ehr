<template>
  <div class="mt-4">
    <div class="bg-gray-50 font-bold">Description</div>
    <div>
      <textarea v-model="transfer.description" class="border min-w-full h-32" required></textarea>
    </div>
  </div>
  <div class="mt-4">
    <div class="bg-gray-50 font-bold">Transfer date</div>
    <div>
      <td><input type="date" v-model="transfer.transferDate" required></td>
    </div>
  </div>
</template>
<script>

export default {
  props: {
    transfer: Object,
    mode: {
      type: String,
      default: 'new'
    }
  },
  emits: ['input'],
  watch: {
    transfer() {
      this.$emit('input', this.transfer)
    }
  }
}
</script>
