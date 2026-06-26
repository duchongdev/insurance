<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageCard from '@/components/PageCard.vue'
import PageToolbar from '@/components/PageToolbar.vue'
import StatusTag from '@/components/StatusTag.vue'
import {
  createChannel,
  deleteChannel,
  fetchChannels,
  fetchHuaAnConfigs,
  updateChannel,
  type Channel,
  type HuaAnConfig,
} from '@/api/admin'
import { maskSecret } from '@/utils/auth'

const loading = ref(false)
const list = ref<Channel[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref<number | undefined>(undefined)
const keyword = ref('')

const huaAnConfigs = ref<HuaAnConfig[]>([])
const dialogVisible = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const form = reactive({
  huaanSettingId: undefined as number | undefined,
  channelCode: '',
  channelName: '',
  channelKey: '',
  callbackUrl: '',
  status: 1,
})

const revealedIds = ref<Set<number>>(new Set())
const filteredList = ref<Channel[]>([])

const linkedHuaAnSettingIds = ref<Set<number>>(new Set())

const huaAnOptions = computed(() =>
  huaAnConfigs.value
    .filter((c) => c.id && (!editingId.value ? !linkedHuaAnSettingIds.value.has(c.id!) : true))
    .map((c) => ({
      value: c.id!,
      label: `${c.channelCode} · ${c.envType === 'prod' ? '生产环境' : '测试环境'}${c.name ? ` · ${c.name}` : ''}`,
      channelCode: c.channelCode,
      envType: c.envType,
    })),
)

function envLabel(envType?: string) {
  return envType === 'prod' ? '生产环境' : '测试环境'
}

function resetForm() {
  form.huaanSettingId = undefined
  form.channelCode = ''
  form.channelName = ''
  form.channelKey = ''
  form.callbackUrl = ''
  form.status = 1
  editingId.value = null
}

function openCreateDialog() {
  resetForm()
  dialogVisible.value = true
  void loadLinkedHuaAnSettings()
}

function openEditDialog(row: Channel) {
  editingId.value = row.id
  form.huaanSettingId = row.huaanSettingId
  form.channelCode = row.channelCode
  form.channelName = row.channelName || ''
  form.channelKey = ''
  form.callbackUrl = row.callbackUrl
  form.status = row.status
  dialogVisible.value = true
}

function isRevealed(id: number) {
  return revealedIds.value.has(id)
}

function toggleReveal(id: number) {
  const next = new Set(revealedIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
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

async function loadLinkedHuaAnSettings() {
  try {
    const { data } = await fetchChannels(1, 200)
    linkedHuaAnSettingIds.value = new Set(
      (data.list || []).filter((c) => c.huaanSettingId).map((c) => c.huaanSettingId!),
    )
  } catch {
    linkedHuaAnSettingIds.value = new Set()
  }
}

async function loadHuaAnConfigs() {
  try {
    const { data } = await fetchHuaAnConfigs()
    huaAnConfigs.value = data.list || []
  } catch {
    huaAnConfigs.value = []
  }
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

async function onSave() {
  if (!form.callbackUrl.trim()) {
    ElMessage.warning('请填写回调地址')
    return
  }
  if (!editingId.value && !form.huaanSettingId) {
    ElMessage.warning('请选择华安配置中的渠道编码')
    return
  }

  saving.value = true
  try {
    if (editingId.value) {
      await updateChannel(editingId.value, {
        channelName: form.channelName.trim(),
        channelKey: form.channelKey.trim() || undefined,
        callbackUrl: form.callbackUrl.trim(),
        status: form.status,
      })
      ElMessage.success('已保存')
    } else {
      const { data } = await createChannel({
        huaanSettingId: form.huaanSettingId!,
        channelName: form.channelName.trim() || undefined,
        channelKey: form.channelKey.trim() || undefined,
        callbackUrl: form.callbackUrl.trim(),
        status: form.status,
      })
      dialogVisible.value = false
      const autoKey = !form.channelKey.trim()
      if (autoKey && data.channelKey) {
        await ElMessageBox.alert(
          `渠道编码：${data.channelCode}\n渠道密钥（已自动生成）：${data.channelKey}`,
          '创建成功，请妥善保存渠道密钥',
          { confirmButtonText: '知道了' },
        )
      } else {
        ElMessage.success('创建成功')
      }
      resetForm()
    }
    if (editingId.value) {
      dialogVisible.value = false
      resetForm()
    }
    await loadData()
    await loadLinkedHuaAnSettings()
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Channel) {
  await ElMessageBox.confirm(`确认删除渠道「${row.channelCode}」？`, '提示', { type: 'warning' })
  await deleteChannel(row.id)
  ElMessage.success('已删除')
  await loadData()
  await loadLinkedHuaAnSettings()
}

function onPageChange(p: number) {
  page.value = p
  loadData()
}

function channelEnvType(row: Channel): string {
  if (row.envType) return row.envType
  const cfg = huaAnConfigs.value.find((c) => c.id === row.huaanSettingId)
  return cfg?.envType || 'test'
}

function linkedHuaAnLabel(row: Channel): string {
  return row.channelCode
}

onMounted(async () => {
  await loadHuaAnConfigs()
  await loadLinkedHuaAnSettings()
  await loadData()
})
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
      <p class="hint">
        渠道编码须从「华安配置」中选择；华安密钥随关联配置自动同步，无需手填。渠道密钥默认脱敏显示，可点击「显示」或「复制」。
      </p>
      <el-table v-loading="loading" :data="filteredList" size="default">
        <el-table-column prop="id" label="编号" width="72" />
        <el-table-column label="渠道编码" min-width="120">
          <template #default="{ row }">{{ linkedHuaAnLabel(row) }}</template>
        </el-table-column>
        <el-table-column label="华安环境" width="100" align="center">
          <template #default="{ row }">
            <StatusTag
              :status="channelEnvType(row) === 'prod' ? 'warning' : 'info'"
              :label="envLabel(channelEnvType(row))"
            />
          </template>
        </el-table-column>
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
        <el-table-column label="操作" width="120" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="openEditDialog(row)">编辑</el-button>
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

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? '编辑渠道' : '新增渠道'"
    width="520px"
    destroy-on-close
    @closed="resetForm"
  >
    <el-form label-width="100px">
      <el-form-item v-if="editingId" label="渠道编码">
        <el-input :model-value="form.channelCode" disabled />
      </el-form-item>
      <el-form-item v-else label="华安配置" required>
        <el-select
          v-model="form.huaanSettingId"
          placeholder="请选择华安配置（生产/测试均可）"
          style="width: 100%"
          filterable
        >
          <el-option
            v-for="opt in huaAnOptions"
            :key="opt.value"
            :label="opt.label"
            :value="opt.value"
          />
        </el-select>
        <p v-if="huaAnOptions.length === 0" class="field-hint">请先在「华安配置」中添加配置，或各配置均已关联渠道</p>
      </el-form-item>
      <el-form-item label="渠道名称">
        <el-input v-model="form.channelName" placeholder="可选" />
      </el-form-item>
      <el-form-item label="回调地址" required>
        <el-input v-model="form.callbackUrl" placeholder="https://渠道域名/投保结果回调" />
      </el-form-item>
      <el-form-item v-if="editingId" label="状态">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">禁用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="渠道密钥">
        <el-input
          v-model="form.channelKey"
          :placeholder="editingId ? '留空则不修改' : '留空则自动生成 32 位十六进制'"
          show-password
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="onSave">确定</el-button>
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

.field-hint {
  margin: 6px 0 0;
  font-size: 12px;
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
