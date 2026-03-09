<script setup>
import { computed } from 'vue'
import StarRating from './StarRating.vue'
import WeatherBadge from './WeatherBadge.vue'

const props = defineProps({ race: Object })

const formattedDate = computed(() => {
  if (!props.race.date) return ''
  return new Date(props.race.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
})
</script>

<template>
  <router-link :to="`/races/${race.id}`" class="block bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow">
    <img v-if="race.photo" :src="race.photo" class="w-full h-32 object-cover" />
    <div class="p-4">
    <div class="flex items-start justify-between mb-2">
      <div>
        <h3 class="font-semibold text-base">{{ race.name }}</h3>
        <p class="text-gray-400 text-xs">{{ formattedDate }}</p>
      </div>
      <span
        class="text-xs font-medium px-2 py-0.5 rounded-full"
        :class="race.is_completed ? 'bg-primary/15 text-primary' : 'bg-amber-100 text-amber-700'"
      >
        {{ race.is_completed ? 'Done' : 'Planned' }}
      </span>
    </div>

    <div class="flex items-center gap-2 flex-wrap text-xs text-gray-500">
      <span class="bg-gray-100 px-2 py-0.5 rounded">{{ race.type }}</span>
      <WeatherBadge v-if="race.conditions" :conditions="race.conditions" />
      <span v-if="race.setup" class="truncate">{{ race.setup.name }}</span>
    </div>

    <div v-if="race.is_completed" class="mt-2 flex items-center justify-between">
      <span v-if="race.result" class="text-sm font-medium">{{ race.result }}</span>
      <StarRating v-if="race.rating" :rating="race.rating" :readonly="true" size="sm" />
    </div>
    </div>
  </router-link>
</template>
