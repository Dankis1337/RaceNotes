import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as racesApi from '../api/races'
import { offlineDb } from '../utils/offlineDb'

export const useRacesStore = defineStore('races', () => {
  const races = ref([])
  const currentRace = ref(null)
  const loading = ref(false)

  async function fetchRaces(params = {}) {
    loading.value = true
    try {
      const { data } = await racesApi.getRaces(params)
      races.value = data || []
      // Cache for offline use
      if (races.value.length) {
        offlineDb.cacheRaces(races.value).catch(() => {})
      }
    } catch {
      // Offline fallback
      try {
        const cached = await offlineDb.getCachedRaces()
        races.value = cached || []
      } catch {
        races.value = []
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchRace(id) {
    loading.value = true
    try {
      const { data } = await racesApi.getRace(id)
      currentRace.value = data
    } catch {
      // Try cached
      try {
        const cached = await offlineDb.getCachedRaces()
        currentRace.value = cached.find(r => r.id == id) || null
      } catch {
        currentRace.value = null
      }
    } finally {
      loading.value = false
    }
  }

  async function createRace(raceData) {
    try {
      const { data } = await racesApi.createRace(raceData)
      return data
    } catch (e) {
      if (!navigator.onLine) {
        await offlineDb.addPendingSync({ type: 'create_race', payload: raceData })
        // Add to local cache with temp id
        const tempRace = { ...raceData, id: `temp_${Date.now()}`, _pending: true }
        const cached = await offlineDb.getCachedRaces().catch(() => [])
        cached.unshift(tempRace)
        await offlineDb.cacheRaces(cached).catch(() => {})
        return tempRace
      }
      throw e
    }
  }

  async function updateRace(id, raceData) {
    try {
      const { data } = await racesApi.updateRace(id, raceData)
      return data
    } catch (e) {
      if (!navigator.onLine) {
        await offlineDb.addPendingSync({ type: 'update_race', entityId: id, payload: raceData })
        return { id, ...raceData, _pending: true }
      }
      throw e
    }
  }

  async function removeRace(id) {
    try {
      await racesApi.deleteRace(id)
    } catch (e) {
      if (!navigator.onLine) {
        await offlineDb.addPendingSync({ type: 'delete_race', entityId: id })
      } else {
        throw e
      }
    }
    races.value = races.value.filter(r => r.id !== id)
  }

  return { races, currentRace, loading, fetchRaces, fetchRace, createRace, updateRace, removeRace }
})
