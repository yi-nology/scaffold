<template>
  <div class="flex flex-col items-center">
    <!-- Spinner rings -->
    <div class="relative w-20 h-20">
      <!-- Outer glow -->
      <div class="absolute inset-0 rounded-full blur-2xl bg-cyber-400/10 animate-glow-pulse"></div>

      <!-- Ring 1: static -->
      <div class="absolute inset-0 rounded-full border border-cyber-400/10"></div>

      <!-- Ring 2: slow spin -->
      <svg class="absolute inset-0 w-full h-full" style="animation: spin-slow 3s linear infinite" viewBox="0 0 80 80">
        <circle cx="40" cy="40" r="36"
          fill="none"
          stroke="rgba(34,211,238,0.18)"
          stroke-width="1"
          stroke-dasharray="4 8"
          stroke-linecap="round"
        />
      </svg>

      <!-- Ring 3: main gradient arc -->
      <svg class="absolute inset-0 w-full h-full" style="animation: spin-main 1.4s linear infinite" viewBox="0 0 80 80">
        <defs>
          <linearGradient id="arc-grad" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stop-color="#22d3ee" stop-opacity="0"/>
            <stop offset="50%" stop-color="#22d3ee" stop-opacity="0.8"/>
            <stop offset="100%" stop-color="#8b5cf6" stop-opacity="0.9"/>
          </linearGradient>
        </defs>
        <circle cx="40" cy="40" r="36"
          fill="none"
          stroke="url(#arc-grad)"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-dasharray="56 170"
          stroke-dashoffset="0"
        />
      </svg>

      <!-- Ring 4: inner counter spin -->
      <svg class="absolute inset-[10px] w-[calc(100%-20px)] h-[calc(100%-20px)]" style="animation: spin-counter 2s linear infinite" viewBox="0 0 60 60">
        <circle cx="30" cy="30" r="26"
          fill="none"
          stroke="rgba(167,139,250,0.25)"
          stroke-width="1"
          stroke-dasharray="3 6"
          stroke-linecap="round"
        />
      </svg>

      <!-- Center icon -->
      <div class="absolute inset-0 flex items-center justify-center">
        <svg class="w-7 h-7 text-cyber-400/50 animate-glow-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"/>
        </svg>
      </div>
    </div>

    <!-- Animated dots -->
    <div class="flex items-center gap-1.5 mt-5">
      <span class="w-1.5 h-1.5 rounded-full bg-cyber-400/60" style="animation: dot-pulse 1.4s ease-in-out 0s infinite"></span>
      <span class="w-1.5 h-1.5 rounded-full bg-cyber-400/40" style="animation: dot-pulse 1.4s ease-in-out 0.2s infinite"></span>
      <span class="w-1.5 h-1.5 rounded-full bg-violet-400/40" style="animation: dot-pulse 1.4s ease-in-out 0.4s infinite"></span>
    </div>
  </div>
</template>

<style scoped>
@keyframes spin-main {
  from { transform: rotate(0deg); }
  to   { transform: rotate(360deg); }
}
@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to   { transform: rotate(360deg); }
}
@keyframes spin-counter {
  from { transform: rotate(0deg); }
  to   { transform: rotate(-360deg); }
}
@keyframes dot-pulse {
  0%, 80%, 100% { transform: scale(1); opacity: 0.4; }
  40%           { transform: scale(1.4); opacity: 1; }
}
</style>
