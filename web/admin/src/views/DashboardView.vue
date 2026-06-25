<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import PageCard from '@/components/PageCard.vue'
import PageToolbar from '@/components/PageToolbar.vue'
import StatCard from '@/components/StatCard.vue'
import StatusTag from '@/components/StatusTag.vue'
import { fetchBanks, fetchChannels, type BankRecord, type Channel } from '@/api/admin'

const router = useRouter()
const loading = ref(false)

const channelTotal = ref(0)
const bankTotal = ref(0)
const activeChannels = ref(0)
const activeBanks = ref(0)
const recentChannels = ref<Channel[]>([])
const recentBanks = ref<BankRecord[]>([])

const todos = computed(() => {
  const items: { text: string; status: 'success' | 'warning' | 'danger' | 'default'; action?: () => void }[] = []
  if (bankTotal.value === 0) {
    items.push({
      text: '银行列表为空，请前往银行列表页获取数据',
      status: 'warning',
      action: () => router.push('/banks'),
    })
  } else {
    items.push({ text: '银行列表数据已就绪', status: 'success' })
  }
  if (channelTotal.value === 0) {
    items.push({
      text: '尚未配置渠道，请新增渠道后再对接',
      status: 'warning',
      action: () => router.push('/channels'),
    })
  } else {
    items.push({ text: `已配置 ${channelTotal.value} 个渠道`, status: 'success' })
  }
  items.push({ text: '本服务当前阶段仅做渠道 API 透明转发', status: 'default' })
  return items
})

async function loadOverview() {
  loading.value = true
  try {
    const [channelRes, bankRes] = await Promise.all([fetchChannels(1, 1), fetchBanks(1, 1)])
    channelTotal.value = channelRes.data.total || 0
    bankTotal.value = bankRes.data.total || 0

    const channelSize = Math.min(Math.max(channelTotal.value, 1), 500)
    const bankSize = Math.min(Math.max(bankTotal.value, 1), 500)
    const [channelListRes, bankListRes] = await Promise.all([
      fetchChannels(1, channelSize),
      fetchBanks(1, bankSize),
    ])

    const channels = channelListRes.data.list || []
    const banks = bankListRes.data.list || []
    recentChannels.value = channels.slice(0, 5)
    recentBanks.value = banks.slice(0, 5)
    activeChannels.value = channels.filter((c) => c.status === 1).length
    activeBanks.value = banks.filter((b) => b.status === 1).length
  } finally {
    loading.value = false
  }
}

onMounted(loadOverview)
</script>

