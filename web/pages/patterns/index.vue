<template>
  <div class="patterns">
    <nuxt-link to="/patterns/create">
      <el-button icon="el-icon-circle-plus" @click="">
        New Pattern
      </el-button>
    </nuxt-link>
    <el-table
      class="pattern-table"
      :data="patterns"
      style="width: 100%"
    >
      <el-table-column
        prop="type"
        label="Type"
        width="100"
      />
      <el-table-column
        prop="description"
        label="Description"
        width="200"
      />
      <el-table-column
        prop="tags"
        label="Tags"
      >
        <template slot-scope="scope">
          <el-tag
            v-for="(tag, i) in scope.row.tags"
            :key="i"
            class="tag"
            size="mini"
            type="info"
          >
            {{ tag.value }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop=""
        width="70"
        label=""
      >
        <template slot-scope="scope">
          <nuxt-link class="icon-btn" :to="`/patterns/${scope.row.id}/edit`">
            <i class="el-icon-edit" />
          </nuxt-link>
          <nuxt-link class="icon-btn" to="#" @click.native="deleteTarget = scope.row">
            <i class="el-icon-delete" />
          </nuxt-link>
        </template>
      </el-table-column>
    </el-table>

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
import Pattern from '@/models/Pattern';

export default {
    middleware: 'authenticated',
    data() {
        return {
            patterns: [],
            deleteTarget: null
        };
    },
    asyncData({ app }) {
        return app.$axios.get('/api/patterns')
            .then(res => ({ patterns: res.data.map(p => new Pattern(p)) }))
            .catch(console.error);
    },
    methods: {
        remove() {
            this.$axios(`/api/patterns/${this.deleteTarget.id}/delete`)
                .then(() => {
                    const target = this.patterns.indexOf(this.deleteTarget);
                    this.patterns.splice(target, 1);
                    this.deleteTarget = null;
                });
        }
    }
};
</script>

<style lang="scss" scoped>
.tag {
    margin: 0 2px;
}
</style>
