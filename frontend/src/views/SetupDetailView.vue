<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSetupsStore } from '../stores/setups'
import { PencilIcon, TrashIcon } from '@heroicons/vue/24/outline'

const route = useRoute()
const router = useRouter()
const setupsStore = useSetupsStore()
const showConfirm = ref(false)

onMounted(() => {
  setupsStore.fetchSetup(route.params.id)
})

const setup = () => setupsStore.currentSetup

async function handleDelete() {
  await setupsStore.removeSetup(route.params.id)
  router.push('/setups')
}
</script>

<template>
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto">
    <div v-if="setupsStore.loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else-if="setup()">
      <div class="flex items-center justify-between mb-4">
        <button @click="router.back()" class="text-gray-400 text-sm">&larr; Back</button>
        <div class="flex gap-2">
          <router-link :to="`/setups/${setup().id}/edit`" class="p-2 text-gray-500 hover:text-primary">
            <PencilIcon class="w-5 h-5" />
          </router-link>
          <button @click="showConfirm = true" class="p-2 text-gray-500 hover:text-red-500">
            <TrashIcon class="w-5 h-5" />
          </button>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-md p-5 space-y-4">
        <h1 class="text-xl font-bold">{{ setup().name }}</h1>

        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p class="text-gray-400 text-xs">Bike</p>
            <p class="font-medium">{{ setup().bike_name }}</p>
          </div>
          <div>
            <p class="text-gray-400 text-xs">Tires</p>
            <p class="font-medium">{{ setup().tires }}</p>
          </div>
        </div>

        <div v-if="setup().components_description">
          <p class="text-gray-400 text-xs">Components</p>
          <p class="text-sm whitespace-pre-line">{{ setup().components_description }}</p>
        </div>
      </div>
    </template>

    <!-- Delete confirm modal -->
    <div v-if="showConfirm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 px-4">
      <div class="bg-white rounded-xl p-5 w-full max-w-sm">
        <h3 class="font-semibold text-lg mb-2">Delete Setup?</h3>
        <p class="text-gray-500 text-sm mb-4">This action cannot be undone.</p>
        <div class="flex gap-3">
          <button @click="showConfirm = false" class="flex-1 py-2 rounded-lg border border-gray-300 text-gray-600">Cancel</button>
          <button @click="handleDelete" class="flex-1 py-2 rounded-lg bg-red-500 text-white">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>
