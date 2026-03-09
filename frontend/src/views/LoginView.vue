<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useI18n } from '../utils/i18n'

const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()

const form = ref({ username: '', password: '' })
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.loginUser(form.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || t('login_failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center px-4 bg-gray-50">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-primary">{{ t('app_name') }}</h1>
        <p class="text-gray-500 mt-1">{{ t('login_subtitle') }}</p>
      </div>

      <form @submit.prevent="handleLogin" class="bg-white rounded-xl shadow-md p-6 space-y-4">
        <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('username') }}</label>
          <input
            v-model="form.username"
            type="text"
            required
            class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none"
            :placeholder="t('enter_username')"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('password') }}</label>
          <input
            v-model="form.password"
            type="password"
            required
            class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none"
            :placeholder="t('enter_password')"
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-primary text-white font-medium py-2.5 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
        >
          {{ loading ? t('signing_in') : t('sign_in') }}
        </button>

        <p class="text-center text-sm text-gray-500">
          {{ t('no_account') }}
          <router-link to="/register" class="text-primary font-medium">{{ t('sign_up') }}</router-link>
        </p>
      </form>
    </div>
  </div>
</template>
