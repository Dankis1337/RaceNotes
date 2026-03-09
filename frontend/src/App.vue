<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Toast from './components/Toast.vue'
import { useToast } from './composables/useToast'
import { useOfflineSync } from './composables/useOfflineSync'
import {
  FlagIcon,
  WrenchScrewdriverIcon,
  CalculatorIcon,
  UserCircleIcon,
} from '@heroicons/vue/24/outline'
import {
  FlagIcon as FlagSolid,
  WrenchScrewdriverIcon as WrenchSolid,
  CalculatorIcon as CalculatorSolid,
  UserCircleIcon as UserSolid,
} from '@heroicons/vue/24/solid'

const route = useRoute()
const { message, type, show, close } = useToast()
const { pendingCount, syncing, isOnline } = useOfflineSync()

const showNav = computed(() => {
  return !['Login', 'Register'].includes(route.name)
})

const navItems = [
  { to: '/', name: 'Races', label: 'Races', icon: FlagIcon, iconActive: FlagSolid },
  { to: '/setups', name: 'Setups', label: 'Setups', icon: WrenchScrewdriverIcon, iconActive: WrenchSolid },
  { to: '/calculator', name: 'Calculator', label: 'Calculator', icon: CalculatorIcon, iconActive: CalculatorSolid },
  { to: '/profile', name: 'Profile', label: 'Profile', icon: UserCircleIcon, iconActive: UserSolid },
]

function isActive(item) {
  if (item.to === '/') return route.path === '/' || route.path.startsWith('/races')
  return route.path.startsWith(item.to)
}
</script>

<template>
  <Toast :message="message" :type="type" :show="show" @close="close" />

  <div v-if="!isOnline || pendingCount > 0" class="fixed top-0 left-0 right-0 z-50 text-center text-xs py-1 font-medium" :class="!isOnline ? 'bg-amber-400 text-amber-900' : 'bg-blue-400 text-white'">
    <template v-if="!isOnline">Offline mode</template>
    <template v-else-if="syncing">Syncing...</template>
    <template v-else>{{ pendingCount }} pending sync</template>
  </div>

  <router-view />

  <nav
    v-if="showNav"
    class="fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 pb-[env(safe-area-inset-bottom)]"
  >
    <div class="flex justify-around items-center h-16 max-w-lg mx-auto">
      <router-link
        v-for="item in navItems"
        :key="item.name"
        :to="item.to"
        class="flex flex-col items-center gap-0.5 text-xs transition-colors"
        :class="isActive(item) ? 'text-primary' : 'text-gray-400'"
      >
        <component :is="isActive(item) ? item.iconActive : item.icon" class="w-6 h-6" />
        <span>{{ item.label }}</span>
      </router-link>
    </div>
  </nav>
</template>
