<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRacesStore } from '../stores/races'
import { useSetupsStore } from '../stores/setups'
import { PlusIcon } from '@heroicons/vue/24/solid'
import RaceCard from '../components/RaceCard.vue'
import FilterChips from '../components/FilterChips.vue'

const auth = useAuthStore()
const racesStore = useRacesStore()
const setupsStore = useSetupsStore()

const statusFilter = ref('')
const typeFilter = ref('')
const setupFilter = ref('')

const setupOptions = computed(() => {
  const opts = [{ label: 'All', value: '' }]
  setupsStore.setups.forEach(s => opts.push({ label: s.name, value: String(s.id) }))
  return opts
})

const statusOptions = [
  { label: 'All', value: '' },
  { label: 'Planned', value: 'false' },
  { label: 'Completed', value: 'true' },
]

const typeOptions = [
  { label: 'All', value: '' },
  { label: 'Road', value: 'Road' },
  { label: 'MTB', value: 'MTB' },
  { label: 'Gravel', value: 'Gravel' },
  { label: 'Cyclocross', value: 'Cyclocross' },
  { label: 'Track', value: 'Track' },
]

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
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto">
    <h1 class="text-2xl font-bold mb-1">
      Hi, {{ auth.user?.name || 'Rider' }}!
    </h1>
    <p class="text-gray-500 text-sm mb-4">Your race diary</p>

    <div
      v-if="yesterdayIncomplete.length"
      class="bg-accent rounded-xl p-3 mb-4 text-sm"
    >
      You have {{ yesterdayIncomplete.length }} race(s) from yesterday without results.
      <router-link
        v-if="yesterdayIncomplete.length === 1"
        :to="`/races/${yesterdayIncomplete[0].id}`"
        class="text-primary font-medium underline ml-1"
      >
        Fill in results
      </router-link>
    </div>

    <div class="space-y-3 mb-4">
      <FilterChips label="Status" :options="statusOptions" :modelValue="statusFilter" @update:modelValue="onStatusChange" />
      <FilterChips label="Type" :options="typeOptions" :modelValue="typeFilter" @update:modelValue="onTypeChange" />
      <FilterChips v-if="setupOptions.length > 1" label="Setup" :options="setupOptions" :modelValue="setupFilter" @update:modelValue="onSetupChange" />
    </div>

    <div v-if="racesStore.loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="racesStore.races.length === 0" class="text-center py-12">
      <p class="text-gray-400 text-lg mb-2">No races yet</p>
      <p class="text-gray-400 text-sm">Tap + to add your first race</p>
    </div>

    <div v-else class="space-y-3">
      <RaceCard v-for="race in racesStore.races" :key="race.id" :race="race" />
    </div>

    <router-link
      to="/races/new"
      class="fixed right-4 w-14 h-14 bg-primary text-white rounded-full shadow-lg flex items-center justify-center hover:bg-primary-dark transition-colors z-30"
      style="bottom: calc(5rem + env(safe-area-inset-bottom, 0px))"
    >
      <PlusIcon class="w-7 h-7" />
    </router-link>
  </div>
</template>
