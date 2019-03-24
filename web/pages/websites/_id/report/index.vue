<template>
  <div class="report">
    <div class="center">
      <h2> {{ report.WebsiteURL }}</h2>
    </div>
    <el-form>
      <el-form-item label="Searched at">
        {{ report.CreatedAt }}
      </el-form-item>
      <el-form-item label="Enviroment loaded in">
        {{ report.StartedIn }}s
      </el-form-item>
      <el-form-item label="Page loaded in">
        {{ report.LoadedIn }}s
      </el-form-item>
      <el-form-item label="Resource patterns matched in">
        {{ report.ResourceCheckIn }}s
      </el-form-item>
      <el-form-item label="HTML patterns matched in">
        {{ report.HTMLCheckIn }}s
      </el-form-item>
      <el-form-item label="Total inpection in">
        {{ report.TotalIn }}s
      </el-form-item>
    </el-form>
    <div v-for="m in report.Matches" :key="m.ID" class="patterns">
      <div class="divider" />
      <div class="pattern-h text-desc">
        <el-tag type="info">
          {{ m.Pattern.Type }}
        </el-tag>
        {{ m.Pattern.Description }}
      </div>
      <el-form>
        <el-form-item label="Pattern">
          {{ m.Pattern.Value }}
        </el-form-item>
        <el-form-item label="Matched">
          {{ m.Value }}
        </el-form-item>
      </el-form>
      <div class="text-desc text-right" />
    </div>
    <div class="divider" />
  </div>
</template>

<script>
export default {
    data() {
        return {
            report: null
        };
    },
    asyncData({ app, params }) {
        return app.$axios.get(`/api/websites/${params.id}/report`)
            .then(res => ({ report: res.data }))
            .catch(console.error);
    }
};
</script>

<style lang="scss" scoped>
.report {
    max-width: 600px;
    margin: auto;
}
.center {
    margin: auto;
    text-align: center;
}
.el-form-item {
    margin: 0;
}
.pattern-h {
    display: flex;
    align-items: center;
    justify-content: space-between;
}
</style>
