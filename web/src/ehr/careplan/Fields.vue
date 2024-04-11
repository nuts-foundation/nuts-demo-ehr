<template>
  <div class="mt-4" v-if="carePlan">
    <div class="mt-6">
      <label>Title</label>
      <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
        <input type="text" v-model="carePlan.title"/>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    carePlan: {
      title: "",
    },
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