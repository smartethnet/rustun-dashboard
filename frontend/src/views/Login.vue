<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1 class="title">🌐 {{ t('login.title') }}</h1>
        <p class="subtitle">{{ t('login.subtitle') }}</p>
        
        <!-- Language Switcher -->
        <div class="language-switcher">
          <el-button
            :type="locale === 'zh' ? 'primary' : ''"
            size="small"
            @click="switchLanguage('zh')"
          >
            中文
          </el-button>
          <el-button
            :type="locale === 'en' ? 'primary' : ''"
            size="small"
            @click="switchLanguage('en')"
          >
            English
          </el-button>
        </div>
      </div>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            :placeholder="t('login.usernamePlaceholder')"
            size="large"
            prefix-icon="User"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('login.passwordPlaceholder')"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            class="login-button"
          >
            {{ loading ? t('login.loggingIn') : t('login.loginButton') }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <p class="hint">{{ t('login.defaultHint') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { authAPI } from '../api'

const router = useRouter()
const { t, locale } = useI18n()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: 'admin',
  password: 'admin123',
})

const rules = computed(() => ({
  username: [
    { required: true, message: t('login.usernameRequired'), trigger: 'blur' },
  ],
  password: [
    { required: true, message: t('login.passwordRequired'), trigger: 'blur' },
  ],
}))

const switchLanguage = (lang) => {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    
    try {
      // Test credentials
      await authAPI.test(form.username, form.password)
      
      // Save credentials
      localStorage.setItem('username', form.username)
      localStorage.setItem('password', form.password)
      
      ElMessage.success(t('login.loginSuccess'))
      router.push('/')
    } catch (error) {
      ElMessage.error(t('login.loginFailed'))
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: fadeIn 0.5s ease;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #909399;
}

.language-switcher {
  margin-top: 16px;
  display: flex;
  gap: 8px;
  justify-content: center;
}

.login-form {
  margin-top: 24px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 600;
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

.hint {
  font-size: 13px;
  color: #909399;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

