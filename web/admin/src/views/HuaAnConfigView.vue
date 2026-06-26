<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageCard from '@/components/PageCard.vue'
import PageToolbar from '@/components/PageToolbar.vue'
import {
  createHuaAnConfig,
  deleteHuaAnConfig,
  fetchHuaAnConfigs,
  testHuaAnConnectivity,
  updateHuaAnConfig,
  type HuaAnConfig,
  type HuaAnConnectivityResult,
  type HuaAnEnvType,
} from '@/api/admin'
import { maskSecret } from '@/utils/auth'

const loading = ref(false)
const saving = ref(false)
const list = ref<HuaAnConfig[]>([])
const dialogVisible = ref(false)
const editingId = ref<number | null>(null)
const testingId = ref<number | null>(null)
const testResults = ref<Record<number, HuaAnConnectivityResult>>({})
const revealedIds = ref<Set<number>>(new Set())

const form = reactive({
  name: '',
  envType: 'test' as HuaAnEnvType,
  baseUrl: '',
  channelCode: '',
  channelSecret: '',
})

function envLabel(envType: HuaAnEnvType): string {
  return envType === 'prod' ? '生产环境' : '测试环境'
}

function envCardClass(envType: HuaAnEnvType): string {
  return envType === 'prod' ? 'config-card--prod' : 'config-card--test'
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

function displaySecret(row: HuaAnConfig): string {
  if (!row.channelSecret) return '—'
  return isRevealed(row.id!) ? row.channelSecret : maskSecret(row.channelSecret)
}

async function copySecret(row: HuaAnConfig) {
  if (!row.channelSecret) {
    ElMessage.warning('华安侧密钥为空')
    return
  }
  try {
    await navigator.clipboard.writeText(row.channelSecret)
    ElMessage.success('已复制华安侧密钥')
  } catch {
    ElMessage.error('复制失败，请手动选择复制')
  }
}

function resetForm() {
  form.name = ''
  form.envType = 'test'
  form.baseUrl = ''
  form.channelCode = ''
  form.channelSecret = ''
  editingId.value = null
}

function openCreateDialog() {
  resetForm()
  dialogVisible.value = true
}

function openEditDialog(row: HuaAnConfig) {
  editingId.value = row.id!
  form.name = row.name || ''
  form.envType = row.envType
  form.baseUrl = row.baseUrl
  form.channelCode = row.channelCode
  form.channelSecret = ''
  dialogVisible.value = true
}

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchHuaAnConfigs()
    list.value = data.list || []
  } finally {
    loading.value = false
  }
}

async function onSave() {
  if (!form.baseUrl.trim() || !form.channelCode.trim()) {
    ElMessage.warning('请填写接口地址与渠道编码')
    return
  }
  if (!editingId.value && !form.channelSecret.trim()) {
    ElMessage.warning('请填写华安侧密钥')
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name.trim(),
      envType: form.envType,
      baseUrl: form.baseUrl.trim(),
      channelCode: form.channelCode.trim(),
      channelSecret: form.channelSecret.trim() || undefined,
    }
    if (editingId.value) {
      await updateHuaAnConfig(editingId.value, payload)
      ElMessage.success('配置已更新')
    } else {
      await createHuaAnConfig(payload)
      ElMessage.success('配置已添加')
    }
    dialogVisible.value = false
    resetForm()
    await loadData()
  } finally {
    saving.value = false
  }
}

async function onDelete(row: HuaAnConfig) {
  await ElMessageBox.confirm(`确认删除「${row.name || envLabel(row.envType)}」？`, '提示', { type: 'warning' })
  await deleteHuaAnConfig(row.id!)
  ElMessage.success('已删除')
  await loadData()
}

async function onTest(row: HuaAnConfig) {
  testingId.value = row.id!
  try {
    const { data } = await testHuaAnConnectivity(row.id!)
    testResults.value = { ...testResults.value, [row.id!]: data }
    if (data.ok) ElMessage.success(data.message)
    else ElMessage.error(data.message)
  } finally {
    testingId.value = null
  }
}

onMounted(loadData)
</script>

