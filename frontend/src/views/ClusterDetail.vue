<template>
  <layout>
    <div class="cluster-detail" v-loading="loading">
      <!-- Breadcrumb -->
      <el-breadcrumb separator="/" class="breadcrumb">
        <el-breadcrumb-item :to="{ path: '/clusters' }">{{ t('cluster.title') }}</el-breadcrumb-item>
        <el-breadcrumb-item>{{ clusterName }}</el-breadcrumb-item>
      </el-breadcrumb>
      
      <!-- Cluster Header -->
      <el-card class="header-card">
        <div class="cluster-header">
          <div class="cluster-info">
            <h1 class="cluster-title">
              <el-icon size="32" color="#409eff"><Collection /></el-icon>
              {{ clusterName }}
            </h1>
            <p class="cluster-subtitle">
              {{ t('cluster.detail.clientsInCluster', { count: clients.length }) }}
            </p>
          </div>
          
          <div class="cluster-actions">
            <el-button type="primary" @click="showAddClientDialog">
              <el-icon><Plus /></el-icon>
              {{ t('cluster.detail.addClient') }}
            </el-button>
            
            <el-button type="danger" @click="handleDeleteCluster">
              <el-icon><Delete /></el-icon>
              {{ t('cluster.detail.deleteCluster') }}
            </el-button>
          </div>
        </div>
      </el-card>
      
      <!-- Clients Table -->
      <el-card class="clients-card">
        <template #header>
          <span>{{ t('cluster.clients') }}</span>
        </template>
        
        <el-table :data="clients" style="width: 100%">
          <el-table-column prop="identity" :label="t('client.identity')" min-width="200" />
          <el-table-column prop="private_ip" :label="t('client.privateIp')" width="150" />
          <el-table-column prop="mask" :label="t('client.mask')" width="150" />
          <el-table-column prop="gateway" :label="t('client.gateway')" width="150" />
          
          <el-table-column :label="t('client.routes')" min-width="200">
            <template #default="{ row }">
              <el-tag
                v-for="(cidr, index) in row.ciders"
                :key="index"
                size="small"
                style="margin-right: 8px; margin-bottom: 4px;"
              >
                {{ cidr }}
              </el-tag>
              <span v-if="!row.ciders || row.ciders.length === 0" class="text-gray-400">
                {{ t('client.noRoutes') }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column :label="t('common.actions')" width="180" align="right" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="handleEdit(row)"
              >
                {{ t('common.edit') }}
              </el-button>
              
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(row)"
              >
                {{ t('common.delete') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <el-empty
          v-if="!loading && clients.length === 0"
          :description="t('dashboard.noClients')"
        >
          <el-button type="primary" @click="showAddClientDialog">
            {{ t('client.addFirstClient') }}
          </el-button>
        </el-empty>
      </el-card>
    </div>
    
    <!-- Add/Edit Client Dialog -->
    <client-dialog
      v-model="dialogVisible"
      :client="selectedClient"
      :cluster-name="clusterName"
      @success="handleDialogSuccess"
    />
  </layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox, ElMessage } from 'element-plus'
import { clusterAPI, clientAPI } from '../api'
import Layout from '../components/Layout.vue'
import ClientDialog from '../components/ClientDialog.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const clusterName = computed(() => route.params.name)
const clusterData = ref(null)
const clients = ref([])
const dialogVisible = ref(false)
const selectedClient = ref(null)

const fetchClusterData = async () => {
  loading.value = true
  try {
    const response = await clusterAPI.getByName(clusterName.value)
    clusterData.value = response.data.cluster
    clients.value = response.data.clients || []
  } catch (error) {
    console.error('Failed to fetch cluster:', error)
    ElMessage.error(t('cluster.deleteFailed'))
  } finally {
    loading.value = false
  }
}

const showAddClientDialog = () => {
  selectedClient.value = null
  dialogVisible.value = true
}

const handleEdit = (client) => {
  selectedClient.value = { ...client }
  dialogVisible.value = true
}

const handleDelete = async (client) => {
  try {
    await ElMessageBox.confirm(
      t('client.deleteConfirm', { identity: client.identity }),
      t('common.confirm'),
      {
        confirmButtonText: t('common.delete'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      }
    )
    
    loading.value = true
    await clientAPI.delete(client.cluster, client.identity)
    ElMessage.success(t('client.deleteSuccess'))
    await fetchClusterData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete client:', error)
    }
  } finally {
    loading.value = false
  }
}

const handleDeleteCluster = async () => {
  try {
    await ElMessageBox.confirm(
      t('cluster.deleteConfirm', { name: clusterName.value, count: clients.value.length }),
      t('common.confirm'),
      {
        confirmButtonText: t('common.delete'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        confirmButtonClass: 'el-button--danger',
      }
    )
    
    loading.value = true
    await clusterAPI.delete(clusterName.value)
    ElMessage.success(t('cluster.deleteSuccess'))
    router.push('/clusters')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete cluster:', error)
    }
    loading.value = false
  }
}

const handleDialogSuccess = () => {
  fetchClusterData()
}

onMounted(() => {
  fetchClusterData()
})
</script>

<style scoped>
.cluster-detail {
  max-width: 1400px;
  margin: 0 auto;
}

.breadcrumb {
  margin-bottom: 20px;
}

.header-card {
  margin-bottom: 20px;
}

.cluster-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.cluster-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 8px 0;
}

.cluster-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.cluster-actions {
  display: flex;
  gap: 12px;
}

@media (max-width: 768px) {
  .cluster-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .cluster-actions {
    width: 100%;
    flex-direction: column;
  }
  
  .cluster-actions .el-button {
    width: 100%;
  }
}
</style>
