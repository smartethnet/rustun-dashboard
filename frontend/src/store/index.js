import { defineStore } from 'pinia'
import { ref } from 'vue'
import { clusterAPI, clientAPI } from '../api'

export const useAppStore = defineStore('app', () => {
  // State
  const clusters = ref([])
  const clients = ref([])
  const loading = ref(false)
  const selectedCluster = ref(null)
  
  // Actions
  const fetchClusters = async () => {
    loading.value = true
    try {
      const response = await clusterAPI.getAll()
      clusters.value = response.data.data || []
      return clusters.value
    } catch (error) {
      console.error('Failed to fetch clusters:', error)
      clusters.value = []
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const fetchClients = async (cluster = '') => {
    loading.value = true
    try {
      const response = await clientAPI.getAll(cluster)
      clients.value = response.data.data || []
      return clients.value
    } catch (error) {
      console.error('Failed to fetch clients:', error)
      clients.value = []
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const createClient = async (data) => {
    loading.value = true
    try {
      await clientAPI.create(data)
      await fetchClients()
      await fetchClusters()
    } finally {
      loading.value = false
    }
  }
  
  const updateClient = async (cluster, identity, data) => {
    loading.value = true
    try {
      await clientAPI.update(cluster, identity, data)
      await fetchClients()
    } finally {
      loading.value = false
    }
  }
  
  const deleteClient = async (cluster, identity) => {
    loading.value = true
    try {
      await clientAPI.delete(cluster, identity)
      await fetchClients()
      await fetchClusters()
    } finally {
      loading.value = false
    }
  }
  
  const deleteCluster = async (name) => {
    loading.value = true
    try {
      await clusterAPI.delete(name)
      await fetchClusters()
      await fetchClients()
    } finally {
      loading.value = false
    }
  }
  
  return {
    // State
    clusters,
    clients,
    loading,
    selectedCluster,
    
    // Actions
    fetchClusters,
    fetchClients,
    createClient,
    updateClient,
    deleteClient,
    deleteCluster,
  }
})
