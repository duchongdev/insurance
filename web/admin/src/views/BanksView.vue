<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken } from '@/utils/auth'
import { refreshBankList } from '@/api/admin'

const channelCode = ref('')
const huaAnKey = ref('')
const loading = ref(false)
const result = ref('')
const errorMsg = ref('')

function formatJson(text: string): string {
  try {
    return JSON.stringify(JSON.parse(text), null, 2)
  } catch {
    return text
  }
}

async function loadCache() {
  errorMsg.value = ''
  result.value = ''
  const code = channelCode.value.trim()
  if (!code) {
    ElMessage.warning('请填写渠道码')
    return
  }
  loading.value = true
  try {
    const res = await axios.get('/admin/api/bank-list', {
      params: { channelCode: code },
      headers: { Authorization: `Bearer ${getToken()}` },
      transformResponse: [(data) => data],
      responseType: 'text',
      validateStatus: () => true,
    })
    if (res.status === 404) {
      errorMsg.value = '暂无缓存，请点击「刷新并缓存」'
      return
    }
    if (res.status >= 400) {
      let msg = res.data
      try {
        msg = JSON.parse(res.data).message
      } catch {
        /* keep raw */
      }
      errorMsg.value = msg || '加载失败'
      return
    }
    result.value = formatJson(res.data)
  } finally {
    loading.value = false
  }
}

async function refreshCache() {
  errorMsg.value = ''
  result.value = ''
  const code = channelCode.value.trim()
  const key = huaAnKey.value.trim()
  if (!code || !key) {
    ElMessage.warning('请填写渠道码和华安密钥')
    return
  }
  loading.value = true
  try {
    const { data } = await refreshBankList(code, key)
    result.value = formatJson(data)
    ElMessage.success('刷新成功')
  } catch {
    /* interceptor shows error */
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <el-card shadow="never" v-loading="loading">
    <template #header>银行列表</template>
    <p class="desc">手动请求华安刷新银行列表并写入 Redis；渠道 getBankList 优先读缓存。</p>
    <el-form label-width="100px" style="max-width: 560px">
      <el-form-item label="渠道码">
        <el-input v-model="channelCode" placeholder="channelCode" />
      </el-form-item>
      <el-form-item label="华安密钥">
        <el-input v-model="huaAnKey" placeholder="huaAnKey" show-password />
      </el-form-item>
      <el-form-item>
        <el-button @click="loadCache">查看缓存</el-button>
        <el-button type="primary" @click="refreshCache">刷新并缓存</el-button>
      </el-form-item>
    </el-form>
    <el-alert v-if="errorMsg" type="warning" :title="errorMsg" show-icon :closable="false" style="margin-bottom: 12px" />
    <el-input v-if="result" v-model="result" type="textarea" :rows="16" readonly class="result" />
  </el-card>
</template>

<style scoped>
.desc {
  color: #909399;
  font-size: 14px;
  margin: 0 0 16px;
}
.result {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
</style>
