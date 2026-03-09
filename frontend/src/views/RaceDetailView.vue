<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRacesStore } from '../stores/races'
import StarRating from '../components/StarRating.vue'
import WeatherBadge from '../components/WeatherBadge.vue'
import { PencilIcon, TrashIcon } from '@heroicons/vue/24/outline'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const racesStore = useRacesStore()
const { toast } = useToast()
const showConfirm = ref(false)

onMounted(() => {
  racesStore.fetchRace(route.params.id)
})

const race = () => racesStore.currentRace

async function handleDelete() {
  await racesStore.removeRace(route.params.id)
  toast('Race deleted')
  router.push('/')
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto">
    <div v-if="racesStore.loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else-if="race()">
      <div class="flex items-center justify-between mb-4">
        <button @click="router.back()" class="text-gray-400 text-sm">&larr; Back</button>
        <div class="flex gap-2">
          <router-link :to="`/races/${race().id}/edit`" class="p-2 text-gray-500 hover:text-primary">
            <PencilIcon class="w-5 h-5" />
          </router-link>
          <button @click="showConfirm = true" class="p-2 text-gray-500 hover:text-red-500">
            <TrashIcon class="w-5 h-5" />
          </button>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-md overflow-hidden space-y-4">
        <img v-if="race().photo" :src="race().photo" class="w-full h-48 object-cover" />
        <div class="px-5" :class="{ 'pt-4': !race().photo }">
          <div class="flex items-center gap-2 mb-1">
            <h1 class="text-xl font-bold">{{ race().name }}</h1>
            <span
              class="text-xs font-medium px-2 py-0.5 rounded-full"
              :class="race().is_completed ? 'bg-primary/15 text-primary' : 'bg-amber-100 text-amber-700'"
            >
              {{ race().is_completed ? 'Completed' : 'Planned' }}
            </span>
          </div>
          <p class="text-gray-500 text-sm">{{ formatDate(race().date) }}</p>
          <span class="inline-block mt-1 bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded">{{ race().type }}</span>
        </div>

        <div class="grid grid-cols-2 gap-3 text-sm px-5">
          <div v-if="race().bike_name">
            <p class="text-gray-400 text-xs">Bike</p>
            <p class="font-medium">{{ race().bike_name }}</p>
          </div>
          <div v-if="race().tires">
            <p class="text-gray-400 text-xs">Tires</p>
            <p class="font-medium">{{ race().tires }}</p>
          </div>
          <div v-if="race().setup">
            <p class="text-gray-400 text-xs">Setup</p>
            <router-link :to="`/setups/${race().setup.id}`" class="font-medium text-primary">{{ race().setup.name }}</router-link>
          </div>
        </div>

        <div v-if="race().tire_pressure_front || race().tire_pressure_rear" class="grid grid-cols-2 gap-3 text-sm px-5">
          <div v-if="race().tire_pressure_front">
            <p class="text-gray-400 text-xs">Front Pressure</p>
            <p class="font-medium">{{ race().tire_pressure_front }} bar</p>
          </div>
          <div v-if="race().tire_pressure_rear">
            <p class="text-gray-400 text-xs">Rear Pressure</p>
            <p class="font-medium">{{ race().tire_pressure_rear }} bar</p>
          </div>
        </div>

        <div v-if="race().other_components" class="px-5">
          <p class="text-gray-400 text-xs">Other Components</p>
          <p class="text-sm whitespace-pre-line">{{ race().other_components }}</p>
        </div>

        <div v-if="race().conditions || race().wind || race().road_conditions || race().temperature !== null" class="space-y-2 px-5">
          <p class="text-gray-400 text-xs">Weather & Conditions</p>
          <div class="flex flex-wrap gap-2">
            <WeatherBadge v-if="race().conditions" :conditions="race().conditions" />
            <span v-if="race().temperature !== null" class="bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded">{{ race().temperature }}°C</span>
            <span v-if="race().wind" class="bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded">Wind: {{ race().wind }}</span>
            <span v-if="race().road_conditions" class="bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded">{{ race().road_conditions }}</span>
          </div>
        </div>

        <div v-if="race().nutrition_plan" class="px-5">
          <p class="text-gray-400 text-xs">Nutrition Plan</p>
          <p class="text-sm whitespace-pre-line">{{ race().nutrition_plan }}</p>
        </div>

        <div v-if="race().is_completed" class="border-t pt-4 space-y-3 px-5 pb-5">
          <h2 class="font-semibold">Results</h2>
          <div v-if="race().result">
            <p class="text-gray-400 text-xs">Result</p>
            <p class="font-medium">{{ race().result }}</p>
          </div>
          <div v-if="race().rating">
            <p class="text-gray-400 text-xs mb-1">Rating</p>
            <StarRating :rating="race().rating" :readonly="true" />
          </div>
          <div v-if="race().feelings">
            <p class="text-gray-400 text-xs">Feelings</p>
            <p class="text-sm whitespace-pre-line">{{ race().feelings }}</p>
          </div>
        </div>

        <div v-if="!race().is_completed" class="border-t pt-4 px-5 pb-5">
          <router-link
            :to="`/races/${race().id}/edit`"
            class="block w-full text-center bg-primary text-white font-medium py-2.5 rounded-lg hover:bg-primary-dark transition-colors"
          >
            Fill in Results
          </router-link>
        </div>
      </div>
    </template>

    <!-- Delete confirm modal -->
    <div v-if="showConfirm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 px-4">
      <div class="bg-white rounded-xl p-5 w-full max-w-sm">
        <h3 class="font-semibold text-lg mb-2">Delete Race?</h3>
        <p class="text-gray-500 text-sm mb-4">This action cannot be undone.</p>
        <div class="flex gap-3">
          <button @click="showConfirm = false" class="flex-1 py-2 rounded-lg border border-gray-300 text-gray-600">Cancel</button>
          <button @click="handleDelete" class="flex-1 py-2 rounded-lg bg-red-500 text-white">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>
