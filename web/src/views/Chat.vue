<script setup lang="ts">
import { ref, onMounted, nextTick, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog, ScrollbarInst } from 'naive-ui'
import api from '../api'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
// import 'highlight.js/styles/atom-one-dark.css' // Removed in favor of CDN link in index.html

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return '<pre class="hljs"><code>' +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }
    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
  }
})

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

interface User {
  username: string
}

interface Conversation {
  conversation_id: string
  title: string
  updated_at: string
}

interface Message {
  role: string
  content: string
  created_at?: string
  streaming?: boolean
}

const user = ref<User>(JSON.parse(localStorage.getItem('user') || '{}'))
const conversations = ref<Conversation[]>([])
const currentConversationId = ref<string | null>(null)
const messages = ref<Message[]>([])
const inputMessage = ref('')
const selectedModel = ref('llama3.2')
const modelOptions = [
  { label: 'Llama 3.2', value: 'llama3.2' },
  { label: 'DeepSeek R1', value: 'deepseek-r1:1.5b' },
  { label: 'Qwen 2.5', value: 'qwen2.5:0.5b' }
]
const streaming = ref(false)
const scrollbarRef = ref<ScrollbarInst | null>(null)

const loadConversations = async () => {
  try {
    const res = await api.get('/v1/conversations?limit=100')
    conversations.value = res.data.conversations || []
  } catch (error) {
    console.error('Failed to load conversations', error)
  }
}

const selectConversation = async (id: string) => {
  if (currentConversationId.value === id) return
  currentConversationId.value = id
  messages.value = []
  try {
    const res = await api.get(`/v1/conversations/${id}/messages?limit=50`)
    messages.value = (res.data.messages || []).reverse()
    scrollToBottom()
  } catch (error) {
    message.error('加载消息失败')
  }
}

const createNewConversation = () => {
  currentConversationId.value = null
  messages.value = []
}

const deleteConversation = async (id: string, e: Event) => {
  e.stopPropagation()
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这个会话吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await api.delete(`/v1/conversations/${id}`)
        message.success('删除成功')
        if (currentConversationId.value === id) {
          createNewConversation()
        }
        loadConversations()
      } catch (error) {
        message.error('删除失败')
      }
    }
  })
}

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}

const scrollToBottom = () => {
  nextTick(() => {
    if (scrollbarRef.value) {
      scrollbarRef.value.scrollTo({ top: 99999, behavior: 'smooth' })
    }
  })
}

const renderMarkdown = (text: string) => {
  return md.render(text || '')
}

const sendMessage = async () => {
  const content = inputMessage.value.trim()
  if (!content || streaming.value) return

  const userMsg = {
    role: 'user',
    content: content,
    created_at: new Date().toISOString()
  }
  messages.value.push(userMsg)
  inputMessage.value = ''
  scrollToBottom()

  const aiMsg = reactive({
    role: 'assistant',
    content: '',
    streaming: true
  })
  messages.value.push(aiMsg)

  let convId = currentConversationId.value
  if (!convId) {
    convId = crypto.randomUUID()
    currentConversationId.value = convId
    conversations.value.unshift({
      conversation_id: convId,
      title: content.slice(0, 20) || 'New Chat',
      updated_at: new Date().toISOString()
    })
  }

  streaming.value = true
  
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/v1/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        conversation_id: convId,
        content: content,
        model: selectedModel.value,
        stream: true
      })
    })

    if (!response.ok) throw new Error(response.statusText)

    const reader = response.body!.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      
      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''
      
      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed || trimmed === 'data: [DONE]') continue

        if (trimmed.startsWith('data:')) {
            const jsonStr = trimmed.replace(/^data:\s*/, '')
            try {
                const data = JSON.parse(jsonStr)
                if (data.content) {
                    aiMsg.content += data.content
                    scrollToBottom()
                }
                if (data.done) {
                    aiMsg.streaming = false
                }
            } catch (e) {
                console.error('Error parsing SSE data', e, trimmed)
            }
        } else {
            console.log('Ignored line:', trimmed)
        }
      }
    }
    
    loadConversations()

  } catch (error) {
    console.error('Chat error', error)
    aiMsg.content += '\n\n[Error: Request failed]'
    message.error('发送失败')
  } finally {
    streaming.value = false
    aiMsg.streaming = false
  }
}

// Helper for reactive inside ref array
import { reactive } from 'vue'

onMounted(() => {
  loadConversations()
})
</script>

