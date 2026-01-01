<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <h2>🌐 Rustun</h2>
        <p class="logo-subtitle">{{ t('nav.subtitle') }}</p>
      </div>
      
      <el-menu
        :default-active="currentRoute"
        class="sidebar-menu"
        @select="handleMenuSelect"
      >
        <el-menu-item index="/topology">
          <el-icon><Share /></el-icon>
          <span>{{ t('nav.topology') }}</span>
        </el-menu-item>
        <el-menu-item index="/ai">
          <el-icon><ChatDotRound /></el-icon>
          <span>{{ t('nav.aiChat') }}</span>
        </el-menu-item>
      </el-menu>
      
      <div class="sidebar-footer">
        <div class="project-info">
          <div class="project-name">Rustun Dashboard</div>
          <div class="project-version">v1.0.0</div>
        </div>
        <div class="project-links">
          <a href="https://github.com/smartethnet/rustun" target="_blank" class="project-link" :title="t('footer.github')">
            <el-icon><Link /></el-icon>
            <span>{{ t('footer.github') }}</span>
          </a>
          <a href="https://github.com/smartethnet/rustun#readme" target="_blank" class="project-link" :title="t('footer.docs')">
            <el-icon><Document /></el-icon>
            <span>{{ t('footer.docs') }}</span>
          </a>
          <a href="https://github.com/smartethnet/rustun/issues" target="_blank" class="project-link" :title="t('footer.issues')">
            <el-icon><ChatDotSquare /></el-icon>
            <span>{{ t('footer.issues') }}</span>
          </a>
        </div>
        <div class="project-copyright">
          © 2026 SmartEthNet
        </div>
      </div>
    </el-aside>
    
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <h3 class="header-title">{{ pageTitle }}</h3>
        </div>
        
        <div class="header-right">
          <!-- Language Switcher -->
          <el-dropdown @command="handleLanguageChange">
            <div class="language-dropdown">
              🌐
              <span>{{ currentLanguage }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="zh" :disabled="locale === 'zh'">
                  中文
                </el-dropdown-item>
                <el-dropdown-item command="en" :disabled="locale === 'en'">
                  English
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <!-- User Dropdown -->
          <el-dropdown @command="handleUserCommand">
            <div class="user-dropdown">
              <el-avatar :size="32" icon="User" />
              <span class="username">{{ username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  {{ t('nav.logout') }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox } from 'element-plus'
import { Link, Document, ChatDotRound, ChatDotSquare, Share, User, ArrowDown, SwitchButton } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()

const currentRoute = computed(() => route.path)
const username = computed(() => localStorage.getItem('username') || 'Admin')

const currentLanguage = computed(() => {
  return locale.value === 'zh' ? '中文' : 'English'
})

const pageTitle = computed(() => {
  if (route.path === '/ai' || route.path === '/ai/') {
    return t('nav.aiChat')
  }
  return t('nav.topology')
})

const handleMenuSelect = (index) => {
  router.push(index)
}

const handleLanguageChange = (lang) => {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

const handleUserCommand = async (command) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm(
        t('logout.confirmMessage'),
        t('logout.confirmTitle'),
        {
          confirmButtonText: t('nav.logout'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        }
      )
      
      localStorage.removeItem('username')
      localStorage.removeItem('password')
      router.push('/login')
    } catch {
      // User cancelled
    }
  }
}
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
}

.sidebar {
  background: #304156;
  color: white;
  display: flex;
  flex-direction: column;
}

.logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: white;
}

.logo-subtitle {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 5px;
}

.sidebar-menu {
  border: none;
  background: transparent;
  flex: 1;
}

.sidebar-menu :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.7);
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: #409eff;
  color: white;
}

.sidebar-footer {
  padding: 16px;
  background: rgba(0, 0, 0, 0.1);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 12px;
}

.project-info {
  margin-bottom: 12px;
}

.project-name {
  font-weight: bold;
  color: white;
  margin-bottom: 4px;
}

.project-version {
  color: rgba(255, 255, 255, 0.5);
}

.project-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.project-link {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  transition: color 0.2s;
}

.project-link:hover {
  color: white;
}

.project-copyright {
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 24px;
}

.header-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.language-dropdown {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background 0.3s;
  font-size: 14px;
  color: #606266;
}

.language-dropdown:hover {
  background: #f5f7fa;
  color: #409eff;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 6px;
  transition: background 0.3s;
}

.user-dropdown:hover {
  background: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.main-content {
  background: #f5f7fa;
  padding: 24px;
  min-height: calc(100vh - 60px);
}
</style>
