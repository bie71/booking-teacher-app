<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-full max-w-4xl mx-4 shadow-lg max-h-[90vh] overflow-y-auto">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold">Manage Schedule for {{ teacher?.name || 'Teacher' }}</h2>
        <button @click="$emit('close')" class="text-gray-500 hover:text-gray-700">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Add New Schedule Form -->
      <div class="mb-6 p-4 bg-gray-50 rounded-lg">
        <h3 class="text-lg font-semibold mb-3">Add New Schedule</h3>
        <form @submit.prevent="addSchedule" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
              <input 
                v-model="newSchedule.date" 
                type="date" 
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select 
                v-model="newSchedule.status" 
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="available">Available</option>
                <option value="booked">Booked</option>
                <option value="unavailable">Unavailable</option>
              </select>
            </div>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Start Time</label>
              <input 
                v-model="newSchedule.startTime" 
                type="time" 
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">End Time</label>
              <input 
                v-model="newSchedule.endTime" 
                type="time" 
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div class="flex justify-end">
            <button 
              type="submit"
              :disabled="isAdding"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {{ isAdding ? 'Adding...' : 'Add Schedule' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Existing Schedules -->
      <div>
        <h3 class="text-lg font-semibold mb-3">Existing Schedules</h3>
        
        <!-- Loading State -->
        <div v-if="isLoading" class="text-center py-4">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p class="mt-2 text-gray-600">Loading schedules...</p>
        </div>

        <!-- Schedules List -->
        <div v-else-if="schedules.length > 0" class="space-y-3">
          <div 
            v-for="schedule in schedules" 
            :key="schedule.id"
            class="p-4 border rounded-lg flex justify-between items-center"
            :class="getScheduleClass(schedule.status)"
          >
            <div>
              <div class="font-medium">{{ formatDate(schedule.date) }}</div>
              <div class="text-sm text-gray-600">
                {{ schedule.start_time }} - {{ schedule.end_time }}
              </div>
              <div class="text-sm">
                Status: <span class="font-medium">{{ schedule.status }}</span>
              </div>
            </div>
            <div class="flex space-x-2">
              <button 
                @click="deleteSchedule(schedule.id)"
                class="text-red-600 hover:text-red-800"
                title="Delete Schedule"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-8 text-gray-500">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p class="mt-2">No schedules found for this teacher</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { teacherService } from '@/services/api'
import { useUIStore } from '@/stores/ui'

const props = defineProps({
  teacher: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])
const uiStore = useUIStore()

// Reactive data
const schedules = ref([])
const isLoading = ref(false)
const isAdding = ref(false)

// New schedule form
const newSchedule = ref({
  teacher_id: props.teacher.id,
  date: '',
  startTime: '',
  endTime: '',
  status: 'available'
})

// Methods
const loadSchedules = async () => {
  isLoading.value = true
  try {
    const response = await teacherService.getTeacherSchedules(props.teacher.id)
    schedules.value = response.data || []
  } catch (error) {
    uiStore.showError('Failed to load schedules: ' + error.message)
    console.error('Error loading schedules:', error)
  } finally {
    isLoading.value = false
  }
}

const addSchedule = async () => {
  isAdding.value = true
  try {
    const scheduleData = {
      teacher_id: props.teacher.id,
      date: newSchedule.value.date,
      start_time: newSchedule.value.startTime,
      end_time: newSchedule.value.endTime,
      status: newSchedule.value.status
    }

    await teacherService.createSchedule(scheduleData)
    uiStore.showSuccess('Schedule added successfully')
    
    // Reset form and reload schedules
    newSchedule.value = {
      teacher_id: props.teacher.id,
      date: '',
      startTime: '',
      endTime: '',
      status: 'available'
    }
    
    await loadSchedules()
  } catch (error) {
    uiStore.showError('Failed to add schedule: ' + error.message)
    console.error('Error adding schedule:', error)
  } finally {
    isAdding.value = false
  }
}

const deleteSchedule = async (scheduleId) => {
  try {
    const confirmed = await uiStore.confirm(
      'Are you sure you want to delete this schedule?',
      'Delete Schedule'
    )
    
    if (confirmed) {
      await teacherService.deleteSchedule(scheduleId)
      uiStore.showSuccess('Schedule deleted successfully')
      await loadSchedules()
    }
  } catch (error) {
    uiStore.showError('Failed to delete schedule: ' + error.message)
    console.error('Error deleting schedule:', error)
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const getScheduleClass = (status) => {
  const baseClass = 'border-l-4 '
  switch (status) {
    case 'available':
      return baseClass + 'border-green-500 bg-green-50'
    case 'booked':
      return baseClass + 'border-red-500 bg-red-50'
    case 'unavailable':
      return baseClass + 'border-gray-500 bg-gray-50'
    default:
      return baseClass + 'border-gray-300'
  }
}

// Load schedules on mount
onMounted(() => {
  loadSchedules()
})
</script>
