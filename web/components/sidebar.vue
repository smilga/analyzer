<template>
    <el-menu
        default-active="2"
        class="sidebar"
        @open="handleOpen"
        @close="handleClose"
        background-color="#ececec"
        text-color="#545c64"
        active-text-color="#000000">

        <nuxt-link to="/">
            <el-menu-item index="2">
                <i class="el-icon-menu"></i>
                <span>Websites</span>
            </el-menu-item>
        </nuxt-link>

        <nuxt-link to="/services">
            <el-menu-item index="4">
                <i class="el-icon-setting"></i>
                <span>Services</span>
            </el-menu-item>
        </nuxt-link>
    </el-menu>
</template>

<script>
export default {
    methods: {
        handleOpen(key, keyPath) {
            console.log(key, keyPath);
        },
        handleClose(key, keyPath) {
            console.log(key, keyPath);
        }
    }
}
</script>

<style lang="scss">
@import '@/assets/scss/_variables.scss';

.el-menu.sidebar {
    width: 220px;
    min-height: 100%;
    @include shadow;
}
</style>
