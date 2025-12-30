import axios from 'axios'
import { ElMessage } from 'element-plus'

// Create axios instance
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add Basic Auth
    const username = localStorage.getItem('username') || 'admin'
    const password = localStorage.getItem('password') || 'admin123'
    const token = btoa(`${username}:${password}`)
    config.headers.Authorization = `Basic ${token}`
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    const message = error.response?.data?.message || error.message || 'Request failed'
    ElMessage.error(message)
    
    if (error.response?.status === 401) {
      // Clear credentials on auth failure
      localStorage.removeItem('username')
      localStorage.removeItem('password')
      window.location.href = '/login'
    }
    
    return Promise.reject(error)
  }
)

// Cluster APIs
export const clusterAPI = {
  // Get all clusters
  getAll() {
    return api.get('/clusters')
  },
  
  // Get cluster by name
  getByName(name) {
    return api.get(`/clusters/${name}`)
  },
  
  // Delete cluster
  delete(name) {
    return api.delete(`/clusters/${name}`)
  },
}

// Client APIs
export const clientAPI = {
  // Get all clients
  getAll(cluster = '') {
    const params = cluster ? { cluster } : {}
    return api.get('/clients', { params })
  },
  
  // Get client by cluster and identity
  getByClusterAndIdentity(cluster, identity) {
    return api.get(`/clients/${cluster}/${identity}`)
  },
  
  // Create client
  create(data) {
    return api.post('/clients', data)
  },
  
  // Update client
  update(cluster, identity, data) {
    return api.put(`/clients/${cluster}/${identity}`, data)
  },
  
  // Delete client
  delete(cluster, identity) {
    return api.delete(`/clients/${cluster}/${identity}`)
  },
}

// Auth API
export const authAPI = {
  // Test credentials
  test(username, password) {
    const token = btoa(`${username}:${password}`)
    return axios.get('/api/clusters', {
      headers: {
        Authorization: `Basic ${token}`,
      },
    })
  },
}

export default api
