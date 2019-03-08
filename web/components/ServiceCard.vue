<template>
    <el-card class="service-card" shadow="hover" :body-style="{ padding: '0px' }">
        <div class="cont">
            <img :src="service.LogoURL" class="image logo">
        </div>
        <div slot="header" class="clearfix card-header">
            <span class="text"><h4>{{ service.Name }}</h4></span>
            <span class="action">
                <nuxt-link class="icon-btn" :to="`/services/${service.ID}/edit`">
                    <i class="el-icon-edit"></i>
                </nuxt-link>
                <nuxt-link  class="icon-btn" to="#"  @click.native="$emit('delete')">
                    <i class="el-icon-delete"></i>
                </nuxt-link>
            </span>
        </div>
    </el-card>
</template>

<script>
export default {
    props: {
        service: {
            required: true,
            type: Object
        }
    }
}
</script>

<style lang="scss" scoped>
@import '@/assets/scss/_variables.scss';

.service-card {
    width: 300px;
}

.card-header {
    display:flex;
    align-items: center;
    justify-content: space-between;
    .text {
        h4 {
            text-transform: uppercase;
            margin: 10px 0;
        }
    }
    .action {
        display: flex;
        a {
            margin: 0 5px;
        }
    }

}
.el-card__header {
    padding: 5px 15px;
}
.logo {
    width: 200px;
    height: 80px;
    object-fit: contain;
}
.cont {
    padding: 0;
    display:flex;
    align-items: center;
    justify-content: center;
}
</style>