<template>
  <div>
    <PageToolbar title="华安配置" subtitle="管理生产/测试环境华安上游连接，可分别做连通性测试">
      <template #actions>
        <el-button type="primary" @click="openCreateDialog">添加配置</el-button>
      </template>
    </PageToolbar>

    <PageCard v-loading="loading">
      <p class="hint">
        每条配置需选择环境类型（生产/测试）。<span class="legend legend--prod">生产环境</span>与
        <span class="legend legend--test">测试环境</span>以颜色区分；渠道管理与银行列表同步时将按需选择具体配置。
      </p>

      <div v-if="!loading && list.length === 0" class="empty">暂无配置，请点击「添加配置」</div>

      <div class="config-list">
        <article
          v-for="row in list"
          :key="row.id"
          class="config-card"
          :class="envCardClass(row.envType)"
        >
          <header class="config-card__header">
            <div class="config-card__title">
              <span class="env-badge" :class="row.envType === 'prod' ? 'env-badge--prod' : 'env-badge--test'">
                {{ envLabel(row.envType) }}
              </span>
              <strong>{{ row.name || row.baseUrl }}</strong>
            </div>
            <div class="config-card__actions">
              <el-button
                size="small"
                :loading="testingId === row.id"
                @click="onTest(row)"
              >
                连通性测试
              </el-button>
              <el-button size="small" link type="primary" @click="openEditDialog(row)">编辑</el-button>
              <el-button size="small" link type="danger" @click="onDelete(row)">删除</el-button>
            </div>
          </header>

          <dl class="config-card__body">
            <div class="field">
              <dt>接口地址</dt>
              <dd>{{ row.baseUrl }}</dd>
            </div>
            <div class="field">
              <dt>渠道编码</dt>
              <dd>{{ row.channelCode }}</dd>
            </div>
            <div class="field">
              <dt>华安侧密钥</dt>
              <dd class="secret-row">
                <span class="secret-text">{{ displaySecret(row) }}</span>
                <el-button type="primary" link @click="toggleReveal(row.id!)">
                  {{ isRevealed(row.id!) ? '隐藏' : '显示' }}
                </el-button>
                <el-button type="primary" link @click="copySecret(row)">复制</el-button>
              </dd>
            </div>
          </dl>

          <div v-if="testResults[row.id!]" class="config-card__test">
            <el-alert
              :type="testResults[row.id!].ok ? 'success' : 'error'"
              :title="testResults[row.id!].message"
              :closable="false"
              show-icon
            />
            <p v-if="testResults[row.id!].upstreamUrl" class="test-meta">
              探测地址：{{ testResults[row.id!].upstreamUrl }}
              <template v-if="testResults[row.id!].ok">
                · 耗时 {{ testResults[row.id!].durationMs }} ms
                <template v-if="testResults[row.id!].httpStatus"> · HTTP {{ testResults[row.id!].httpStatus }}</template>
              </template>
            </p>
          </div>
        </article>
      </div>
    </PageCard>
  </div>

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? '编辑华安配置' : '添加华安配置'"
    width="520px"
    destroy-on-close
    @closed="resetForm"
  >
    <el-form label-width="100px">
      <el-form-item label="环境类型" required>
        <el-radio-group v-model="form.envType">
          <el-radio value="test">
            <span class="legend legend--test">测试环境</span>
          </el-radio>
          <el-radio value="prod">
            <span class="legend legend--prod">生产环境</span>
          </el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="配置名称">
        <el-input v-model="form.name" placeholder="可选，便于识别" />
      </el-form-item>
      <el-form-item label="接口地址" required>
        <el-input v-model="form.baseUrl" placeholder="https://ins.api.hahealth.ink/" />
      </el-form-item>
      <el-form-item label="渠道编码" required>
        <el-input v-model="form.channelCode" placeholder="channel_code" />
      </el-form-item>
      <el-form-item label="华安侧密钥" :required="!editingId">
        <el-input
          v-model="form.channelSecret"
          :placeholder="editingId ? '留空则不修改' : 'channel_secret'"
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
.hint {
  margin: 0 0 20px;
  font-size: 13px;
  color: var(--color-text-muted);
  line-height: 1.6;
}

.legend {
  display: inline-block;
  padding: 0 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.legend--prod {
  background: #fef0f0;
  color: #c45656;
  border: 1px solid #fbc4c4;
}

.legend--test {
  background: #ecf5ff;
  color: #337ecc;
  border: 1px solid #b3d8ff;
}

.empty {
  text-align: center;
  color: var(--color-text-muted);
  padding: 32px 0;
}

.config-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-card {
  border-radius: 8px;
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.config-card--prod {
  border-left: 4px solid #e6a23c;
  background: linear-gradient(to right, #fffbf5 0%, var(--color-bg-card) 48px);
}

.config-card--test {
  border-left: 4px solid #409eff;
  background: linear-gradient(to right, #f5f9ff 0%, var(--color-bg-card) 48px);
}

.config-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-border);
  flex-wrap: wrap;
}

.config-card__title {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 0;
}

.config-card__title strong {
  font-size: 14px;
  word-break: break-all;
}

.config-card__actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.env-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
  white-space: nowrap;
}

.env-badge--prod {
  background: #fdf6ec;
  color: #b88230;
}

.env-badge--test {
  background: #ecf5ff;
  color: #337ecc;
}

.config-card__body {
  margin: 0;
  padding: 12px 16px 16px;
  display: grid;
  gap: 10px;
}

.field {
  display: grid;
  grid-template-columns: 88px 1fr;
  gap: 8px;
  font-size: 13px;
}

.field dt {
  margin: 0;
  color: var(--color-text-muted);
}

.field dd {
  margin: 0;
  color: var(--color-text-primary);
  word-break: break-all;
}

.secret-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.secret-text {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  margin-right: 4px;
  color: var(--color-text-secondary);
}

.config-card__test {
  padding: 0 16px 16px;
}

.test-meta {
  margin: 8px 0 0;
  font-size: 12px;
  color: var(--color-text-muted);
}
</style>
