
<template>
  <div class="flex items-center justify-between border-t border-gray-200 dark:border-gray-700 px-4 py-3 sm:px-6">
    <div class="flex flex-1 justify-between sm:hidden">
      <button @click="$emit('update:page', Math.max(1, page-1))"
        class="relative inline-flex items-center rounded-md px-4 py-2 text-sm font-medium bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
        Previous
      </button>
      <button @click="$emit('update:page', Math.min(totalPages, page+1))"
        class="relative ml-3 inline-flex items-center rounded-md px-4 py-2 text-sm font-medium bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
        Next
      </button>
    </div>
    <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
      <p class="text-sm text-gray-700 dark:text-gray-300">
        Showing
        <span class="font-medium">{{ startItem }}</span>
        to
        <span class="font-medium">{{ endItem }}</span>
        of
        <span class="font-medium">{{ total }}</span>
        results
      </p>
      <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
        <button
          @click="$emit('update:page', Math.max(1, page-1))"
          :disabled="page === 1"
          class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 focus:outline-none disabled:opacity-40">
          <span class="sr-only">Previous</span>
          ‹
        </button>
        <button
          v-for="p in pagesToShow"
          :key="p"
          @click="$emit('update:page', p)"
          :class="[
            'relative inline-flex items-center px-4 py-2 text-sm font-semibold ring-1 ring-inset ring-gray-300 focus:outline-none',
            p === page ? 'z-10 bg-indigo-600 text-white' : 'bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800'
          ]">
          {{ p }}
        </button>
        <button
          @click="$emit('update:page', Math.min(totalPages, page+1))"
          :disabled="page === totalPages"
          class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 focus:outline-none disabled:opacity-40">
          <span class="sr-only">Next</span>
          ›
        </button>
      </nav>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({
  page: { type: Number, default: 1 },
  limit: { type: Number, default: 10 },
  total: { type: Number, default: 0 }
})

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)))
const startItem = computed(() => props.total === 0 ? 0 : ((props.page - 1) * props.limit) + 1)
const endItem = computed(() => Math.min(props.total, props.page * props.limit))

const pagesToShow = computed(() => {
  const pages = []
  const start = Math.max(1, props.page - 2)
  const end = Math.min(totalPages.value, props.page + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})
</script>
