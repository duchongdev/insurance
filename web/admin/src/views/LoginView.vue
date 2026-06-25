<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function onSubmit() {
  if (!username.value || !password.value) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-panel">
      <div class="login-brand">
        <span class="login-brand__mark">HA</span>
        <div>
          <h1 class="login-brand__title">华安保险渠道对接平台</h1>
          <p class="login-brand__desc">管理后台登录</p>
        </div>
      </div>

      <div class="login-card">
        <h2 class="login-card__title">管理员登录</h2>
        <el-form label-position="top" @submit.prevent="onSubmit">
          <el-form-item label="用户名">
            <el-input v-model="username" autocomplete="username" placeholder="请输入用户名" size="large" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input
              v-model="password"
              type="password"
              autocomplete="current-password"
              placeholder="请输入密码"
              show-password
              size="large"
            />
          </el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" size="large" class="login-btn">
            登录
          </el-button>
        </el-form>
        <p class="login-hint">首次部署账号见 config/config.yaml 中 admin 段（仅表为空时自动创建）。</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-page);
  padding: 24px;
}

.login-panel {
  width: 100%;
  max-width: 420px;
}

.login-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.login-brand__mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 4px;
  background: var(--color-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 700;
}

.login-brand__title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.login-brand__desc {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--color-text-muted);
}

.login-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: 28px 32px;
}

.login-card__title {
  margin: 0 0 24px;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.login-btn {
  width: 100%;
  margin-top: 8px;
}

.login-hint {
  margin: 16px 0 0;
  font-size: 12px;
  color: var(--color-text-muted);
  line-height: 1.6;
}

.login-card :deep(.el-form-item__label) {
  color: var(--color-text-secondary);
  font-weight: 500;
}
</style>
