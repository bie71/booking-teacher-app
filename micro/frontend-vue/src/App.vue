<script setup>
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
import NotificationContainer from '@/components/ui/NotificationContainer.vue'
import LoadingOverlay from '@/components/ui/LoadingOverlay.vue'

// Initialize stores with error handling
let authStore = null
let uiStore = null
const storesInitialized = ref(false)

try {
  authStore = useAuthStore()
  uiStore = useUIStore()
  storesInitialized.value = true
} catch (error) {
  console.error('Failed to initialize stores:', error)
}

onMounted(async () => {
  if (!storesInitialized.value) {
    console.warn('Stores not initialized, skipping initialization')
    return
  }

  try {
    // Initialize theme
    if (uiStore?.initializeTheme) {
      uiStore.initializeTheme()
    }
    
    // Initialize keyboard shortcuts
    if (uiStore?.initializeKeyboardShortcuts) {
      uiStore.initializeKeyboardShortcuts()
    }
    
    // Initialize auth state
    if (authStore?.initializeAuth) {
      await authStore.initializeAuth()
    }
  } catch (error) {
    console.error('Error during app initialization:', error)
  }
})
</script>

<template>
  <div id="app" class="min-h-screen flex flex-col bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
    <!-- Header -->
    <AppHeader v-if="storesInitialized" />
    
    <!-- Main Content -->
    <main class="flex-1 pt-16">
      <router-view />
    </main>
    
    <!-- Footer -->
    <AppFooter v-if="storesInitialized" />
    
    <!-- Global Components -->
    <NotificationContainer v-if="storesInitialized" />
    <LoadingOverlay v-if="storesInitialized && uiStore?.isLoading" />
    
    <!-- Error fallback when stores fail to initialize -->
    <div v-if="!storesInitialized" class="flex-1 flex items-center justify-center bg-gray-50">
      <div class="text-center">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">Loading Application...</h1>
        <p class="text-gray-600">Please wait while we initialize the application.</p>
      </div>
    </div>
  </div>
</template>

<style>
/* Global styles */
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Japanese font support */
.japanese {
  font-family: 'Noto Sans JP', 'Hiragino Sans', 'Hiragino Kaku Gothic ProN', 'Yu Gothic', 'Meiryo', sans-serif;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* Dark mode scrollbar */
.dark ::-webkit-scrollbar-track {
  background: #374151;
}

.dark ::-webkit-scrollbar-thumb {
  background: #6b7280;
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

/* Focus styles */
.focus-ring {
  @apply focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from {
  transform: translateX(-100%);
}

.slide-leave-to {
  transform: translateX(100%);
}

/* Button styles */
.btn {
  @apply inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors duration-200;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500;
}

.btn-success {
  @apply bg-green-600 text-white hover:bg-green-700 focus:ring-green-500;
}

.btn-danger {
  @apply bg-red-600 text-white hover:bg-red-700 focus:ring-red-500;
}

.btn-outline {
  @apply bg-transparent border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-gray-500;
}

/* Card styles */
.card {
  @apply bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700;
}

.card-header {
  @apply px-6 py-4 border-b border-gray-200 dark:border-gray-700;
}

.card-body {
  @apply px-6 py-4;
}

.card-footer {
  @apply px-6 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50;
}

/* Form styles */
.form-group {
  @apply mb-4;
}

.form-label {
  @apply block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2;
}

.form-input {
  @apply block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
}

.form-select {
  @apply block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
}

  .form-textarea {
    @apply block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
    resize: vertical;
  }

.form-error {
  @apply text-red-600 dark:text-red-400 text-sm mt-1;
}

/* Badge styles */
.badge {
  @apply inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium;
}

.badge-success {
  @apply bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100;
}

.badge-warning {
  @apply bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100;
}

.badge-error {
  @apply bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100;
}

.badge-info {
  @apply bg-blue-100 text-blue-800 dark:bg-blue-800 dark:text-blue-100;
}

.badge-gray {
  @apply bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100;
}

/* Loading animation */
.loading-spinner {
  @apply animate-spin rounded-full border-2 border-gray-300 border-t-blue-600;
}

/* Responsive utilities */
.container {
  @apply max-w-7xl mx-auto px-4 sm:px-6 lg:px-8;
}

.section {
  @apply py-12 sm:py-16 lg:py-20;
}

/* Print styles */
@media print {
  .no-print {
    display: none !important;
  }
}
</style>
