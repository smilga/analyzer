<template>
    <div class="services">
        <div class="header">
            <nuxt-link class="add-btn" to="/services/create">
                <el-button icon="el-icon-circle-plus">Add Service</el-button>
            </nuxt-link>
        </div>
        <div class="service-list">
            <service-card v-for="(s, i) in services" @delete="deleteTarget = s" :key="i" :service="s"></service-card>
        </div>

        <el-dialog
            title="Delete confirmation"
            :visible="deleteTarget !== null"
            width="30%"
            >
            <span>Are you sure want to delete?</span>
            <span slot="footer" class="dialog-footer">
                <el-button @click="deleteTarget = null">Cancel</el-button>
                <el-button type="error" @click="remove">Confirm</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import ServiceCard from '@/components/ServiceCard';

export default {
    middleware: 'authenticated',
    components: {
        ServiceCard
    },
    asyncData ({ app }) {
        return app.$axios.get('/api/services')
            .then(res => ({ services: res.data }))
    },
    data() {
        return {
            services: [],
            deleteTarget: null,
        };
    },
    methods: {
        remove() {
            this.$axios(`/api/services/${this.deleteTarget.ID}/delete`)
                .then(() => {
                    const target = this.services.indexOf(this.deleteTarget);
                    this.services.splice(target, 1);
                    this.deleteTarget = null;
                })
        }
    },
}
</script>

<style lang="scss" scoped>
.service-list {
    display: flex;
    flex-wrap: wrap;
    .service-card {
        margin: 0 10px 10px 0;
    }
}
.header {
    display: flex;
    align-items: center;
}
.add-btn {
    margin-bottom: 20px;
}
</style>
