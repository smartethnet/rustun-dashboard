<template>
  <div class="ai-chat-container">
    <el-card class="chat-card" shadow="never">
      <!-- Messages Area -->
      <div class="messages-container" ref="messagesContainer">
        <!-- Welcome Message -->
        <div v-if="messages.length === 0" class="welcome-message">
          <el-icon class="welcome-icon"><ChatDotRound /></el-icon>
          <h2>{{ t('ai.welcome') }}</h2>
          <p>{{ t('ai.welcomeDesc') }}</p>
          
          <!-- Quick Actions -->
          <div class="quick-actions">
            <h4>{{ t('ai.quickActions') }}</h4>
            <el-space wrap>
              <el-button 
                v-for="action in quickActions" 
                :key="action.text"
                size="small"
                @click="sendQuickAction(action.text)"
              >
                {{ action.text }}
              </el-button>
            </el-space>
          </div>
        </div>

        <!-- Chat Messages -->
        <div v-else class="messages-list">
          <div 
            v-for="(message, index) in messages" 
            :key="index"
            :class="['message-item', message.role]"
          >
            <!-- User Message -->
            <div v-if="message.role === 'user'" class="message-bubble user-message">
              <div class="message-content">{{ message.content }}</div>
            </div>

            <!-- Assistant Message -->
            <div v-else class="message-bubble assistant-message">
              <!-- Tool Calls Badge -->
              <div v-if="message.toolCalls && message.toolCalls.length > 0" class="tool-calls">
                <div class="tool-calls-header">
                  <el-icon><Tools /></el-icon>
                  <span>{{ t('ai.toolsUsed') }} ({{ message.toolCalls.length }})</span>
                </div>
                <div class="tool-calls-list">
                  <el-collapse accordion>
                    <el-collapse-item 
                      v-for="(tool, idx) in message.toolCalls" 
                      :key="idx"
                      :name="idx"
                    >
                      <template #title>
                        <div class="tool-title">
                          <el-tag 
                            size="small"
                            type="success"
                            effect="plain"
                          >
                            {{ getToolName(tool.tool) }}
                          </el-tag>
                          <el-icon class="expand-icon"><ArrowRight /></el-icon>
                        </div>
                      </template>
                      <div class="tool-details">
                        <div class="tool-section" v-if="tool.arguments">
                          <div class="tool-section-title">{{ t('ai.parameters') }}:</div>
                          <pre class="tool-code">{{ formatJSON(tool.arguments) }}</pre>
                        </div>
                        <div class="tool-section" v-if="tool.result">
                          <div class="tool-section-title">{{ t('ai.result') }}:</div>
                          <pre class="tool-code">{{ formatJSON(tool.result) }}</pre>
                        </div>
                      </div>
                    </el-collapse-item>
                  </el-collapse>
                </div>
              </div>
              
              <!-- Message Content -->
              <div class="message-content" v-html="formatMessage(message.content)"></div>
            </div>
          </div>

          <!-- Loading -->
          <div v-if="loading" class="message-item assistant">
            <div class="message-bubble assistant-message loading">
              <el-icon class="is-loading"><Loading /></el-icon>
              {{ t('ai.thinking') }}
            </div>
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="input-area">
        <el-input
          v-model="inputMessage"
          :placeholder="t('ai.inputPlaceholder')"
          type="textarea"
          :rows="2"
          :maxlength="1000"
          show-word-limit
          @keydown.enter.ctrl="sendMessage"
          :disabled="loading"
        />
        <el-button
          type="primary"
          :icon="Promotion"
          @click="sendMessage"
          :loading="loading"
          :disabled="!inputMessage.trim()"
        >
          {{ t('ai.send') }}
        </el-button>
      </div>

      <!-- Tips -->
      <div class="tips" v-if="messages.length === 0">
        <el-alert 
          :title="t('ai.tips')" 
          type="info" 
          :closable="false"
          show-icon
        >
          <ul>
            <li>{{ t('ai.tip1') }}</li>
            <li>{{ t('ai.tip2') }}</li>
            <li>{{ t('ai.tip3') }}</li>
          </ul>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatDotRound, Delete, Tools, Loading, Promotion, ArrowRight } from '@element-plus/icons-vue'
import { agentApi } from '@/api/agent'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const { t } = useI18n()

// Configure marked
marked.setOptions({
  breaks: true,
  gfm: true,
})

// State
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const messagesContainer = ref(null)

