<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
      <p class="text-gray-600 dark:text-gray-400">Loading your dashboard...</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

onMounted(() => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }

  const userRole = authStore.userRole
  console.log('User role:', userRole)

  switch (userRole) {
    case 'admin':
      router.replace('/admin/dashboard')
      break
    case 'teacher':
      router.replace('/teacher/dashboard')
      break
    case 'user':
    default:
      router.replace('/student/dashboard')
      break
  }
})
</script>
