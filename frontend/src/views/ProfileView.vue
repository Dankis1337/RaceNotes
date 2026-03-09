<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const editing = ref(false)
const loading = ref(false)
const error = ref('')
const success = ref('')

const form = ref({
  name: '',
  email: '',
  height: null,
  weight: null,
})

onMounted(async () => {
  if (!authStore.user) await authStore.fetchProfile()
  if (authStore.user) {
    form.value = {
      name: authStore.user.name,
      email: authStore.user.email,
      height: authStore.user.height,
      weight: authStore.user.weight,
    }
  }
})

async function handleSave() {
  error.value = ''
  success.value = ''
  loading.value = true
  try {
    const payload = { ...form.value }
    if (payload.height) payload.height = Number(payload.height)
    if (payload.weight) payload.weight = Number(payload.weight)
    await authStore.updateUser(payload)
    editing.value = false
    success.value = 'Profile updated'
    setTimeout(() => { success.value = '' }, 3000)
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to update profile'
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto">
    <h1 class="text-2xl font-bold mb-4">Profile</h1>

    <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3 mb-4">{{ error }}</div>
    <div v-if="success" class="bg-green-50 text-green-600 text-sm rounded-lg p-3 mb-4">{{ success }}</div>

    <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
      <div v-if="!editing">
        <div class="space-y-3">
          <div>
            <p class="text-gray-400 text-xs">Name</p>
            <p class="font-medium">{{ authStore.user?.name }}</p>
          </div>
          <div>
            <p class="text-gray-400 text-xs">Username</p>
            <p class="font-medium">{{ authStore.user?.username }}</p>
          </div>
          <div>
            <p class="text-gray-400 text-xs">Email</p>
            <p class="font-medium">{{ authStore.user?.email }}</p>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <p class="text-gray-400 text-xs">Height</p>
              <p class="font-medium">{{ authStore.user?.height ? authStore.user.height + ' cm' : '—' }}</p>
            </div>
            <div>
              <p class="text-gray-400 text-xs">Weight</p>
              <p class="font-medium">{{ authStore.user?.weight ? authStore.user.weight + ' kg' : '—' }}</p>
            </div>
          </div>
        </div>

        <button
          @click="editing = true"
          class="w-full mt-4 bg-primary text-white font-medium py-2.5 rounded-lg hover:bg-primary-dark transition-colors"
        >
          Edit Profile
        </button>
      </div>

      <form v-else @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input v-model="form.name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input :value="authStore.user?.username" type="text" disabled class="w-full border border-gray-200 rounded-lg px-3 py-2.5 bg-gray-50 text-gray-400 outline-none" />
          <p class="text-xs text-gray-400 mt-1">Username cannot be changed</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input v-model="form.email" type="email" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Height (cm)</label>
            <input v-model="form.height" type="number" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Weight (kg)</label>
            <input v-model="form.weight" type="number" step="0.1" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
          </div>
        </div>

        <div class="flex gap-3">
          <button type="button" @click="editing = false" class="flex-1 py-2.5 rounded-lg border border-gray-300 text-gray-600">
            Cancel
          </button>
          <button type="submit" :disabled="loading" class="flex-1 py-2.5 rounded-lg bg-primary text-white font-medium disabled:opacity-50">
            {{ loading ? 'Saving...' : 'Save' }}
          </button>
        </div>
      </form>
    </div>

    <button
      @click="handleLogout"
      class="w-full mt-4 py-2.5 rounded-xl border border-red-300 text-red-500 font-medium hover:bg-red-50 transition-colors"
    >
      Log Out
    </button>
  </div>
</template>
