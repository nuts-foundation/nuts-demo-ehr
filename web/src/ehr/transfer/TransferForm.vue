<template>
  <div class="mt-4" v-if="transfer">
    <div class="bg-gray-50 font-bold">Problems</div>
    <button @click="transfer.carePlan.patientProblems.push({problem: {name: ''}, interventions: []})" class="btn btn-secondary">Add
      problem
    </button>
    <ul>
      <li v-for="patientProblem in transfer.carePlan.patientProblems" class="border pl-4">
        <b>Problem:</b>
        <div>
          <textarea v-model="patientProblem.problem.name" class="border min-w-full h-8" required></textarea>
          <button @click="patientProblem.interventions.push({ name: ''})" class="btn btn-secondary">Add intervention</button>
          <div v-for="intervention in patientProblem.interventions" class="border">
            <b>Intervention:</b>
            <textarea v-model="intervention.comment" class="border min-w-full h-8"></textarea>
          </div>
        </div>
      </li>
    </ul>
  </div>
  <div class="mt-4">
    <div class="bg-gray-50 font-bold">Transfer date</div>
    <div>
      <td><input type="date" v-model="transfer.transferDate" required></td>
    </div>
  </div>
</template>
<script>

export default {
  props: {
    transfer: {
      carePlan: {
        patientProblems: [{
          problem: {name: String},
          interventions: [
            {comment: String}
          ]
        }]
      }
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
