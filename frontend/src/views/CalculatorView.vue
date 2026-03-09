<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { calculatePressure } from '../api/calculator'
import { calculateTirePressureOffline } from '../utils/tirePressureCalc'

const authStore = useAuthStore()

const form = ref({
  rider_weight: 70,
  bike_weight: 8,
  tire_width: 28,
  tire_type: 'clincher',
  surface: 'road',
  conditions: 'dry',
})

const result = ref(null)
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  if (!authStore.user) await authStore.fetchProfile()
  if (authStore.user?.weight) {
    form.value.rider_weight = authStore.user.weight
  }
})

async function handleCalculate() {
  error.value = ''
  loading.value = true
  result.value = null
  try {
    const { data } = await calculatePressure(form.value)
    result.value = data
  } catch (e) {
    // Fallback to offline calculation
    try {
      result.value = calculateTirePressureOffline(form.value)
    } catch {
      error.value = 'Calculation failed'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto">
    <h1 class="text-2xl font-bold mb-4">Tire Pressure Calculator</h1>

    <form @submit.prevent="handleCalculate" class="space-y-5">
      <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Rider Weight (kg)</label>
            <input v-model.number="form.rider_weight" type="number" step="0.1" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Bike Weight (kg)</label>
            <input v-model.number="form.bike_weight" type="number" step="0.1" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Tire Width (mm)</label>
          <input v-model.number="form.tire_width" type="number" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Tire Type</label>
          <select v-model="form.tire_type" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option value="clincher">Clincher</option>
            <option value="tubeless">Tubeless</option>
            <option value="tubular">Tubular</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Surface</label>
          <select v-model="form.surface" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option value="road">Road</option>
            <option value="gravel">Gravel</option>
            <option value="mixed">Mixed</option>
            <option value="cobblestone">Cobblestone</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Conditions</label>
          <select v-model="form.conditions" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option value="dry">Dry</option>
            <option value="wet">Wet</option>
            <option value="mud">Mud</option>
            <option value="snow">Snow</option>
          </select>
        </div>
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="w-full bg-primary text-white font-medium py-3 rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 shadow-md"
      >
        {{ loading ? 'Calculating...' : 'Calculate' }}
      </button>
    </form>

    <div v-if="result" class="mt-5 bg-white rounded-xl shadow-md p-5 space-y-4">
      <h2 class="font-semibold text-lg">Recommended Pressure</h2>

      <div class="grid grid-cols-2 gap-4 text-center">
        <div class="bg-green-50 rounded-xl p-4">
          <p class="text-xs text-gray-500 mb-1">Front</p>
          <p class="text-3xl font-bold text-primary">{{ result.front_pressure }}</p>
          <p class="text-xs text-gray-400">{{ result.unit }}</p>
        </div>
        <div class="bg-green-50 rounded-xl p-4">
          <p class="text-xs text-gray-500 mb-1">Rear</p>
          <p class="text-3xl font-bold text-primary">{{ result.rear_pressure }}</p>
          <p class="text-xs text-gray-400">{{ result.unit }}</p>
        </div>
      </div>

      <div v-if="result.recommendations?.length" class="space-y-2">
        <h3 class="text-sm font-medium text-gray-600">Tips</h3>
        <p v-for="(rec, i) in result.recommendations" :key="i" class="text-sm text-gray-500 leading-relaxed">
          {{ rec }}
        </p>
      </div>
    </div>
  </div>
</template>
