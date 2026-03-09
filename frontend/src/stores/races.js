import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as racesApi from '../api/races'

export const useRacesStore = defineStore('races', () => {
  const races = ref([])
  const currentRace = ref(null)
  const loading = ref(false)

  async function fetchRaces(params = {}) {
    loading.value = true
    try {
      const { data } = await racesApi.getRaces(params)
      races.value = data || []
    } finally {
      loading.value = false
    }
  }

  async function fetchRace(id) {
    loading.value = true
    try {
      const { data } = await racesApi.getRace(id)
      currentRace.value = data
    } finally {
      loading.value = false
    }
  }

  async function createRace(raceData) {
    const { data } = await racesApi.createRace(raceData)
    return data
  }

  async function updateRace(id, raceData) {
    const { data } = await racesApi.updateRace(id, raceData)
    return data
  }

  async function removeRace(id) {
    await racesApi.deleteRace(id)
    races.value = races.value.filter(r => r.id !== id)
  }

  return { races, currentRace, loading, fetchRaces, fetchRace, createRace, updateRace, removeRace }
})
