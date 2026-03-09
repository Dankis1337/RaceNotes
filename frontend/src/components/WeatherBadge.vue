<script setup>
import { computed } from 'vue'

const props = defineProps({ conditions: String })

const badge = computed(() => {
  const map = {
    Sunny: { icon: 'sun', bg: 'bg-yellow-100 text-yellow-700' },
    Cloudy: { icon: 'cloud', bg: 'bg-gray-100 text-gray-600' },
    Rain: { icon: 'rain', bg: 'bg-blue-100 text-blue-600' },
    Snow: { icon: 'snow', bg: 'bg-sky-100 text-sky-600' },
  }
  return map[props.conditions] || { icon: '', bg: 'bg-gray-100 text-gray-500' }
})

const icons = { sun: '\u2600\uFE0F', cloud: '\u2601\uFE0F', rain: '\uD83C\uDF27\uFE0F', snow: '\u2744\uFE0F' }
</script>

<template>
  <span :class="['px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1', badge.bg]">
    <span v-if="badge.icon">{{ icons[badge.icon] }}</span>
    {{ conditions }}
  </span>
</template>
