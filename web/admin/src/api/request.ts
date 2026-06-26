import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { clearToken, getToken } from '@/utils/auth'
import router from '@/router'

const request = axios.create({
  baseURL: '/admin/api',
  timeout: 60000,
})

request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => response,
  (error: AxiosError<{ message?: string }>) => {
    const status = error.response?.status
    const message = error.response?.data?.message || error.message || '请求失败'

    if (status === 401) {
      if (router.currentRoute.value.path === '/login') {
        ElMessage.error('账号或密码错误')
        return Promise.reject(error)
      }
      clearToken()
      ElMessage.error('登录已过期，请重新登录')
      router.push({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } })
      return Promise.reject(error)
    }

    ElMessage.error(message)
    return Promise.reject(error)
  },
)

export default request
