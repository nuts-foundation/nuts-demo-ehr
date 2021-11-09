<template>
  <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-md" v-if="errors.length">
    <label class="text-red-500">Please correct the following error{{
        errors.length === 0 ? '' : 's'
      }}:</label>

    <ul class="text-sm">
      <li v-for="error in errors">— {{ error }}</li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    errors: []
  }
}
</script>