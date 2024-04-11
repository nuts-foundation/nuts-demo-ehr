<template>
  <h1>New Shared Care Plan</h1>
  <form @submit.stop.prevent="create" novalidate>
    <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-md shadow-sm" v-if="formErrors.length">
      <label class="text-red-500">Please correct the following error{{ formErrors.length === 0 ? '' : 's' }}:</label>
      <ul class="text-sm">
        <li v-for="error in formErrors">— {{ error }}</li>
      </ul>
    </div>

    <fields :care-plan="carePlan"></fields>

    <div class="mt-6">
      <button type="submit" class="btn btn-primary mr-4" :class="{'btn-loading': loading}">
        Create
      </button>
    </div>
  </form>
</template>
<script>
import Fields from "./Fields.vue";

export default {
  components: {Fields},
  data() {
    return {
      loading: false,
      formErrors: [],
      carePlan: {
        title: "",
      }
    }
  },
  methods: {
    checkForm() {
      this.formErrors.length = 0

      if (!this.carePlan.title) {
        this.formErrors.push("Care plan title is required.")
      }

      return this.formErrors.length === 0
    },
    create() {
      if (!this.checkForm()) {
        return;
      }

      this.loading = true

      this.$api.createDossier(null, {patientID: this.$route.params.id, name: 'Shared Care Plan - ' + this.carePlan.title})
          .then(result => this.$api.createCarePlan(null, {dossierID: result.data.id, title: this.carePlan.title}))
          .then(result => {
            return this.$router.push({
              name: 'ehr.patient.careplan.edit',
              params: {dossierID: result.data.id}
            })
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    },
  },
}
</script>