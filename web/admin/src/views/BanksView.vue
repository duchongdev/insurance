<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PageCard from '@/components/PageCard.vue'
import PageToolbar from '@/components/PageToolbar.vue'
import StatusTag from '@/components/StatusTag.vue'
import { fetchBanks, refreshBankList, type BankRecord } from '@/api/admin'

const DEFAULT_CHANNEL_CODE = 'I8p0Wn'
const DEFAULT_HUAAN_KEY = 'd158da0dbddf4bec8d4f703fcd6a82cb'

const loading = ref(false)
const list = ref<BankRecord[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref<number | undefined>(undefined)
const keyword = ref('')

const channelCode = ref(DEFAULT_CHANNEL_CODE)
const huaAnKey = ref(DEFAULT_HUAAN_KEY)
const fetching = ref(false)

const filteredList = ref<BankRecord[]>([])

function applyFilter() {
  let rows = list.value
  if (statusFilter.value !== undefined) {
    rows = rows.filter((r) => r.status === statusFilter.value)
  }
  const kw = keyword.value.trim().toLowerCase()
  if (kw) {
    rows = rows.filter(
      (r) => r.bankCode.toLowerCase().includes(kw) || r.bankName.toLowerCase().includes(kw),
    )
  }
  filteredList.value = rows
}

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchBanks(page.value, pageSize.value, statusFilter.value)
    list.value = data.list || []
    total.value = data.total || 0
    applyFilter()
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

function onKeywordChange() {
  applyFilter()
}

onMounted(loadData)
</script>

<template>
  <div>
    <PageToolbar title="银行列表" subtitle="华安 getBankList 全量同步与本地缓存">
      <template #filters>
        <el-input
          v-model="keyword"
          placeholder="搜索银行编码/名称"
          clearable
          style="width: 200px"
          @input="onKeywordChange"
          @clear="onKeywordChange"
        />
        <el-select
          v-model="statusFilter"
          placeholder="全部状态"
          clearable
          style="width: 120px"
          @change="onStatusChange"
        >
          <el-option label="启用" :value="1" />
          <el-option label="停用" :value="2" />
        </el-select>
      </template>
      <template #actions>
        <el-button type="primary" :loading="fetching" @click="onFetchBankList">获取银行列表</el-button>
      </template>
    </PageToolbar>

    <PageCard title="同步配置">
      <p class="desc">
        请求华安 getBankList，全量覆盖更新 bank_info_t 并刷新 Redis 缓存。下游渠道调用本服务时优先读缓存，未命中则读库，库中无数据再请求华安。
      </p>
      <el-form inline class="sync-form">
        <el-form-item label="渠道码">
          <el-input v-model="channelCode" placeholder="channelCode" style="width: 180px" />
        </el-form-item>
        <el-form-item label="华安密钥">
          <el-input v-model="huaAnKey" placeholder="huaAnKey" show-password style="width: 280px" />
        </el-form-item>
      </el-form>
    </PageCard>

    <PageCard title="银行数据">
      <el-table v-loading="loading" :data="filteredList" size="default">
        <el-table-column prop="id" label="编号" width="72" />
        <el-table-column prop="bankCode" label="银行编码" min-width="120" />
        <el-table-column prop="bankName" label="银行名称" min-width="160" />
        <el-table-column label="储蓄卡" width="90" align="center">
          <template #default="{ row }">
            <StatusTag
              :status="row.debitCard === 1 ? 'success' : 'default'"
              :label="row.debitCard === 1 ? '支持' : '不支持'"
            />
          </template>
        </el-table-column>
        <el-table-column label="信用卡" width="90" align="center">
          <template #default="{ row }">
            <StatusTag
              :status="row.creditCard === 1 ? 'success' : 'default'"
              :label="row.creditCard === 1 ? '支持' : '不支持'"
            />
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <StatusTag :status="row.status === 1 ? 'success' : 'danger'" :label="row.status === 1 ? '启用' : '停用'" />
          </template>
        </el-table-column>
        <el-table-column prop="updateTime" label="更新时间" min-width="170">
          <template #default="{ row }">
            <span class="text-muted">{{ row.updateTime || '—' }}</span>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="!loading && filteredList.length === 0" class="empty">暂无银行数据，请先点击「获取银行列表」</div>
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
    </PageCard>
  </div>
</template>

<style scoped>
.desc {
  color: var(--color-text-muted);
  font-size: 13px;
  margin: 0 0 16px;
  line-height: 1.6;
}

.sync-form {
  margin-bottom: 0;
}

.sync-form :deep(.el-form-item__label) {
  color: var(--color-text-secondary);
}

.text-muted {
  color: var(--color-text-muted);
  font-size: 13px;
}

.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.empty {
  text-align: center;
  color: var(--color-text-muted);
  padding: 24px 0 8px;
  font-size: 14px;
}
</style>
