<template>
  <h1>New Episode</h1>
  <form @submit.stop.prevent="createDossier" novalidate>
    <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-md shadow-sm" v-if="formErrors.length">
      <label class="text-red-500">Please correct the following error{{ formErrors.length === 0 ? '' : 's' }}:</label>
      <ul class="text-sm">
        <li v-for="error in formErrors">— {{ error }}</li>
      </ul>
    </div>

    <episode-fields :episode="episode"></episode-fields>

    <div class="mt-6">
      <button type="submit" class="btn btn-primary mr-4" :class="{'btn-loading': loading}">
        Create
      </button>
    </div>
  </form>
</template>
<script>
import EpisodeFields from "./EpisodeFields.vue";

export default {
  components: {EpisodeFields},
  data() {
    return {
      loading: false,
      formErrors: [],
      episode: {
        period: {
          start: new Date().toISOString().split('T')[0],
          end: null,
        },
        diagnosis: ""
      }
    }
  },
  emits: ['statusUpdate'],
  methods: {
    checkForm() {
      this.formErrors.length = 0

      if (!this.episode.period.start) {
        this.formErrors.push("Episode start date is a required field")
      }
      if (!this.episode.diagnosis) {
        this.formErrors.push("Diagnosis is a required field")
      }

      return this.formErrors.length === 0
    },
    createDossier() {
      if (!this.checkForm()) {
        return;
      }

      this.loading = true

      this.$api.createDossier({body: {patientID: this.$route.params.id, name: 'Episode'}})
          .then(dossier => this.createEpisode(dossier.id))
          .then(episode => {
            return this.$router.push({
              name: 'ehr.patient.episode.edit',
              params: {episodeID: episode.id}
            })
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    },
    createEpisode(dossierID) {
      return this.$api.createEpisode({
        body: {
          dossierID,
          ...this.episode,
        }
      })
    }
  },
}
</script>