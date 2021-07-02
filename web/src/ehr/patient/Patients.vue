<template>
  <div>
    <h1 class="text-2xl">Patients</h1>

    <p v-if="!!error" class="m-4">Error: {{ error }}</p>

    <button
        @click="$router.push({name: 'ehr.patients.new'})"
        class="float-right inline-flex items-center"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
      </svg>
      Add new Patient</button>
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th>Name</th>
        <th>Birth date</th>
        <th>Gender</th>
      </tr>
      </thead>
      <tbody>
      <tr class="hover:bg-gray-100 cursor-pointer"
          v-for="patient in patients"
          @click="$router.push({name: 'ehr.patient', params: {id: patient.PatientID}})">
        <td>{{ patient.surname }}, {{ patient.firstName }}</td>
        <td>{{ patient.dob }}</td>
        <td>{{ patient.gender }}</td>
      </tr>
      </tbody>
    </table>
  </div>

</template>
<script>
export default {
  data() {
    return {
      patients: [],
      error: null,
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    list() {
      this.$api.get("web/private/patients")
          .then((response) => this.patients = response)
          .catch((reason) => this.error = reason)
    },
  },
}
</script>