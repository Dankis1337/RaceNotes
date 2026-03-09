<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Toast from './components/Toast.vue'
import { useToast } from './composables/useToast'
import { useOfflineSync } from './composables/useOfflineSync'
import { useI18n } from './utils/i18n'
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
const { t } = useI18n()

const showNav = computed(() => {
  return !['Login', 'Register'].includes(route.name)
})

const navItems = [
  { to: '/', name: 'Races', labelKey: 'nav_races', icon: FlagIcon, iconActive: FlagSolid },
  { to: '/setups', name: 'Setups', labelKey: 'nav_setups', icon: WrenchScrewdriverIcon, iconActive: WrenchSolid },
  { to: '/calculator', name: 'Calculator', labelKey: 'nav_calculator', icon: CalculatorIcon, iconActive: CalculatorSolid },
  { to: '/profile', name: 'Profile', labelKey: 'nav_profile', icon: UserCircleIcon, iconActive: UserSolid },
]

function isActive(item) {
  if (item.to === '/') return route.path === '/' || route.path.startsWith('/races')
  return route.path.startsWith(item.to)
}
</script>

<template>
  <Toast :message="message" :type="type" :show="show" @close="close" />

  <div v-if="!isOnline || pendingCount > 0" class="fixed top-0 left-0 right-0 z-50 text-center text-xs py-1 font-medium" :class="!isOnline ? 'bg-amber-400 text-amber-900' : 'bg-blue-400 text-white'">
    <template v-if="!isOnline">{{ t('offline_mode') }}</template>
    <template v-else-if="syncing">{{ t('syncing') }}</template>
    <template v-else>{{ t('pending_sync', { n: pendingCount }) }}</template>
  </div>

  <div class="lg:flex lg:min-h-screen">
    <!-- Desktop sidebar -->
    <aside v-if="showNav" class="hidden lg:flex lg:flex-col lg:w-56 lg:fixed lg:inset-y-0 bg-white border-r border-gray-200 z-40">
      <div class="px-5 py-6">
        <h1 class="text-xl font-bold text-primary">RaceNotes</h1>
      </div>
      <nav class="flex-1 px-3 space-y-1">
        <router-link
          v-for="item in navItems"
          :key="item.name"
          :to="item.to"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
          :class="isActive(item) ? 'bg-primary/10 text-primary' : 'text-gray-600 hover:bg-gray-50'"
        >
          <component :is="isActive(item) ? item.iconActive : item.icon" class="w-5 h-5" />
          <span>{{ t(item.labelKey) }}</span>
        </router-link>
      </nav>
    </aside>

    <!-- Main content -->
    <main class="flex-1 lg:ml-56">
      <router-view />
    </main>
  </div>

  <!-- Mobile bottom nav -->
  <nav
    v-if="showNav"
    class="fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 pb-[env(safe-area-inset-bottom)] lg:hidden"
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
        <span>{{ t(item.labelKey) }}</span>
      </router-link>
    </div>
  </nav>
</template>
