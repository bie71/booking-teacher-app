<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 ">
    <div class="bg-white rounded-lg shadow-lg max-w-4xl w-full p-6 relative dark:bg-gray-800">
      
      <!-- Close Button
      <button @click="$emit('close')" class="absolute top-2 right-2 text-gray-500 hover:text-gray-700 dark:text-gray-400">
        ✕
      </button> -->

      <!-- Header -->
      <div class="mb-4 border-b pb-2">
        <slot name="header" />
      </div>

      <!-- Body -->
      <div class="mb-4 max-h-[80vh] overflow-y-auto">
        <slot name="body" />
      </div>

      <!-- Footer -->
      <div class="pt-2 border-t">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: {
    type: Boolean,
    required: true
  }
})
</script>
