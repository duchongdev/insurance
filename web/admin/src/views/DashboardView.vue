<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchStats, type Stats } from '@/api/admin'

const loading = ref(false)
const channelCode = ref('')
const stats = ref<Stats>({
  policyCount: 0,
  userCount: 0,
  signCount: 0,
  apiLogCount: 0,
})

async function loadStats() {
  loading.value = true
  try {
    const { data } = await fetchStats(channelCode.value.trim() || undefined)
    stats.value = data
  } finally {
    loading.value = false
  }
}

onMounted(loadStats)
</script>

<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <span>数据概览</span>
        <div class="filters">
          <el-input v-model="channelCode" placeholder="渠道编码（可选）" clearable style="width: 200px" />
          <el-button type="primary" :loading="loading" @click="loadStats">查询</el-button>
        </div>
      </div>
    </template>
    <el-row :gutter="16" v-loading="loading">
      <el-col :span="6">
        <el-statistic title="保单" :value="stats.policyCount" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="用户" :value="stats.userCount" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="签约" :value="stats.signCount" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="接口日志" :value="stats.apiLogCount" />
      </el-col>
    </el-row>
  </el-card>
</template>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.filters {
  display: flex;
  gap: 8px;
}
</style>
