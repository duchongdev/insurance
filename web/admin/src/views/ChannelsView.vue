<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageCard from '@/components/PageCard.vue'
import PageToolbar from '@/components/PageToolbar.vue'
import StatusTag from '@/components/StatusTag.vue'
import { createChannel, deleteChannel, fetchChannels, type Channel } from '@/api/admin'
import { maskSecret } from '@/utils/auth'

const loading = ref(false)
const list = ref<Channel[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref<number | undefined>(undefined)
const keyword = ref('')

const createVisible = ref(false)
const creating = ref(false)
const form = reactive({
  channelCode: '',
  channelName: '',
  huaAnKey: '',
  channelKey: '',
  callbackUrl: '',
})

const revealedIds = ref<Set<number>>(new Set())

const filteredList = ref<Channel[]>([])

function resetForm() {
  form.channelCode = ''
  form.channelName = ''
  form.huaAnKey = ''
  form.channelKey = ''
  form.callbackUrl = ''
}

function openCreateDialog() {
  resetForm()
  createVisible.value = true
}

function isRevealed(id: number) {
  return revealedIds.value.has(id)
}

function toggleReveal(id: number) {
  const next = new Set(revealedIds.value)
  if (next.has(id)) {
    next.delete(id)
  } else {
    next.add(id)
  }
  revealedIds.value = next
}

function displaySecret(row: Channel, field: 'channelKey' | 'huaAnKey') {
  const value = row[field]
  if (!value) return '—'
  return isRevealed(row.id) ? value : maskSecret(value)
}

async function copySecret(value: string | undefined, label: string) {
  if (!value) {
    ElMessage.warning(`${label}为空`)
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success(`已复制${label}`)
  } catch {
    ElMessage.error('复制失败，请手动选择复制')
  }
}

function applyFilter() {
  let rows = list.value
  if (statusFilter.value !== undefined) {
    rows = rows.filter((r) => r.status === statusFilter.value)
  }
  const kw = keyword.value.trim().toLowerCase()
  if (kw) {
    rows = rows.filter(
      (r) =>
        r.channelCode.toLowerCase().includes(kw) ||
        (r.channelName && r.channelName.toLowerCase().includes(kw)),
    )
  }
  filteredList.value = rows
}

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchChannels(page.value, pageSize.value)
    list.value = data.list || []
    total.value = data.total || 0
    applyFilter()
  } finally {
    loading.value = false
  }
}

function onFilterChange() {
  applyFilter()
}

async function onCreate() {
  if (!form.channelCode.trim() || !form.huaAnKey.trim() || !form.callbackUrl.trim()) {
    ElMessage.warning('请填写渠道编码、华安密钥和回调地址')
    return
  }
  creating.value = true
  try {
    const { data } = await createChannel({
      channelCode: form.channelCode.trim(),
      channelName: form.channelName.trim(),
      huaAnKey: form.huaAnKey.trim(),
      channelKey: form.channelKey.trim() || undefined,
      callbackUrl: form.callbackUrl.trim(),
      status: 1,
    })
    createVisible.value = false
    const autoKey = !form.channelKey.trim()
    if (autoKey && data.channelKey) {
      await ElMessageBox.alert(
        `渠道编码：${data.channelCode}\n渠道密钥（已自动生成）：${data.channelKey}\n华安密钥：${data.huaAnKey}`,
        '创建成功，请妥善保存密钥',
        { confirmButtonText: '知道了' },
      )
    } else {
      ElMessage.success('创建成功')
    }
    resetForm()
    await loadData()
  } finally {
    creating.value = false
  }
}

async function onDelete(row: Channel) {
  await ElMessageBox.confirm(`确认删除渠道「${row.channelCode}」？`, '提示', { type: 'warning' })
  await deleteChannel(row.id)
  ElMessage.success('已删除')
  await loadData()
}

function onPageChange(p: number) {
  page.value = p
  loadData()
}

onMounted(loadData)
</script>

<template>
  <div>
    <PageToolbar title="渠道管理" subtitle="维护对接渠道编码、密钥与投保结果回调地址">
      <template #filters>
        <el-input
          v-model="keyword"
          placeholder="搜索编码/名称"
          clearable
          style="width: 200px"
          @input="onFilterChange"
          @clear="onFilterChange"
        />
        <el-select
          v-model="statusFilter"
          placeholder="全部状态"
          clearable
          style="width: 120px"
          @change="onFilterChange"
        >
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </template>
      <template #actions>
        <el-button type="primary" @click="openCreateDialog">新增渠道</el-button>
      </template>
    </PageToolbar>

    <PageCard>
      <template #header>
        <h3 class="card-title">渠道列表</h3>
      </template>
      <p class="hint">密钥在库内为明文存储，列表默认脱敏显示；点击「显示」可查看完整内容，「复制」可写入剪贴板。</p>
      <el-table v-loading="loading" :data="filteredList" size="default">
        <el-table-column prop="id" label="编号" width="72" />
        <el-table-column prop="channelCode" label="渠道编码" min-width="120" />
        <el-table-column prop="channelName" label="渠道名称" min-width="120">
          <template #default="{ row }">{{ row.channelName || '—' }}</template>
        </el-table-column>
        <el-table-column prop="callbackUrl" label="回调地址" min-width="200" show-overflow-tooltip />
        <el-table-column label="渠道密钥" min-width="240">
          <template #default="{ row }">
            <span class="secret-text">{{ displaySecret(row, 'channelKey') }}</span>
            <el-button type="primary" link @click="toggleReveal(row.id)">
              {{ isRevealed(row.id) ? '隐藏' : '显示' }}
            </el-button>
            <el-button type="primary" link @click="copySecret(row.channelKey, '渠道密钥')">复制</el-button>
          </template>
        </el-table-column>
        <el-table-column label="华安密钥" min-width="240">
          <template #default="{ row }">
            <span class="secret-text">{{ displaySecret(row, 'huaAnKey') }}</span>
            <el-button type="primary" link @click="toggleReveal(row.id)">
              {{ isRevealed(row.id) ? '隐藏' : '显示' }}
            </el-button>
            <el-button type="primary" link @click="copySecret(row.huaAnKey, '华安密钥')">复制</el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <StatusTag :status="row.status === 1 ? 'success' : 'default'" :label="row.status === 1 ? '启用' : '禁用'" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="danger" link @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
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
    </PageCard>
  </div>

  <el-dialog v-model="createVisible" title="新增渠道" width="520px" destroy-on-close @closed="resetForm">
    <el-form label-width="100px">
      <el-form-item label="渠道编码" required>
        <el-input v-model="form.channelCode" placeholder="channelCode" />
      </el-form-item>
      <el-form-item label="渠道名称">
        <el-input v-model="form.channelName" placeholder="可选" />
      </el-form-item>
      <el-form-item label="华安密钥" required>
        <el-input v-model="form.huaAnKey" placeholder="华安分配的 huaAnKey" show-password />
      </el-form-item>
      <el-form-item label="回调地址" required>
        <el-input v-model="form.callbackUrl" placeholder="https://渠道域名/投保结果回调" />
      </el-form-item>
      <el-form-item label="渠道密钥">
        <el-input v-model="form.channelKey" placeholder="留空则自动生成 32 位十六进制" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="createVisible = false">取消</el-button>
      <el-button type="primary" :loading="creating" @click="onCreate">确定</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.card-title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
}

.hint {
  margin: 0 0 16px;
  font-size: 13px;
  color: var(--color-text-muted);
}

.secret-text {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  margin-right: 4px;
  word-break: break-all;
  color: var(--color-text-secondary);
}

.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
