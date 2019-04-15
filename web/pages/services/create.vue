<template>
  <div class="create-service">
    <create-edit-service class="service" :service="service" @save="save" />
  </div>
</template>

<script>
import CreateEditService from '@/components/CreateEditService';
import Feature from '@/models/Feature';
import Service from '@/models/Service';

export default {
    middleware: 'authenticated',
    components: {
        CreateEditService
    },
    data() {
        return {
            service: new Service()
        };
    },
    methods: {
        save() {
            this.fixfeatures();

            this.$axios.post('/api/services', this.service)
                .then(() => this.$router.push({ path: '/services' }))
                .catch(console.warn);
        },
        fixfeatures() {
            this.service.features = this.service.features.map((t) => {
                if (typeof t === 'string') {
                    return new Feature({ value: t });
                }
                return t;
            });
        }
    }
};

</script>
