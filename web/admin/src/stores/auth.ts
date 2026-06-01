import { defineStore } from 'pinia'
import { ref } from 'vue'
import { clearToken, getToken, setToken } from '@/utils/auth'
import { login as loginApi } from '@/api/admin'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(getToken())

  function isLoggedIn() {
    return !!token.value
  }

  async function login(username: string, password: string) {
    const { data } = await loginApi(username, password)
    token.value = data.token
    setToken(data.token)
  }

  function logout() {
    token.value = ''
    clearToken()
  }

  return { token, isLoggedIn, login, logout }
})
