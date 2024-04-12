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
          <div v-if="carePlan.participants.length === 0">
            No organizations involved yet.
          </div>
          <div v-else>
            {{ carePlan.participants.map(p => p.name).join(', ') }}
          </div>
        </div>
      </div>
      <div class="mt-6">
        <label>Activities</label>
        <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
          <div v-if="carePlan.fhirCarePlan.activity && carePlan.fhirCarePlan.activity.length === 0">
            No activities added yet.
          </div>
          <div v-else>
            <table>
              <thead>
              <tr>
                <th>Code</th>
                <th>Display</th>
                <th>Status</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(activity, i) in carePlan.fhirCarePlan.activity" :key="'carePlanActivity' + i">
                <td>{{ carePlan.fhirActivityTasks[activity.reference.reference] ? carePlan.fhirActivityTasks[activity.reference.reference].code.coding[0].code : '' }}</td>
                <td>{{ carePlan.fhirActivityTasks[activity.reference.reference] ? carePlan.fhirActivityTasks[activity.reference.reference].code.coding[0].display : '' }}</td>
                <td>{{ carePlan.fhirActivityTasks[activity.reference.reference] ? carePlan.fhirActivityTasks[activity.reference.reference].status : '' }}</td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <div class="mt-6">
        <label>Add Activity</label>
        <div class="bg-white p-5 shadow-sm rounded-lg mb-3">
          <p>
            Creates a new activity for this care plan. If the activity is requested from another organization,
            they will be notified and added to the care plan. They can then collaborate on the care plan.
          </p>
          <div class="mb-2 mt-2">
            Activity:
            <select id="addActivitySelect" v-model="newActivity.code">
              <option :value="availableActivity.code" v-for="availableActivity in availableActivities"
                      :key="availableActivity.code.coding[0].system + '|' + availableActivity.code.coding[0].code">
                {{ availableActivity.code.coding[0].code }} - {{ availableActivity.code.coding[0].display }}
              </option>
            </select>
          </div>
          <div class="mb-2">
            Request from organization:
            <div v-if="newActivity.owner === null">
              <auto-complete
                  :items="availableOrganizations"
                  @selected="(item) => newActivity.owner = item"
                  @search="searchOrganizations"
                  v-slot="slotProps">
                {{ organizationName(slotProps.item) }}
              </auto-complete>
            </div>
            <div v-else>
              {{ organizationName(newActivity.owner) }}
            </div>
          </div>
          <button id="add-activity-button"
                  type="submit"
                  class="btn btn-primary mr-4"
                  @click="addNewActivity">
            Add/Request Activity
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<style>
td {
  border: 1px solid #cccccc;
  padding: 5px;
}
</style>
<script>
import CarePlanFields from "./Fields.vue";
import AutoComplete from "../../components/Autocomplete.vue"
import AvailableActivities from "./Activities"

export default {
  components: {CarePlanFields, AutoComplete},
  data() {
    return {
      carePlan: null,
      availableActivities: AvailableActivities(),
      availableOrganizations: [],
      newActivity: {
        code: null,
        owner: null,
      },
    }
  },
  methods: {
    fetchCarePlan(dossierID) {
      this.$api.getCarePlan({dossierID})
          .then(result => this.carePlan = result.data)
          .catch(e => this.$status.error(e))
    },
    searchOrganizations(query) {
      this.$api.searchOrganizations(null, {
        query: {"credentialSubject.organization.name": query + '*'},
        excludeOwn: false
      })
          .then((result) => this.availableOrganizations = Object.values(result.data))
          .catch(error => this.$status.error(error))
    },
    organizationName(organization) {
      return organization.name + ' (' + Object.keys(organization.identifiers).map(i => i + ': ' + organization.identifiers[i]).join(', ') + ')'
    },
    addNewActivity() {
      this.$api.createCarePlanActivity(
          {dossierID: this.$route.params.dossierID},
          {
            code: this.newActivity.code,
            owner: {
              system: "http://fhir.nl/fhir/NamingSystem/ura",
              value: this.newActivity.owner.identifiers.ura
            }
          },
      )
        .then((result) => {
          this.carePlan = result.data
          this.newActivity = {
            code: null,
            owner: null,
          }
        })
        .catch(error => this.$status.error(error))
    }
  },
  mounted() {
    this.fetchCarePlan(this.$route.params.dossierID)
  }
}
</script>