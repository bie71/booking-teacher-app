<template>
  <div>
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900">
        Teacher Schedule Management
      </h2>
      <button 
        @click="$emit('add-schedule')"
        class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
      >
        Add New Schedule
      </button>
    </div>

    <!-- Search and Filters -->
    <div class="mb-4 flex space-x-4">
      <div class="flex-1">
        <input 
          v-model="searchQuery"
          type="text"
          placeholder="Search schedules..."
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <select 
        v-model="levelFilter"
        class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">All Levels</option>
        <option v-for="level in teacherLevels" :key="level.value" :value="level.value">
          {{ level.label }}
        </option>
      </select>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
      <p class="mt-2 text-gray-600">Loading schedules...</p>
    </div>

    <!-- Schedules Table -->
    <div v-else class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Schedule
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Level
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Price/Hour
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="schedule in filteredSchedules" :key="schedule.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">
                {{ schedule.day ? schedule.day.charAt(0).toUpperCase() + schedule.day.slice(1) : 'N/A' }}
              </div>
              <div class="text-sm text-gray-500">
                {{ schedule.start_time }} - {{ schedule.end_time }}
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                {{ schedule.max_students || 0 }} students
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              ${{ schedule.price_per_hour || 0 }}/hour
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getStatusClass(schedule.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ schedule.status || 'Available' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex space-x-2">
                <button 
                  @click="editSchedule(schedule)"
                  class="text-blue-600 hover:text-blue-900"
                  title="Edit Schedule"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button 
                  @click="deleteSchedule(schedule)"
                  class="text-red-600 hover:text-red-900"
                  title="Delete Schedule"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Empty State -->
      <div v-if="filteredSchedules.length === 0" class="text-center py-8">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No schedules found</h3>
        <p class="mt-1 text-sm text-gray-500">Get started by adding a new schedule.</p>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="schedules.length > 0" class="mt-4 flex justify-between items-center">
      <div class="text-sm text-gray-700">
        Showing {{ (pagination.page - 1) * pagination.limit + 1 }} to 
        {{ Math.min(pagination.page * pagination.limit, pagination.total) }} of 
        {{ pagination.total }} results
      </div>
      <div class="flex space-x-2">
        <button 
          @click="previousPage"
          :disabled="pagination.page <= 1"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
        >
          Previous
        </button>
        <button 
          @click="nextPage"
          :disabled="pagination.page >= Math.ceil(pagination.total / pagination.limit)"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useTeachersStore } from '@/stores/teachers'
import { useUIStore } from '@/stores/ui'
import { teacherService } from '@/services/api'

const emit = defineEmits(['edit-schedule', 'delete-schedule', 'add-schedule'])

const teachersStore = useTeachersStore()
const uiStore = useUIStore()

const schedules = ref([])
const pagination = ref({
  page: 1,
  limit: 10,
  total: 0
})

const searchQuery = ref('')
const levelFilter = ref('')
const teacherLevels = computed(() => teachersStore.teacherLevels)
const isLoading = ref(false)

const filteredSchedules = computed(() => {
  let filtered = schedules.value

  // Search filter
  if (searchQuery.value) {
    const searchTerm = searchQuery.value.toLowerCase()
    filtered = filtered.filter(schedule => {
      const name = (schedule.name || '').toLowerCase()
      const email = (schedule.email || '').toLowerCase()
      return name.includes(searchTerm) || email.includes(searchTerm)
    })
  }

  // Level filter
  if (levelFilter.value) {
    filtered = filtered.filter(schedule => {
      const level = schedule.language_level || schedule.LanguageLevel
      return level === levelFilter.value
    })
  }

  return filtered
})

// Methods
const editSchedule = (schedule) => {
  emit('edit-schedule', schedule)
}

const deleteSchedule = (schedule) => {
  emit('delete-schedule', schedule)
}

const loadSchedules = async () => {
  try {
    isLoading.value = true
    const response = await teacherService.getSchedules({
      page: pagination.value.page,
      limit: pagination.value.limit
    })
    
    schedules.value = response.schedules || response.data || response
    if (response.pagination) {
      pagination.value = {
        page: response.pagination.current_page || pagination.value.page,
        limit: response.pagination.limit || pagination.value.limit,
        total: response.pagination.total_data || response.pagination.total || 0
      }
    }
  } catch (error) {
    uiStore.showError('Failed to load schedules')
  } finally {
    isLoading.value = false
  }
}

const previousPage = () => {
  if (pagination.value.page > 1) {
    pagination.value.page -= 1
    loadSchedules()
  }
}

const nextPage = () => {
  const totalPages = Math.ceil(pagination.value.total / pagination.value.limit)
  if (pagination.value.page < totalPages) {
    pagination.value.page += 1
    loadSchedules()
  }
}

onMounted(() => {
  loadSchedules()
})

// Watch for changes
watch([searchQuery, levelFilter], () => {
  teachersStore.setPage(1)
  loadSchedules()
})
</script>
