<template>
  <div class="flex flex-row gap-4 m4">
    <div class="w-24 h-24 border"></div>
    <div v-if="!!patient">
      <div class="flex">
        <h1 class="text-2xl mb-2 mr-4">{{ patient.surname }}, {{ patient.firstName }}</h1>
        <button
            @click="$router.push({name: 'ehr.patient.edit', params: {id: patient.ObjectID}})"
            class="float-right inline-flex items-center"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
          </svg>
          Edit
        </button>
      </div>
      <div class="grid grid-cols-2 gap-x-6">
        <div><span class="text-sm font-bold">SSN</span>: {{ patient.ssn }}</div>
        <div><span class="text-sm font-bold">Gender</span>: {{ patient.gender }}</div>
        <div><span class="text-sm font-bold">Birth date</span>: {{ patient.dob }}</div>
        <div><span class="text-sm font-bold">Patient number</span>: {{ patient.id }}</div>
        <div><span class="text-sm font-bold">E-mail</span>: {{ patient.email }}</div>
        <div><span class="text-sm font-bold">Zipcode</span>: {{ patient.zipcode }}</div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  props: {
    patient: Object,
  },
}
</script>