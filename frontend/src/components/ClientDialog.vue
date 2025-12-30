<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? t('client.editClient') : t('client.addClient')"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item :label="t('client.form.cluster')" prop="cluster">
        <el-input
          v-model="form.cluster"
          :placeholder="t('client.form.clusterPlaceholder')"
          :disabled="isEdit"
        />
      </el-form-item>
      
      <el-form-item :label="t('client.form.identity')" prop="identity">
        <el-input
          v-model="form.identity"
          :placeholder="t('client.form.identityPlaceholder')"
          :disabled="isEdit"
        />
      </el-form-item>
      
      <el-form-item :label="t('client.form.privateIp')" prop="private_ip">
        <el-input v-model="form.private_ip" :placeholder="t('client.form.privateIpPlaceholder')" />
      </el-form-item>
      
      <el-form-item :label="t('client.form.mask')" prop="mask">
        <el-input v-model="form.mask" :placeholder="t('client.form.maskPlaceholder')" />
      </el-form-item>
      
      <el-form-item :label="t('client.form.gateway')" prop="gateway">
        <el-input v-model="form.gateway" :placeholder="t('client.form.gatewayPlaceholder')" />
      </el-form-item>
      
      <el-form-item :label="t('client.form.routes')">
        <el-select
          v-model="form.ciders"
          multiple
          filterable
          allow-create
          default-first-option
          :reserve-keyword="false"
          :placeholder="t('client.form.routesPlaceholder')"
          style="width: 100%"
        >
          <el-option
            v-for="item in commonCIDRs"
            :key="item"
            :label="item"
            :value="item"
          />
        </el-select>
        <div class="form-tip">
          {{ t('client.form.routesTip') }}
        </div>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ isEdit ? t('common.update') : t('common.create') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useAppStore } from '../store'

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  client: {
    type: Object,
    default: null,
  },
  clusterName: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue', 'success'])

const store = useAppStore()
const { t } = useI18n()
const formRef = ref(null)
const loading = ref(false)

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const isEdit = computed(() => !!props.client)

const form = ref({
  cluster: '',
  identity: '',
  private_ip: '',
  mask: '255.255.255.0',
  gateway: '',
  ciders: [],
})

const commonCIDRs = [
  '192.168.1.0/24',
  '192.168.100.0/24',
  '10.0.0.0/8',
  '172.16.0.0/12',
]

const rules = computed(() => ({
  cluster: [
    { required: true, message: t('client.form.clusterRequired'), trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: t('client.form.clusterPattern'), trigger: 'blur' },
  ],
  identity: [
    { required: true, message: t('client.form.identityRequired'), trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: t('client.form.identityPattern'), trigger: 'blur' },
  ],
  private_ip: [
    { required: true, message: t('client.form.privateIpRequired'), trigger: 'blur' },
    { pattern: /^(\d{1,3}\.){3}\d{1,3}$/, message: t('client.form.privateIpPattern'), trigger: 'blur' },
  ],
  mask: [
    { required: true, message: t('client.form.maskRequired'), trigger: 'blur' },
    { pattern: /^(\d{1,3}\.){3}\d{1,3}$/, message: t('client.form.maskPattern'), trigger: 'blur' },
  ],
  gateway: [
    { required: true, message: t('client.form.gatewayRequired'), trigger: 'blur' },
    { pattern: /^(\d{1,3}\.){3}\d{1,3}$/, message: t('client.form.gatewayPattern'), trigger: 'blur' },
  ],
}))

watch(() => props.modelValue, (val) => {
  if (val) {
    if (props.client) {
      // Edit mode
      form.value = {
        cluster: props.client.cluster,
        identity: props.client.identity,
        private_ip: props.client.private_ip,
        mask: props.client.mask,
        gateway: props.client.gateway,
        ciders: props.client.ciders || [],
      }
    } else {
      // Add mode
      form.value = {
        cluster: props.clusterName || '',
        identity: '',
        private_ip: '',
        mask: '255.255.255.0',
        gateway: '',
        ciders: [],
      }
    }
    
    // Reset validation
    if (formRef.value) {
      formRef.value.clearValidate()
    }
  }
})

const handleClose = () => {
  visible.value = false
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    
    try {
      if (isEdit.value) {
        // Update client
        await store.updateClient(
          props.client.cluster,
          props.client.identity,
          form.value
        )
        ElMessage.success(t('client.updateSuccess'))
      } else {
        // Create client
        await store.createClient(form.value)
        ElMessage.success(t('client.createSuccess'))
      }
      
      emit('success')
      handleClose()
    } catch (error) {
      console.error('Failed to save client:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
