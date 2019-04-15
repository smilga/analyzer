<template>
  <div class="stats">
    <span class="stats">Queue: {{ queued }} websites</span>
    <span class="stats">Timeouted: {{ timeouted }} websites</span>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
    middleware: 'authenticated',
    computed: {
        ...mapState({
            queued: state => state.websites.queued,
            timeouted: state => state.websites.timeouted
        })
    }
};
</script>
