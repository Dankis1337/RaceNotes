<script setup>
import { ref, watch } from 'vue'
import { CheckCircleIcon, ExclamationCircleIcon, XMarkIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  message: { type: String, default: '' },
  type: { type: String, default: 'success' },
  show: { type: Boolean, default: false },
})

const emit = defineEmits(['close'])

watch(() => props.show, (val) => {
  if (val) {
    setTimeout(() => emit('close'), 3000)
  }
})
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="translate-y-2 opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-2 opacity-0"
  >
    <div
      v-if="show"
      class="fixed top-4 left-4 right-4 z-[60] max-w-lg mx-auto"
    >
      <div
        class="flex items-center gap-2 rounded-xl px-4 py-3 shadow-lg"
        :class="type === 'error' ? 'bg-red-500 text-white' : 'bg-primary text-white'"
      >
        <CheckCircleIcon v-if="type === 'success'" class="w-5 h-5 shrink-0" />
        <ExclamationCircleIcon v-else class="w-5 h-5 shrink-0" />
        <span class="text-sm flex-1">{{ message }}</span>
        <button @click="emit('close')" class="shrink-0">
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
  </Transition>
</template>
