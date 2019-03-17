<template>
  <div class="create-pattern">
    <create-edit-pattern class="pattern" :pattern="pattern" @save="save" />
  </div>
</template>

<script>
import CreateEditPattern from '@/components/CreateEditPattern';
import Tag from '@/models/Tag';

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
            this.fixTags();

            this.$axios.post('/api/patterns', this.pattern)
                .then(() => this.$router.push({ path: '/patterns' }))
                .catch(console.warn);
        },
        fixTags() {
            this.pattern.tags.map((t) => {
                if (typeof t === 'string') {
                    return new Tag({ Value: t });
                }
            });
        }
    }
};
</script>
