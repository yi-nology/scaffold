<script setup lang="ts">
import { onMounted, ref, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useScaffoldStore } from '../stores/scaffold'

const { t } = useI18n()
const store = useScaffoldStore()

onMounted(() => {
  store.fetchTemplates()
})

const emit = defineEmits<{
  (e: 'select', templateId: string): void
}>()

// ── Access Key Management ─────────────────────────────────────────────
const accessKey = ref<string>('')
const showAccessKeyInput = ref<boolean>(false)
const showAccessKeyPassword = ref<boolean>(false) // 控制密码可见性

// Load access key from localStorage on mount
onMounted(() => {
  const savedKey = localStorage.getItem('scaffold-access-key')
  if (savedKey) {
    accessKey.value = savedKey
  }
})

function saveAccessKey(key: string) {
  accessKey.value = key
  localStorage.setItem('scaffold-access-key', key)
}

function clearAccessKey() {
  accessKey.value = ''
  localStorage.removeItem('scaffold-access-key')
}

function formatDate(date: Date | string): string {
  const d = new Date(date)
  return d.toLocaleTimeString('zh-CN', { 
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// ── Add template modal ──────────────────────────────────────────
const showAddModal  = ref(false)
const formId        = ref('')
const formRepo      = ref('')
const formKey       = ref('')
const formError     = ref('')
const submitting    = ref(false)

function openAddModal() {
  formId.value    = ''
  formRepo.value  = ''
  formKey.value   = ''
  formError.value = ''
  showAddModal.value = true
}

function closeAddModal() { showAddModal.value = false }

// ── Async Task Management ─────────────────────────────────────────────
const activeTasks = ref<Record<string, any>>({})
const pollingInterval = ref<number | null>(null)

async function handleAdd() {
  formError.value = ''
  if (!formId.value.trim())   { formError.value = 'Template ID 不能为空'; return }
  if (!formRepo.value.trim()) { formError.value = 'Repository URL 不能为空'; return }
  if (!accessKey.value.trim()) { formError.value = '请输入访问密钥'; return }

  submitting.value = true
  try {
    // 提交异步任务
    const result = await store.addTemplate(formId.value.trim(), formRepo.value.trim(), accessKey.value.trim())
    
    // 开始轮询任务状态
    startPollingTask(result.taskId, formId.value.trim())
    
    // 关闭模态框但不清空表单，让用户可以看到进度
    showAddModal.value = false
    
  } catch (e: any) {
    formError.value = e.message || '添加失败，请检查仓库地址是否正确'
    submitting.value = false
  }
}

function startPollingTask(taskId: string, templateId: string) {
  // 存储任务信息
  activeTasks.value[taskId] = {
    id: taskId,
    templateId: templateId,
    status: 'pending',
    message: '任务已提交',
    createdAt: new Date()
  }

  // 开始轮询
  const poll = async () => {
    try {
      const task = await store.getTaskStatus(taskId)
      activeTasks.value[taskId].status = task.status
      activeTasks.value[taskId].message = task.message
      activeTasks.value[taskId].updatedAt = task.updated_at

      // 如果任务完成或失败，停止轮询
      if (task.status === 'completed' || task.status === 'failed') {
        stopPolling()
        if (task.status === 'completed') {
          // 重新加载模板列表
          await store.fetchTemplates()
          // 清除任务
          delete activeTasks.value[taskId]
        }
      }
    } catch (error) {
      console.error('Polling error:', error)
    }
  }

  // 立即执行一次
  poll()
  // 设置定时轮询
  pollingInterval.value = window.setInterval(poll, 2000) as unknown as number
}

function stopPolling() {
  if (pollingInterval.value) {
    clearInterval(pollingInterval.value)
    pollingInterval.value = null
  }
}

// 组件卸载时清理轮询
onUnmounted(() => {
  stopPolling()
})

// ── Delete template modal ──────────────────────────────────────
const showDelModal    = ref(false)
const delTargetId     = ref('')
const delTargetName   = ref('')
const delKey          = ref('')
const delError        = ref('')
const deleting        = ref(false)

function openDelModal(id: string, name: string) {
  delTargetId.value   = id
  delTargetName.value = name
  delKey.value        = ''
  delError.value      = ''
  showDelModal.value  = true
}

function closeDelModal() { showDelModal.value = false }

async function handleDelete() {
  delError.value = ''
  if (!accessKey.value.trim()) { delError.value = '请输入访问密钥'; return }

  deleting.value = true
  try {
    await store.deleteTemplate(delTargetId.value, accessKey.value.trim())
    closeDelModal()
  } catch (e: any) {
    delError.value = e.message || '删除失败'
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <div>
    <!-- ── Header ── -->
    <div class="flex items-center justify-between mb-7">
      <div>
        <h2 class="text-xl font-semibold text-white">{{ t('template.choose') }}</h2>
        <p v-if="store.templates.length > 0" class="text-xs font-mono text-gray-600 mt-1">
          <span class="text-cyber-400/50">{{ store.templates.length }}</span>
          <span class="ml-1 opacity-60">templates</span>
        </p>
      </div>
      <div class="flex items-center gap-2">
        <!-- Refresh -->
        <button
          v-if="!store.loading"
          @click="store.fetchTemplates"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[11px] font-mono text-gray-600
                 hover:text-cyber-400 border border-white/[0.06] hover:border-cyber-400/20
                 bg-white/[0.02] hover:bg-cyber-400/5 transition-all duration-200"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          {{ t('common.refresh') }}
        </button>

        <!-- Add template -->
        <button
          @click="openAddModal"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[11px] font-mono
                 text-cyber-400/70 hover:text-cyber-400
                 border border-cyber-400/15 hover:border-cyber-400/35
                 bg-cyber-400/[0.04] hover:bg-cyber-400/8
                 transition-all duration-200"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"/>
          </svg>
          Add
        </button>
      </div>
    </div>

    <!-- ── Loading ── -->
    <div v-if="store.loading" class="glass-card rounded-2xl p-16 text-center">
      <div class="relative w-10 h-10 mx-auto">
        <div class="absolute inset-0 border border-cyber-400/15 rounded-full"></div>
        <div class="absolute inset-0 border border-t-cyber-400 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin"></div>
      </div>
      <p class="text-gray-600 mt-5 text-sm font-light">{{ t('template.loading') }}</p>
    </div>

    <!-- ── Error ── -->
    <div v-else-if="store.error" class="glass-card rounded-2xl p-5 border-red-500/20 bg-red-500/[0.04]">
      <div class="flex items-start gap-3">
        <div class="w-8 h-8 rounded-lg bg-red-500/10 border border-red-500/20 flex items-center justify-center flex-shrink-0 mt-0.5">
          <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <div>
          <p class="text-sm font-medium text-red-400 mb-0.5">{{ t('common.error') }}</p>
          <p class="text-xs text-red-400/70">{{ store.error }}</p>
        </div>
      </div>
    </div>

    <!-- ── Empty ── -->
    <div v-else-if="store.templates.length === 0" class="glass-card glow-card rounded-3xl p-16 text-center">
      <div class="w-16 h-16 mx-auto mb-6 relative">
        <div class="absolute inset-0 bg-midnight-700/80 rounded-2xl border border-white/[0.06]"></div>
        <div class="relative w-full h-full flex items-center justify-center">
          <svg class="w-8 h-8 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
          </svg>
        </div>
      </div>
      <p class="text-gray-400 text-base mb-1.5">{{ t('template.empty') }}</p>
      <p class="text-gray-600 text-xs font-mono">{{ t('template.emptyHint') }}</p>
    </div>

    <!-- ── Active Tasks ── -->
    <div v-if="Object.keys(activeTasks).length > 0" class="mb-6">
      <div class="glass-card rounded-2xl p-5">
        <h3 class="text-sm font-semibold text-white mb-4 flex items-center gap-2">
          <svg class="w-4 h-4 text-cyber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          正在处理的任务
        </h3>
        <div class="space-y-3">
          <div 
            v-for="task in activeTasks" 
            :key="task.id"
            class="flex items-center gap-3 p-3 rounded-xl bg-midnight-800/30 border"
            :class="{
              'border-cyber-400/20': task.status === 'pending' || task.status === 'running',
              'border-green-400/20': task.status === 'completed',
              'border-red-400/20': task.status === 'failed'
            }"
          >
            <div class="flex-shrink-0">
              <div 
                class="w-3 h-3 rounded-full"
                :class="{
                  'bg-cyber-400 animate-pulse': task.status === 'pending' || task.status === 'running',
                  'bg-green-400': task.status === 'completed',
                  'bg-red-400': task.status === 'failed'
                }"
              ></div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-xs font-mono text-gray-300 truncate">
                {{ task.templateId }}
              </div>
              <div class="text-[10px] text-gray-500 mt-1">
                {{ task.message }}
              </div>
            </div>
            <div class="text-[10px] text-gray-500 whitespace-nowrap">
              {{ formatDate(task.createdAt) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Template Cards ── -->
    <div v-else class="space-y-3">
      <div
        v-for="(template, index) in store.templates"
        :key="template.id"
        :style="{ '--i': index }"
        @click="emit('select', template.id)"
        class="card-stagger template-card rounded-2xl p-5 group cursor-pointer"
      >
        <div class="flex items-start gap-4">
          <!-- Index -->
          <div class="flex-shrink-0 w-8 text-right pt-0.5">
            <span class="num-badge">{{ String(index + 1).padStart(2, '0') }}</span>
          </div>
          <!-- Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-3 mb-2">
              <h3 class="text-[15px] font-semibold text-white group-hover:text-cyber-400 transition-colors duration-200 leading-snug">
                {{ template.name }}
              </h3>
              <div class="flex items-center gap-1.5 flex-shrink-0 mt-0.5">
                <!-- Delete button -->
                <button
                  @click.stop="openDelModal(template.id, template.name)"
                  class="w-7 h-7 rounded-lg flex items-center justify-center text-transparent
                         group-hover:text-red-400/50 hover:!text-red-400
                         border border-transparent group-hover:border-red-500/15
                         hover:!border-red-500/30 hover:!bg-red-500/8
                         transition-all duration-200"
                  title="Delete template"
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                </button>
                <!-- Arrow -->
                <div class="w-7 h-7 rounded-lg flex items-center justify-center
                            border border-transparent group-hover:border-cyber-400/20
                            group-hover:bg-cyber-400/8 transition-all duration-200">
                  <svg class="w-3.5 h-3.5 text-gray-600 group-hover:text-cyber-400 transition-all duration-200 transform group-hover:translate-x-0.5"
                       fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                  </svg>
                </div>
              </div>
            </div>
            <p class="text-gray-500 text-[13px] leading-relaxed line-clamp-2 mb-3">
              {{ template.description }}
            </p>
            <div class="flex items-center gap-2.5 flex-wrap">
              <span class="flex items-center gap-1 text-[11px] font-mono text-gray-600">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                {{ template.author || t('template.unknown') }}
              </span>
              <span v-if="template.tags?.length" class="text-gray-700 text-xs">·</span>
              <div v-if="template.tags?.length" class="flex gap-1.5 flex-wrap">
                <span v-for="tag in template.tags.slice(0, 4)" :key="tag" class="tag-chip px-2 py-0.5 rounded-md">
                  {{ tag }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ══════════════════════════════════════════
         Add Template Modal
         ══════════════════════════════════════════ -->
    <Teleport to="body">
      <Transition enter-active-class="transition-all duration-300 ease-out" leave-active-class="transition-all duration-200 ease-in" enter-from-class="opacity-0" leave-to-class="opacity-0">
        <div v-if="showAddModal" class="fixed inset-0 z-[200] flex items-center justify-center p-4" @click.self="closeAddModal">
          <div class="absolute inset-0 bg-midnight-950/80 backdrop-blur-sm"></div>
          <Transition enter-active-class="transition-all duration-300 ease-out" leave-active-class="transition-all duration-200 ease-in" enter-from-class="opacity-0 scale-95 translate-y-2" leave-to-class="opacity-0 scale-95 translate-y-2">
            <div v-if="showAddModal" class="relative z-10 w-full max-w-md glass-card glow-card rounded-2xl overflow-hidden">
              <!-- Header -->
              <div class="px-6 py-4 border-b border-white/[0.05] flex items-center justify-between">
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-cyber-400/10 border border-cyber-400/20 flex items-center justify-center">
                    <svg class="w-3.5 h-3.5 text-cyber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"/>
                    </svg>
                  </div>
                  <span class="text-sm font-semibold text-white">Add Template</span>
                </div>
                <button @click="closeAddModal" class="w-7 h-7 rounded-lg flex items-center justify-center text-gray-600 hover:text-gray-300 hover:bg-white/[0.06] transition-all duration-150">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
              <!-- Body -->
              <div class="px-6 py-5 space-y-4">
                <div>
                  <label class="block text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.1em] mb-2">Template ID</label>
                  <input v-model="formId" type="text" placeholder="e.g. my-template" class="input-field w-full rounded-xl px-4 py-2.5 text-sm font-mono" @keyup.enter="handleAdd"/>
                </div>
                <div>
                  <label class="block text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.1em] mb-2">Repository URL</label>
                  <input v-model="formRepo" type="url" placeholder="https://github.com/user/repo" class="input-field w-full rounded-xl px-4 py-2.5 text-sm font-mono" @keyup.enter="handleAdd"/>
                </div>
                <div class="bg-midnight-800/30 rounded-xl p-4 border border-white/[0.05]">
                  <div class="flex items-center justify-between mb-3">
                    <label class="block text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.1em]">
                      <span class="flex items-center gap-1.5">
                        <svg class="w-3 h-3 text-copper-400/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
                        </svg>
                        Access Key
                      </span>
                    </label>
                    <button 
                      @click="showAccessKeyInput = !showAccessKeyInput"
                      class="text-xs text-cyber-400 hover:text-cyber-300 transition-colors"
                    >
                      {{ showAccessKeyInput ? 'Hide' : 'Show' }}
                    </button>
                  </div>
                  <div v-if="showAccessKeyInput" class="mb-3 relative">
                    <input 
                      v-model="accessKey" 
                      :type="showAccessKeyPassword ? 'text' : 'password'" 
                      placeholder="Enter access key" 
                      class="input-field w-full rounded-xl px-4 py-2.5 pr-12 text-sm font-mono tracking-widest"
                      @blur="saveAccessKey(accessKey)"
                    />
                    <button
                      type="button"
                      @click="showAccessKeyPassword = !showAccessKeyPassword"
                      class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 transition-colors"
                      title="Toggle password visibility"
                    >
                      <svg v-if="showAccessKeyPassword" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                      </svg>
                      <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                      </svg>
                    </button>
                  </div>
                  <div v-else class="text-xs text-gray-500 truncate">
                    {{ accessKey ? '••••••••••••' : 'No access key set' }}
                  </div>
                  <div v-if="accessKey" class="mt-2 flex items-center gap-2">
                    <button 
                      @click="clearAccessKey"
                      class="text-xs text-red-400 hover:text-red-300 transition-colors"
                    >
                      Clear Key
                    </button>
                  </div>
                </div>
                <Transition enter-active-class="transition-all duration-200 ease-out" leave-active-class="transition-all duration-150 ease-in" enter-from-class="opacity-0 -translate-y-1" leave-to-class="opacity-0">
                  <div v-if="formError" class="flex items-center gap-2 px-3.5 py-2.5 rounded-xl bg-red-500/[0.06] border border-red-500/20">
                    <svg class="w-3.5 h-3.5 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    <span class="text-xs text-red-400">{{ formError }}</span>
                  </div>
                </Transition>
              </div>
              <!-- Footer -->
              <div class="px-6 py-4 border-t border-white/[0.05] flex items-center justify-end gap-2.5">
                <button @click="closeAddModal" class="btn-ghost px-5 py-2.5 rounded-xl text-sm font-medium" :disabled="submitting">Cancel</button>
                <button @click="handleAdd" :disabled="submitting" class="btn-primary px-5 py-2.5 rounded-xl text-sm font-semibold text-midnight-950 disabled:opacity-40 disabled:cursor-not-allowed disabled:transform-none">
                  <span class="flex items-center gap-2">
                    <svg v-if="submitting" class="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
                    </svg>
                    <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"/>
                    </svg>
                    {{ submitting ? 'Adding...' : 'Add Template' }}
                  </span>
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </Transition>
    </Teleport>

    <!-- ══════════════════════════════════════════
         Delete Confirm Modal
         ══════════════════════════════════════════ -->
    <Teleport to="body">
      <Transition enter-active-class="transition-all duration-300 ease-out" leave-active-class="transition-all duration-200 ease-in" enter-from-class="opacity-0" leave-to-class="opacity-0">
        <div v-if="showDelModal" class="fixed inset-0 z-[200] flex items-center justify-center p-4" @click.self="closeDelModal">
          <div class="absolute inset-0 bg-midnight-950/80 backdrop-blur-sm"></div>
          <Transition enter-active-class="transition-all duration-300 ease-out" leave-active-class="transition-all duration-200 ease-in" enter-from-class="opacity-0 scale-95 translate-y-2" leave-to-class="opacity-0 scale-95 translate-y-2">
            <div v-if="showDelModal" class="relative z-10 w-full max-w-sm glass-card rounded-2xl overflow-hidden border border-red-500/20">
              <!-- Header -->
              <div class="px-6 py-4 border-b border-red-500/10 flex items-center justify-between">
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-red-500/10 border border-red-500/20 flex items-center justify-center">
                    <svg class="w-3.5 h-3.5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </div>
                  <span class="text-sm font-semibold text-white">Delete Template</span>
                </div>
                <button @click="closeDelModal" class="w-7 h-7 rounded-lg flex items-center justify-center text-gray-600 hover:text-gray-300 hover:bg-white/[0.06] transition-all duration-150">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
              <!-- Body -->
              <div class="px-6 py-5 space-y-4">
                <!-- Warning -->
                <div class="flex items-start gap-3 px-3.5 py-3 rounded-xl bg-red-500/[0.05] border border-red-500/15">
                  <svg class="w-4 h-4 text-red-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                  </svg>
                  <div>
                    <p class="text-xs text-red-400/90 font-medium">即将删除以下模板</p>
                    <p class="text-[11px] font-mono text-red-400/60 mt-0.5 truncate">{{ delTargetName }}</p>
                  </div>
                </div>
                <!-- Access Key Section -->
                <div class="bg-midnight-800/30 rounded-xl p-4 border border-white/[0.05]">
                  <div class="flex items-center justify-between mb-3">
                    <label class="block text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.1em]">
                      <span class="flex items-center gap-1.5">
                        <svg class="w-3 h-3 text-copper-400/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
                        </svg>
                        Access Key Required
                      </span>
                    </label>
                    <button 
                      @click="showAccessKeyInput = !showAccessKeyInput"
                      class="text-xs text-cyber-400 hover:text-cyber-300 transition-colors"
                    >
                      {{ showAccessKeyInput ? 'Hide' : 'Show' }}
                    </button>
                  </div>
                  <div v-if="showAccessKeyInput" class="mb-3 relative">
                    <input 
                      v-model="accessKey" 
                      :type="showAccessKeyPassword ? 'text' : 'password'" 
                      placeholder="Enter access key" 
                      class="input-field w-full rounded-xl px-4 py-2.5 pr-12 text-sm font-mono tracking-widest"
                      @blur="saveAccessKey(accessKey)"
                    />
                    <button
                      type="button"
                      @click="showAccessKeyPassword = !showAccessKeyPassword"
                      class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 transition-colors"
                      title="Toggle password visibility"
                    >
                      <svg v-if="showAccessKeyPassword" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                      </svg>
                      <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                      </svg>
                    </button>
                  </div>
                  <div v-else class="text-xs text-gray-500 truncate">
                    {{ accessKey ? '••••••••••••' : 'No access key set' }}
                  </div>
                  <div v-if="accessKey" class="mt-2 flex items-center gap-2">
                    <button 
                      @click="clearAccessKey"
                      class="text-xs text-red-400 hover:text-red-300 transition-colors"
                    >
                      Clear Key
                    </button>
                  </div>
                  <div v-if="!accessKey" class="mt-2 text-xs text-yellow-400">
                    Please set an access key to delete templates
                  </div>
                </div>
                <Transition enter-active-class="transition-all duration-200 ease-out" leave-active-class="transition-all duration-150 ease-in" enter-from-class="opacity-0 -translate-y-1" leave-to-class="opacity-0">
                  <div v-if="delError" class="flex items-center gap-2 px-3.5 py-2.5 rounded-xl bg-red-500/[0.06] border border-red-500/20">
                    <svg class="w-3.5 h-3.5 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    <span class="text-xs text-red-400">{{ delError }}</span>
                  </div>
                </Transition>
              </div>
              <!-- Footer -->
              <div class="px-6 py-4 border-t border-white/[0.05] flex items-center justify-end gap-2.5">
                <button @click="closeDelModal" class="btn-ghost px-5 py-2.5 rounded-xl text-sm font-medium" :disabled="deleting">Cancel</button>
                <button @click="handleDelete" :disabled="deleting"
                  class="flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-semibold bg-red-500/80 hover:bg-red-500 text-white border border-red-500/30 hover:border-red-400/50 shadow-lg shadow-red-500/15 hover:shadow-red-500/25 transition-all duration-200 disabled:opacity-40 disabled:cursor-not-allowed"
                >
                  <svg v-if="deleting" class="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
                  </svg>
                  <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                  {{ deleting ? 'Deleting...' : 'Confirm Delete' }}
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>