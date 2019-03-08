<template>
    <div class="create-service">
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
    data() {
        return {
            service: new Service({Name: 'Untitled'}),
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
