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

const form = reactive({
  channelCode: '',
  channelName: '',
  huaAnKey: '',
  channelKey: '',
})

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
  if (!form.channelCode || !form.huaAnKey) {
    ElMessage.warning('请填写渠道编码和华安密钥')
    return
  }
  loading.value = true
  try {
    await createChannel({
      channelCode: form.channelCode,
      channelName: form.channelName,
      huaAnKey: form.huaAnKey,
      channelKey: form.channelKey || undefined,
      status: 1,
    })
    ElMessage.success('创建成功')
    form.channelCode = ''
    form.channelName = ''
    form.huaAnKey = ''
    form.channelKey = ''
    await loadData()
  } finally {
    loading.value = false
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
    <template #header>新增渠道</template>
    <el-form inline>
      <el-form-item label="渠道编码">
        <el-input v-model="form.channelCode" placeholder="channelCode" />
      </el-form-item>
      <el-form-item label="渠道名称">
        <el-input v-model="form.channelName" placeholder="名称" />
      </el-form-item>
      <el-form-item label="华安密钥">
        <el-input v-model="form.huaAnKey" placeholder="huaAnKey" />
      </el-form-item>
      <el-form-item label="渠道密钥">
        <el-input v-model="form.channelKey" placeholder="留空自动生成" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="onCreate">新增</el-button>
      </el-form-item>
    </el-form>
  </el-card>

  <el-card shadow="never" style="margin-top: 16px">
    <template #header>渠道列表</template>
    <el-table v-loading="loading" :data="list" stripe border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="channelCode" label="编码" min-width="120" />
      <el-table-column prop="channelName" label="名称" min-width="120" />
      <el-table-column label="渠道密钥" min-width="140">
        <template #default="{ row }">{{ maskSecret(row.channelKey) }}</template>
      </el-table-column>
      <el-table-column label="华安密钥" min-width="140">
        <template #default="{ row }">{{ maskSecret(row.huaAnKey) }}</template>
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
</template>

<style scoped>
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