// Quick Actions
const quickActions = computed(() => [
  { text: t('ai.action1') },
  { text: t('ai.action2') },
  { text: t('ai.action3') },
  { text: t('ai.action4') },
])

// Tool name mapping
const toolNameMap = {
  'list_clusters': t('ai.tools.listClusters'),
  'list_clients': t('ai.tools.listClients'),
  'get_client': t('ai.tools.getClient'),
  'create_client': t('ai.tools.createClient'),
  'update_client': t('ai.tools.updateClient'),
  'delete_client': t('ai.tools.deleteClient'),
}

const getToolName = (tool) => {
  return toolNameMap[tool] || tool
}

// Format message with markdown
const formatMessage = (content) => {
  const rawHtml = marked.parse(content)
  return DOMPurify.sanitize(rawHtml)
}

// Format JSON for display
const formatJSON = (jsonString) => {
  try {
    const obj = typeof jsonString === 'string' ? JSON.parse(jsonString) : jsonString
    return JSON.stringify(obj, null, 2)
  } catch {
    return jsonString
  }
}

// Send message
const sendMessage = async () => {
  if (!inputMessage.value.trim() || loading.value) return

  const userMessage = inputMessage.value.trim()
  inputMessage.value = ''

  // Add user message
  messages.value.push({
    role: 'user',
    content: userMessage
  })

  // Scroll to bottom
  await nextTick()
  scrollToBottom()

  loading.value = true

  // Add placeholder assistant message for streaming
  const assistantMessageIndex = messages.value.length
  messages.value.push({
    role: 'assistant',
    content: '',
    toolCalls: []
  })

  try {
    // Build history for API
    const history = messages.value.slice(0, -1).map(msg => ({
      role: msg.role,
      content: msg.content
    }))

    // Call streaming API
    await agentApi.chatStream(userMessage, history, (chunk) => {
      if (chunk.type === 'content') {
        // Update message content
        messages.value[assistantMessageIndex].content = chunk.fullMessage
        
        // Scroll to bottom
        nextTick(() => scrollToBottom())
      } else if (chunk.type === 'tool_call') {
        // Add tool call
        if (!messages.value[assistantMessageIndex].toolCalls) {
          messages.value[assistantMessageIndex].toolCalls = []
        }
        messages.value[assistantMessageIndex].toolCalls.push(chunk.toolCall)
        
        // Scroll to bottom
        nextTick(() => scrollToBottom())
      } else if (chunk.type === 'done') {
        // Ensure final state
        messages.value[assistantMessageIndex].content = chunk.fullMessage
        if (chunk.toolCalls && chunk.toolCalls.length > 0) {
          messages.value[assistantMessageIndex].toolCalls = chunk.toolCalls
        }
      }
    })

  } catch (error) {
    console.error('Chat error:', error)
    ElMessage.error(t('ai.error') + ': ' + (error.message || 'Unknown error'))
    
    // Remove the placeholder message if error
    messages.value.splice(assistantMessageIndex, 1)
    // Re-add user message input
    inputMessage.value = userMessage
  } finally {
    loading.value = false
    await nextTick()
    scrollToBottom()
  }
}

// Send quick action
const sendQuickAction = (text) => {
  inputMessage.value = text
  sendMessage()
}

// Clear history
const clearHistory = async () => {
  try {
    await ElMessageBox.confirm(
      t('ai.clearConfirmMessage'),
      t('ai.clearConfirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    
    messages.value = []
    ElMessage.success(t('ai.cleared'))
  } catch {
    // User cancelled
  }
}

// Scroll to bottom
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.ai-chat-container {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
}

.chat-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
  overflow: hidden;
}

/* Header */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.robot-icon {
  font-size: 32px;
  color: #409eff;
}

.chat-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: #909399;
}

/* Messages Container */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #f5f7fa;
}

/* Welcome Message */
.welcome-message {
  text-align: center;
  padding: 60px 20px;
}

.welcome-icon {
  font-size: 64px;
  color: #409eff;
  margin-bottom: 20px;
}

.welcome-message h2 {
  margin: 0 0 12px;
  font-size: 24px;
  color: #303133;
}

.welcome-message p {
  margin: 0 0 40px;
  font-size: 14px;
  color: #606266;
}

.quick-actions {
  max-width: 600px;
  margin: 0 auto;
}

