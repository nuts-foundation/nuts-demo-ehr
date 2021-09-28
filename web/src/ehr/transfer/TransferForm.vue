<template>
  <div class="mt-4" v-if="transfer">
    <div class="mt-6">
      <label>Transfer date</label>

      <div>
        <input type="date" v-model="transfer.transferDate" required>
      </div>
    </div>

    <div class="flex justify-between items-center mt-8 mb-4">
      <h2>Problems</h2>

      <button
          v-if="mode === 'new'"
          @click="transfer.carePlan.patientProblems.push({problem: {name: ''}, interventions: [{comment: ''}]})"
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>

    <div v-for="patientProblem in transfer.carePlan.patientProblems" class="bg-white p-5 shadow-sm rounded-lg mb-3">
      <label>Problem</label>

      <div>
        <textarea
            v-if="mode === 'new'"
            placeholder="The problem.."
            v-model="patientProblem.problem.name"
            class="min-w-full border"
            required
        ></textarea>
        <p v-else>{{patientProblem.problem.name}}</p>

        <div v-for="(intervention, i) in patientProblem.interventions" class="mt-3">
          <label>
            Intervention
            <span v-if="patientProblem.interventions.length > 1">{{ i + 1 }}</span>
          </label>

          <textarea
              v-if="mode === 'new'"
              placeholder="The intervention.."
              @input="e => addOrRemoveIntervention(e, patientProblem)"
              v-model="intervention.comment"
              class="min-w-full border"></textarea>
          <p v-else>{{intervention.comment}}</p>
        </div>
      </div>
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
  methods: {
    addOrRemoveIntervention(e, patientProblem) {
      const isEmpty = value => (value || '').trim().length === 0;

      if (isEmpty(e.target.value) ||
          isEmpty(patientProblem.interventions[patientProblem.interventions.length - 1].comment)) {
        if (isEmpty(patientProblem.interventions[patientProblem.interventions.length - 2].comment)) {
          patientProblem.interventions.pop();
        }

        return;
      }

      patientProblem.interventions.push({
        comment: ''
      });
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
