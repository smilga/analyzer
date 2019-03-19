<template>
  <div class="create-pattern">
    <create-edit-pattern class="pattern" :pattern="pattern" @save="save" />
  </div>
</template>

<script>
import CreateEditPattern from '@/components/CreateEditPattern';
import Tag from '@/models/Tag';
import Pattern, { TYPE } from '@/models/Pattern';

export default {
    middleware: 'authenticated',
    components: {
        CreateEditPattern
    },
    data() {
        return {
            pattern: new Pattern({ Type: TYPE.RESOURCE })
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
            this.pattern.tags = this.pattern.tags.map((t) => {
                if (typeof t === 'string') {
                    return new Tag({ Value: t });
                }
                return t;
            });
        }
    }
};
</script>
