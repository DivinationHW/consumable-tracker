import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getMe } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<{ user_id: number; username: string; role: string } | null>(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setLogin(t: string, u: { user_id: number; username: string; role: string }) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
    localStorage.setItem('user', JSON.stringify(u))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchMe() {
    try {
      const res = await getMe()
      user.value = res.data
    } catch {
      logout()
    }
  }

  const saved = localStorage.getItem('user')
  if (saved) {
    try { user.value = JSON.parse(saved) } catch {}
  }

  return { user, token, isLoggedIn, isAdmin, setLogin, logout, fetchMe }
})
