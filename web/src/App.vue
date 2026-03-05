<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useScaffoldStore } from './stores/scaffold'
import TemplateList from './components/TemplateList.vue'
import ConfigForm from './components/ConfigForm.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

const { t, locale } = useI18n()
const store = useScaffoldStore()
const currentStep = ref<'select' | 'config' | 'generating' | 'done'>('select')

// ── Theme management ───────────────────────────────────────────────
type Theme = 'dark' | 'light' | 'system'
const currentTheme = ref<Theme>('dark')
let mediaQuery: MediaQueryList | null = null

function applyTheme(t: Theme) {
  const root = document.documentElement
  if (t === 'system') {
    root.setAttribute('data-theme', mediaQuery?.matches ? 'dark' : 'light')
  } else {
    root.setAttribute('data-theme', t)
  }
}

function cycleTheme() {
  const order: Theme[] = ['dark', 'light', 'system']
  const next = order[(order.indexOf(currentTheme.value) + 1) % order.length]
  currentTheme.value = next
  localStorage.setItem('scaffold-theme', next)
  applyTheme(next)
}

function onSystemChange() {
  if (currentTheme.value === 'system') applyTheme('system')
}

onMounted(() => {
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', onSystemChange)
  const saved = localStorage.getItem('scaffold-theme') as Theme | null
  currentTheme.value = saved || 'dark'
  applyTheme(currentTheme.value)
})

onUnmounted(() => {
  mediaQuery?.removeEventListener('change', onSystemChange)
})

// ── Theme icon helpers ─────────────────────────────────────────────
const themeLabel = computed(() => ({
  dark: 'Dark', light: 'Light', system: 'Auto'
}[currentTheme.value]))

// ── Lang toggle ───────────────────────────────────────────────────
const toggleLocale = () => {
  const newLocale = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = newLocale
  localStorage.setItem('scaffold-locale', newLocale)
}

// ── Page flow ─────────────────────────────────────────────────────
const handleTemplateSelect = async (templateId: string) => {
  await store.selectTemplate(templateId)
  currentStep.value = 'config'
}

const handleGenerate = async () => {
  currentStep.value = 'generating'
  try {
    await store.generateProject()
    currentStep.value = 'done'
  } catch (e) {
    console.error(e)
    currentStep.value = 'config'
  }
}

const handleReset = () => {
  store.reset()
  currentStep.value = 'select'
}

const stepIndicator = computed(() => {
  const _ = locale.value
  const steps = [t('steps.select'), t('steps.config'), t('steps.generate')]
  return steps.map((step, i) => ({
    name: step,
    active: currentStep.value === ['select', 'config', 'generating'][i] ||
            (currentStep.value === 'done' && i === 2),
    completed: ['config', 'generating', 'done'].indexOf(currentStep.value) > i
  }))
})
</script>

