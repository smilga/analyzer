<template>
  <div class="create-edit-service">
    <el-form ref="form" :model="service" label-width="120px">
      <el-form-item label="Name">
        <el-input v-model="service.name" placeholder="Maxtraffic" />
      </el-form-item>

      <el-form-item label="Features">
        <el-select
          v-model="service.features"
          class="features-list"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="Select features"
        >
          <el-option
            v-for="item in features"
            :key="item.value"
            :label="item.value"
            :value="item"
          />
        </el-select>
      </el-form-item>
    </el-form>

    <el-button icon="el-icon-plus" @click="$emit('save')">
      Save
    </el-button>
  </div>
</template>

<script>
import Service from '@/models/Service';
import Feature from '@/models/Feature';

export default {
    props: {
        service: {
            required: true,
            type: Object
        }
    },
    data() {
        return {
            newFeature: null,
            features: []
        };
    },
    mounted() {
        this.$axios.get('/api/features')
            .then(res => this.features = res.data.map(t => new Feature(t)));
    },
    methods: {
        saveFeature() {
            this.$axios.post('/api/features', this.newFeature)
                .then(res => this.service.features.push(new Feature(res.data)));
        }
    }
};
</script>

<style lang="scss" scoped>
@import '@/assets/scss/_variables.scss';

.create-edit-service {
    max-width: 600px;
    margin: auto;
}

.el-tag + .el-tag {
    margin-left: 10px;
}
.button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
}
.input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
}
.features-list {
    width: 100%;
}
</style>
