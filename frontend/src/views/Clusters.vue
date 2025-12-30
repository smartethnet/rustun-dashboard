<template>
  <layout>
    <div class="clusters-page">
      <div class="page-header">
        <h1 class="page-title">{{ t('cluster.title') }}</h1>
        <p class="page-description">{{ t('cluster.description') }}</p>
      </div>
      
      <!-- Cluster Cards -->
      <el-row :gutter="20" v-loading="store.loading">
        <el-col
          v-for="cluster in clusters"
          :key="cluster.name"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card class="cluster-card hover-card" @click="goToDetail(cluster.name)">
            <div class="cluster-icon">
              <el-icon size="48" color="#409eff"><Collection /></el-icon>
            </div>
            
            <h3 class="cluster-name">{{ cluster.name }}</h3>
            
            <div class="cluster-stats">
              <div class="stat-item">
                <span class="stat-label">{{ t('cluster.clients') }}</span>
                <span class="stat-value">{{ cluster.client_count }}</span>
              </div>
            </div>
            
            <div class="cluster-actions">
              <el-button
                type="primary"
                size="small"
                @click.stop="goToDetail(cluster.name)"
              >
                {{ t('common.viewDetails') }}
              </el-button>
              
              <el-button
                type="danger"
                size="small"
                @click.stop="handleDelete(cluster)"
              >
                {{ t('common.delete') }}
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-empty
        v-if="!store.loading && clusters.length === 0"
        :description="t('dashboard.noClusters')"
        :image-size="200"
      >
        <el-button type="primary" @click="goToClients">
          {{ t('cluster.addFirstClient') }}
        </el-button>
      </el-empty>
    </div>
  </layout>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useAppStore } from '../store'
import Layout from '../components/Layout.vue'

const router = useRouter()
const store = useAppStore()
const { t } = useI18n()

const clusters = computed(() => store.clusters)

const goToDetail = (name) => {
  router.push(`/clusters/${name}`)
}

const goToClients = () => {
  router.push('/clients')
}

const handleDelete = async (cluster) => {
  try {
    await ElMessageBox.confirm(
      t('cluster.deleteConfirm', { name: cluster.name, count: cluster.client_count }),
      t('common.confirm'),
      {
        confirmButtonText: t('common.delete'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        confirmButtonClass: 'el-button--danger',
      }
    )
    
    await store.deleteCluster(cluster.name)
    ElMessage.success(t('cluster.deleteSuccess'))
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete cluster:', error)
    }
  }
}

onMounted(() => {
  store.fetchClusters()
})
</script>

<style scoped>
.clusters-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
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

.cluster-card {
  margin-bottom: 20px;
  cursor: pointer;
  text-align: center;
  transition: all 0.3s;
}

.cluster-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.cluster-icon {
  margin: 20px 0;
}

.cluster-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 16px 0;
}

.cluster-stats {
  margin: 20px 0;
  padding: 16px 0;
  border-top: 1px solid #e4e7ed;
  border-bottom: 1px solid #e4e7ed;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #409eff;
}

.cluster-actions {
  margin-top: 20px;
  display: flex;
  gap: 8px;
  justify-content: center;
}
</style>
