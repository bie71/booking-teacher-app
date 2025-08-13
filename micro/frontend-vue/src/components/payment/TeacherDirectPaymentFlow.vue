<template>
  <div class="payment-flow-container">
    <!-- Payment Steps Progress (2 steps only) -->
    <div class="mb-6">
      <div class="flex items-center space-x-4">
        <div :class="['step', paymentStep >= 1 ? 'active' : '']">1</div>
        <div class="flex-1 h-1 bg-gray-200 dark:bg-gray-700">
          <div :class="['h-full bg-blue-600 transition-all', paymentStep >= 2 ? 'w-full' : 'w-0']"></div>
        </div>
        <div :class="['step', paymentStep >= 2 ? 'active' : '']">2</div>
      </div>
      <div class="flex justify-between mt-2 text-sm text-gray-600 dark:text-gray-400">
        <span>Booking Details</span>
        <span>Payment</span>
      </div>
    </div>

    <!-- Step 1: Booking Details (since schedule is pre-selected) -->
    <div v-if="paymentStep === 1" class="space-y-4">
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
              <span>{{ Utils.formatDate(selectedSchedule.date) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Time:</span>
              <span>{{ selectedSchedule.start_time }} - {{ selectedSchedule.end_time }}</span>
            </div>
            <div class="flex justify-between font-semibold">
              <span>Total:</span>
              <span>${{ teacher.price_per_hour }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="flex justify-end">
        <button @click="createBooking" class="btn btn-primary">
          Create Booking & Proceed to Payment
        </button>
      </div>
    </div>

    <!-- Step 2: Payment -->
    <div v-if="paymentStep === 2 && createdBooking">
      <PaymentProcessor
        :booking-data="{
          booking_id: createdBooking.id,
          teacher_name: teacher.name,
          date: selectedSchedule.date,
          start_time: selectedSchedule.start_time,
          end_time: selectedSchedule.end_time,
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
  selectedSchedule: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['payment-success', 'payment-failed', 'cancel'])

const authStore = useAuthStore()
const bookingsStore = useBookingsStore()
const uiStore = useUIStore()
const paymentsStore = usePaymentsStore()

const paymentStep = ref(1)
const bookingNotes = ref('')
const createdBooking = ref(null)
const loading = ref(false)

const createBooking = async () => {
  if (!authStore.currentUser) {
    uiStore.showError('Please log in to book a lesson')
    return
  }

  loading.value = true
  try {
    const bookingData = {
      user_id: authStore.currentUser.id,
      schedule_id: props.selectedSchedule.id,
      note: bookingNotes.value || 'Booked via website',
      total_price: props.teacher.price_per_hour
    }

    const result = await bookingsStore.createBooking(bookingData)
    
    if (result.success) {
      createdBooking.value = result.data
      paymentStep.value = 2 // Proceed to payment
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

.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  resize: vertical;
}
</style>
