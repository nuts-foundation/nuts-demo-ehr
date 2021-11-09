<template>
  <form @submit.stop.prevent="createDossier" novalidate>
    <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-md shadow-sm" v-if="formErrors.length">
      <label class="text-red-500">Please correct the following error{{ formErrors.length === 0 ? '' : 's' }}:</label>

      <ul class="text-sm">
        <li v-for="error in formErrors">— {{ error }}</li>
      </ul>
    </div>

    <div class="mt-6">
      <button type="submit" class="btn btn-primary mr-4" :class="{'btn-loading': loading}">
        Create
      </button>
    </div>
  </form>
</template>
<script>
export default {
  data() {
    return {
      loading: false,
      formErrors: [],
    }
  },
  emits: ['statusUpdate'],
  methods: {
    checkForm() {
      this.formErrors.length = 0

      return this.formErrors.length === 0
    },
    createDossier() {
      if (!this.checkForm()) {
        return;
      }

      this.loading = true

      this.$api.createDossier({body: {patientID: this.$route.params.id, name: 'Collaboration'}})
          .then(dossier => this.createCollaboration(dossier.id))
          .then(collaboration => {
            return this.$router.push({
              name: 'ehr.patient.collaboration.edit',
              params: {collaborationID: collaboration.id}
            })
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    },
    createCollaboration(dossierID) {
      return this.$api.createCollaboration({body: {dossierID, name: 'Collaboration'}})
    }
  },
}
</script>