<template>
  <div v-loading="loading">
    <PageToolbar title="数据概览" subtitle="渠道与银行配置运行状态">
      <template #actions>
        <el-button type="primary" @click="loadOverview">刷新</el-button>
      </template>
    </PageToolbar>

    <el-row :gutter="16" class="stat-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard title="渠道总数" :value="channelTotal" trend="neutral" trend-text="对接渠道配置" />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard title="银行总数" :value="bankTotal" trend="neutral" trend-text="华安银行列表缓存" />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          title="启用渠道"
          :value="activeChannels"
          trend="up"
          :trend-text="channelTotal ? `${Math.round((activeChannels / channelTotal) * 100)}% 启用率` : '暂无渠道'"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          title="启用银行"
          :value="activeBanks"
          trend="up"
          :trend-text="bankTotal ? `${Math.round((activeBanks / bankTotal) * 100)}% 启用率` : '暂无银行'"
        />
      </el-col>
    </el-row>

    <PageCard title="待办事项">
      <ul class="todo-list">
        <li v-for="(item, idx) in todos" :key="idx" class="todo-item">
          <StatusTag :status="item.status" :label="item.status === 'success' ? '正常' : item.status === 'warning' ? '待处理' : '说明'" />
          <span class="todo-item__text">{{ item.text }}</span>
          <el-button v-if="item.action" type="primary" link @click="item.action">去处理</el-button>
        </li>
      </ul>
    </PageCard>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <PageCard title="最近渠道">
          <el-table :data="recentChannels" size="default" empty-text="暂无渠道数据">
            <el-table-column prop="channelCode" label="渠道编码" min-width="110" />
            <el-table-column prop="channelName" label="渠道名称" min-width="100">
              <template #default="{ row }">{{ row.channelName || '—' }}</template>
            </el-table-column>
            <el-table-column label="状态" width="80" align="center">
              <template #default="{ row }">
                <StatusTag :status="row.status === 1 ? 'success' : 'default'" :label="row.status === 1 ? '启用' : '禁用'" />
              </template>
            </el-table-column>
          </el-table>
          <div class="card-footer">
            <el-button type="primary" link @click="router.push('/channels')">查看全部</el-button>
          </div>
        </PageCard>
      </el-col>
      <el-col :xs="24" :lg="12">
        <PageCard title="最近银行">
          <el-table :data="recentBanks" size="default" empty-text="暂无银行数据">
            <el-table-column prop="bankCode" label="银行编码" min-width="110" />
            <el-table-column prop="bankName" label="银行名称" min-width="120" />
            <el-table-column label="状态" width="80" align="center">
              <template #default="{ row }">
                <StatusTag :status="row.status === 1 ? 'success' : 'danger'" :label="row.status === 1 ? '启用' : '停用'" />
              </template>
            </el-table-column>
          </el-table>
          <div class="card-footer">
            <el-button type="primary" link @click="router.push('/banks')">查看全部</el-button>
          </div>
        </PageCard>
      </el-col>
    </el-row>

    <PageCard title="数据统计">
      <el-row :gutter="24">
        <el-col :xs="24" :md="8">
          <div class="panel-item">
            <div class="panel-item__label">渠道启用率</div>
            <div class="panel-item__bar">
              <div
                class="panel-item__fill panel-item__fill--primary"
                :style="{ width: channelTotal ? `${(activeChannels / channelTotal) * 100}%` : '0%' }"
              />
            </div>
            <div class="panel-item__meta">{{ activeChannels }} / {{ channelTotal }}</div>
          </div>
        </el-col>
        <el-col :xs="24" :md="8">
          <div class="panel-item">
            <div class="panel-item__label">银行启用率</div>
            <div class="panel-item__bar">
              <div
                class="panel-item__fill panel-item__fill--success"
                :style="{ width: bankTotal ? `${(activeBanks / bankTotal) * 100}%` : '0%' }"
              />
            </div>
            <div class="panel-item__meta">{{ activeBanks }} / {{ bankTotal }}</div>
          </div>
        </el-col>
        <el-col :xs="24" :md="8">
          <div class="panel-item">
            <div class="panel-item__label">服务状态</div>
            <div class="panel-item__status">
              <StatusTag status="success" label="运行中" />
              <span class="panel-item__hint">API 透明转发模式</span>
            </div>
          </div>
        </el-col>
      </el-row>
    </PageCard>
  </div>
</template>

<style scoped>
.stat-row {
  margin-bottom: 16px;
}

.stat-row .el-col {
  margin-bottom: 16px;
}

.todo-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.todo-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid var(--color-border);
}

.todo-item:last-child {
  border-bottom: none;
}

.todo-item__text {
  flex: 1;
  color: var(--color-text-primary);
}

.card-footer {
  margin-top: 12px;
  text-align: right;
}

.panel-item {
  padding: 8px 0;
}

.panel-item__label {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-bottom: 10px;
}

.panel-item__bar {
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.panel-item__fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.2s ease;
}

.panel-item__fill--primary {
  background: var(--color-primary);
}

.panel-item__fill--success {
  background: var(--color-success);
}

.panel-item__meta {
  margin-top: 8px;
  font-size: 13px;
  color: var(--color-text-muted);
}

.panel-item__status {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 4px;
}

.panel-item__hint {
  font-size: 13px;
  color: var(--color-text-muted);
}
</style>
