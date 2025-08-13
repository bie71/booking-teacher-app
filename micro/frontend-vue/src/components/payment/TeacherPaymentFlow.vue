<template>
  <div class="payment-flow-container">
    <!-- Payment Steps Progress -->
    <div class="mb-6">
      <div class="flex items-center space-x-4">
        <div :class="['step', paymentStep >= 1 ? 'active' : '']">1</div>
        <div class="flex-1 h-1 bg-gray-200 dark:bg-gray-700">
          <div :class="['h-full bg-blue-600 transition-all', paymentStep >= 2 ? 'w-full' : 'w-0']"></div>
        </div>
        <div :class="['step', paymentStep >= 2 ? 'active' : '']">2</div>
        <div class="flex-1 h-1 bg-gray-200 dark:bg-gray-700">
          <div :class="['h-full bg-blue-600 transition-all', paymentStep >= 3 ? 'w-full' : 'w-0']"></div>
        </div>
        <div :class="['step', paymentStep >= 3 ? 'active' : '']">3</div>
      </div>
      <div class="flex justify-between mt-2 text-sm text-gray-600 dark:text-gray-400">
        <span>Select Schedule</span>
        <span>Booking Details</span>
        <span>Payment</span>
      </div>
    </div>

    <!-- Step 1: Schedule Selection -->
    <div v-if="paymentStep === 1" class="space-y-4">
      <h3 class="text-lg font-semibold">Select Available Schedule</h3>
      <div v-if="availableSchedules.length === 0" class="text-center py-8">
        <p class="text-gray-600 dark:text-gray-400">No available schedules found</p>
      </div>
      <div v-else class="grid gap-3">
        <div
          v-for="schedule in availableSchedules"
          :key="schedule.id"
          @click="selectSchedule(schedule)"
          :class="[
            'p-4 border rounded-lg cursor-pointer transition-all',
            selectedSchedule?.id === schedule.id
              ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
              : 'border-gray-200 dark:border-gray-700 hover:border-blue-300'
          ]"
        >
          <div class="flex justify-between items-center">
            <div>
              <div class="font-medium">{{ Utils.formatDate(schedule.date) }}</div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                {{ schedule.start_time }} - {{ schedule.end_time }}
              </div>
            </div>
            <div class="text-lg font-semibold text-blue-600">
              ${{ teacher.price_per_hour }}
            </div>
          </div>
        </div>
      </div>
      <div class="flex justify-end">
        <button
          @click="nextStep"
          :disabled="!selectedSchedule"
          class="btn btn-primary"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Step 2: Booking Details -->
    <div v-if="paymentStep === 2" class="space-y-4">
      <h3 class="text-lg font-semibold">Booking Details</h3>
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Special Notes (Optional)
          </label>
          <textarea
            v-model="bookingNotes"
            rows="3"
            class="form-textarea"
            placeholder="Any special requests or notes for the teacher..."
          ></textarea>
        </div>
        
        <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg">
          <h4 class="font-medium mb-2">Booking Summary</h4>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span>Teacher:</span>
              <span>{{ teacher.name }}</span>
            </div>
            <div class="flex justify-between">
              <span>Date:</span>
              <span>{{ Utils.formatDate(selectedSchedule?.date) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Time:</span>
              <span>{{ selectedSchedule?.start_time }} - {{ selectedSchedule?.end_time }}</span>
            </div>
            <div class="flex justify-between font-semibold">
              <span>Total:</span>
              <span>${{ teacher.price_per_hour }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="flex justify-between">
        <button @click="prevStep" class="btn btn-outline">
          Back
        </button>
        <button @click="createBooking" class="btn btn-primary">
          Create Booking
        </button>
      </div>
    </div>

    <!-- Step 3: Payment -->
    <div v-if="paymentStep === 3 && createdBooking">
      <PaymentProcessor
        :booking-data="{
          booking_id: createdBooking.id,
          teacher_name: teacher.name,
          date: selectedSchedule?.date,
          start_time: selectedSchedule?.start_time,
          end_time: selectedSchedule?.end_time,
          amount: teacher.price_per_hour * 100,
          duration: 60
        }"
        @payment-success="handlePaymentSuccess"
        @payment-failed="handlePaymentFailed"
        @cancel="handleCancel"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useBookingsStore } from '@/stores/bookings'
import { useUIStore } from '@/stores/ui'
import { usePaymentsStore } from '@/stores/payments'
import { Utils } from '@/utils'
import PaymentProcessor from '@/components/PaymentProcessor.vue'

const props = defineProps({
  teacher: {
    type: Object,
    required: true
  },
  availableSchedules: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['payment-success', 'payment-failed', 'cancel'])

const authStore = useAuthStore()
const bookingsStore = useBookingsStore()
const uiStore = useUIStore()
const paymentsStore = usePaymentsStore()

const paymentStep = ref(1)
const selectedSchedule = ref(null)
const bookingNotes = ref('')
const createdBooking = ref(null)
const loading = ref(false)

const selectSchedule = (schedule) => {
  selectedSchedule.value = schedule
}

const nextStep = () => {
  if (paymentStep.value < 3) {
    paymentStep.value++
  }
}

const prevStep = () => {
  if (paymentStep.value > 1) {
    paymentStep.value--
  }
}

const createBooking = async () => {
  if (!selectedSchedule.value || !authStore.currentUser) {
    uiStore.showError('Missing required booking information')
    return
  }

  loading.value = true
  try {
    const bookingData = {
      user_id: authStore.currentUser.id,
      schedule_id: selectedSchedule.value.id,
      note: bookingNotes.value || 'Booked via website',
      total_price: props.teacher.price_per_hour
    }

    const result = await bookingsStore.createBooking(bookingData)
    
    if (result.success) {
      createdBooking.value = result.data
      nextStep()
    } else {
      uiStore.showError(result.message || 'Failed to create booking')
    }
  } catch (error) {
    uiStore.showError('Error creating booking: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handlePaymentSuccess = (data) => {
  uiStore.showSuccess('Payment completed successfully!')
  emit('payment-success', data)
}

const handlePaymentFailed = (data) => {
  uiStore.showError('Payment failed. Please try again.')
  emit('payment-failed', data)
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.payment-flow-container {
  width: 100%;
}

.step {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 500;
  background-color: #e5e7eb;
  color: #6b7280;
}

.step.active {
  background-color: #3b82f6;
  color: white;
}

.btn {
  padding: 8px 16px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn-outline {
  border: 1px solid #d1d5db;
  color: #374151;
}

.btn-outline:hover {
  background-color: #f9fafb;
}

.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  resize: vertical;
}
</style>
