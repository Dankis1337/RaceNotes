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
import { useI18n } from '../utils/i18n'

const { t } = useI18n()
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

const conditionKeys = { Sunny: 'sunny', Cloudy: 'cloudy', Rain: 'rain', Snow: 'snow' }
const windKeys = { None: 'wind_none', Light: 'wind_light', Moderate: 'wind_moderate', Strong: 'wind_strong' }
const roadKeys = { Dry: 'road_dry', Wet: 'road_wet', Mud: 'road_mud' }
const typeKeys = { Road: 'type_road', MTB: 'type_mtb', Gravel: 'type_gravel', Cyclocross: 'type_cyclocross', Track: 'type_track' }

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
  if (['MTB', 'Gravel', 'Cyclocross'].includes(form.value.type)) {
    calcForm.value.surface = 'gravel'
  } else {
    calcForm.value.surface = 'road'
  }
  if (form.value.road_conditions) {
    const mapping = { Dry: 'dry', Wet: 'wet', Mud: 'mud' }
    calcForm.value.conditions = mapping[form.value.road_conditions] || 'dry'
  }
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
    try {
      calcResult.value = calculateTirePressureOffline(calcForm.value)
    } catch {
      error.value = t('calculator_error')
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

    if (payload.tire_pressure_front) payload.tire_pressure_front = Number(payload.tire_pressure_front)
    if (payload.tire_pressure_rear) payload.tire_pressure_rear = Number(payload.tire_pressure_rear)
    if (payload.temperature !== null && payload.temperature !== '') payload.temperature = Number(payload.temperature)
    else payload.temperature = null
    if (payload.rating) payload.rating = Number(payload.rating)

    if (useSetup.value && payload.setup_id) {
      payload.setup_id = Number(payload.setup_id)
      delete payload.bike_name
      delete payload.tires
    } else {
      payload.setup_id = null
    }

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
      toast(t('race_updated'))
    } else {
      await racesStore.createRace(payload)
      toast(t('race_created'))
    }
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || t('failed_save_race')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="px-4 pt-6 pb-24 max-w-lg mx-auto lg:max-w-2xl">
    <div class="flex items-center justify-between mb-4">
      <button @click="router.back()" class="text-gray-400 text-sm">&larr; {{ t('back') }}</button>
      <h1 class="text-lg font-bold">{{ isEdit ? t('edit_race') : t('new_race') }}</h1>
      <div class="w-10"></div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-5">
      <div v-if="error" class="bg-red-50 text-red-600 text-sm rounded-lg p-3">{{ error }}</div>

      <!-- Basic Info -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('basic_info') }}</h2>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('race_name') }} *</label>
          <input v-model="form.name" type="text" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('race_name_placeholder')" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('date') }} *</label>
          <input v-model="form.date" type="date" required class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('type') }} *</label>
          <select v-model="form.type" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option v-for="tp in raceTypes" :key="tp" :value="tp">{{ t(typeKeys[tp]) }}</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('photo') }}</label>
          <PhotoUpload v-model="form.photo" />
        </div>
      </div>

      <!-- Setup / Bike -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('bike_setup') }}</h2>

        <label class="flex items-center gap-2 cursor-pointer">
          <input v-model="useSetup" type="checkbox" class="w-4 h-4 text-primary rounded border-gray-300 focus:ring-primary" />
          <span class="text-sm text-gray-700">{{ t('use_saved_setup') }}</span>
        </label>

        <div v-if="useSetup">
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('select_setup') }}</label>
          <select v-model="form.setup_id" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
            <option :value="null" disabled>{{ t('choose_setup') }}</option>
            <option v-for="s in setupsStore.setups" :key="s.id" :value="s.id">{{ s.name }} — {{ s.bike_name }}</option>
          </select>
          <p v-if="selectedSetup" class="text-xs text-gray-400 mt-1">
            {{ selectedSetup.bike_name }} &middot; {{ selectedSetup.tires }}
          </p>
        </div>

        <template v-if="needsManualBike">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('bike_name') }} *</label>
            <input v-model="form.bike_name" type="text" :required="needsManualBike" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('bike_name_placeholder')" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('tires') }} *</label>
            <input v-model="form.tires" type="text" :required="needsManualBike" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('tires_placeholder')" />
          </div>
        </template>
      </div>

      <!-- Tire Pressure -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('tire_pressure') }}</h2>
          <button type="button" @click="openCalculator" class="flex items-center gap-1 text-primary text-sm font-medium hover:underline">
            <CalculatorIcon class="w-4 h-4" />
            {{ t('calculate') }}
          </button>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('front_bar') }}</label>
            <input v-model="form.tire_pressure_front" type="number" step="0.01" min="0" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="5.5" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('rear_bar') }}</label>
            <input v-model="form.tire_pressure_rear" type="number" step="0.01" min="0" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="6.0" />
          </div>
        </div>
      </div>

      <!-- Other Components -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('other_components_section') }}</h2>
        <textarea v-model="form.other_components" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" :placeholder="t('other_components_placeholder')"></textarea>
      </div>

      <!-- Weather -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('weather_conditions') }}</h2>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('temperature') }}</label>
          <input v-model="form.temperature" type="number" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" placeholder="20" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">{{ t('weather') }}</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="c in conditionOptions"
              :key="c"
              type="button"
              @click="form.conditions = form.conditions === c ? null : c"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="form.conditions === c ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-300'"
            >{{ t(conditionKeys[c]) }}</button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">{{ t('wind_label') }}</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="w in windOptions"
              :key="w"
              type="button"
              @click="form.wind = form.wind === w ? null : w"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="form.wind === w ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-300'"
            >{{ t(windKeys[w]) }}</button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">{{ t('road') }}</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="r in roadOptions"
              :key="r"
              type="button"
              @click="form.road_conditions = form.road_conditions === r ? null : r"
              class="px-3 py-1.5 rounded-full text-sm border transition-colors"
              :class="form.road_conditions === r ? 'bg-primary text-white border-primary' : 'bg-white text-gray-600 border-gray-300'"
            >{{ t(roadKeys[r]) }}</button>
          </div>
        </div>
      </div>

      <!-- Nutrition -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('nutrition_plan_section') }}</h2>
        <textarea v-model="form.nutrition_plan" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" :placeholder="t('nutrition_placeholder')"></textarea>
      </div>

      <!-- Results (for completed races) -->
      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-semibold text-sm text-gray-400 uppercase tracking-wide">{{ t('results') }}</h2>
          <label class="flex items-center gap-2 cursor-pointer">
            <input v-model="form.is_completed" type="checkbox" class="w-4 h-4 text-primary rounded border-gray-300 focus:ring-primary" />
            <span class="text-sm text-gray-700">{{ t('is_completed') }}</span>
          </label>
        </div>

        <template v-if="form.is_completed">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('result') }}</label>
            <input v-model="form.result" type="text" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none" :placeholder="t('result_placeholder')" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('rating') }}</label>
            <StarRating :rating="form.rating || 0" @update:rating="form.rating = $event" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('feelings') }}</label>
            <textarea v-model="form.feelings" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:ring-2 focus:ring-primary focus:border-transparent outline-none resize-none" :placeholder="t('feelings_placeholder')"></textarea>
          </div>
        </template>
      </div>

      <!-- Submit -->
      <button
        type="submit"
        :disabled="loading"
        class="w-full bg-primary text-white font-medium py-3 rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 shadow-md"
      >
        {{ loading ? t('saving') : (isEdit ? t('save_changes') : t('create_race')) }}
      </button>
    </form>

    <!-- Calculator Modal -->
    <div v-if="showCalculator" class="fixed inset-0 bg-black/40 flex items-end sm:items-center justify-center z-50 px-4 pb-4">
      <div class="bg-white rounded-xl p-5 w-full max-w-sm max-h-[85vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-lg">{{ t('tire_pressure_calculator') }}</h3>
          <button @click="showCalculator = false" class="text-gray-400 text-xl leading-none">&times;</button>
        </div>

        <div class="space-y-3">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('rider_weight') }}</label>
              <input v-model.number="calcForm.rider_weight" type="number" step="0.1" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('bike_weight') }}</label>
              <input v-model.number="calcForm.bike_weight" type="number" step="0.1" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('tire_width') }}</label>
            <input v-model.number="calcForm.tire_width" type="number" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none" />
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('tire_type') }}</label>
            <select v-model="calcForm.tire_type" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
              <option value="clincher">{{ t('clincher') }}</option>
              <option value="tubeless">{{ t('tubeless') }}</option>
              <option value="tubular">{{ t('tubular') }}</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('surface') }}</label>
            <select v-model="calcForm.surface" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
              <option value="road">{{ t('surface_road') }}</option>
              <option value="gravel">{{ t('surface_gravel') }}</option>
              <option value="mixed">{{ t('surface_mixed') }}</option>
              <option value="cobblestone">{{ t('surface_cobblestone') }}</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('conditions') }}</label>
            <select v-model="calcForm.conditions" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:border-transparent outline-none bg-white">
              <option value="dry">{{ t('cond_dry') }}</option>
              <option value="wet">{{ t('cond_wet') }}</option>
              <option value="mud">{{ t('cond_mud') }}</option>
              <option value="snow">{{ t('cond_snow') }}</option>
            </select>
          </div>

          <button
            type="button"
            @click="runCalculator"
            :disabled="calcLoading"
            class="w-full bg-primary text-white font-medium py-2.5 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
          >
            {{ calcLoading ? t('calculating') : t('calculate') }}
          </button>

          <div v-if="calcResult" class="bg-green-50 rounded-lg p-4 space-y-2">
            <div class="grid grid-cols-2 gap-3 text-center">
              <div>
                <p class="text-xs text-gray-500">{{ t('front') }}</p>
                <p class="text-2xl font-bold text-primary">{{ calcResult.front_pressure }}</p>
                <p class="text-xs text-gray-400">bar</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">{{ t('rear') }}</p>
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
              {{ t('apply_to_race') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
