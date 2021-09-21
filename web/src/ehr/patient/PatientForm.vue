<template>
  <form class="space-y-5">
    <div>
      <label for="ssn">SSN</label>
      <input class="alt" id="ssn" :disabled="mode === 'edit'" type="text" v-model="value.ssn">
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="firstname">First name</label>
        <input class="alt" id="firstname" type="text" v-model="value.firstName">
      </div>

      <div>
        <label for="surname">Surname</label>
        <input class="alt" id="surname" type="text" v-model="value.surname">
      </div>
    </div>

    <div>
      <label for="gender">Gender</label>
      <div class="custom-select">
        <select class="alt" id="gender" name="" v-model="value.gender">
          <option value="unknown" selected>Unknown</option>
          <option value="male">Male</option>
          <option value="female">Female</option>
          <option value="other">Other</option>
        </select>

        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#444">
          <path d="M24 24H0V0h24v24z" fill="none" opacity=".87"/>
          <path d="M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6-1.41-1.41z"/>
        </svg>
      </div>
    </div>

    <div>
      <label for="dob">Date of birth</label>
      <input class="alt" id="dob" type="date" v-model="value.dob">
    </div>

    <div>
      <label for="zipcode">Zipcode</label>
      <input class="alt" id="zipcode" type="text" v-model="value.zipcode">
    </div>

    <div>
      <label for="email">Email</label>
      <input class="alt" id="email" type="email" v-model="value.email">
    </div>
  </form>
</template>
<script>
export default {
  props: {
    value: Object,
    mode: {
      type: String,
      default: 'new'
    }
  },
  emits: ['input'],
  watch: {
    value() {
      this.$emit('input', this.value)
    }
  }
}
</script>
