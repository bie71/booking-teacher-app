<template>
  <!-- Mobile: Full width at top, Desktop: Fixed position at top-right -->
  <div class="fixed top-0 left-0 right-0 z-50 p-2 sm:top-4 sm:right-4 sm:left-auto sm:max-w-lg sm:w-auto space-y-2">
    <transition-group
      name="notification"
      tag="div"
      class="space-y-2"
    >
      <div
        v-for="notification in uiStore.activeNotifications"
        :key="notification.id"
        :class="getNotificationClasses(notification.type)"
        class="w-full sm:w-auto bg-white dark:bg-gray-800 shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden"
      >
        <div class="p-3 sm:p-4">
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <!-- Success Icon -->
              <svg
                v-if="notification.type === 'success'"
                class="h-5 w-5 sm:h-6 sm:w-6 text-green-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              
              <!-- Error Icon -->
              <svg
                v-else-if="notification.type === 'error'"
                class="h-5 w-5 sm:h-6 sm:w-6 text-red-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              
              <!-- Warning Icon -->
              <svg
                v-else-if="notification.type === 'warning'"
                class="h-5 w-5 sm:h-6 sm:w-6 text-yellow-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
                />
              </svg>
              
              <!-- Info Icon -->
              <svg
                v-else
                class="h-5 w-5 sm:h-6 sm:w-6 text-blue-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            
            <div class="ml-3 flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 dark:text-gray-100 break-words">
                {{ notification.message }}
              </p>
            </div>
            
            <div class="ml-2 flex-shrink-0 flex">
              <button
                @click="uiStore.removeNotification(notification.id)"
                class="bg-white dark:bg-gray-800 rounded-md inline-flex text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <span class="sr-only">Close</span>
                <svg class="h-4 w-4 sm:h-5 sm:w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
        
        <!-- Progress Bar -->
        <div
          v-if="notification.duration > 0"
          class="h-1 bg-gray-200 dark:bg-gray-700"
        >
          <div
            :class="getProgressBarClasses(notification.type)"
            class="h-full transition-all duration-100 ease-linear"
            :style="{ width: getProgressWidth(notification) + '%' }"
          ></div>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useUIStore } from '@/stores/ui'

const uiStore = useUIStore()
const progressIntervals = ref(new Map())

const getNotificationClasses = (type) => {
  const baseClasses = 'transform transition-all duration-300 ease-in-out'
  
  switch (type) {
    case 'success':
      return `${baseClasses} border-l-4 border-green-400`
    case 'error':
      return `${baseClasses} border-l-4 border-red-400`
    case 'warning':
      return `${baseClasses} border-l-4 border-yellow-400`
    default:
      return `${baseClasses} border-l-4 border-blue-400`
  }
}

const getProgressBarClasses = (type) => {
  switch (type) {
    case 'success':
      return 'bg-green-400'
    case 'error':
      return 'bg-red-400'
    case 'warning':
      return 'bg-yellow-400'
    default:
      return 'bg-blue-400'
  }
}

const getProgressWidth = (notification) => {
  if (notification.duration <= 0) return 100
  
  const elapsed = Date.now() - notification.timestamp.getTime()
  const progress = Math.max(0, 100 - (elapsed / notification.duration) * 100)
  
  return progress
}

const startProgressTracking = (notification) => {
  if (notification.duration <= 0) return
  
  const interval = setInterval(() => {
    const elapsed = Date.now() - notification.timestamp.getTime()
    if (elapsed >= notification.duration) {
      clearInterval(interval)
      progressIntervals.value.delete(notification.id)
    }
  }, 100)
  
  progressIntervals.value.set(notification.id, interval)
}

const stopProgressTracking = (notificationId) => {
  const interval = progressIntervals.value.get(notificationId)
  if (interval) {
    clearInterval(interval)
    progressIntervals.value.delete(notificationId)
  }
}

// Watch for new notifications
const unwatchNotifications = uiStore.$subscribe((mutation, state) => {
  if (mutation.type === 'direct' && mutation.events?.some(e => e.key === 'notifications')) {
    // Start progress tracking for new notifications
    state.notifications.forEach(notification => {
      if (!progressIntervals.value.has(notification.id)) {
        startProgressTracking(notification)
      }
    })
    
    // Clean up intervals for removed notifications
    progressIntervals.value.forEach((interval, id) => {
      if (!state.notifications.some(n => n.id === id)) {
        stopProgressTracking(id)
      }
    })
  }
})

onMounted(() => {
  // Start progress tracking for existing notifications
  uiStore.activeNotifications.forEach(notification => {
    startProgressTracking(notification)
  })
})

onUnmounted(() => {
  // Clean up all intervals
  progressIntervals.value.forEach(interval => clearInterval(interval))
  progressIntervals.value.clear()
  
  // Unwatch notifications
  unwatchNotifications()
})
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.notification-move {
  transition: transform 0.3s ease;
}

/* Mobile responsive adjustments */
@media (max-width: 640px) {
  .notification-enter-from {
    opacity: 0;
    transform: translateY(-100%);
  }

  .notification-leave-to {
    opacity: 0;
    transform: translateY(-100%);
  }
}
</style>
