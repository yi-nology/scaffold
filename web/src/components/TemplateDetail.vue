<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { TemplateMeta, Variable } from '../stores/scaffold'

const { t } = useI18n()

const props = withDefaults(defineProps<{
  template: TemplateMeta
  variables?: Variable[]
}>(), {
  variables: () => []
})

const isExpanded = ref(false)

const hasLinks = computed(() => 
  props.template.repository || 
  props.template.homepage || 
  props.template.bugs
)

const hasKeywords = computed(() => 
  props.template.keywords && props.template.keywords.length > 0
)

const hasTags = computed(() => 
  props.template.tags && props.template.tags.length > 0
)

const hasVariables = computed(() => props.variables.length > 0)

const groupedVariables = computed(() => {
  const groups: Record<string, Variable[]> = {}
  for (const v of props.variables) {
    const group = v.group || 'Other'
    if (!groups[group]) groups[group] = []
    groups[group].push(v)
  }
  return Object.entries(groups)
})

const typeColorClass = (type: string) => {
  switch (type) {
    case 'string': return 'bg-gray-500/15 text-gray-400 border-gray-500/20'
    case 'number': return 'bg-cyber-500/15 text-cyber-400 border-cyber-400/20'
    case 'boolean': return 'bg-violet-500/15 text-violet-400 border-violet-400/20'
    case 'enum': return 'bg-copper-400/15 text-copper-400 border-copper-400/20'
    default: return 'bg-gray-500/15 text-gray-400 border-gray-500/20'
  }
}

const groupDotClass = (index: number) => {
  const classes = [
    'bg-cyber-400 shadow-[0_0_6px_rgba(34,211,238,0.6)]',
    'bg-violet-400 shadow-[0_0_6px_rgba(167,139,250,0.6)]',
    'bg-copper-400 shadow-[0_0_6px_rgba(245,158,11,0.6)]',
    'bg-emerald-400 shadow-[0_0_6px_rgba(52,211,153,0.6)]'
  ]
  return classes[index % classes.length]
}

function formatDefault(val: any): string {
  if (val === null || val === undefined) return '—'
  if (typeof val === 'boolean') return val ? 'true' : 'false'
  return String(val)
}
</script>

