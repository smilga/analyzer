<template>
  <div class="create-edit-pattern">
    <el-form ref="form" :model="pattern" label-width="120px">
      <el-form-item label="Type">
        <el-select v-model="pattern.type" placeholder="js_source">
          <el-option
            v-for="t in Object.values(TYPE)"
            :key="t"
            :label="t"
            :value="t"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="Description">
        <el-input v-model="pattern.description" placeholder="Maxtraffic script" />
      </el-form-item>

      <el-form-item label="What to match">
        <el-input v-model="pattern.value" placeholder="*mt.js*" />
      </el-form-item>

      <el-form-item label="Tags">
        <el-select
          v-model="pattern.tags"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="Select tags"
        >
          <el-option
            v-for="item in tags"
            :key="item.value"
            :label="item.value"
            :value="item"
          />
        </el-select>
        <!--
        <el-tag
          v-for="tag in pattern.tags"
          :key="tag.id"
          closable
          :disable-transitions="false"
          @close="removeTag(tag)"
        >
          {{ tag.value }}
        </el-tag>
        <el-input
          v-if="newTag !== null"
          ref="saveTagInput"
          v-model="newTag.value"
          class="input-new-tag"
          size="mini"
          @keyup.enter.native="saveTag"
          @blur="saveTag"
        />
        <el-button v-else class="button-new-tag" size="small" @click="newTag = new Tag">
          + New Tag
        </el-button>
          -->
      </el-form-item>
    </el-form>

    <el-button icon="el-icon-plus" @click="$emit('save')">
      Save
    </el-button>
  </div>
</template>

<script>
import Pattern, { TYPE } from '@/models/Pattern';
import Tag from '@/models/Tag';

export default {
    props: {
        pattern: {
            required: true,
            type: Object
        }
    },
    data() {
        return {
            newTag: null,
            tags: [],
            TYPE
        };
    },
    mounted() {
        this.$axios.get('/api/tags/')
            .then(res => this.tags = res.data.map(t => new Tag(t)));
    },
    methods: {
        translate(type) {
            switch (type) {
            case TYPE.RESOURCE:
                return 'Page loads resource from url';
            case TYPE.JS_SOURCE:
                return 'Page scripts contains';
            case TYPE.HTML:
                return 'Page has HTML tag';
            default:
                return 'Err no name';
            }
        },
        placeholderByType(type) {
            switch (type) {
            case TYPE.RESOURCE:
                return '*mt.js (accepts wildcard)';
            case TYPE.JS_SOURCE:
                return '*jQuery* (accepts wildcard)';
            case TYPE.HTML:
                return 'div.opt-in or div#someID (use CSS selectors)';
            default:
                return 'Err no name';
            }
        },
        saveTag() {
            this.$axios.post('/api/tags', this.newTag)
                .then(res => this.patter.tags.push(new Tag(res.data)));
        }
    }
};
</script>

<style lang="scss" scoped>
@import '@/assets/scss/_variables.scss';

.create-edit-pattern {
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
</style>
