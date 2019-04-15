<template>
  <div class="services" @keyup.enter="fetch">
    <div class="actions">
      <nuxt-link to="/services/create">
        <el-button icon="el-icon-circle-plus">
          New Service
        </el-button>
      </nuxt-link>
    </div>
    <div class="table-actions">
      <el-select
        v-model="selected"
        class="filter-list"
        multiple
        filterable
        default-first-option
        placeholder="Filter services"
      >
        <el-option
          v-for="item in features"
          :key="item.value"
          :label="item.value"
          :value="item.id"
        />
      </el-select>
      <el-input v-model="pagination.search" class="search-input" placeholder="Search" />
      <el-button
        type="primary"
        icon="el-icon-search"
        @click="fetch"
      />
    </div>

    <div class="services-list">
      <el-table
        class="services-table"
        :data="services"
        style="width: 100%"
      >
        <el-table-column
          prop="name"
          label="Name"
          width="200"
        />
        <el-table-column
          prop="features"
          label="Features"
        >
          <template slot-scope="scope">
            <el-tag
              v-for="(feature, i) in scope.row.features"
              :key="i"
              class="tag"
              size="mini"
              type="warning"
            >
              {{ feature.value }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop=""
          width="70"
          label=""
        >
          <template slot-scope="scope">
            <nuxt-link class="icon-btn" :to="`/services/${scope.row.id}/edit`">
              <i class="el-icon-edit" />
            </nuxt-link>
            <nuxt-link class="icon-btn" to="#" @click.native="deleteTarget = scope.row">
              <i class="el-icon-delete" />
            </nuxt-link>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          layout="prev, pager, next, total"
          :total="total"
          :current-page.sync="pagination.page"
          :page-size="pagination.limit"
          @current-change="fetch"
        />
      </div>
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
import Service from '@/models/Service';
import Feature from '@/models/Feature';
import { mapState } from 'vuex';

export default {
    middleware: 'authenticated',
    data() {
        return {
            deleteTarget: null,
            features: [],
            selected: [],
            fetching: false,
            pagination: {
                page: 1,
                limit: 10,
                search: ''
            }
        };
    },
    computed: {
        ...mapState({
            services: state => state.services.list,
            total: state => state.services.total
        })
    },
    mounted() {
        this.fetch();
        this.$axios.get('/api/features')
            .then(res => this.features = res.data.map(p => new Feature(p)));
    },
    methods: {
        fetch() {
            this.fetching = true;
            this.$store.dispatch('services/fetch', {
                features: this.selected.join(','),
                pagination: this.pagination
            }).then(this.fetching = false);
        },
        remove() {
            this.$store.dispatch('services/delete', this.deleteTarget.id)
                .then(() => this.deleteTarget = null);
        }
    }
};
</script>

<style lang="scss" scoped>
.tag {
    margin: 0 2px;
}

.services-table {
    margin-top: 20px;
}

.cell {
    display: flex !important;
}

.features {
    margin: 50px 0 0 0;
    .el-select {
        width: 100%;
    }
}
.tag-preview {
    margin: 0 2px;
}
.actions {
    display: flex;
    justify-content: space-between;
}
.pagination {
    text-align: center;
    margin: 20px;
}
.table-actions {
    margin-top: 10px;
    display: flex;
    justify-content: space-between;
}

span.el-pagination__total{
    display: block !important;
}
.search-input {
    max-width: 250px;
    margin: 0 5px;
    max-width: 250px;
}
.features-list {
    flex-grow: 1;
}
</style>
