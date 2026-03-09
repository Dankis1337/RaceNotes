<script setup>
import { ref } from 'vue'
import { uploadFile } from '../api/upload'
import { CameraIcon, XMarkIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  modelValue: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

const uploading = ref(false)
const error = ref('')

async function handleFileChange(e) {
  const file = e.target.files?.[0]
  if (!file) return

  if (file.size > 10 * 1024 * 1024) {
    error.value = 'File must be under 10 MB'
    return
  }

  error.value = ''
  uploading.value = true
  try {
    const { data } = await uploadFile(file)
    emit('update:modelValue', data.url)
  } catch (err) {
    error.value = err.response?.data?.error || 'Upload failed'
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

function removePhoto() {
  emit('update:modelValue', '')
}
</script>

<template>
  <div>
    <div v-if="modelValue" class="relative">
      <img :src="modelValue" class="w-full h-40 object-cover rounded-lg" @error="removePhoto" />
      <button
        type="button"
        @click="removePhoto"
        class="absolute top-2 right-2 bg-black/50 text-white rounded-full p-1 hover:bg-black/70 transition-colors"
      >
        <XMarkIcon class="w-5 h-5" />
      </button>
    </div>

    <label
      v-else
      class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:border-primary transition-colors"
      :class="{ 'opacity-50 pointer-events-none': uploading }"
    >
      <CameraIcon class="w-8 h-8 text-gray-400 mb-1" />
      <span class="text-sm text-gray-500">{{ uploading ? 'Uploading...' : 'Tap to add photo' }}</span>
      <input
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        class="hidden"
        @change="handleFileChange"
        :disabled="uploading"
      />
    </label>

    <p v-if="error" class="text-red-500 text-xs mt-1">{{ error }}</p>
  </div>
</template>
