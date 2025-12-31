<template>
  <layout>
    <div class="topology-page">
      <!-- Operation Tips -->
      <el-alert
        type="info"
        :closable="false"
        class="tips-alert"
      >
        <div class="tips-content">
          <div class="tip-item">
            <el-icon color="#409eff"><Mouse /></el-icon>
            <span>{{ t('topology.tip1') }}</span>
          </div>
          <div class="tip-item">
            <el-icon color="#67c23a"><Connection /></el-icon>
            <span>{{ t('topology.tip2') }}</span>
          </div>
          <div class="tip-item">
            <el-icon color="#e6a23c"><Edit /></el-icon>
            <span>{{ t('topology.tip3') }}</span>
          </div>
        </div>
      </el-alert>
      
      <!-- Topology Visualization -->
      <el-card class="topology-card" v-loading="store.loading">
        <template #header>
          <div class="card-header">
            <span>{{ t('topology.title') }}</span>
            <div style="display: flex; gap: 8px;">
              <el-button @click="openAddDialog" type="primary" size="small">
                <el-icon><Plus /></el-icon>
                {{ t('client.addClient') }}
              </el-button>
              <el-button @click="refreshTopology" size="small">
                <el-icon><Refresh /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
          </div>
        </template>
        
        <div class="network-wrapper">
          <div ref="networkContainer" class="network-container"></div>
          <div
            v-if="!store.loading && clients.length === 0"
            class="empty-state"
          >
            <el-empty :description="t('topology.emptyHint')">
              <el-button type="primary" size="large" @click="openAddDialog">
                <el-icon><Plus /></el-icon>
                {{ t('topology.createFirstClient') }}
              </el-button>
            </el-empty>
          </div>
        </div>
      </el-card>
      
      <!-- Legend -->
      <el-card class="legend-card">
        <template #header>
          {{ t('topology.legend') }}
        </template>
        
        <div class="legend-items">
          <div class="legend-item">
            <div class="legend-node cluster-node"></div>
            <span>{{ t('topology.cluster') }}</span>
          </div>
          
          <div class="legend-item">
            <div class="legend-node client-node"></div>
            <span>{{ t('topology.client') }}</span>
          </div>
          
          <div class="legend-item">
            <div class="legend-edge"></div>
            <span>{{ t('client.routes') }}</span>
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- Client Management Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'add' ? t('client.addClient') : t('client.editClient')"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="clientForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item :label="t('client.cluster')" prop="cluster">
          <el-input
            v-model="clientForm.cluster"
            :placeholder="t('client.clusterPlaceholder')"
            :disabled="dialogMode === 'edit'"
          >
            <template #append>
              <el-dropdown
                v-if="dialogMode === 'add' && clusters.length > 0"
                trigger="click"
                @command="selectCluster"
              >
                <el-button>
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item
                      v-for="cluster in clusters"
                      :key="cluster.name"
                      :command="cluster.name"
                    >
                      {{ cluster.name }} ({{ cluster.client_count }} {{ t('cluster.clients') }})
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item :label="t('client.name')">
          <el-input
            v-model="clientForm.name"
            :placeholder="t('client.namePlaceholder')"
          />
          <template #extra>
            <span style="font-size: 12px; color: #909399;">{{ t('client.nameHint') }}</span>
          </template>
        </el-form-item>
        
        <el-form-item :label="t('client.identity')" v-if="dialogMode === 'edit'">
          <el-input
            v-model="clientForm.identity"
            disabled
          >
            <template #prepend>
              <span style="font-size: 12px;">UUID</span>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item :label="t('client.privateIP')" v-if="dialogMode === 'edit'">
          <el-input
            v-model="clientForm.private_ip"
            disabled
          />
        </el-form-item>
        
        <el-form-item :label="t('client.mask')" v-if="dialogMode === 'edit'">
          <el-input
            v-model="clientForm.mask"
            disabled
          />
        </el-form-item>
        
        <el-form-item :label="t('client.gateway')" v-if="dialogMode === 'edit'">
          <el-input
            v-model="clientForm.gateway"
            disabled
          />
        </el-form-item>
        
        <el-form-item :label="t('client.routes')">
          <div class="routes-container">
            <div
              v-for="(route, index) in clientForm.ciders"
              :key="index"
              class="route-item"
            >
              <el-input
                v-model="clientForm.ciders[index]"
                :placeholder="t('client.routePlaceholder')"
              />
              <el-button
                type="danger"
                :icon="Delete"
                circle
                @click="removeRoute(index)"
              />
            </div>
            <el-button @click="addRoute" style="width: 100%; margin-top: 8px">
              <el-icon><Plus /></el-icon>
              {{ t('client.addRoute') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <el-button
              v-if="dialogMode === 'edit'"
              type="danger"
              @click="handleDeleteFromDialog"
            >
              <el-icon><Delete /></el-icon>
              {{ t('common.delete') }}
            </el-button>
          </div>
          <div>
            <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitting">
              {{ t('common.save') }}
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
    
    <!-- Cluster Actions Menu -->
    <el-dialog
      v-model="clusterMenuVisible"
      :title="selectedCluster"
      width="350px"
    >
      <div style="text-align: center; padding: 20px;">
        <el-button type="primary" size="large" @click="handleAddClientToCluster" style="width: 100%;">
          <el-icon><Plus /></el-icon>
          {{ t('client.addClient') }}
        </el-button>
      </div>
    </el-dialog>
  </layout>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, ArrowDown, Mouse, Connection, View } from '@element-plus/icons-vue'
import { Network } from 'vis-network/standalone'
import { useAppStore } from '../store'
import { createClient, updateClient, deleteClient } from '../api/clients'
import Layout from '../components/Layout.vue'

const { t } = useI18n()
const store = useAppStore()

const networkContainer = ref(null)
const network = ref(null)

const clusters = computed(() => store.clusters)
const clients = computed(() => store.clients)

// Dialog state
const dialogVisible = ref(false)
const clusterMenuVisible = ref(false)
const dialogMode = ref('add') // 'add' or 'edit'
const submitting = ref(false)
const formRef = ref(null)
const selectedClient = ref(null)
const selectedCluster = ref(null)

// Form data
const clientForm = ref({
  cluster: '',
  name: '', // Optional friendly name (only for creation)
  identity: '', // Only used for editing
  private_ip: '',
  mask: '255.255.255.0',
  gateway: '',
  ciders: [],
})

// Form validation rules
const formRules = {
  cluster: [
    { required: true, message: t('client.clusterRequired'), trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: t('client.clusterPattern'), trigger: 'blur' },
  ],
}

// Generate network data from clients
const generateNetworkData = () => {
  const nodes = []
  const edges = []
  const clusterGroups = {}
  
  // Group clients by cluster
  clients.value.forEach(client => {
    if (!clusterGroups[client.cluster]) {
      clusterGroups[client.cluster] = []
    }
    clusterGroups[client.cluster].push(client)
  })
  
  // Create cluster nodes (center nodes)
  Object.keys(clusterGroups).forEach((clusterName) => {
    const clusterId = `cluster-${clusterName}`
    const clientCount = clusterGroups[clusterName].length
    
    nodes.push({
      id: clusterId,
      label: `${clusterName}\n${clientCount} clients`,
      shape: 'box',
      color: {
        background: '#667eea',
        border: '#5a67d8',
        highlight: {
          background: '#5a67d8',
          border: '#4c51bf',
        },
      },
      font: {
        color: '#ffffff',
        size: 16,
        bold: true,
      },
      borderWidth: 3,
      borderWidthSelected: 4,
      margin: 15,
      level: 0,
      clusterName: clusterName, // Store cluster name for click event
    })
    
    // Create client nodes and connect directly to cluster
    clusterGroups[clusterName].forEach((client) => {
      const clientId = `client-${client.cluster}-${client.identity}`
      // Display name and IP (or just IP if no name)
      const displayName = client.name 
        ? `${client.name}\n${client.private_ip}` 
        : client.private_ip
      
      nodes.push({
        id: clientId,
        label: displayName,
        shape: 'box',
        color: {
          background: '#48bb78',
          border: '#38a169',
          highlight: {
            background: '#38a169',
            border: '#2f855a',
          },
        },
        font: {
          size: 13,
          color: '#ffffff',
          face: 'Monaco, Consolas, monospace',
          bold: true,
        },
        borderWidth: 2,
        borderWidthSelected: 3,
        margin: 12,
        level: 1,
        title: `${client.name ? `Name: ${client.name}\n` : ''}IP: ${client.private_ip}\nIdentity: ${client.identity}\nCluster: ${client.cluster}\nMask: ${client.mask}\nGateway: ${client.gateway}\nRoutes: ${client.ciders?.length || 0}`,
        // Store client data for click event
        cluster: client.cluster,
        identity: client.identity,
      })
      
      // Connect client directly to cluster
      edges.push({
        from: clusterId,
        to: clientId,
        color: { 
          color: '#48bb78',
          highlight: '#38a169',
          hover: '#38a169',
        },
        width: 2.5,
        arrows: { to: false },
        smooth: {
          type: 'cubicBezier',
          roundness: 0.5,
        },
      })
      
      // Create route nodes if client has routes
      if (client.ciders && client.ciders.length > 0) {
        client.ciders.forEach((cidr, cidx) => {
          const routeId = `route-${clientId}-${cidx}`
          nodes.push({
            id: routeId,
            label: cidr,
            shape: 'box',
            color: {
              background: '#edf2f7',
              border: '#cbd5e0',
              highlight: {
                background: '#e2e8f0',
                border: '#a0aec0',
              },
            },
            font: {
              size: 10,
              color: '#4a5568',
              face: 'Monaco, Consolas, monospace',
            },
            borderWidth: 1,
            margin: 8,
            level: 2,
          })
          
          // Connect client to route
          edges.push({
            from: clientId,
            to: routeId,
            color: { 
              color: '#a0aec0',
              highlight: '#718096',
            },
            width: 1.5,
            dashes: [5, 5],
            arrows: { 
              to: {
                enabled: true,
                scaleFactor: 0.5,
              },
            },
            smooth: {
              type: 'curvedCW',
              roundness: 0.2,
            },
          })
        })
      }
    })
  })
  
  return { nodes, edges }
}

// Initialize network
const initNetwork = () => {
  if (!networkContainer.value) return
  
  const data = generateNetworkData()
  
  const options = {
    nodes: {
      shadow: {
        enabled: true,
        color: 'rgba(0,0,0,0.15)',
        size: 8,
        x: 2,
        y: 2,
      },
      font: {
        size: 14,
        face: '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
      },
      shapeProperties: {
        borderRadius: 6,
      },
    },
    edges: {
      shadow: {
        enabled: true,
        color: 'rgba(0,0,0,0.1)',
        size: 5,
        x: 1,
        y: 1,
      },
      smooth: {
        enabled: true,
        type: 'continuous',
      },
    },
    layout: {
      hierarchical: {
        enabled: true,
        direction: 'UD',
        sortMethod: 'directed',
        levelSeparation: 140,
        nodeSpacing: 180,
        treeSpacing: 220,
      },
    },
    physics: {
      enabled: false,
    },
    interaction: {
      hover: true,
      hoverConnectedEdges: true,
      tooltipDelay: 200,
      zoomView: true,
      dragView: true,
      selectable: true,
      selectConnectedEdges: false,
      navigationButtons: false,
      keyboard: {
        enabled: false,
      },
    },
    tooltips: {
      enabled: true,
    },
  }
  
  network.value = new Network(networkContainer.value, data, options)
  
  // Manual double-click detection
  let lastClickTime = 0
  let lastClickedNodeId = null
  const DOUBLE_CLICK_THRESHOLD = 300 // milliseconds
  
  network.value.on('click', (params) => {
    const currentTime = new Date().getTime()
    const timeDiff = currentTime - lastClickTime
    const currentNodeId = params.nodes.length > 0 ? params.nodes[0] : null
    
    // Check if this is a double click (same target within threshold)
    if (timeDiff < DOUBLE_CLICK_THRESHOLD && timeDiff > 0 && lastClickedNodeId === currentNodeId) {
      // This is a double click
      if (params.event && params.event.preventDefault) {
        params.event.preventDefault()
      }
      
      if (currentNodeId) {
        // Double-clicked on a node
        const nodeId = currentNodeId
        
        // Handle cluster nodes
        if (nodeId.startsWith('cluster-')) {
          const nodeData = network.value.body.data.nodes.get(nodeId)
          if (nodeData && nodeData.clusterName) {
            selectedCluster.value = nodeData.clusterName
            clusterMenuVisible.value = true
          }
        }
        // Handle client nodes
        else if (nodeId.startsWith('client-')) {
          const nodeData = network.value.body.data.nodes.get(nodeId)
          
          if (nodeData && nodeData.cluster && nodeData.identity) {
            const client = clients.value.find(
              c => c.cluster === nodeData.cluster && c.identity === nodeData.identity
            )
            
            if (client) {
              handleEditClient(client)
            }
          }
        }
      } else {
        // Double-clicked on canvas (blank area) - open add client dialog
        openAddDialog()
      }
      
      // Reset to prevent triple-click
      lastClickTime = 0
      lastClickedNodeId = null
    } else {
      // This is a single click, just record it
      lastClickTime = currentTime
      lastClickedNodeId = currentNodeId
    }
  })
}

// Dialog handlers
const openAddDialog = () => {
  dialogMode.value = 'add'
  resetForm()
  dialogVisible.value = true
}

const handleAddClientToCluster = () => {
  clusterMenuVisible.value = false
  dialogMode.value = 'add'
  resetForm()
  clientForm.value.cluster = selectedCluster.value
  dialogVisible.value = true
}

const handleEditClient = (client) => {
  selectedClient.value = client
  dialogMode.value = 'edit'
  clientForm.value = {
    cluster: client.cluster,
    identity: client.identity,
    name: client.name || '',
    private_ip: client.private_ip,
    mask: client.mask,
    gateway: client.gateway,
    ciders: client.ciders ? [...client.ciders] : [],
  }
  dialogVisible.value = true
}

const handleDeleteClient = async (client) => {
  try {
    await ElMessageBox.confirm(
      t('client.deleteConfirm'),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      }
    )
    
    await store.deleteClient(client.cluster, client.identity)
    ElMessage.success(t('client.deleteSuccess'))
    await refreshTopology()
  } catch (error) {
    // User cancelled the confirmation
    if (error === 'cancel') {
      throw error // Re-throw to let caller know user cancelled
    }
    // Actual deletion error
    console.error('Failed to delete client:', error)
    ElMessage.error(t('client.deleteFailed'))
    throw error
  }
}

