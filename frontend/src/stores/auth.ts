import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/api/types'
import { authApi } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(
    localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!) : null
  )

  const isLoggedIn = computed(() => !!user.value)
  const isSeller = computed(() => user.value?.role === 'seller')
  const isBuyer = computed(() => user.value?.role === 'buyer')

  async function register(email: string, password: string, name: string, role: string) {
    const { data } = await authApi.register(email, password, name, role)
    setAuth(data.user)
  }

  async function login(email: string, password: string) {
    const { data } = await authApi.login(email, password)
    setAuth(data.user)
  }

  function setAuth(u: User) {
    user.value = u
    localStorage.setItem('user', JSON.stringify(u))
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // Cookie may already be expired
    }
    user.value = null
    localStorage.removeItem('user')
  }

  return { user, isLoggedIn, isSeller, isBuyer, register, login, logout }
})
