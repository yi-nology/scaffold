<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useScaffoldStore } from '../stores/scaffold'
import TemplateDetail from './TemplateDetail.vue'

const { t } = useI18n()
const store = useScaffoldStore()

const emit = defineEmits<{
  (e: 'generate'): void
  (e: 'back'): void
}>()

const groupedVariables = computed(() => {
  const groups: Record<string, typeof store.variables> = {}
  for (const v of store.variables) {
    const group = v.group || 'Other'
    if (!groups[group]) groups[group] = []
    groups[group].push(v)
  }
  return Object.entries(groups).filter(([_, vars]) => vars.length > 0)
})

const getGroupLabel = (group: string) => {
  const key = `form.groups.${group.toLowerCase()}`
  const translated = t(key)
  return translated === key ? group : translated
}

const groupDotClass = (index: number) => {
  const classes = ['group-dot-0', 'group-dot-1', 'group-dot-2', 'group-dot-3']
  return classes[index % classes.length]
}

const selectedTagInfo = computed(() => {
  if (!store.selectedVersion) {
    return store.availableTags.length > 0 ? store.availableTags[0] : null
  }
  return store.availableTags.find(tag => tag.name === store.selectedVersion)
})

const latestVersion = computed(() => {
  return store.availableTags.length > 0 ? store.availableTags[0].name : ''
})

function getValue(name: string): any {
  return store.values[name]
}

function setValue(name: string, value: any) {
  store.setValue(name, value)
}

function toggleBoolean(name: string) {
  store.setValue(name, !store.values[name])
}
</script>

