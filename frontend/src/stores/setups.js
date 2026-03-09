import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as setupsApi from '../api/setups'
import { offlineDb } from '../utils/offlineDb'

export const useSetupsStore = defineStore('setups', () => {
  const setups = ref([])
  const currentSetup = ref(null)
  const loading = ref(false)

  async function fetchSetups() {
    loading.value = true
    try {
      const { data } = await setupsApi.getSetups()
      setups.value = data || []
      if (setups.value.length) {
        offlineDb.cacheSetups(setups.value).catch(() => {})
      }
    } catch {
      try {
        const cached = await offlineDb.getCachedSetups()
        setups.value = cached || []
      } catch {
        setups.value = []
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchSetup(id) {
    loading.value = true
    try {
      const { data } = await setupsApi.getSetup(id)
      currentSetup.value = data
    } catch {
      try {
        const cached = await offlineDb.getCachedSetups()
        currentSetup.value = cached.find(s => s.id == id) || null
      } catch {
        currentSetup.value = null
      }
    } finally {
      loading.value = false
    }
  }

  async function createSetup(setupData) {
    try {
      const { data } = await setupsApi.createSetup(setupData)
      return data
    } catch (e) {
      if (!navigator.onLine) {
        await offlineDb.addPendingSync({ type: 'create_setup', payload: setupData })
        const tempSetup = { ...setupData, id: `temp_${Date.now()}`, _pending: true }
        const cached = await offlineDb.getCachedSetups().catch(() => [])
        cached.unshift(tempSetup)
        await offlineDb.cacheSetups(cached).catch(() => {})
        return tempSetup
      }
      throw e
    }
  }

  async function updateSetup(id, setupData) {
    try {
      const { data } = await setupsApi.updateSetup(id, setupData)
      return data
    } catch (e) {
      if (!navigator.onLine) {
        await offlineDb.addPendingSync({ type: 'update_setup', entityId: id, payload: setupData })
        return { id, ...setupData, _pending: true }
      }
      throw e
    }
  }

  async function removeSetup(id) {
    try {
      await setupsApi.deleteSetup(id)
    } catch (e) {
      if (!navigator.onLine) {
        await offlineDb.addPendingSync({ type: 'delete_setup', entityId: id })
      } else {
        throw e
      }
    }
    setups.value = setups.value.filter(s => s.id !== id)
  }

  return { setups, currentSetup, loading, fetchSetups, fetchSetup, createSetup, updateSetup, removeSetup }
})
