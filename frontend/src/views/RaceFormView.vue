<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRacesStore } from '../stores/races'
import { useSetupsStore } from '../stores/setups'
import { useAuthStore } from '../stores/auth'
import { calculatePressure } from '../api/calculator'
import { calculateTirePressureOffline } from '../utils/tirePressureCalc'
import StarRating from '../components/StarRating.vue'
import { CalculatorIcon } from '@heroicons/vue/24/outline'
import { useToast } from '../composables/useToast'
import PhotoUpload from '../components/PhotoUpload.vue'

const route = useRoute()
const router = useRouter()
const racesStore = useRacesStore()
const setupsStore = useSetupsStore()
const authStore = useAuthStore()

const { toast } = useToast()
const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const error = ref('')
const showCalculator = ref(false)

const raceTypes = ['Road', 'MTB', 'Gravel', 'Cyclocross', 'Track']
const conditionOptions = ['Sunny', 'Cloudy', 'Rain', 'Snow']
const windOptions = ['None', 'Light', 'Moderate', 'Strong']
const roadOptions = ['Dry', 'Wet', 'Mud']

const form = ref({
  name: '',
  date: '',
  type: 'Road',
  photo: '',
  setup_id: null,
  bike_name: '',
  tires: '',
  tire_pressure_front: null,
  tire_pressure_rear: null,
  other_components: '',
  temperature: null,
  conditions: null,
  wind: null,
  road_conditions: null,
  nutrition_plan: '',
  result: '',
  rating: null,
  feelings: '',
  is_completed: false,
})

const useSetup = ref(false)

const needsManualBike = computed(() => !useSetup.value || !form.value.setup_id)

const selectedSetup = computed(() => {
  if (!form.value.setup_id) return null
  return setupsStore.setups.find(s => s.id === form.value.setup_id)
})

// Calculator form
const calcForm = ref({
  rider_weight: 70,
  bike_weight: 8,
  tire_width: 28,
  tire_type: 'clincher',
  surface: 'road',
  conditions: 'dry',
})
const calcResult = ref(null)
const calcLoading = ref(false)

onMounted(async () => {
  await setupsStore.fetchSetups()

  if (!authStore.user) {
    await authStore.fetchProfile()
  }
  if (authStore.user?.weight) {
    calcForm.value.rider_weight = authStore.user.weight
  }

  if (isEdit.value) {
    await racesStore.fetchRace(route.params.id)
    const race = racesStore.currentRace
    if (race) {
      form.value = {
        name: race.name,
        date: race.date ? race.date.substring(0, 10) : '',
        type: race.type,
        photo: race.photo || '',
        setup_id: race.setup_id || null,
        bike_name: race.bike_name || '',
        tires: race.tires || '',
        tire_pressure_front: race.tire_pressure_front,
        tire_pressure_rear: race.tire_pressure_rear,
        other_components: race.other_components || '',
        temperature: race.temperature,
        conditions: race.conditions || null,
        wind: race.wind || null,
        road_conditions: race.road_conditions || null,
        nutrition_plan: race.nutrition_plan || '',
        result: race.result || '',
        rating: race.rating || null,
        feelings: race.feelings || '',
        is_completed: race.is_completed,
      }
      useSetup.value = !!race.setup_id
    }
  }
})

watch(useSetup, (val) => {
  if (!val) {
    form.value.setup_id = null
  }
})

function openCalculator() {
  calcResult.value = null
  // Pre-fill surface from race type
  if (['MTB', 'Gravel', 'Cyclocross'].includes(form.value.type)) {
    calcForm.value.surface = 'gravel'
  } else {
    calcForm.value.surface = 'road'
  }
  // Pre-fill conditions from weather
  if (form.value.road_conditions) {
    const mapping = { Dry: 'dry', Wet: 'wet', Mud: 'mud' }
    calcForm.value.conditions = mapping[form.value.road_conditions] || 'dry'
  }
  // Pre-fill tire width from tires string or selected setup
  const tiresStr = useSetup.value && selectedSetup.value ? selectedSetup.value.tires : form.value.tires
  if (tiresStr) {
    const match = tiresStr.match(/(\d{2,3})(?:mm|c|C)?/)
    if (match) calcForm.value.tire_width = parseInt(match[1])
  }
  showCalculator.value = true
}

async function runCalculator() {
  calcLoading.value = true
  try {
    const { data } = await calculatePressure(calcForm.value)
    calcResult.value = data
  } catch (e) {
    // Fallback to offline calculation
    try {
      calcResult.value = calculateTirePressureOffline(calcForm.value)
    } catch {
      error.value = 'Calculator error'
    }
  } finally {
    calcLoading.value = false
  }
}

