<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 z-100 flex items-center justify-center">
    <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-xl">
      <div class="flex items-center space-x-4">
        <div class="loading-spinner w-8 h-8"></div>
        <div class="text-gray-900 dark:text-gray-100">
          <p class="font-medium">{{ message }}</p>
          <p v-if="subMessage" class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {{ subMessage }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: String,
    default: 'Loading...'
  },
  subMessage: {
    type: String,
    default: ''
  }
})
</script>
