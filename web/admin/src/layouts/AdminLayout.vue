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
  { path: '/policies', title: '保单', icon: 'Document' },
  { path: '/users', title: '用户', icon: 'User' },
  { path: '/signs', title: '签约', icon: 'EditPen' },
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
      <div class="logo">华安渠道对接</div>
      <el-menu :default-active="activeMenu" router background-color="#001529" text-color="#ffffffa6" active-text-color="#fff">
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span class="title">华安保险渠道对接平台</span>
        <el-button type="danger" plain @click="logout">退出</el-button>
      </el-header>
      <el-main class="main">
        <el-breadcrumb separator="/" class="breadcrumb">
          <el-breadcrumb-item>管理后台</el-breadcrumb-item>
          <el-breadcrumb-item>{{ route.meta.title as string }}</el-breadcrumb-item>
        </el-breadcrumb>
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
  background: #001529;
}
.logo {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  padding: 20px 16px;
  border-bottom: 1px solid #ffffff1a;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
}
.title {
  font-size: 18px;
  font-weight: 600;
}
.main {
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}
.breadcrumb {
  margin-bottom: 16px;
}
</style>

<style>
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
</style>
