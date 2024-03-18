<template>
  <div class="flex justify-between items-center mb-3">
    <h2>Observations</h2>
  </div>

  <table class="reports-list min-w-full divide-y divide-gray-200">
    <thead>
    <tr>
      <th>Type</th>
      <th>Date</th>
      <th>Value</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="observation in observations">
      <td>{{ observation.code.coding[0].display }}</td>
      <td>{{ observation.effectiveDateTime }}</td>
      <td>{{ observation.valueQuantity.value }}</td>
    </tr>
    </tbody>
  </table>
</template>
<script>
export default {
  props: {
    observations: []
  }
}
</script>
