import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: () => import('@/views/DashboardView.vue'), meta: { title: '概览' } },
        { path: 'channels', name: 'channels', component: () => import('@/views/ChannelsView.vue'), meta: { title: '渠道管理' } },
        { path: 'policies', name: 'policies', component: () => import('@/views/PoliciesView.vue'), meta: { title: '保单' } },
        { path: 'users', name: 'users', component: () => import('@/views/UsersView.vue'), meta: { title: '用户' } },
        { path: 'signs', name: 'signs', component: () => import('@/views/SignsView.vue'), meta: { title: '签约' } },
        { path: 'banks', name: 'banks', component: () => import('@/views/BanksView.vue'), meta: { title: '银行列表' } },
      ],
    },
  ],
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  if (!getToken()) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router