const handleDeleteFromDialog = async () => {
  try {
    await handleDeleteClient(selectedClient.value)
    dialogVisible.value = false
  } catch (error) {
    // User cancelled confirmation dialog, do nothing
    if (error === 'cancel') {
      console.log('Delete cancelled by user')
    }
    // Error already handled in handleDeleteClient
  }
}

const resetForm = () => {
  clientForm.value = {
    cluster: '',
    name: '',
    identity: '',
    private_ip: '',
    mask: '255.255.255.0',
    gateway: '',
    ciders: [],
  }
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const addRoute = () => {
  clientForm.value.ciders.push('')
}

const removeRoute = (index) => {
  clientForm.value.ciders.splice(index, 1)
}

const selectCluster = (cluster) => {
  clientForm.value.cluster = cluster
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      submitting.value = true
      
      if (dialogMode.value === 'add') {
        // For creation, backend will auto-generate IP, mask, gateway, and identity
        const data = {
          cluster: clientForm.value.cluster,
          name: clientForm.value.name || undefined,
          ciders: clientForm.value.ciders.filter(r => r.trim() !== ''),
        }
        await createClient(data)
        ElMessage.success(t('client.createSuccess'))
      } else {
        // For update, include identity and name
        const data = {
          cluster: clientForm.value.cluster,
          identity: clientForm.value.identity,
          name: clientForm.value.name || '',
          private_ip: clientForm.value.private_ip,
          mask: clientForm.value.mask,
          gateway: clientForm.value.gateway,
          ciders: clientForm.value.ciders.filter(r => r.trim() !== ''),
        }
        await updateClient(data.cluster, data.identity, data)
        ElMessage.success(t('client.updateSuccess'))
      }
      
      dialogVisible.value = false
      resetForm()
      await refreshTopology()
    } catch (error) {
      ElMessage.error(
        dialogMode.value === 'add' 
          ? t('client.createFailed') 
          : t('client.updateFailed')
      )
    } finally {
      submitting.value = false
    }
  })
}

