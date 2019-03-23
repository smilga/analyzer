<template>
  <div class="layout">
    <toolbar />
    <nuxt class="content" />
  </div>
</template>

<script>
import sidebar from '@/components/sidebar.vue';
import toolbar from '@/components/toolbar.vue';

export default {
    components: {
        sidebar,
        toolbar
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

@font-face {
    font-family: 'Product Sans';
    src: url('../assets/fonts/ProductSansRegular.ttf') format('truetype');
}

</style>
