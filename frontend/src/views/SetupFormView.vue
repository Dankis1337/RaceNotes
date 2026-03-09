<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSetupsStore } from '../stores/setups'
import { useToast } from '../composables/useToast'
import PhotoUpload from '../components/PhotoUpload.vue'
import { useI18n } from '../utils/i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const setupsStore = useSetupsStore()

const { toast } = useToast()
const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const error = ref('')

const form = ref({
  name: '',
  photo: '',
  bike_name: '',
  tires: '',
  components_description: '',
})

onMounted(async () => {
  if (isEdit.value) {
    await setupsStore.fetchSetup(route.params.id)
    const s = setupsStore.currentSetup
    if (s) {
      form.value = {
        name: s.name,
        photo: s.photo || '',
        bike_name: s.bike_name,
        tires: s.tires,
        components_description: s.components_description || '',
      }
    }
  }
})

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    const payload = { ...form.value }
    if (!payload.photo) payload.photo = null
    if (!payload.components_description) payload.components_description = null

    if (isEdit.value) {
      await setupsStore.updateSetup(route.params.id, payload)
      toast(t('setup_updated'))
    } else {
      await setupsStore.createSetup(payload)
      toast(t('setup_created'))
    }
    router.push('/setups')
  } catch (e) {
    error.value = e.response?.data?.error || t('failed_save_setup')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="px-4 pt-6 pb-24 max-w-lg mx-auto lg:max-w-2xl">
    <div class="flex items-center justify-between mb-4">
      <button @click="router.back()" class="text-gray-400 text-sm">&larr; {{ t('back') }}</button>
      <h1 class="text-lg font-bold">{{ isEdit ? t('edit_setup') : t('new_setup') }}</h1>
      <div class="w-10"></div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-5">
      <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('setup_name') }} *</label>
          <input v-model="form.name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('setup_name_placeholder')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('photo') }}</label>
          <PhotoUpload v-model="form.photo" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('bike_name') }} *</label>
          <input v-model="form.bike_name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('bike_name_setup_placeholder')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('tires') }} *</label>
          <input v-model="form.tires" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('tires_placeholder')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('components_description') }}</label>
          <textarea v-model="form.components_description" rows="4" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" :placeholder="t('components_placeholder')"></textarea>
        </div>
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="w-full bg-primary text-white font-medium py-3 rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 shadow-md"
      >
        {{ loading ? t('saving') : (isEdit ? t('save_changes') : t('create_setup')) }}
      </button>
    </form>
  </div>
</template>
