<template>
  <form class="space-y-3">
    <label>SSN:
      <input type="text" v-model="value.ssn"></label>
    <label>First name:
      <input type="text" v-model="value.firstName"></label>
    <label>Surname:
      <input type="text" v-model="value.surname"></label>
    <label>Gender:
      <select name="" id="" v-model="value.gender">
        <option value="unknown" selected>Unknown</option>
        <option value="male">Male</option>
        <option value="female">Female</option>
        <option value="other">Other</option>
      </select>
      </label>
    <label>Date of birth:
      <input type="date" v-model="value.dob"></label>
    <label>Zipcode:
      <input type="text" v-model="value.zipcode"></label>
    <label>Email:
      <input type="email" v-model="value.email"></label>

  </form>
</template>
<script>
export default {
  props: {
    value: Object,
  },
  emits: ['input'],
  watch: {
    value() {
      this.$emit('input', this.value)
    }
  }
}
</script>
