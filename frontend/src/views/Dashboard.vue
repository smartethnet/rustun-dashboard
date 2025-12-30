<template>
  <layout>
    <div class="dashboard">
      <h1 class="page-title">{{ t('dashboard.title') }}</h1>
      
      <!-- Statistics Cards -->
      <el-row :gutter="20" class="stats-row">
        <el-col :xs="24" :sm="12" :md="8">
          <el-card class="stat-card hover-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #409eff20;">
                <el-icon size="32" color="#409eff"><Collection /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ clusters.length }}</div>
                <div class="stat-label">{{ t('dashboard.totalClusters') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8">
          <el-card class="stat-card hover-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #67c23a20;">
                <el-icon size="32" color="#67c23a"><Monitor /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ clients.length }}</div>
                <div class="stat-label">{{ t('dashboard.totalClients') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8">
          <el-card class="stat-card hover-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #e6a23c20;">
                <el-icon size="32" color="#e6a23c"><Connection /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ activeConnections }}</div>
                <div class="stat-label">{{ t('dashboard.activeConnections') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <!-- Clusters Overview -->
      <el-card class="section-card">
        <template #header>
          <div class="card-header">
            <span>{{ t('dashboard.clustersOverview') }}</span>
            <el-button type="primary" size="small" @click="goToClusters">
              {{ t('common.viewAll') }}
            </el-button>
          </div>
        </template>
        
        <el-table
          :data="clusters"
          style="width: 100%"
          v-loading="store.loading"
        >
          <el-table-column prop="name" :label="t('cluster.name')" min-width="200">
            <template #default="{ row }">
              <el-link
                type="primary"
                @click="goToClusterDetail(row.name)"
              >
                <el-icon><Collection /></el-icon>
                {{ row.name }}
              </el-link>
            </template>
          </el-table-column>
          
          <el-table-column prop="client_count" :label="t('cluster.clients')" width="120">
            <template #default="{ row }">
              <el-tag>{{ row.client_count }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column :label="t('common.status')" width="120">
            <template #default>
              <el-tag type="success" effect="light">{{ t('common.active') }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column :label="t('common.actions')" width="150" align="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="goToClusterDetail(row.name)"
              >
                {{ t('common.details') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <el-empty
          v-if="!store.loading && clusters.length === 0"
          :description="t('dashboard.noClusters')"
        />
      </el-card>
      
      <!-- Recent Clients -->
      <el-card class="section-card">
        <template #header>
          <div class="card-header">
            <span>{{ t('dashboard.recentClients') }}</span>
            <el-button type="primary" size="small" @click="goToClients">
              {{ t('common.viewAll') }}
            </el-button>
          </div>
        </template>
        
        <el-table
          :data="recentClients"
          style="width: 100%"
          v-loading="store.loading"
        >
          <el-table-column prop="cluster" :label="t('client.cluster')" width="150" />
          <el-table-column prop="identity" :label="t('client.identity')" min-width="200" />
          <el-table-column prop="private_ip" :label="t('client.privateIp')" width="150" />
          <el-table-column prop="gateway" :label="t('client.gateway')" width="150" />
          
          <el-table-column :label="t('client.routes')" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.ciders && row.ciders.length > 0">
                {{ row.ciders.length }}
              </el-tag>
              <span v-else class="text-gray-400">-</span>
            </template>
          </el-table-column>
        </el-table>
        
        <el-empty
          v-if="!store.loading && clients.length === 0"
          :description="t('dashboard.noClients')"
        />
      </el-card>
    </div>
  </layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../store'
import Layout from '../components/Layout.vue'

const { t } = useI18n()

const router = useRouter()
const store = useAppStore()

const clusters = computed(() => store.clusters)
const clients = computed(() => store.clients)
const recentClients = computed(() => clients.value.slice(0, 5))
const activeConnections = computed(() => clients.value.length)

const goToClusters = () => {
  router.push('/clusters')
}

const goToClients = () => {
  router.push('/clients')
}

const goToClusterDetail = (name) => {
  router.push(`/clusters/${name}`)
}

onMounted(async () => {
  await store.fetchClusters()
  await store.fetchClients()
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 24px;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.section-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}
</style>