<template>
  <div class="template-detail glass-card rounded-2xl overflow-hidden">
    <!-- Toggle Header -->
    <button
      @click="isExpanded = !isExpanded"
      class="w-full px-5 py-3.5 flex items-center justify-between
             hover:bg-white/[0.02] transition-colors duration-200"
    >
      <div class="flex items-center gap-2.5">
        <div class="w-1.5 h-1.5 rounded-full bg-cyber-400 shadow-[0_0_6px_rgba(34,211,238,0.6)]"></div>
        <h3 class="text-[11px] font-mono font-medium text-gray-500 uppercase tracking-[0.12em]">
          {{ t('template.detail.projectInfo') }}
        </h3>
        <!-- Meta badges inline -->
        <span v-if="template.version"
              class="px-2 py-0.5 rounded-md text-[10px] font-mono
                     bg-cyber-500/10 text-cyber-400/70 border border-cyber-400/15">
          v{{ template.version }}
        </span>
        <span v-if="template.license"
              class="px-2 py-0.5 rounded-md text-[10px] font-mono
                     bg-violet-500/10 text-violet-400/70 border border-violet-400/15">
          {{ template.license }}
        </span>
        <span v-if="template.author"
              class="px-2 py-0.5 rounded-md text-[10px] font-mono
                     bg-white/[0.04] text-gray-500 border border-white/[0.06]">
          {{ template.author }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-[10px] text-gray-600 font-mono">
          {{ isExpanded ? t('template.detail.collapse') : t('template.detail.expand') }}
        </span>
        <svg
          :class="['w-3.5 h-3.5 text-gray-600 transition-transform duration-300', isExpanded ? 'rotate-180' : '']"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
        </svg>
      </div>
    </button>

    <!-- Expandable Content -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      leave-active-class="transition-all duration-200 ease-in"
      enter-from-class="opacity-0 max-h-0"
      enter-to-class="opacity-100 max-h-[2000px]"
      leave-from-class="opacity-100 max-h-[2000px]"
      leave-to-class="opacity-0 max-h-0"
    >
      <div v-show="isExpanded" class="overflow-hidden">
        <div class="px-5 pb-5 space-y-5 border-t border-white/[0.04] pt-5">

          <!-- Description -->
          <div v-if="template.description" class="text-sm text-gray-400 leading-relaxed">
            {{ template.description }}
          </div>

          <!-- Meta Row: Author + License + Version -->
          <div class="flex flex-wrap items-center gap-x-5 gap-y-2 text-sm">
            <div v-if="template.author" class="flex items-center gap-2">
              <svg class="w-3.5 h-3.5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
              </svg>
              <span class="text-gray-600 text-xs">{{ t('template.detail.author') }}</span>
              <span class="text-gray-300 text-xs font-medium">{{ template.author }}</span>
            </div>
            <div v-if="template.license" class="flex items-center gap-2">
              <svg class="w-3.5 h-3.5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span class="text-gray-600 text-xs">{{ t('template.detail.license') }}</span>
              <span class="text-cyber-400 text-xs font-mono font-medium">{{ template.license }}</span>
            </div>
            <div v-if="template.version" class="flex items-center gap-2">
              <svg class="w-3.5 h-3.5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
              </svg>
              <span class="text-gray-600 text-xs">{{ t('template.detail.version') }}</span>
              <span class="text-gray-300 text-xs font-mono font-medium">{{ template.version }}</span>
            </div>
          </div>

          <!-- Technology Tags -->
          <div v-if="hasTags">
            <h4 class="text-[10px] font-mono text-gray-600 uppercase tracking-[0.1em] mb-2">
              {{ t('template.detail.techTags') }}
            </h4>
            <div class="flex flex-wrap gap-1.5">
              <span 
                v-for="tag in template.tags" 
                :key="tag"
                class="px-2.5 py-1 rounded-lg text-[11px] font-mono
                       bg-cyber-500/8 text-cyber-400/80 border border-cyber-400/15
                       hover:bg-cyber-500/15 hover:border-cyber-400/25 transition-all"
              >
                {{ tag }}
              </span>
            </div>
          </div>

          <!-- Keywords -->
          <div v-if="hasKeywords">
            <h4 class="text-[10px] font-mono text-gray-600 uppercase tracking-[0.1em] mb-2">
              {{ t('template.detail.keywords') }}
            </h4>
            <div class="flex flex-wrap gap-1.5">
              <span 
                v-for="keyword in template.keywords" 
                :key="keyword"
                class="px-2 py-0.5 rounded-md text-[10px] font-mono
                       bg-white/[0.04] text-gray-500 border border-white/[0.06]
                       hover:bg-white/[0.06] transition-all"
              >
                {{ keyword }}
              </span>
            </div>
          </div>

          <!-- Links -->
          <div v-if="hasLinks">
            <h4 class="text-[10px] font-mono text-gray-600 uppercase tracking-[0.1em] mb-2">
              {{ t('template.detail.links') }}
            </h4>
            <div class="flex flex-wrap gap-2">
              <a v-if="template.repository" 
                 :href="template.repository" 
                 target="_blank"
                 rel="noopener noreferrer"
                 class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                        bg-midnight-800/50 border border-white/[0.06]
                        hover:border-cyber-400/30 hover:bg-cyber-400/5 hover:text-cyber-400
                        text-gray-500 transition-all group"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
                </svg>
                {{ t('template.detail.viewSource') }}
                <svg class="w-2.5 h-2.5 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                </svg>
              </a>

              <a v-if="template.homepage" 
                 :href="template.homepage" 
                 target="_blank"
                 rel="noopener noreferrer"
                 class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                        bg-midnight-800/50 border border-white/[0.06]
                        hover:border-violet-400/30 hover:bg-violet-400/5 hover:text-violet-400
                        text-gray-500 transition-all group"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
                </svg>
                {{ t('template.detail.homepage') }}
                <svg class="w-2.5 h-2.5 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                </svg>
              </a>

              <a v-if="template.bugs" 
                 :href="template.bugs" 
                 target="_blank"
                 rel="noopener noreferrer"
                 class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs
                        bg-midnight-800/50 border border-white/[0.06]
                        hover:border-amber-400/30 hover:bg-amber-400/5 hover:text-amber-400
                        text-gray-500 transition-all group"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                {{ t('template.detail.reportBug') }}
                <svg class="w-2.5 h-2.5 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                </svg>
              </a>
            </div>
          </div>

          <!-- Variables Overview -->
          <div v-if="hasVariables">
            <h4 class="text-[10px] font-mono text-gray-600 uppercase tracking-[0.1em] mb-3">
              {{ t('template.detail.variables') }}
              <span class="ml-1.5 text-gray-700">
                ({{ t('template.detail.variableCount', { count: variables.length }) }})
              </span>
            </h4>
            <div class="space-y-3">
              <div v-for="([groupName, vars], groupIndex) in groupedVariables" :key="groupName">
                <!-- Group Label -->
                <div class="flex items-center gap-2 mb-2">
                  <div :class="['w-1 h-1 rounded-full', groupDotClass(groupIndex)]"></div>
                  <span class="text-[10px] font-mono text-gray-500 uppercase tracking-wider">
                    {{ groupName }}
                  </span>
                  <div class="flex-1 h-px bg-white/[0.04]"></div>
                </div>
                <!-- Variable Items -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-1.5">
                  <div
                    v-for="v in vars"
                    :key="v.name"
                    class="flex items-center gap-2 px-3 py-2 rounded-lg
                           bg-midnight-800/30 border border-white/[0.03]
                           hover:border-white/[0.06] transition-all"
                  >
                    <!-- Variable Name -->
                    <span class="text-xs font-mono text-gray-300 truncate min-w-0 flex-shrink">
                      {{ v.name }}
                    </span>
                    <!-- Required indicator -->
                    <span v-if="v.required" class="text-copper-400 text-[9px] flex-shrink-0">*</span>
                    <!-- Spacer -->
                    <span class="flex-1"></span>
                    <!-- Type Badge -->
                    <span :class="['px-1.5 py-0.5 rounded text-[9px] font-mono border flex-shrink-0', typeColorClass(v.type)]">
                      {{ v.type }}
                    </span>
                    <!-- Default Value -->
                    <span v-if="v.default !== undefined && v.default !== ''"
                          class="text-[10px] text-gray-600 font-mono truncate max-w-[100px] flex-shrink-0"
                          :title="formatDefault(v.default)"
                    >
                      {{ formatDefault(v.default) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.glass-card {
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.06);
}
</style>
