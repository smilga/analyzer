<template>
  <div class="layout">
    <toolbar />
    <nuxt class="content" />
  </div>
</template>

<script>
import sidebar from '@/components/sidebar.vue';
import toolbar from '@/components/toolbar.vue';
import WSocket from '@/models/WSocket';
import mapState from 'vuex';

export default {
    components: {
        sidebar,
        toolbar
    },
    mounted() {
        if (this.$store.state.auth.user) {
            new WSocket({ store: this.$store, notify: this.$toasted });
        }
    }
};
</script>

<style lang="scss">
@import '@/assets/scss/_variables.scss';

body {
    margin: 0;
}

* {
    font-family: 'Product Sans';
}

.divider {
    margin: 10px 0;
    border-bottom: 1px solid #dcdfe6;
}

.text-desc {
    font-size: 14px;
    color: #606269;
}
.text-right {
    text-align: right;
}

.icon-btn {
    border-radius: 50%;
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: .3s;
    cursor: pointer;
    &:hover {
        background: $themeColor;
    }
}

.layout {
    min-height: 100vh;
}

.el-card__header {
    padding: 5px 15px;
}

.content {
    padding: 30px;
}

a, .nuxt-link-active, .nuxt-link, .nuxt-link-exact-active {
    text-decoration: none;
    color: inherit;
}

.toasted.bubble {
    padding: 14px 26px 14px 13px !important;
    border-radius: 8px !important;
    border: 1px solid #ebeef5 !important;
    background-color: #fff !important;
    box-shadow: 0 2px 12px 0 rgba(0,0,0,.1) !important;
    color: #606266 !important;
}

@font-face {
    font-family: 'Product Sans';
    src: url('../assets/fonts/ProductSansRegular.ttf') format('truetype');
}

</style>
