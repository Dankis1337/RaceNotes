import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as setupsApi from '../api/setups'

export const useSetupsStore = defineStore('setups', () => {
  const setups = ref([])
  const currentSetup = ref(null)
  const loading = ref(false)

  async function fetchSetups() {
    loading.value = true
    try {
      const { data } = await setupsApi.getSetups()
      setups.value = data || []
    } finally {
      loading.value = false
    }
  }

  async function fetchSetup(id) {
    loading.value = true
    try {
      const { data } = await setupsApi.getSetup(id)
      currentSetup.value = data
    } finally {
      loading.value = false
    }
  }

  async function createSetup(setupData) {
    const { data } = await setupsApi.createSetup(setupData)
    return data
  }

  async function updateSetup(id, setupData) {
    const { data } = await setupsApi.updateSetup(id, setupData)
    return data
  }

  async function removeSetup(id) {
    await setupsApi.deleteSetup(id)
    setups.value = setups.value.filter(s => s.id !== id)
  }

  return { setups, currentSetup, loading, fetchSetups, fetchSetup, createSetup, updateSetup, removeSetup }
})
