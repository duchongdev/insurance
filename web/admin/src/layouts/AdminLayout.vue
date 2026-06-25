<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const activeMenu = computed(() => route.path)

const menuItems = [
  { path: '/dashboard', title: '概览', icon: 'Odometer' },
  { path: '/channels', title: '渠道管理', icon: 'Connection' },
  { path: '/banks', title: '银行列表', icon: 'Money' },
]

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <span class="logo__mark">HA</span>
        <span class="logo__text">华安渠道对接</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        class="sidebar-menu"
        background-color="transparent"
        text-color="#ffffff"
        active-text-color="#ffffff"
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-container">
      <el-header class="top-bar" height="48px">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item>管理后台</el-breadcrumb-item>
          <el-breadcrumb-item>{{ route.meta.title as string }}</el-breadcrumb-item>
        </el-breadcrumb>
        <el-button class="logout-btn" text @click="logout">退出登录</el-button>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  min-height: 100vh;
}

.aside {
  background: var(--color-bg-sidebar);
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 100;
  overflow-y: auto;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo__mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  background: var(--color-primary);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
}

.logo__text {
  color: #fff;
  font-size: 15px;
  font-weight: 600;
}

.sidebar-menu {
  border-right: none;
  padding: 8px;
}

.sidebar-menu :deep(.el-menu-item) {
  height: 40px;
  line-height: 40px;
  border-radius: 4px;
  margin-bottom: 4px;
  color: rgba(255, 255, 255, 0.85);
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.08);
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: var(--color-primary);
  color: #fff;
}

.main-container {
  margin-left: 220px;
  min-height: 100vh;
  background: var(--color-bg-page);
}

.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
}

.top-bar :deep(.el-breadcrumb__inner) {
  color: var(--color-text-secondary);
  font-weight: 400;
}

.top-bar :deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: var(--color-text-primary);
  font-weight: 500;
}

.logout-btn {
  color: var(--color-text-secondary);
}

.logout-btn:hover {
  color: var(--color-primary);
}

.main {
  padding: 16px 20px 24px;
  min-height: calc(100vh - 48px);
}
</style>
