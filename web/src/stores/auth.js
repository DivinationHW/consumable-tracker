import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const role = ref(localStorage.getItem('role') || '')

  const isAdmin = computed(() => role.value === 'admin')
  const isReadonly = computed(() => role.value === 'readonly')
  const loggedIn = computed(() => !!token.value)

  function setToken(t, r) {
    token.value = t
    role.value = r
    localStorage.setItem('token', t)
    localStorage.setItem('role', r)
  }

  function logout() {
    token.value = ''
    role.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('role')
  }

  return { token, role, isAdmin, isReadonly, loggedIn, setToken, logout }
})
