<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchPolicies, type PolicyRecord } from '@/api/admin'

const loading = ref(false)
const list = ref<PolicyRecord[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const channelCode = ref('')

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchPolicies(page.value, pageSize.value, channelCode.value.trim() || undefined)
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
        <span>保单列表</span>
        <div class="filters">
          <el-input v-model="channelCode" placeholder="渠道编码" clearable style="width: 180px" />
          <el-button type="primary" :loading="loading" @click="onSearch">查询</el-button>
        </div>
      </div>
    </template>
    <el-table v-loading="loading" :data="list" stripe border>
      <el-table-column prop="policyId" label="policyId" min-width="160" />
      <el-table-column prop="channelCode" label="channelCode" min-width="120" />
      <el-table-column prop="policyStatus" label="policyStatus" width="120" />
      <el-table-column prop="productCode" label="productCode" min-width="120" />
      <el-table-column prop="policyNo" label="policyNo" min-width="140" />
      <el-table-column prop="createdAt" label="createdAt" min-width="180" />
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