<template>
  <div class="space-y-5">

    <!-- ── Header ── -->
    <div class="glass-card rounded-2xl p-5">
      <div class="flex items-center gap-3">
        <button
          @click="emit('back')"
          class="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0
                 btn-ghost transition-all duration-200"
          :title="t('common.back')"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
        </button>
        <div class="flex-1 min-w-0">
          <h2 class="text-base font-semibold text-white truncate">{{ store.selectedTemplate?.name }}</h2>
          <p class="text-gray-600 text-xs mt-0.5 line-clamp-1">{{ store.selectedTemplate?.description }}</p>
        </div>
      </div>
    </div>

    <!-- ── Template Detail ── -->
    <TemplateDetail
      v-if="store.selectedTemplate"
      :template="store.selectedTemplate"
      :variables="store.variables"
    />

    <!-- ── Version Selector ── -->
    <div v-if="store.availableTags.length > 0" class="glass-card rounded-2xl overflow-hidden">
      <div class="px-5 py-3.5 border-b border-white/[0.04] flex items-center gap-2">
        <div class="w-1.5 h-1.5 rounded-full bg-copper-400 shadow-[0_0_6px_rgba(245,158,11,0.6)]"></div>
        <h3 class="text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.12em]">
          {{ t('form.version') }}
        </h3>
      </div>
      <div class="p-5">
        <div class="flex flex-wrap gap-2">
          <!-- Latest -->
          <button
            type="button"
            @click="store.setVersion('')"
            :class="[
              'px-3.5 py-2 rounded-xl text-xs font-medium transition-all duration-250',
              store.selectedVersion === ''
                ? 'bg-cyber-500/15 border border-cyber-400/50 text-cyber-400 shadow-[0_0_12px_rgba(34,211,238,0.15)]'
                : 'btn-ghost text-gray-500 hover:text-gray-300'
            ]"
          >
            {{ t('form.latestVersion') }}
            <span v-if="latestVersion" class="ml-1.5 font-mono text-[10px] opacity-60">({{ latestVersion }})</span>
          </button>
          <!-- Tags -->
          <button
            v-for="tag in store.availableTags"
            :key="tag.name"
            type="button"
            @click="store.setVersion(tag.name)"
            :class="[
              'px-3.5 py-2 rounded-xl text-xs font-mono transition-all duration-250',
              store.selectedVersion === tag.name
                ? 'bg-cyber-500/15 border border-cyber-400/50 text-cyber-400 shadow-[0_0_12px_rgba(34,211,238,0.15)]'
                : 'btn-ghost text-gray-500 hover:text-gray-300'
            ]"
          >
            {{ tag.name }}
          </button>
        </div>

        <!-- Tag message -->
        <Transition
          enter-active-class="transition-all duration-300 ease-out"
          leave-active-class="transition-all duration-200 ease-in"
          enter-from-class="opacity-0 -translate-y-1"
          leave-to-class="opacity-0 -translate-y-1"
        >
          <div v-if="selectedTagInfo?.message" class="mt-4 p-3.5 rounded-xl bg-cyber-400/[0.04] border border-cyber-400/15">
            <div class="flex items-start gap-2.5">
              <svg class="w-3.5 h-3.5 text-cyber-400/60 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <p class="text-xs text-gray-400 leading-relaxed">{{ selectedTagInfo.message }}</p>
            </div>
          </div>
        </Transition>
      </div>
    </div>

    <!-- ── Configuration Groups ── -->
    <div
      v-for="([groupName, vars], groupIndex) in groupedVariables"
      :key="groupName"
      class="glass-card rounded-2xl overflow-hidden"
    >
      <!-- Group header -->
      <div class="px-5 py-3.5 border-b border-white/[0.04] flex items-center gap-2.5">
        <div :class="['w-1.5 h-1.5 rounded-full', groupDotClass(groupIndex)]"></div>
        <h3 class="text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.12em]">
          {{ getGroupLabel(groupName) }}
        </h3>
      </div>

      <!-- Fields -->
      <div class="p-5">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <div
            v-for="v in vars"
            :key="v.name"
            :class="[
              v.type === 'string' && v.name.toLowerCase().includes('description') ? 'sm:col-span-2' : '',
              v.type === 'enum' && (v.options?.length ?? 0) > 3 ? 'sm:col-span-2' : '',
            ]"
          >
            <!-- Label -->
            <label class="flex items-center gap-1.5 text-[12px] font-medium text-gray-400 mb-2">
              {{ v.prompt || v.name }}
              <span v-if="v.required" class="text-copper-400 text-[10px]">*</span>
              <span v-if="v.description" class="text-gray-600 font-normal text-[11px] truncate ml-1">
                — {{ v.description }}
              </span>
            </label>

            <!-- Textarea -->
            <textarea
              v-if="v.type === 'string' && v.name.toLowerCase().includes('description')"
              :value="getValue(v.name)"
              @input="setValue(v.name, ($event.target as HTMLTextAreaElement).value)"
              :placeholder="String(v.default || '')"
              rows="3"
              class="input-field w-full rounded-xl px-4 py-3 text-sm placeholder-gray-600 font-light resize-none"
            />

            <!-- Text input -->
            <input
              v-else-if="v.type === 'string'"
              type="text"
              :value="getValue(v.name)"
              @input="setValue(v.name, ($event.target as HTMLInputElement).value)"
              :placeholder="String(v.default || '')"
              class="input-field w-full rounded-xl px-4 py-2.5 text-sm"
            />

            <!-- Number input -->
            <input
              v-else-if="v.type === 'number'"
              type="number"
              :value="getValue(v.name)"
              @input="setValue(v.name, Number(($event.target as HTMLInputElement).value))"
              :placeholder="String(v.default || '')"
              class="input-field w-full rounded-xl px-4 py-2.5 text-sm font-mono"
            />

            <!-- Boolean toggle -->
            <div v-else-if="v.type === 'boolean'" class="flex items-center gap-3">
              <button
                type="button"
                @click="toggleBoolean(v.name)"
                :class="[
                  'relative w-12 h-6 rounded-full transition-all duration-300 focus:outline-none',
                  getValue(v.name)
                    ? 'bg-cyber-500/30 border border-cyber-400/50 shadow-[0_0_10px_rgba(34,211,238,0.2)]'
                    : 'bg-midnight-700 border border-white/[0.08]'
                ]"
              >
                <span :class="[
                  'absolute top-[3px] w-[18px] h-[18px] rounded-full shadow-sm transition-all duration-300',
                  getValue(v.name)
                    ? 'left-[23px] bg-cyber-400'
                    : 'left-[3px] bg-gray-500'
                ]"/>
              </button>
              <span class="text-xs font-mono" :class="getValue(v.name) ? 'text-cyber-400' : 'text-gray-600'">
                {{ getValue(v.name) ? 'ON' : 'OFF' }}
              </span>
            </div>

            <!-- Enum options -->
            <div v-else-if="v.type === 'enum'" class="flex flex-wrap gap-2">
              <button
                v-for="opt in v.options"
                :key="opt"
                type="button"
                @click="setValue(v.name, opt)"
                :class="[
                  'px-3.5 py-2 rounded-xl text-xs font-mono transition-all duration-250',
                  getValue(v.name) === opt
                    ? 'bg-cyber-500/15 border border-cyber-400/50 text-cyber-400 shadow-[0_0_10px_rgba(34,211,238,0.12)]'
                    : 'btn-ghost text-gray-500 hover:text-gray-300'
                ]"
              >
                {{ opt }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Generate Button ── -->
    <div class="flex justify-end pt-1 pb-2">
      <button
        @click="emit('generate')"
        :disabled="store.loading"
        class="btn-primary px-8 py-3.5 rounded-xl font-semibold text-midnight-950 text-sm
               disabled:opacity-40 disabled:cursor-not-allowed disabled:transform-none"
      >
        <span class="flex items-center gap-2">
          <svg v-if="!store.loading" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
          <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
          </svg>
          {{ store.loading ? t('common.generating') : t('common.generate') }}
        </span>
      </button>
    </div>

  </div>
</template>
