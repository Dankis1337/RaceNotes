<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSetupsStore } from '../stores/setups'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const setupsStore = useSetupsStore()

const { toast } = useToast()
const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const error = ref('')

const form = ref({
  name: '',
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
    if (!payload.components_description) payload.components_description = null

    if (isEdit.value) {
      await setupsStore.updateSetup(route.params.id, payload)
      toast('Setup updated')
    } else {
      await setupsStore.createSetup(payload)
      toast('Setup created')
    }
    router.push('/setups')
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to save setup'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="px-4 pt-6 pb-24 max-w-lg mx-auto">
    <div class="flex items-center justify-between mb-4">
      <button @click="router.back()" class="text-gray-400 text-sm">&larr; Back</button>
      <h1 class="text-lg font-bold">{{ isEdit ? 'Edit Setup' : 'New Setup' }}</h1>
      <div class="w-10"></div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-5">
      <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Setup Name *</label>
          <input v-model="form.name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. Road Race Setup" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Bike Name *</label>
          <input v-model="form.bike_name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. Canyon Aeroad CF SLX" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Tires *</label>
          <input v-model="form.tires" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. Continental GP5000 28mm" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Components Description</label>
          <textarea v-model="form.components_description" rows="4" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" placeholder="Shimano Ultegra Di2, Zipp 303 wheels..."></textarea>
        </div>
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="w-full bg-primary text-white font-medium py-3 rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 shadow-md"
      >
        {{ loading ? 'Saving...' : (isEdit ? 'Save Changes' : 'Create Setup') }}
      </button>
    </form>
  </div>
</template>