.quick-actions h4 {
  margin: 0 0 16px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

/* Messages List */
.messages-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-item.user {
  justify-content: flex-end;
}

.message-item.assistant {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  word-wrap: break-word;
}

.user-message {
  background: #409eff;
  color: white;
  border-bottom-right-radius: 4px;
}

.assistant-message {
  background: white;
  color: #303133;
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.assistant-message.loading {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #909399;
}

.message-content {
  line-height: 1.6;
  font-size: 14px;
}

.message-content :deep(strong) {
  font-weight: 600;
}

.message-content :deep(p) {
  margin: 8px 0;
}

.message-content :deep(p:first-child) {
  margin-top: 0;
}

.message-content :deep(p:last-child) {
  margin-bottom: 0;
}

.message-content :deep(ul),
.message-content :deep(ol) {
  margin: 8px 0;
  padding-left: 24px;
}

.message-content :deep(li) {
  margin: 4px 0;
}

.message-content :deep(code) {
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 13px;
}

.message-content :deep(pre) {
  background: rgba(0, 0, 0, 0.05);
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 8px 0;
}

.message-content :deep(pre code) {
  background: none;
  padding: 0;
}

.message-content :deep(blockquote) {
  border-left: 3px solid #409eff;
  padding-left: 12px;
  margin: 8px 0;
  color: #606266;
}

.message-content :deep(a) {
  color: #409eff;
  text-decoration: none;
}

.message-content :deep(a:hover) {
  text-decoration: underline;
}

.message-content :deep(h1),
.message-content :deep(h2),
.message-content :deep(h3),
.message-content :deep(h4),
.message-content :deep(h5),
.message-content :deep(h6) {
  margin: 16px 0 8px;
  font-weight: 600;
}

.message-content :deep(h1:first-child),
.message-content :deep(h2:first-child),
.message-content :deep(h3:first-child),
.message-content :deep(h4:first-child),
.message-content :deep(h5:first-child),
.message-content :deep(h6:first-child) {
  margin-top: 0;
}

.message-content :deep(hr) {
  border: none;
  border-top: 1px solid #e4e7ed;
  margin: 16px 0;
}

.message-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 8px 0;
}

.message-content :deep(table th),
.message-content :deep(table td) {
  border: 1px solid #e4e7ed;
  padding: 8px 12px;
  text-align: left;
}

.message-content :deep(table th) {
  background: #f5f7fa;
  font-weight: 600;
}

.message-content :deep(br) {
  margin: 4px 0;
}

.tool-calls {
  margin-bottom: 12px;
  padding: 12px;
  background: #f0f9ff;
  border: 1px solid #b3d8ff;
  border-radius: 8px;
}

.tool-calls-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 13px;
  font-weight: 600;
  color: #409eff;
}

.tool-calls-list {
  margin-top: 8px;
}

.tool-calls-list :deep(.el-collapse) {
  border: none;
}

.tool-calls-list :deep(.el-collapse-item__header) {
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 8px 12px;
  margin-bottom: 8px;
  height: auto;
  line-height: 1.5;
}

.tool-calls-list :deep(.el-collapse-item__wrap) {
  border: none;
  background: white;
  border-radius: 6px;
  margin-bottom: 8px;
}

.tool-calls-list :deep(.el-collapse-item__content) {
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-top: none;
  border-radius: 0 0 6px 6px;
}

.tool-title {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.expand-icon {
  margin-left: auto;
  transition: transform 0.3s;
}

.tool-calls-list :deep(.el-collapse-item.is-active .expand-icon) {
  transform: rotate(90deg);
}

.tool-details {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tool-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.tool-section-title {
  font-size: 12px;
  font-weight: 600;
  color: #606266;
}

.tool-code {
  background: #f5f7fa;
  padding: 8px 12px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #303133;
  overflow-x: auto;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.tool-calls .el-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Input Area */
.input-area {
  display: flex;
  gap: 12px;
  padding: 20px;
  background: white;
  border-top: 1px solid #e4e7ed;
}

.input-area .el-textarea {
  flex: 1;
}

.input-area .el-button {
  align-self: flex-end;
}

/* Tips */
.tips {
  padding: 0 20px 20px;
  background: white;
}

.tips ul {
  margin: 8px 0 0;
  padding-left: 20px;
}

.tips li {
  margin: 4px 0;
  font-size: 13px;
  color: #606266;
}

/* Scrollbar */
.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: #f5f7fa;
}

.messages-container::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

/* Responsive */
@media (max-width: 768px) {
  .message-bubble {
    max-width: 85%;
  }
  
  .input-area {
    flex-direction: column;
  }
  
  .input-area .el-button {
    width: 100%;
  }
}
</style>

