<template>
    <div class="edit-service">
        <create-edit-service @save="save" class="service" :service="service"></create-edit-service>
    </div>
</template>

<script>
import CreateEditService from '@/components/CreateEditService';
import Service from '@/models/Service';

export default {
    middleware: 'authenticated',
    components: {
        CreateEditService
    },
    asyncData({ app, params }) {
        return app.$axios.get(`/api/services/${params.id}`)
            .then(res => ({ service: new Service(res.data) }))
            .catch(err => console.log(err))
    },
    data() {
        return {
            service: null
        }
    },
    methods: {
        save() {
            this.$axios.post('/api/services', this.service)
                .then(() => this.$router.push({path: '/services'}))
                .catch(console.warn);
        }
    },
}

</script>

