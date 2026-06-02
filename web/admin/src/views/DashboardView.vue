<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchBanks, fetchChannels } from '@/api/admin'

const loading = ref(false)
const channelCount = ref(0)
const bankCount = ref(0)

async function loadOverview() {
  loading.value = true
  try {
    const [channels, banks] = await Promise.all([
      fetchChannels(1, 1),
      fetchBanks(1, 1),
    ])
    channelCount.value = channels.data.total || 0
    bankCount.value = banks.data.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(loadOverview)
</script>

<template>
  <el-card shadow="never">
    <template #header>
      <span>概览</span>
    </template>
    <el-row :gutter="16" v-loading="loading">
      <el-col :span="8">
        <el-statistic title="渠道数量" :value="channelCount" />
      </el-col>
      <el-col :span="8">
        <el-statistic title="银行数量" :value="bankCount" />
      </el-col>
    </el-row>
    <p class="hint">本服务当前阶段仅做渠道 API 透明转发；渠道配置与银行列表可在左侧菜单维护。</p>
  </el-card>
</template>

<style scoped>
.hint {
  margin-top: 24px;
  color: #909399;
  font-size: 14px;
}
</style>
