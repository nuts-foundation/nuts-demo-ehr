<template>
  <div class="inline-block relative">
    <div class="h-5 w-5 mr-3">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
        <path d="M0 0h24v24H0V0z" fill="none"/>
        <path
            d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5v-3h3.56c.69 1.19 1.97 2 3.45 2s2.75-.81 3.45-2H19v3zm0-5h-4.99c0 1.1-.9 2-2 2s-2-.9-2-2H5V5h14v9z"/>
      </svg>
    </div>

    <span class="count" v-if="!!inbox && inbox.messageCount > 0" :data-inbox-message-count="inbox.messageCount">{{ inbox.messageCount }}</span>
  </div>
</template>

<style>
.count {
  background-color: #fa3e3e;
  border-radius: 2px;
  color: white;

  padding: 1px 5px;
  font-size: 10px;
  line-height: normal;

  position: absolute; /* Position the badge within the relatively positioned button */
  top: -6px;
  right: 4px;
}
</style>
<script>
export default {
  name: 'inbox-badge',
  props: {inbox: Object},
}
</script>
