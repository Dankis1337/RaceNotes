import { ref, onMounted, onUnmounted } from 'vue'
import { offlineDb } from '../utils/offlineDb'
import api from '../api/axios'

const pendingCount = ref(0)
const syncing = ref(false)
const isOnline = ref(navigator.onLine)

let syncInterval = null

async function refreshPendingCount() {
  try {
    const items = await offlineDb.getPendingSync()
    pendingCount.value = items.length
  } catch {
    pendingCount.value = 0
  }
}

async function syncPendingItems() {
  if (!navigator.onLine || syncing.value) return
  syncing.value = true

  try {
    const items = await offlineDb.getPendingSync()
    for (const item of items) {
      try {
        switch (item.type) {
          case 'create_race':
            await api.post('/races', item.payload)
            break
          case 'update_race':
            await api.put(`/races/${item.entityId}`, item.payload)
            break
          case 'delete_race':
            await api.delete(`/races/${item.entityId}`)
            break
          case 'create_setup':
            await api.post('/setups', item.payload)
            break
          case 'update_setup':
            await api.put(`/setups/${item.entityId}`, item.payload)
            break
          case 'delete_setup':
            await api.delete(`/setups/${item.entityId}`)
            break
        }
        await offlineDb.removePendingItem(item.localId)
      } catch (e) {
        // If server rejects (4xx), remove from queue; if network error, stop
        if (e.response && e.response.status >= 400 && e.response.status < 500) {
          await offlineDb.removePendingItem(item.localId)
        } else {
          break
        }
      }
    }
  } finally {
    syncing.value = false
    await refreshPendingCount()
  }
}

export function useOfflineSync() {
  function handleOnline() {
    isOnline.value = true
    syncPendingItems()
  }

  function handleOffline() {
    isOnline.value = false
  }

  onMounted(() => {
    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
    refreshPendingCount()

    // Auto-sync every 30s when online
    syncInterval = setInterval(() => {
      if (navigator.onLine && pendingCount.value > 0) {
        syncPendingItems()
      }
    }, 30000)
  })

  onUnmounted(() => {
    window.removeEventListener('online', handleOnline)
    window.removeEventListener('offline', handleOffline)
    if (syncInterval) clearInterval(syncInterval)
  })

  return { pendingCount, syncing, isOnline, syncPendingItems, refreshPendingCount }
}