<template>
  <n-layout has-sider class="full-height">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="260"
      :native-scrollbar="false"
      style="background-color: #18181c;"
    >
      <div style="padding: 16px; display: flex; flex-direction: column; height: 100%;">
        <n-button type="primary" dashed block @click="createNewConversation">
          + 新建会话
        </n-button>
        <div style="margin-top: 20px; flex: 1; overflow-y: auto;">
          <n-list hoverable clickable>
            <n-list-item v-for="conv in conversations" :key="conv.conversation_id" 
                @click="selectConversation(conv.conversation_id)"
                :style="currentConversationId === conv.conversation_id ? 'background-color: #2d2d30' : ''">
                <n-thing :title="conv.title || '无标题会话'" content-style="margin-top: 4px;">
                    <template #description>
                        <n-text depth="3" style="font-size: 12px">
                            {{ new Date(conv.updated_at).toLocaleString() }}
                        </n-text>
                    </template>
                    <template #header-extra>
                        <n-button size="tiny" text type="error" @click="(e) => deleteConversation(conv.conversation_id, e)">x</n-button>
                    </template>
                </n-thing>
            </n-list-item>
          </n-list>
        </div>
      </div>
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered style="padding: 16px; display: flex; align-items: center; justify-content: space-between;">
        <div style="font-weight: bold; font-size: 16px;">
            {{ currentConversationId ? '对话中' : '新会话' }}
        </div>
        <n-space align="center">
            <n-avatar round size="small">{{ user.username?.[0]?.toUpperCase() }}</n-avatar>
            <n-text>{{ user.username }}</n-text>
            <n-button size="tiny" secondary type="error" @click="handleLogout">退出</n-button>
        </n-space>
      </n-layout-header>

      <n-layout-content content-style="padding: 0; display: flex; flex-direction: column; height: calc(100vh - 60px);">
        <div style="flex: 1; overflow: hidden; position: relative;">
            <n-scrollbar ref="scrollbarRef" style="height: 100%; padding: 20px;">
                <div v-if="messages.length === 0" class="flex-center full-height" style="flex-direction: column; color: #666; height: 100%; display: flex; justify-content: center; align-items: center;">
                    <h1>Krag AI</h1>
                    <p>开始一个新的对话吧</p>
                </div>
                <div v-else>
                    <div v-for="(msg, index) in messages" :key="index" 
                         style="display: flex; margin-bottom: 24px; flex-direction: column;">
                         
                        <div v-if="msg.role === 'user'" style="display: flex; justify-content: flex-end; align-items: flex-start;">
                            <div class="message-bubble message-user">
                                {{ msg.content }}
                            </div>
                            <n-avatar round size="small" style="margin-left: 12px; background-color: #2080f0; flex-shrink: 0;">U</n-avatar>
                        </div>

                        <div v-else style="display: flex; justify-content: flex-start; align-items: flex-start;">
                            <n-avatar round size="small" style="margin-right: 12px; background-color: #18a058; flex-shrink: 0;">AI</n-avatar>
                            <div class="message-bubble message-ai">
                                <div class="markdown-body" v-html="renderMarkdown(msg.content)"></div>
                                <span v-if="msg.streaming" class="blinking-cursor">|</span>
                            </div>
                        </div>
                    </div>
                </div>
            </n-scrollbar>
        </div>

        <div style="padding: 20px; border-top: 1px solid #333; background-color: #18181c;">
            <n-input
                v-model:value="inputMessage"
                type="textarea"
                placeholder="输入消息，按 Enter 发送 (Shift+Enter 换行)..."
                :autosize="{ minRows: 2, maxRows: 6 }"
                @keydown.enter.prevent="(e) => { if(!e.shiftKey) sendMessage() }"
                :disabled="streaming"
            />
            <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 12px;">
                <div style="width: 200px;">
                    <n-select v-model:value="selectedModel" :options="modelOptions" size="small" placeholder="选择模型" />
                </div>
                <n-button type="primary" :loading="streaming" :disabled="!inputMessage.trim() && !streaming" @click="sendMessage">
                    发送
                </n-button>
            </div>
        </div>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style>
.full-height {
  height: 100vh;
}
.message-bubble {
  padding: 12px 16px;
  border-radius: 12px;
  max-width: 85%;
  overflow-wrap: anywhere;
}
.message-user {
  background-color: #2080f0;
  color: white;
  align-self: flex-end;
  border-bottom-right-radius: 2px;
}
.message-ai {
  background-color: #2d2d30;
  color: #e5e7eb;
  align-self: flex-start;
  border-bottom-left-radius: 2px;
  border: 1px solid #3d3d40;
}
.markdown-body {
  background-color: transparent !important;
  font-size: 14px;
}
.markdown-body pre {
  background-color: #161b22 !important;
  border-radius: 6px;
}

.blinking-cursor {
  animation: blink 1s step-end infinite;
  margin-left: 4px;
}
@keyframes blink {
  50% { opacity: 0; }
}
</style>
