<template>
<!-- Payment Flow Modal -->
<div v-if="showPaymentModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center">
  <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-4xl shadow-lg max-h-[90vh] overflow-y-auto">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-bold">Book Lesson with {{ teacher.name }}</h2>
      <button @click="cancelPaymentFlow" class="text-gray-400 hover:text-gray-600">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
    
    <TeacherDirectPaymentFlow
      :teacher="teacher"
      :selected-schedule="selectedSchedule"
      @payment-success="handlePaymentSuccess"
      @payment-failed="handlePaymentFailed"
      @cancel="cancelPaymentFlow"
    />
  </div>
</div>


  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="container py-8">
      <div v-if="teachersStore.isLoading" class="animate-pulse">
        <div class="skeleton h-64 w-full mb-8 rounded-lg"></div>
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div class="lg:col-span-2">
            <div class="skeleton h-8 w-3/4 mb-4"></div>
            <div class="skeleton h-4 w-full mb-2"></div>
            <div class="skeleton h-4 w-2/3 mb-4"></div>
          </div>
          <div>
            <div class="skeleton h-48 w-full rounded-lg"></div>
          </div>
        </div>
      </div>

      <div v-else-if="teacher" class="max-w-6xl mx-auto">
        <!-- Teacher Header -->
        <div class="card mb-8">
          <div class="card-body">
            <div class="flex flex-col md:flex-row items-start space-y-4 md:space-y-0 md:space-x-6">
              <div class="w-32 h-32 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden flex-shrink-0">
                <img
                  v-if="teacher.profile_image"
                  :src="teacher.profile_image"
                  :alt="teacher.name"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-3xl font-bold text-gray-500">
                  {{ teacher.name.charAt(0) }}
                </div>
              </div>
              
              <div class="flex-1">
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
                  {{ teacher.name }}
                </h1>
                <div class="flex items-center space-x-4 mb-4">
                  <span :class="getLevelBadgeClass(teacher.language_level)" class="badge">
                    {{ Utils.capitalize(teacher.language_level) }}
                  </span>
                  <span class="text-sm text-gray-500 dark:text-gray-400">
                    Native Japanese Speaker
                  </span>
                </div>
                <p class="text-gray-600 dark:text-gray-400 mb-4">
                  {{ teacher.bio }}
                </p>
                <div class="flex items-center justify-between">
                  <div class="price">
                    {{ Utils.formatCurrency(teacher.price_per_hour) }}/hour
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Available Schedules -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Available Schedules</h2>
          </div>
          <div class="card-body">
            <div v-if="teachersStore.schedules.length === 0" class="text-center py-8">
              <p class="text-gray-500 dark:text-gray-400">
                No available schedules found for this teacher
              </p>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div 
                v-for="schedule in teachersStore.schedules" 
                :key="schedule.id"
                class="border rounded-lg p-4 hover:shadow-md transition-shadow"
                :class="schedule.status === 'available' ? 'border-green-200 bg-green-50 dark:bg-green-900/20' : 'border-gray-200 opacity-50'"
              >
                <div class="flex justify-between items-start mb-2">
                  <span class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ Utils.formatDate(schedule.date) }}
                  </span>
                  <span 
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      schedule.status === 'available' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                    ]"
                  >
                    {{ schedule.status }}
                  </span>
                </div>
                <div class="text-sm text-gray-600 dark:text-gray-400 mb-2">
                  {{ schedule.start_time }} - {{ schedule.end_time }}
                </div>
                <div class="text-sm font-semibold text-gray-900 dark:text-white">
                  {{ Utils.formatCurrency(schedule.price_per_hour || teacher.price_per_hour) }}
                </div>
              <button 
                v-if="schedule.status === 'available'"
                @click="selectSchedule(schedule)"
                class="mt-3 w-full btn btn-primary btn-sm"
              >
                Select & Book
              </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-12">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">
          Teacher not found
        </h2>
        <router-link to="/teachers" class="btn btn-primary">
          Browse Teachers
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTeachersStore } from '@/stores/teachers'
import { useUIStore } from '@/stores/ui'
import { handleApiCall, Utils } from '@/utils';

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const teachersStore = useTeachersStore()
const uiStore = useUIStore()

const teacher = ref(null)
const loading = ref(true)

const getLevelBadgeClass = (level) => {
  return `level-${level}`
}

import { useBookingsStore } from '@/stores/bookings'
import { usePaymentsStore } from '@/stores/payments'
// import TeacherPaymentFlow from '@/components/payment/TeacherPaymentFlow.vue'
import TeacherDirectPaymentFlow from '@/components/payment/TeacherDirectPaymentFlow.vue'

const bookingsStore = useBookingsStore()
const paymentsStore = usePaymentsStore()

const selectedSchedule = ref(null)
const bookingNote = ref('')
const showPaymentModal = ref(false)
const availableSchedules = ref([])

const selectSchedule = (schedule) => {
  if (schedule.status !== 'available') {
    uiStore.showError('This schedule is not available for booking')
    return
  }
  
  selectedSchedule.value = schedule
  showPaymentModal.value = true
}

const bookLesson = async () => {
  if (!authStore.isAuthenticated) {
    uiStore.showError('Please log in to book a lesson')
    router.push('/login')
    return
  }
  
  // Get available schedules for this teacher
  availableSchedules.value = teachersStore.schedules.filter(
    schedule => schedule.status === 'available'
  )
  
  showPaymentModal.value = true
}

const handlePaymentSuccess = (data) => {
  uiStore.showSuccess('Payment completed successfully!')
  showPaymentModal.value = false
  router.push('/bookings')
}

const handlePaymentFailed = (data) => {
  uiStore.showError('Payment failed. Please try again.')
}

const cancelPaymentFlow = () => {
  showPaymentModal.value = false
  selectedSchedule.value = null
  bookingNote.value = ''
}

const fetchTeacher = async () => {
  const result = await handleApiCall(() => teachersStore.fetchTeacher(route.params.id));
  if (result && result.success) {
    teacher.value = result.data
    uiStore.setPageTitle(`${teacher.value.name} - Teacher Profile`)
  }
}

onMounted(() => {
  fetchTeacher()
  // Fetch teacher schedules
  teachersStore.fetchTeacherSchedules(route.params.id)
})
</script>
