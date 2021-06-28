<template>
  <h1 class="text-xl">Patients</h1>

  <p v-if="!!error" class="m-4">Error: {{ error }}</p>

  <table class="min-w-full divide-y divide-gray-200">
    <thead class="bg-gray-50">
    <tr>
      <th>Name</th>
      <th>Birth date</th>
      <th>Gender</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="patient in patients">
      <td>{{ patient.surname }}, {{ patient.firstName }}</td>
      <td>{{ patient.dob }}</td>
      <td>{{ patient.gender }}</td>
    </tr>
    </tbody>
  </table>

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
      this.$api.get("web/patients")
          .then((response) => this.patients = response)
          .catch((reason) => this.error = reason)
    },
  },
}
</script>