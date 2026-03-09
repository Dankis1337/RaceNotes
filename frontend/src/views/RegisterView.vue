<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useI18n } from '../utils/i18n'

const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()

const form = ref({
  name: '',
  username: '',
  email: '',
  password: '',
  height: null,
  weight: null,
})
const error = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    const payload = { ...form.value }
    if (payload.height) payload.height = Number(payload.height)
    if (payload.weight) payload.weight = Number(payload.weight)
    await auth.registerUser(payload)
    if (!auth.isAuthenticated) {
      await auth.loginUser({ username: form.value.username, password: form.value.password })
    }
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || t('registration_failed')
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
        <p class="text-gray-500 mt-1">{{ t('create_account') }}</p>
      </div>

      <form @submit.prevent="handleRegister" class="bg-white rounded-xl shadow-md p-6 space-y-4">
        <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('name') }}</label>
          <input v-model="form.name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('your_name')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('username') }}</label>
          <input v-model="form.username" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('choose_username')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('email') }}</label>
          <input v-model="form.email" type="email" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('email_placeholder')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('password') }}</label>
          <input v-model="form.password" type="password" required minlength="6" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('min_6_chars')" />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('height_cm') }}</label>
            <input v-model="form.height" type="number" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="175" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('weight_kg') }}</label>
            <input v-model="form.weight" type="number" step="0.1" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="70" />
          </div>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-primary text-white font-medium py-2.5 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
        >
          {{ loading ? t('creating_account') : t('sign_up') }}
        </button>

        <p class="text-center text-sm text-gray-500">
          {{ t('already_have_account') }}
          <router-link to="/login" class="text-primary font-medium">{{ t('sign_in') }}</router-link>
        </p>
      </form>
    </div>
  </div>
</template>
