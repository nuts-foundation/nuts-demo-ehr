<template>
  <div class="inline-block relative">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
    </svg>
    <span class="count" v-if="!!inbox">{{ inbox.messageCount }}</span>
  </div>
</template>

<style>
.count {
  background-color: #fa3e3e;
  border-radius: 2px;
  color: white;

  padding: 1px 3px;
  font-size: 10px;
  line-height: normal;

  position: absolute; /* Position the badge within the relatively positioned button */
  top: -5px;
  right: 0;
}
</style>
<script>
export default {
  name: 'inbox-badge',
  props: {inbox: Object},
}
</script>