<script setup>
import { onMounted } from 'vue'
import { useSetupsStore } from '../stores/setups'
import { PlusIcon, WrenchScrewdriverIcon } from '@heroicons/vue/24/solid'

const setupsStore = useSetupsStore()

onMounted(() => {
  setupsStore.fetchSetups()
})
</script>

<template>
  <div class="px-4 pt-6 pb-20 max-w-lg mx-auto">
    <h1 class="text-2xl font-bold mb-4">My Setups</h1>

    <div v-if="setupsStore.loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="setupsStore.setups.length === 0" class="text-center py-12">
      <WrenchScrewdriverIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-gray-400 text-lg mb-1">No setups yet</p>
      <p class="text-gray-400 text-sm">Create your first bike setup</p>
    </div>

    <div v-else class="space-y-3">
      <router-link
        v-for="setup in setupsStore.setups"
        :key="setup.id"
        :to="`/setups/${setup.id}`"
        class="block bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow"
      >
        <img v-if="setup.photo" :src="setup.photo" class="w-full h-32 object-cover" />
        <div class="p-4">
        <h3 class="font-semibold text-base mb-1">{{ setup.name }}</h3>
        <div class="flex items-center gap-2 text-sm text-gray-500">
          <span>{{ setup.bike_name }}</span>
          <span class="text-gray-300">&middot;</span>
          <span>{{ setup.tires }}</span>
        </div>
        <p v-if="setup.components_description" class="text-xs text-gray-400 mt-1 line-clamp-1">
          {{ setup.components_description }}
        </p>
        </div>
      </router-link>
    </div>

    <router-link
      to="/setups/new"
      class="fixed right-4 w-14 h-14 bg-primary text-white rounded-full shadow-lg flex items-center justify-center hover:bg-primary-dark transition-colors z-30"
      style="bottom: calc(5rem + env(safe-area-inset-bottom, 0px))"
    >
      <PlusIcon class="w-7 h-7" />
    </router-link>
  </div>
</template>
