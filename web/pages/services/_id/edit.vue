<template>
  <div class="edit-service">
    <create-edit-service class="service" :service="service" @save="save" />
  </div>
</template>

<script>
import CreateEditService from '@/components/CreateEditService';
import Service from '@/models/Service';
import Feature from '@/models/Feature';

export default {
    middleware: 'authenticated',
    components: {
        CreateEditService
    },
    data() {
        return {
            service: null
        };
    },
    asyncData({ app, params }) {
        return app.$axios.get(`/api/services/${params.id}`)
            .then(res => ({ service: new Service(res.data) }));
    },
    methods: {
        save() {
            this.fixFeatures();

            this.$axios.post('/api/services', this.service)
                .then(() => this.$router.push({ path: '/services' }));
        },
        fixFeatures() {
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
