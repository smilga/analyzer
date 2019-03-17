<template>
  <div class="edit-pattern">
    <create-edit-pattern class="pattern" :pattern="pattern" @save="save" />
  </div>
</template>

<script>
import CreateEditPattern from '@/components/CreateEditPattern';
import Pattern from '@/models/Pattern';

export default {
    middleware: 'authenticated',
    components: {
        CreateEditPattern
    },
    data() {
        return {
            pattern: null
        };
    },
    asyncData({ app, params }) {
        return app.$axios.get(`/api/patterns/${params.id}`)
            .then(res => ({ pattern: new Pattern(res.data) }));
    },
    methods: {
        save() {
            this.$axios.post('/api/patterns', this.pattern)
                .then(() => this.$router.push({ path: '/patterns' }));
        }
    }
};

</script>
