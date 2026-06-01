<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchLogs, type LogRecord } from '@/api/admin'

const loading = ref(false)
const list = ref<LogRecord[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const channelCode = ref('')
const apiPath = ref('')

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchLogs(
      page.value,
      pageSize.value,
      channelCode.value.trim() || undefined,
      apiPath.value.trim() || undefined,
    )
    list.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function onPageChange(p: number) {
  page.value = p
  loadData()
}

function onSearch() {
  page.value = 1
  loadData()
}

onMounted(loadData)
</script>

<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <span>接口日志</span>
        <div class="filters">
          <el-input v-model="channelCode" placeholder="渠道编码" clearable style="width: 160px" />
          <el-input v-model="apiPath" placeholder="接口路径" clearable style="width: 160px" />
          <el-button type="primary" :loading="loading" @click="onSearch">查询</el-button>
        </div>
      </div>
    </template>
    <el-table v-loading="loading" :data="list" stripe border>
      <el-table-column prop="traceId" label="链路ID" min-width="200" />
      <el-table-column prop="channelCode" label="渠道编码" min-width="120" />
      <el-table-column prop="apiPath" label="接口路径" min-width="160" />
      <el-table-column prop="huaanCode" label="响应码" width="100" />
      <el-table-column prop="durationMS" label="耗时(ms)" width="110" />
      <el-table-column prop="createdAt" label="创建时间" min-width="180" />
    </el-table>
    <div class="pager">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="onPageChange"
      />
    </div>
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
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
