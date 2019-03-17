<template>
  <div class="create-pattern">
    <create-edit-pattern class="pattern" :pattern="pattern" @save="save" />
  </div>
</template>

<script>
import CreateEditPattern from '@/components/CreateEditPattern';

export default {
    middleware: 'authenticated',
    components: {
        CreateEditPattern
    },
    data() {
        return {
        };
    },
    methods: {
        save() {
            this.$axios.post('/api/patterns', this.pattern)
                .then(() => this.$router.push({ path: '/patterns' }))
                .catch(console.warn);
        }
    }
};
</script>
