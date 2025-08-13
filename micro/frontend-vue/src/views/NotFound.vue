<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
    <div class="text-center">
      <div class="mb-8">
        <h1 class="text-9xl font-bold text-gray-300 dark:text-gray-700">404</h1>
      </div>
      
      <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          Page Not Found
        </h2>
        <p class="text-xl text-gray-600 dark:text-gray-400 mb-2">
          The page you're looking for doesn't exist.
        </p>
        <p class="text-japanese text-gray-500 dark:text-gray-500">
          お探しのページが見つかりません。
        </p>
      </div>
      
      <div class="space-y-4">
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <router-link to="/" class="btn btn-primary">
            Go Home
          </router-link>
          <router-link to="/teachers" class="btn btn-outline">
            Browse Teachers
          </router-link>
        </div>
        
        <button
          @click="goBack"
          class="btn-link"
        >
          ← Go Back
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUIStore } from '@/stores/ui'

const router = useRouter()
const uiStore = useUIStore()

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}

onMounted(() => {
  uiStore.setPageTitle('Page Not Found')
})
</script>