function applyPressure() {
  if (calcResult.value) {
    form.value.tire_pressure_front = calcResult.value.front_pressure
    form.value.tire_pressure_rear = calcResult.value.rear_pressure
    showCalculator.value = false
  }
}

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    const payload = { ...form.value }

    // Convert types
    if (payload.tire_pressure_front) payload.tire_pressure_front = Number(payload.tire_pressure_front)
    if (payload.tire_pressure_rear) payload.tire_pressure_rear = Number(payload.tire_pressure_rear)
    if (payload.temperature !== null && payload.temperature !== '') payload.temperature = Number(payload.temperature)
    else payload.temperature = null
    if (payload.rating) payload.rating = Number(payload.rating)

    // Handle setup vs manual
    if (useSetup.value && payload.setup_id) {
      payload.setup_id = Number(payload.setup_id)
      delete payload.bike_name
      delete payload.tires
    } else {
      payload.setup_id = null
    }

    // Clean empty strings to null for optional fields
    if (!payload.photo) payload.photo = null
    if (!payload.other_components) payload.other_components = null
    if (!payload.nutrition_plan) payload.nutrition_plan = null
    if (!payload.result) payload.result = null
    if (!payload.feelings) payload.feelings = null
    if (!payload.conditions) payload.conditions = null
    if (!payload.wind) payload.wind = null
    if (!payload.road_conditions) payload.road_conditions = null

    if (isEdit.value) {
      await racesStore.updateRace(route.params.id, payload)
      toast('Race updated')
    } else {
      await racesStore.createRace(payload)
      toast('Race created')
    }
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to save race'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="px-4 pt-6 pb-24 max-w-lg mx-auto">
    <div class="flex items-center justify-between mb-4">
      <button @click="router.back()" class="text-gray-400 text-sm">&larr; Back</button>
      <h1 class="text-lg font-bold">{{ isEdit ? 'Edit Race' : 'New Race' }}</h1>
      <div class="w-10"></div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-5">
      <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

      <!-- Basic Info -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Basic Info</h2>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Race Name *</label>
          <input v-model="form.name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. Spring Classic 2026" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date *</label>
          <input v-model="form.date" type="date" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type *</label>
          <select v-model="form.type" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option v-for="t in raceTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Photo</label>
          <PhotoUpload v-model="form.photo" />
        </div>
      </div>

      <!-- Setup / Bike -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Bike Setup</h2>

        <label class="flex items-center gap-2 cursor-pointer">
          <input v-model="useSetup" type="checkbox" class="w-4 h-4 text-primary rounded border-gray-300 focus:ring-primary" />
          <span class="text-sm text-gray-700">Use saved setup</span>
        </label>

        <div v-if="useSetup">
          <label class="block text-sm font-medium text-gray-700 mb-1">Select Setup</label>
          <select v-model="form.setup_id" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option :value="null" disabled>Choose a setup...</option>
            <option v-for="s in setupsStore.setups" :key="s.id" :value="s.id">{{ s.name }} — {{ s.bike_name }}</option>
          </select>
          <p v-if="selectedSetup" class="text-xs text-gray-400 mt-1">
            {{ selectedSetup.bike_name }} &middot; {{ selectedSetup.tires }}
          </p>
        </div>

        <template v-if="needsManualBike">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Bike Name *</label>
            <input v-model="form.bike_name" type="text" :required="needsManualBike" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. Canyon Aeroad" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tires *</label>
            <input v-model="form.tires" type="text" :required="needsManualBike" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. Continental GP5000 28mm" />
          </div>
        </template>
      </div>

      <!-- Tire Pressure -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Tire Pressure</h2>
          <button type="button" @click="openCalculator" class="flex items-center gap-1 text-primary text-sm font-medium hover:underline">
            <CalculatorIcon class="w-4 h-4" />
            Calculate
          </button>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Front (bar)</label>
            <input v-model="form.tire_pressure_front" type="number" step="0.01" min="0" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="5.5" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Rear (bar)</label>
            <input v-model="form.tire_pressure_rear" type="number" step="0.01" min="0" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="6.0" />
          </div>
        </div>
      </div>

      <!-- Other Components -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Other Components</h2>
        <textarea v-model="form.other_components" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" placeholder="e.g. Saddle height, handlebar position, gearing..."></textarea>
      </div>

      <!-- Weather -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Weather & Conditions</h2>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Temperature (°C)</label>
          <input v-model="form.temperature" type="number" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="20" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Weather</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="c in conditionOptions"
              :key="c"
              type="button"
              @click="form.conditions = form.conditions === c ? null : c"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="form.conditions === c ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-300'"
            >{{ c }}</button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Wind</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="w in windOptions"
              :key="w"
              type="button"
              @click="form.wind = form.wind === w ? null : w"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="form.wind === w ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-300'"
            >{{ w }}</button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Road</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="r in roadOptions"
              :key="r"
              type="button"
              @click="form.road_conditions = form.road_conditions === r ? null : r"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="form.road_conditions === r ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-300'"
            >{{ r }}</button>
          </div>
        </div>
      </div>

      <!-- Nutrition -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Nutrition Plan</h2>
        <textarea v-model="form.nutrition_plan" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" placeholder="e.g. Gel every 30 min, 750ml water per hour..."></textarea>
      </div>

      <!-- Results (for completed races) -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">Results</h2>
          <label class="flex items-center gap-2 cursor-pointer">
            <input v-model="form.is_completed" type="checkbox" class="w-4 h-4 text-primary rounded border-gray-300 focus:ring-primary" />
            <span class="text-sm text-gray-700">Completed</span>
          </label>
        </div>

        <template v-if="form.is_completed">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Result</label>
            <input v-model="form.result" type="text" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="e.g. 3rd place, 2:45:30" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Rating</label>
            <StarRating :rating="form.rating || 0" @update:rating="form.rating = $event" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Feelings</label>
            <textarea v-model="form.feelings" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" placeholder="How did you feel during the race?"></textarea>
          </div>
        </template>
      </div>

      <!-- Submit -->
      <button
        type="submit"
        :disabled="loading"
        class="w-full bg-primary text-white font-medium py-3 rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 shadow-md"
      >
        {{ loading ? 'Saving...' : (isEdit ? 'Save Changes' : 'Create Race') }}
      </button>
    </form>

    <!-- Calculator Modal -->
    <div v-if="showCalculator" class="fixed inset-0 bg-black/40 flex items-end sm:items-center justify-center z-50 px-4 pb-4">
      <div class="bg-white rounded-xl p-5 w-full max-w-sm max-h-[85vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-lg">Tire Pressure Calculator</h3>
          <button @click="showCalculator = false" class="text-gray-400 text-xl leading-none">&times;</button>
        </div>

        <div class="space-y-3">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Rider Weight (kg)</label>
              <input v-model.number="calcForm.rider_weight" type="number" step="0.1" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Bike Weight (kg)</label>
              <input v-model.number="calcForm.bike_weight" type="number" step="0.1" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Tire Width (mm)</label>
            <input v-model.number="calcForm.tire_width" type="number" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Tire Type</label>
            <select v-model="calcForm.tire_type" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
              <option value="clincher">Clincher</option>
              <option value="tubeless">Tubeless</option>
              <option value="tubular">Tubular</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Surface</label>
            <select v-model="calcForm.surface" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
              <option value="road">Road</option>
              <option value="gravel">Gravel</option>
              <option value="mixed">Mixed</option>
              <option value="cobblestone">Cobblestone</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Conditions</label>
            <select v-model="calcForm.conditions" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
              <option value="dry">Dry</option>
              <option value="wet">Wet</option>
              <option value="mud">Mud</option>
              <option value="snow">Snow</option>
            </select>
          </div>

          <button
            type="button"
            @click="runCalculator"
            :disabled="calcLoading"
            class="w-full bg-primary text-white font-medium py-2.5 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
          >
            {{ calcLoading ? 'Calculating...' : 'Calculate' }}
          </button>

          <div v-if="calcResult" class="bg-green-50 rounded-lg p-4 space-y-2">
            <div class="grid grid-cols-2 gap-3 text-center">
              <div>
                <p class="text-xs text-gray-500">Front</p>
                <p class="text-2xl font-bold text-primary">{{ calcResult.front_pressure }}</p>
                <p class="text-xs text-gray-400">bar</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">Rear</p>
                <p class="text-2xl font-bold text-primary">{{ calcResult.rear_pressure }}</p>
                <p class="text-xs text-gray-400">bar</p>
              </div>
            </div>

            <div v-if="calcResult.recommendations?.length" class="border-t border-green-200 pt-2">
              <p v-for="(rec, i) in calcResult.recommendations" :key="i" class="text-xs text-gray-600 leading-relaxed">{{ rec }}</p>
            </div>

            <button
              type="button"
              @click="applyPressure"
              class="w-full bg-primary text-white font-medium py-2 rounded-lg hover:bg-primary-dark transition-colors mt-2"
            >
              Apply to Race
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
