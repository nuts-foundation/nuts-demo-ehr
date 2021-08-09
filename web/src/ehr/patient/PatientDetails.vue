<template>
  <div class="flex flex-row gap-4 m4">
    <div class="w-24 h-24 border">
      <img v-if="patient.avatar_url" :src="patient.avatar_url" alt="avatar">
      <svg v-else-if="patient.gender === 'female'" id="Layer_1" data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 116.79 122.88"><path class="cls-1" d="M75.4,74c-6.51,11.9-26,12.82-35,.82h0c2.81-1,4.86-2.43,5.53-4.32a19.85,19.85,0,0,0,11.3,3.83c4.6,0,9.15-1.91,13.24-5,.46,2,2.31,3.57,5,4.72ZM30.92,76.52C24.24,78,18,79.88,13.67,82.39c-10,5.85-9.45,19.77-11.86,31-.68,3.17-1.36,6.33-1.81,9.5H22.89V115.2a2.31,2.31,0,1,1,4.61,0v7.68H89.23V115.2a2.31,2.31,0,1,1,4.62,0v7.68h22.94c-.39-2.06-.84-4.09-1.18-6.19-1.88-11.73-1.89-28.4-13-34.23-4.94-2.6-12.08-4.63-19-6.42,10.25,1.28,23-.69,24.66-5.44-15.07-1.22-19.8-24.35-19.7-41.19C86.11,4.64,63.88-4.72,46.68,2.23,25,11,27.87,30.12,25.82,48.11,24.47,59.87,20.46,69.72,9.56,70.6c1.53,4.31,12,6.43,21.36,5.92ZM64.64,18.3c-10.12,10.24-20.79,14.25-32,17.09a6.74,6.74,0,0,0-.11,1.81c0,15.48,11.68,34.06,24.71,34.06S83.4,52.68,83.4,37.2a21.28,21.28,0,0,0-1.1-6.93c-9.84-.11-13.67-5.06-17.66-12Z"/></svg>
      <svg v-else id="Layer_1" data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.88 116.06"><path d="M0,116.06C3.81,84.35,2.87,87.21,40.09,77.24c4.77,21.43,41.49,21.59,42.21,0C121.73,87.81,120,84,122.88,116.06ZM28,22c16.1-19.9,34.67-30.73,48.61-13,16.8.88,23.53,25.55,10.14,36.79,0,.22,0,.44-.08.66a7.31,7.31,0,0,1,1.55-.3,4.86,4.86,0,0,1,2.74.52,3.87,3.87,0,0,1,1.91,2.31c.7,2.19-1.65,11.62-3.51,13.09a3.89,3.89,0,0,1-2.85.69L86,62.64c-.13,6.62-3.37,8.69-7.61,12.67-.55.52-1.12,1-1.63,1.55C72.42,81,67.73,83,63,83s-9.71-2.07-14.22-6.08c-.68-.61-1.14-1-1.59-1.39C43,72,39.52,69.94,38.86,62.66A5,5,0,0,1,36.05,62c-3.47-2-4.72-11.25-3.53-14.13.87-2.42,2.2-3.21,4-3.16C32.26,41.92,39.15,25.33,28,22ZM39.3,48.66c-6.82-2.38-4.58,12.13.54,10.6a1.5,1.5,0,0,1,.47-.08A1.59,1.59,0,0,1,42,60.73c.17,7.24,3.38,9,7.28,12.38.6.52,1.23,1,1.63,1.42,3.93,3.49,8,5.29,12.1,5.29s7.87-1.72,11.53-5.26c4.94-4.77,8.84-6.41,8.19-14.18a1.62,1.62,0,0,1,.26-1,1.6,1.6,0,0,1,2.22-.45c4.23,2.8,6-9.85,3.32-9.64-1.85.14-2.91,1.7-4.09,1.5a1.6,1.6,0,0,1-1.31-1.84c1.5-8.79.82-14.51-1.06-18.41a15.66,15.66,0,0,0-7.09-7c-10.49,8-20.09,4-28.17,10.74-3.42,2.85-5.2,7.06-3.91,13a1.6,1.6,0,0,1,0,1c-.63,1.73-2.35.83-3.51.42Z"/></svg>
    </div>
    <div v-if="!!patient">
      <div class="flex">
        <router-link
            :to="{name: 'ehr.patient.overview', params: {id: patient.ObjectID}}"
            class="text-2xl mb-2 mr-4 hover:cursor-pointer hover:underline"
            v-if="patient.surname || patient.firstName"
        >{{ patient.surname }}, {{ patient.firstName }}</router-link>
        <h1 v-else class="text-2xl mb-2 mr-4">Unknown patient</h1>
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
        <div><span class="text-sm font-bold">Birth date</span>: {{ patient.dob ? patient.dob : "unknown" }}</div>
        <div v-if="patient.email"><span class="text-sm font-bold">E-mail</span>: {{ patient.email }}</div>
        <div v-if="patient.zipcode"><span class="text-sm font-bold">Zipcode</span>: {{ patient.zipcode }}</div>
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