const refreshTopology = async () => {
  await store.fetchClusters()
  await store.fetchClients()
  await nextTick()
  if (network.value) {
    network.value.destroy()
  }
  initNetwork()
}

onMounted(async () => {
  await store.fetchClusters()
  await store.fetchClients()
  await nextTick()
  initNetwork()
})

onBeforeUnmount(() => {
  if (network.value) {
    network.value.destroy()
  }
})
</script>

<style scoped>
.topology-page {
  max-width: 1600px;
  margin: 0 auto;
}

.stats-row {
  margin-bottom: 16px;
}

.tips-alert {
  margin-bottom: 16px;
}

.tips-content {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-top: 8px;
}

.tip-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.stat-card {
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

.topology-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.network-wrapper {
  width: 100%;
  height: 500px;
  position: relative;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

.network-container {
  width: 100%;
  height: 100%;
  border-radius: 8px;
}

.empty-state {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  text-align: center;
  z-index: 10;
  pointer-events: all;
}

.empty-state :deep(.el-empty) {
  background: rgba(255, 255, 255, 0.95);
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-width: 400px;
  margin: 0 auto;
}

.empty-state :deep(.el-button) {
  margin-top: 16px;
}

.legend-card {
  margin-bottom: 16px;
}

.legend-items {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-node {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 2px solid;
}

.cluster-node {
  background: #667eea;
  border-color: #5a67d8;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.client-node {
  background: #48bb78;
  border-color: #38a169;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.legend-edge {
  width: 40px;
  height: 2px;
  background: #48bb78;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.routes-container {
  width: 100%;
}

.route-item {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.client-info {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.client-info p {
  margin: 8px 0;
  color: #606266;
}

@media (max-width: 768px) {
  .tips-content {
    flex-direction: column;
    gap: 12px;
  }
  
  .network-wrapper {
    height: 400px;
  }
}
</style>

