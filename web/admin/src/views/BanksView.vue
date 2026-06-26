<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PageCard from '@/components/PageCard.vue'
import PageToolbar from '@/components/PageToolbar.vue'
import StatusTag from '@/components/StatusTag.vue'
import { fetchBanks, fetchHuaAnConfigs, refreshBankList, type BankRecord, type HuaAnConfig } from '@/api/admin'

const loading = ref(false)
const list = ref<BankRecord[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref<number | undefined>(undefined)
const keyword = ref('')

const huaAnConfigs = ref<HuaAnConfig[]>([])
const selectedSettingId = ref<number | undefined>(undefined)
const fetching = ref(false)

const filteredList = ref<BankRecord[]>([])

const huaAnOptions = computed(() =>
  huaAnConfigs.value
    .filter((c) => c.id)
    .map((c) => ({
      value: c.id!,
      label: `${c.envType === 'prod' ? '生产环境' : '测试环境'} · ${c.channelCode}${c.name ? ` · ${c.name}` : ''}`,
    })),
)

const selectedConfig = computed(() =>
  huaAnConfigs.value.find((c) => c.id === selectedSettingId.value),
)

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

async function loadHuaAnConfigs() {
  try {
    const { data } = await fetchHuaAnConfigs()
    huaAnConfigs.value = data.list || []
    if (!selectedSettingId.value && huaAnConfigs.value.length > 0) {
      selectedSettingId.value = huaAnConfigs.value[0].id
    }
  } catch {
    huaAnConfigs.value = []
  }
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
  if (!selectedSettingId.value) {
    ElMessage.warning('请先在「华安配置」中添加配置，并选择要使用的华安配置')
    return
  }
  fetching.value = true
  try {
    await refreshBankList(selectedSettingId.value)
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

onMounted(async () => {
  await loadHuaAnConfigs()
  await loadData()
})
</script>

<template>
  <div>
    <PageToolbar title="银行列表" subtitle="华安 getBankList 全量同步与本地缓存">
      <template #filters>
        <el-select
          v-model="selectedSettingId"
          placeholder="选择华安配置"
          style="width: 280px"
          filterable
        >
          <el-option
            v-for="opt in huaAnOptions"
            :key="opt.value"
            :label="opt.label"
            :value="opt.value"
          />
        </el-select>
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

    <PageCard title="同步说明">
      <p class="desc">
        请求华安 getBankList，全量覆盖更新 bank_info_t 并刷新 Redis 缓存。请在上方选择
        <router-link to="/huaan-config">华安配置</router-link>
        后再点击「获取银行列表」；当前选择：
        <strong v-if="selectedConfig">
          {{ selectedConfig.envType === 'prod' ? '生产环境' : '测试环境' }} · {{ selectedConfig.channelCode }}
        </strong>
        <strong v-else>未选择</strong>。
      </p>
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
      <div v-if="!loading && filteredList.length === 0" class="empty">暂无银行数据，请先选择华安配置并点击「获取银行列表」</div>
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
  margin: 0;
  line-height: 1.6;
}

.desc a {
  color: var(--color-primary);
  text-decoration: none;
}

.desc a:hover {
  text-decoration: underline;
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
