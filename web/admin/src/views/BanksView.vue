<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchBanks, refreshBankList, type BankRecord } from '@/api/admin'

const DEFAULT_CHANNEL_CODE = 'I8p0Wn'
const DEFAULT_HUAAN_KEY = 'd158da0dbddf4bec8d4f703fcd6a82cb'

const loading = ref(false)
const list = ref<BankRecord[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref<number | undefined>(undefined)

const channelCode = ref(DEFAULT_CHANNEL_CODE)
const huaAnKey = ref(DEFAULT_HUAAN_KEY)
const fetching = ref(false)

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchBanks(page.value, pageSize.value, statusFilter.value)
    list.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

async function onFetchBankList() {
  const code = channelCode.value.trim()
  const key = huaAnKey.value.trim()
  if (!code || !key) {
    ElMessage.warning('请填写渠道码和华安密钥')
    return
  }
  fetching.value = true
  try {
    await refreshBankList(code, key)
    ElMessage.success('银行列表已获取并更新')
    page.value = 1
    await loadData()
  } finally {
    fetching.value = false
  }
}

function onPageChange(p: number) {
  page.value = p
  loadData()
}

function onStatusChange() {
  page.value = 1
  loadData()
}

onMounted(loadData)
</script>

<template>
  <el-card shadow="never">
    <template #header>获取银行列表</template>
    <p class="desc">
      请求华安 getBankList，全量覆盖更新 bank_info_t 并刷新 Redis 缓存。下游渠道调用本服务时优先读缓存，未命中则读库，库中无数据再请求华安。
    </p>
    <el-form label-width="100px" style="max-width: 560px">
      <el-form-item label="渠道码">
        <el-input v-model="channelCode" placeholder="channelCode" />
      </el-form-item>
      <el-form-item label="华安密钥">
        <el-input v-model="huaAnKey" placeholder="huaAnKey" show-password />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="fetching" @click="onFetchBankList">获取银行列表</el-button>
      </el-form-item>
    </el-form>
  </el-card>

  <el-card shadow="never" style="margin-top: 16px" v-loading="loading">
    <template #header>
      <div class="header-row">
        <span>银行列表</span>
        <el-select
          v-model="statusFilter"
          placeholder="全部状态"
          clearable
          style="width: 140px"
          @change="onStatusChange"
        >
          <el-option label="启用" :value="1" />
          <el-option label="停用" :value="2" />
        </el-select>
      </div>
    </template>
    <el-table :data="list" stripe border>
      <el-table-column prop="id" label="编号" width="80" />
      <el-table-column prop="bankCode" label="银行编码" min-width="120" />
      <el-table-column prop="bankName" label="银行名称" min-width="160" />
      <el-table-column label="储蓄卡" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.debitCard === 1 ? 'success' : 'info'" size="small">
            {{ row.debitCard === 1 ? '支持' : '不支持' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="信用卡" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.creditCard === 1 ? 'success' : 'info'" size="small">
            {{ row.creditCard === 1 ? '支持' : '不支持' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updateTime" label="更新时间" min-width="170" />
    </el-table>
    <div v-if="!loading && list.length === 0" class="empty">暂无银行数据，请先点击「获取银行列表」</div>
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
.desc {
  color: #909399;
  font-size: 14px;
  margin: 0 0 16px;
}
.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.empty {
  text-align: center;
  color: #909399;
  padding: 24px 0 8px;
  font-size: 14px;
}
</style>