<template>
  <div class="min-h-screen relative overflow-x-hidden">

    <!-- Grid -->
    <div class="fixed inset-0 grid-pattern opacity-70 pointer-events-none"></div>

    <!-- Ambient orbs -->
    <div class="orb fixed top-[-8%] left-[10%] w-[560px] h-[560px] bg-cyber-500/5 rounded-full blur-[130px] animate-glow-pulse pointer-events-none"></div>
    <div class="orb fixed bottom-[-8%] right-[10%] w-[560px] h-[560px] bg-violet-500/4 rounded-full blur-[130px] animate-glow-pulse pointer-events-none" style="animation-delay:2s"></div>
    <div class="orb fixed top-[45%] right-[5%] w-[280px] h-[280px] bg-copper-400/3 rounded-full blur-[100px] animate-glow-pulse pointer-events-none" style="animation-delay:4s"></div>

    <div class="relative z-10">

      <!-- ── Top Navigation ── -->
      <nav class="sticky top-0 z-50 border-b border-white/5 backdrop-blur-2xl bg-midnight-950/70">
        <div class="max-w-4xl mx-auto px-5 sm:px-8 h-[52px] flex items-center justify-between">

          <!-- Logo -->
          <div class="flex items-center gap-2.5">
            <div class="w-[26px] h-[26px] bg-gradient-to-br from-cyber-400 to-violet-500 rounded-[7px] flex items-center justify-center shadow-lg shadow-cyber-500/20 flex-shrink-0">
              <svg class="w-[14px] h-[14px] text-midnight-950" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"/>
              </svg>
            </div>
            <span class="font-semibold text-white text-[13px] tracking-wide">Scaffold</span>
            <span class="version-chip hidden sm:inline-flex items-center px-1.5 py-0.5 rounded-md text-[10px] font-mono text-cyber-400/60 bg-cyber-400/6 border border-cyber-400/12 tracking-wider">
              v2.0
            </span>
          </div>

          <!-- Right controls -->
          <div class="flex items-center gap-2">

            <!-- Theme cycle button -->
            <button
              @click="cycleTheme"
              :title="themeLabel"
              class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-[11px] font-mono
                     text-gray-600 hover:text-cyber-400
                     border border-white/[0.06] hover:border-cyber-400/20
                     bg-white/[0.02] hover:bg-cyber-400/5 transition-all duration-200"
            >
              <!-- Moon — dark -->
              <svg v-if="currentTheme === 'dark'" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
              </svg>
              <!-- Sun — light -->
              <svg v-else-if="currentTheme === 'light'" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364-.707-.707M6.343 6.343l-.707-.707m12.728 0-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 1 1-8 0 4 4 0 0 1 8 0z"/>
              </svg>
              <!-- Monitor — system -->
              <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17 9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2z"/>
              </svg>
              <span class="hidden sm:inline">{{ themeLabel }}</span>
            </button>

            <!-- Lang switcher -->
            <button
              @click="toggleLocale"
              class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-[11px] font-mono
                     text-gray-600 hover:text-cyber-400
                     border border-white/[0.06] hover:border-cyber-400/20
                     bg-white/[0.02] hover:bg-cyber-400/5 transition-all duration-200"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 0 1 6.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"/>
              </svg>
              {{ locale === 'zh-CN' ? 'EN' : '中文' }}
            </button>
          </div>
        </div>
      </nav>

      <!-- ── Page content ── -->
      <div class="max-w-4xl mx-auto px-5 sm:px-8 py-10 sm:py-14">

        <!-- Hero -->
        <header class="text-center mb-12">
          <div class="hero-badge status-badge inline-flex items-center gap-2 px-3 py-1 rounded-full border mb-7">
            <span class="badge-dot w-1.5 h-1.5 rounded-full animate-pulse"></span>
            <span class="badge-text text-[10px] font-mono tracking-[0.18em] uppercase">System Online</span>
          </div>

          <h1 class="hero-title text-[3.2rem] sm:text-[4rem] font-extrabold text-gradient tracking-tight mb-5 leading-none">
            {{ t('common.scaffold') }}
          </h1>

          <div class="hero-sub terminal-pill inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-midnight-850/80 border border-white/[0.05] shadow-inner">
            <span class="font-mono text-cyber-400/50 text-sm select-none">$</span>
            <span class="font-mono text-gray-400 text-sm">{{ t('common.subtitle') }}</span>
            <span class="cursor-blink"></span>
          </div>
        </header>

        <!-- Step indicator -->
        <div class="hero-steps flex justify-center mb-10">
          <div class="glass-card rounded-2xl px-5 py-3 inline-flex items-center gap-1.5">
            <template v-for="(step, i) in stepIndicator" :key="i">
              <div class="flex items-center gap-2">
                <div :class="[
                  'w-7 h-7 rounded-lg flex items-center justify-center transition-all duration-500 flex-shrink-0',
                  step.active
                    ? 'bg-cyber-500/15 border border-cyber-400/50 shadow-[0_0_10px_rgba(34,211,238,0.25)] step-ring-pulse'
                    : step.completed
                      ? 'bg-emerald-500/12 border border-emerald-400/35'
                      : 'bg-midnight-700/60 border border-white/[0.07]'
                ]">
                  <svg v-if="step.completed" class="w-3.5 h-3.5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                  </svg>
                  <span v-else :class="['text-[11px] font-mono', step.active ? 'text-cyber-400' : 'text-gray-600']">
                    {{ i + 1 }}
                  </span>
                </div>
                <span :class="[
                  'text-sm font-medium transition-colors duration-300 leading-none',
                  step.active ? 'text-white' : step.completed ? 'text-emerald-400/60' : 'text-gray-600'
                ]">{{ step.name }}</span>
              </div>
              <div v-if="i < stepIndicator.length - 1" class="w-10 sm:w-14 h-px step-line mx-1"></div>
            </template>
          </div>
        </div>

        <!-- Main content -->
        <main>
          <Transition
            mode="out-in"
            enter-active-class="transition-all duration-400 ease-out"
            leave-active-class="transition-all duration-250 ease-in"
            enter-from-class="opacity-0 translate-y-4"
            leave-to-class="opacity-0 -translate-y-3"
          >
            <TemplateList v-if="currentStep === 'select'" @select="handleTemplateSelect" />
            <ConfigForm  v-else-if="currentStep === 'config'" @generate="handleGenerate" @back="currentStep = 'select'" />

            <!-- Generating -->
            <div v-else-if="currentStep === 'generating'" class="glass-card glow-card rounded-3xl p-16 text-center">
              <LoadingSpinner />
              <p class="text-gray-500 mt-8 text-sm font-light">{{ t('done.generatingMessage') }}</p>
            </div>

            <!-- Done -->
            <div v-else-if="currentStep === 'done'" class="glass-card glow-card rounded-3xl p-14 text-center">
              <div class="w-16 h-16 mx-auto mb-7 relative animate-float">
                <div class="absolute inset-0 bg-emerald-400/15 rounded-2xl blur-2xl"></div>
                <div class="relative w-full h-full bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-2xl flex items-center justify-center shadow-xl shadow-emerald-500/25">
                  <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                  </svg>
                </div>
              </div>
              <h2 class="text-2xl font-bold text-white mb-2">{{ t('done.title') }}</h2>
              <p class="text-gray-500 mb-9 font-light text-sm">{{ t('done.message') }}</p>
              <div class="flex justify-center gap-3">
                <button @click="store.downloadZip()" class="btn-primary px-7 py-3 rounded-xl font-semibold text-midnight-950 text-sm">
                  <span class="flex items-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
                    </svg>
                    {{ t('common.download') }}
                  </span>
                </button>
                <button @click="handleReset" class="btn-ghost px-7 py-3 rounded-xl font-medium text-sm">
                  {{ t('common.another') }}
                </button>
              </div>
            </div>
          </Transition>
        </main>

        <!-- Footer -->
        <footer class="mt-14 text-center">
          <div class="cyber-divider mb-5 max-w-xs mx-auto"></div>
          <p class="text-gray-700 text-[11px] font-mono tracking-wide">
            design by
            <a href="https://murphyyi.com" target="_blank" rel="noopener noreferrer"
               class="text-cyber-400/45 hover:text-cyber-400 transition-colors duration-300 ml-1 underline underline-offset-2">
              murphyyi
            </a>
          </p>
        </footer>

      </div>
    </div>
  </div>
</template>
