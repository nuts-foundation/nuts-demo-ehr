<template>
  <div>
    <h1>Care Plan - {{carePlan && carePlan.fhirCarePlan ? carePlan.fhirCarePlan.title : '' }}</h1>
    <div v-if="carePlan && carePlan.fhirCarePlan">
      <div class="mt-6">
        <label>Status</label>
        <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
          {{ carePlan.fhirCarePlan.status }}
        </div>
      </div>
      <div class="mt-6">
        <label>Involved Organizations</label>
        <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
          {{ carePlan.participants.join(', ') }}
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import CarePlanFields from "./Fields.vue";
import AutoComplete from "../../components/Autocomplete.vue"

export default {
  components: {CarePlanFields, AutoComplete},
  data() {
    return {
      carePlan: null,
    }
  },
  methods: {
    truncate(str, n) {
      return (str.length > n) ? str.substr(0, n - 1) + '...' : str
    },
    fetchCarePlan(dossierID) {
      this.$api.getCarePlan({dossierID})
          .then(result => this.carePlan = result.data)
          .catch(e => this.$status.error(e))
    },
  },
  mounted() {
    this.fetchCarePlan(this.$route.params.dossierID)
  }
}
</script>