import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)

  const isAuthenticated = computed(() => !!token.value)

  async function loginUser(credentials) {
    const { data } = await authApi.login(credentials)
    token.value = data.token
    localStorage.setItem('token', data.token)
  }

  async function registerUser(userData) {
    const { data } = await authApi.register(userData)
    if (data.token) {
      token.value = data.token
      localStorage.setItem('token', data.token)
    }
  }

  async function fetchProfile() {
    const { data } = await authApi.getProfile()
    user.value = data
  }

  async function updateUser(updates) {
    const { data } = await authApi.updateProfile(updates)
    user.value = data
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { token, user, isAuthenticated, loginUser, registerUser, fetchProfile, updateUser, logout }
})
