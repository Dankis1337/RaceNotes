<script setup>
defineProps({
  label: String,
  options: Array,
  modelValue: String,
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div>
    <span class="text-xs font-medium text-gray-500 mr-2">{{ label }}:</span>
    <div class="inline-flex flex-wrap gap-1.5">
      <button
        v-for="opt in options"
        :key="opt.value"
        type="button"
        @click="emit('update:modelValue', opt.value)"
        class="px-3 py-1 rounded-full text-xs font-medium transition-colors"
        :class="modelValue === opt.value
          ? 'bg-primary text-white'
          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
      >
        {{ opt.label }}
      </button>
    </div>
  </div>
</template>
