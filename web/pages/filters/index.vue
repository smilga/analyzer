<template>
  <div class="filters">
    <el-button icon="el-icon-circle-plus" @click="filter = new Filter()">
      Add filter
    </el-button>
    <el-table
      :data="filters"
      style="width: 100%"
    >
      <el-table-column
        prop="name"
        label="Name"
        width="200"
      />
      <el-table-column
        prop="tags"
        label="Tags"
      >
        <template slot-scope="scope">
          <el-select v-model="scope.row.tags" class="tags-list" multiple placeholder="Tags...">
            <el-option
              v-for="item in tags"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      title="Add Filter"
      :visible="filter !== null"
      width="30%"
    >
      <span>Filter name</span>
      <el-input v-if="filter" v-model="filter.name" placeholder="Marketing tools" />
      <span slot="footer" class="dialog-footer">
        <el-button @click="filter = null">Cancel</el-button>
        <el-button @click="saveFilter">Confirm</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import Tag from '@/models/Tag';
import Filter from '@/models/Filter';

export default {
    data() {
        return {
            tags: [],
            filters: [],
            filter: null,
            Filter
        };
    },
    mounted() {
        this.$axios('/api/filters')
            .then(res => this.filters = res.data.map(f => new Filter(f)));

        this.$axios('/api/tags')
            .then(res => this.tags = res.data.map(t => new Tag(t)));
    },
    methods: {
        saveFilter() {
            this.$axios.post('/api/filters', this.filter)
                .then((res) => {
                    this.filters.push(new Filter(res.data));
                    this.filter = null;
                })
                .catch((e) => {
                    this.$notify.error({
                        title: 'Error',
                        message: e.response.data.error,
                        position: 'bottom-right'
                    });
                });
        }
    }
};
</script>

<style lang="scss" scoped>
.tags-list {
    width: 100%;
}
</style>
