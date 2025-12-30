<template>
  <layout>
    <div class="clients-page">
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ t('client.title') }}</h1>
          <p class="page-description">{{ t('client.description') }}</p>
        </div>
        
        <el-button type="primary" size="large" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          {{ t('client.addClient') }}
        </el-button>
      </div>
      
      <!-- Filters -->
      <el-card class="filter-card">
        <el-form :inline="true">
          <el-form-item :label="t('client.cluster')">
            <el-select
              v-model="selectedCluster"
              :placeholder="t('client.allClusters')"
              clearable
              style="width: 200px"
              @change="handleFilterChange"
            >
              <el-option
                v-for="cluster in clusters"
                :key="cluster.name"
                :label="cluster.name"
                :value="cluster.name"
              />
            </el-select>
          </el-form-item>
          
          <el-form-item :label="t('common.search')">
            <el-input
              v-model="searchText"
              :placeholder="t('client.searchPlaceholder')"
              clearable
              style="width: 250px"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </el-card>
      
      <!-- Clients Table -->
      <el-card>
        <el-table
          :data="filteredClients"
          style="width: 100%"
          v-loading="store.loading"
        >
          <el-table-column prop="cluster" :label="t('client.cluster')" width="150">
            <template #default="{ row }">
              <el-link type="primary" @click="goToCluster(row.cluster)">
                {{ row.cluster }}
              </el-link>
            </template>
          </el-table-column>
          
          <el-table-column prop="identity" :label="t('client.identity')" min-width="200" />
          
          <el-table-column prop="private_ip" :label="t('client.privateIp')" width="150">
            <template #default="{ row }">
              <el-tag>{{ row.private_ip }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="mask" :label="t('client.mask')" width="150" />
          <el-table-column prop="gateway" :label="t('client.gateway')" width="150" />
          
          <el-table-column :label="t('client.routes')" min-width="250">
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
              <el-button type="primary" size="small" @click="handleEdit(row)">
                {{ t('common.edit') }}
              </el-button>
              <el-button type="danger" size="small" @click="handleDelete(row)">
                {{ t('common.delete') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <el-empty
          v-if="!store.loading && filteredClients.length === 0"
          :description="t('dashboard.noClients')"
        >
          <el-button type="primary" @click="showAddDialog">
            {{ t('client.addFirstClient') }}
          </el-button>
        </el-empty>
      </el-card>
    </div>
    
    <!-- Add/Edit Client Dialog -->
    <client-dialog
      v-model="dialogVisible"
      :client="selectedClient"
      @success="handleDialogSuccess"
    />
  </layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useAppStore } from '../store'
import Layout from '../components/Layout.vue'
import ClientDialog from '../components/ClientDialog.vue'

const router = useRouter()
const store = useAppStore()
const { t } = useI18n()

const selectedCluster = ref('')
const searchText = ref('')
const dialogVisible = ref(false)
const selectedClient = ref(null)

const clusters = computed(() => store.clusters)
const clients = computed(() => store.clients)

const filteredClients = computed(() => {
  let result = clients.value
  
  // Filter by cluster
  if (selectedCluster.value) {
    result = result.filter(c => c.cluster === selectedCluster.value)
  }
  
  // Filter by search text
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(c =>
      c.identity.toLowerCase().includes(search) ||
      c.private_ip.toLowerCase().includes(search)
    )
  }
  
  return result
})

const handleFilterChange = async () => {
  await store.fetchClients(selectedCluster.value)
}

const showAddDialog = () => {
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
    
    await store.deleteClient(client.cluster, client.identity)
    ElMessage.success(t('client.deleteSuccess'))
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete client:', error)
    }
  }
}

const goToCluster = (name) => {
  router.push(`/clusters/${name}`)
}

const handleDialogSuccess = () => {
  store.fetchClients(selectedCluster.value)
  store.fetchClusters()
}

onMounted(async () => {
  await store.fetchClusters()
  await store.fetchClients()
})
</script>

<style scoped>
.clients-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  gap: 20px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
