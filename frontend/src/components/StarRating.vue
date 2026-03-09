<script setup>
import { StarIcon } from '@heroicons/vue/24/solid'
import { StarIcon as StarOutline } from '@heroicons/vue/24/outline'

const props = defineProps({
  rating: { type: Number, default: 0 },
  readonly: { type: Boolean, default: false },
  size: { type: String, default: 'md' },
})

const emit = defineEmits(['update:rating'])

const sizeClass = { sm: 'w-4 h-4', md: 'w-6 h-6' }

function select(val) {
  if (!props.readonly) emit('update:rating', val)
}
</script>

<template>
  <div class="flex gap-0.5">
    <button
      v-for="i in 5"
      :key="i"
      type="button"
      :disabled="readonly"
      @click="select(i)"
      class="disabled:cursor-default"
    >
      <component
        :is="i <= rating ? StarIcon : StarOutline"
        :class="[sizeClass[size], i <= rating ? 'text-amber-400' : 'text-gray-300']"
      />
    </button>
  </div>
</template>
