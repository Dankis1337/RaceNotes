<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRacesStore } from '../stores/races'
import { useSetupsStore } from '../stores/setups'
import { PlusIcon } from '@heroicons/vue/24/solid'
import RaceCard from '../components/RaceCard.vue'
import FilterChips from '../components/FilterChips.vue'
import { useI18n } from '../utils/i18n'

const { t } = useI18n()
const auth = useAuthStore()
const racesStore = useRacesStore()
const setupsStore = useSetupsStore()

const statusFilter = ref('')
const typeFilter = ref('')
const setupFilter = ref('')
const dismissedReminder = ref(false)

const setupOptions = computed(() => {
  const opts = [{ label: t('all'), value: '' }]
  setupsStore.setups.forEach(s => opts.push({ label: s.name, value: String(s.id) }))
  return opts
})

const statusOptions = computed(() => [
  { label: t('all'), value: '' },
  { label: t('planned'), value: 'false' },
  { label: t('completed'), value: 'true' },
])

const typeOptions = computed(() => [
  { label: t('all'), value: '' },
  { label: t('type_road'), value: 'Road' },
  { label: t('type_mtb'), value: 'MTB' },
  { label: t('type_gravel'), value: 'Gravel' },
  { label: t('type_cyclocross'), value: 'Cyclocross' },
  { label: t('type_track'), value: 'Track' },
])

const yesterdayIncomplete = computed(() => {
  const yesterday = new Date()
  yesterday.setDate(yesterday.getDate() - 1)
  const yStr = yesterday.toISOString().split('T')[0]
  return racesStore.races.filter(r => !r.is_completed && r.date === yStr)
})

async function loadRaces() {
  const params = {}
  if (statusFilter.value) params.is_completed = statusFilter.value
  if (typeFilter.value) params.type = typeFilter.value
  if (setupFilter.value) params.setup_id = setupFilter.value
  await racesStore.fetchRaces(params)
}

onMounted(async () => {
  await auth.fetchProfile()
  await setupsStore.fetchSetups()
  await loadRaces()
})

function onStatusChange(val) {
  statusFilter.value = val
  loadRaces()
}

function onTypeChange(val) {
  typeFilter.value = val
  loadRaces()
}

function onSetupChange(val) {
  setupFilter.value = val
  loadRaces()
}
</script>

<template>
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto lg:max-w-4xl">
    <h1 class="text-2xl font-bold mb-1">
      {{ t('hi_rider', { name: auth.user?.name || 'Rider' }) }}
    </h1>
    <p class="text-gray-500 text-sm mb-4">{{ t('your_race_diary') }}</p>

    <div
      v-if="yesterdayIncomplete.length && !dismissedReminder"
      class="bg-accent rounded-xl p-3 mb-4 text-sm"
    >
      {{ t('yesterday_races', { n: yesterdayIncomplete.length }) }}
      <div class="flex gap-3 mt-2">
        <router-link
          v-if="yesterdayIncomplete.length === 1"
          :to="`/races/${yesterdayIncomplete[0].id}`"
          class="text-primary font-medium underline"
        >
          {{ t('fill_results') }}
        </router-link>
        <button @click="dismissedReminder = true" class="text-gray-500 font-medium">
          {{ t('later') }}
        </button>
      </div>
    </div>

    <div class="space-y-3 mb-4">
      <FilterChips :label="t('status')" :options="statusOptions" :modelValue="statusFilter" @update:modelValue="onStatusChange" />
      <FilterChips :label="t('type')" :options="typeOptions" :modelValue="typeFilter" @update:modelValue="onTypeChange" />
      <FilterChips v-if="setupOptions.length > 1" :label="t('setup')" :options="setupOptions" :modelValue="setupFilter" @update:modelValue="onSetupChange" />
    </div>

    <div v-if="racesStore.loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="racesStore.races.length === 0" class="text-center py-12">
      <p class="text-gray-400 text-lg mb-2">{{ t('no_races_yet') }}</p>
      <p class="text-gray-400 text-sm">{{ t('tap_plus_race') }}</p>
    </div>

    <div v-else class="space-y-3 lg:grid lg:grid-cols-2 lg:gap-4 lg:space-y-0">
      <RaceCard v-for="race in racesStore.races" :key="race.id" :race="race" />
    </div>

    <router-link
      to="/races/new"
      class="fixed right-4 w-14 h-14 bg-primary text-white rounded-full shadow-lg flex items-center justify-center hover:bg-primary-dark transition-colors z-30 lg:right-8"
      style="bottom: calc(5rem + env(safe-area-inset-bottom, 0px))"
    >
      <PlusIcon class="w-7 h-7" />
    </router-link>
  </div>
</template>
