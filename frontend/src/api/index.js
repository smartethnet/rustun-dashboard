import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Create axios instance with basic auth
const createAuthAxios = () => {
  const username = localStorage.getItem('username')
  const password = localStorage.getItem('password')

  return axios.create({
    baseURL: API_BASE_URL,
    auth: {
      username,
      password
    }
  })
}

// Cluster API
export const clusterAPI = {
  getAll: async () => {
    const api = createAuthAxios()
    return await api.get('/api/clusters')
  },

  getByName: async (name) => {
    const api = createAuthAxios()
    return await api.get(`/api/clusters/${name}`)
  },

  delete: async (name) => {
    const api = createAuthAxios()
    return await api.delete(`/api/clusters/${name}`)
  }
}

// Client API
export const clientAPI = {
  getAll: async (cluster = '') => {
    const api = createAuthAxios()
    const params = cluster ? { cluster } : {}
    return await api.get('/api/clients', { params })
  },

  getByClusterAndIdentity: async (cluster, identity) => {
    const api = createAuthAxios()
    return await api.get(`/api/clients/${cluster}/${identity}`)
  },

  create: async (data) => {
    const api = createAuthAxios()
    return await api.post('/api/clients', data)
  },

  update: async (cluster, identity, data) => {
    const api = createAuthAxios()
    return await api.put(`/api/clients/${cluster}/${identity}`, data)
  },

  delete: async (cluster, identity) => {
    const api = createAuthAxios()
    return await api.delete(`/api/clients/${cluster}/${identity}`)
  }
}

