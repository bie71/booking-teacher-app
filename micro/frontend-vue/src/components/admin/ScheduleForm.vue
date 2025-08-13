<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900">
        {{ mode === 'create' ? 'Add New Schedule' : 'Edit Schedule' }}
      </h2>
      <button 
        @click="$emit('cancel')"
        class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
      >
        Cancel
      </button>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Teacher Selection -->
      <div v-if="mode === 'create'">
        <label for="teacher" class="block text-sm font-medium text-gray-700">Teacher *</label>
        <select 
          v-model="form.teacher_id"
          id="teacher"
          required
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">Select Teacher</option>
          <option v-for="teacher in teachers" :key="teacher.id" :value="teacher.id">
            {{ teacher.name }} - {{ teacher.language_level }}
          </option>
        </select>
      </div>

      <!-- Schedule Details -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label for="day" class="block text-sm font-medium text-gray-700">Day *</label>
          <select 
            v-model="form.day"
            id="day"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Select Day</option>
            <option value="monday">Monday</option>
            <option value="tuesday">Tuesday</option>
            <option value="wednesday">Wednesday</option>
            <option value="thursday">Thursday</option>
            <option value="friday">Friday</option>
            <option value="saturday">Saturday</option>
            <option value="sunday">Sunday</option>
          </select>
        </div>

        <div>
          <label for="start_time" class="block text-sm font-medium text-gray-700">Start Time *</label>
          <input 
            v-model="form.start_time"
            type="time"
            id="start_time"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <div>
          <label for="end_time" class="block text-sm font-medium text-gray-700">End Time *</label>
          <input 
            v-model="form.end_time"
            type="time"
            id="end_time"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <div>
          <label for="max_students" class="block text-sm font-medium text-gray-700">Max Students *</label>
          <input 
            v-model.number="form.max_students"
            type="number"
            id="max_students"
            min="1"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
      </div>

      <!-- Status -->
      <div>
        <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
        <select 
          v-model="form.status"
          id="status"
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="available">Available</option>
          <option value="booked">Booked</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      <!-- Submit Button -->
      <div class="flex justify-end space-x-3">
        <button 
          type="button"
          @click="$emit('cancel')"
          class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
        >
          Cancel
        </button>
        <button 
          type="submit"
          :disabled="isSubmitting"
          class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
        >
          {{ isSubmitting ? 'Saving...' : 'Save Schedule' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTeachersStore } from '@/stores/teachers'
import { useUIStore } from '@/stores/ui'

const props = defineProps({
  teacher: {
    type: Object,
    default: null
  },
  schedule: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

const teachersStore = useTeachersStore()
const uiStore = useUIStore()

const isSubmitting = ref(false)
const mode = computed(() => props.schedule ? 'edit' : 'create')

const form = ref({
  teacher_id: '',
  day: '',
  start_time: '',
  end_time: '',
  max_students: 1,
  status: 'available'
})

const teachers = computed(() => teachersStore.teachers)

// Initialize form
const initializeForm = () => {
  if (props.schedule) {
    form.value = {
      teacher_id: props.schedule.teacher_id || '',
      day: props.schedule.day || '',
      start_time: props.schedule.start_time || '',
      end_time: props.schedule.end_time || '',
      max_students: props.schedule.max_students || 1,
      status: props.schedule.status || 'available'
    }
  } else if (props.teacher) {
    form.value.teacher_id = props.teacher.id
  }
}

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    
    const scheduleData = {
      ...form.value,
      max_students: Number(form.value.max_students)
    }

    if (mode.value === 'create') {
      await teachersStore.createSchedule(scheduleData)
      uiStore.showSuccess('Schedule created successfully!')
    } else {
      await teachersStore.updateSchedule(props.schedule.id, scheduleData)
      uiStore.showSuccess('Schedule updated successfully!')
    }
    
    emit('submit')
  } catch (error) {
    uiStore.showError(error.message || 'Failed to save schedule')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  initializeForm()
})
</script>
