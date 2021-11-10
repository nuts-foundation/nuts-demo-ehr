<template>
  <div class="mt-4" v-if="episode">
    <div class="mt-6">
      <label>Start date</label>

      <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
        <input :disabled="mode!=='new'" type="date" v-model="episode.period.start" :min="new Date().toISOString().split('T')[0]" required>
      </div>
    </div>

    <div class="mt-6">
      <label>Diagnosis</label>
      <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
        <input type="text" v-model="episode.diagnosis">
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    episode: {
      period: {
        start: String,
        end: String,
      },
      diagnosis: null,
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