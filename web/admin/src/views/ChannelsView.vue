<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createChannel, deleteChannel, fetchChannels, type Channel } from '@/api/admin'
import { maskSecret } from '@/utils/auth'

const loading = ref(false)
const list = ref<Channel[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const createVisible = ref(false)
const creating = ref(false)
const form = reactive({
  channelCode: '',
  channelName: '',
  huaAnKey: '',
  channelKey: '',
})

/** 列表中已展开明文的行 id */
const revealedIds = ref<Set<number>>(new Set())

function resetForm() {
  form.channelCode = ''
  form.channelName = ''
  form.huaAnKey = ''
  form.channelKey = ''
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

async function loadData() {
  loading.value = true
  try {
    const { data } = await fetchChannels(page.value, pageSize.value)
    list.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

async function onCreate() {
  if (!form.channelCode.trim() || !form.huaAnKey.trim()) {
    ElMessage.warning('请填写渠道编码和华安密钥')
    return
  }
  creating.value = true
  try {
    const { data } = await createChannel({
      channelCode: form.channelCode.trim(),
      channelName: form.channelName.trim(),
      huaAnKey: form.huaAnKey.trim(),
      channelKey: form.channelKey.trim() || undefined,
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
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <span>渠道列表</span>
        <el-button type="primary" @click="openCreateDialog">新增渠道</el-button>
      </div>
    </template>
    <p class="hint">密钥在库内为明文存储，列表默认脱敏显示；点击「显示」可查看完整内容，「复制」可写入剪贴板。</p>
    <el-table v-loading="loading" :data="list" stripe border>
      <el-table-column prop="id" label="编号" width="80" />
      <el-table-column prop="channelCode" label="渠道编码" min-width="120" />
      <el-table-column prop="channelName" label="渠道名称" min-width="120" />
      <el-table-column label="渠道密钥" min-width="220">
        <template #default="{ row }">
          <span class="secret-text">{{ displaySecret(row, 'channelKey') }}</span>
          <el-button type="primary" link @click="toggleReveal(row.id)">
            {{ isRevealed(row.id) ? '隐藏' : '显示' }}
          </el-button>
          <el-button type="primary" link @click="copySecret(row.channelKey, '渠道密钥')">复制</el-button>
        </template>
      </el-table-column>
      <el-table-column label="华安密钥" min-width="220">
        <template #default="{ row }">
          <span class="secret-text">{{ displaySecret(row, 'huaAnKey') }}</span>
          <el-button type="primary" link @click="toggleReveal(row.id)">
            {{ isRevealed(row.id) ? '隐藏' : '显示' }}
          </el-button>
          <el-button type="primary" link @click="copySecret(row.huaAnKey, '华安密钥')">复制</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
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
  </el-card>

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
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.hint {
  margin: 0 0 12px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
.secret-text {
  font-family: ui-monospace, monospace;
  font-size: 13px;
  margin-right: 4px;
  word-break: break-all;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
