<template>
  <div class="payment-processor">
    <div class="card">
      <div class="card-body">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
          Complete Your Payment
        </h2>
        
        <!-- Booking Summary -->
        <div v-if="bookingData" class="mb-6 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <h3 class="font-semibold text-gray-900 dark:text-white mb-3">
            Booking Summary
          </h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Teacher:</span>
              <span class="font-medium">{{ bookingData.teacher_name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Date:</span>
              <span class="font-medium">{{ Utils.formatDate(bookingData.date) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Time:</span>
              <span class="font-medium">{{ bookingData.start_time }} - {{ bookingData.end_time }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Duration:</span>
              <span class="font-medium">{{ bookingData.duration || 60 }} minutes</span>
            </div>
            <hr class="my-3 border-gray-200 dark:border-gray-700">
            <div class="flex justify-between text-lg font-semibold">
              <span>Total Amount:</span>
              <span class="text-blue-600">{{ paymentsStore.formatPaymentAmount(bookingData.amount) }}</span>
            </div>
          </div>
        </div>
        
        <!-- Payment Method Selection -->
        <div class="mb-6">
          <PaymentMethodSelector
            v-model="selectedPaymentMethod"
            @method-selected="onPaymentMethodSelected"
          />
        </div>
        
        <!-- Payment Actions -->
        <div class="flex flex-col sm:flex-row gap-4">
          <button
            @click="$emit('cancel')"
            :disabled="paymentsStore.isProcessingPayment"
            class="btn btn-outline flex-1"
          >
            Cancel
          </button>
          <button
            @click="processPayment"
            :disabled="!selectedPaymentMethod || paymentsStore.isProcessingPayment"
            class="btn btn-primary flex-1"
          >
            <svg v-if="paymentsStore.isProcessingPayment" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ paymentsStore.isProcessingPayment ? 'Processing...' : 'Pay Now' }}
          </button>
        </div>
        
        <!-- Payment Security Info -->
        <div class="mt-6 p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
          <div class="flex items-start space-x-3">
            <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <div>
              <h5 class="font-medium text-green-900 dark:text-green-100">
                Secure Payment
              </h5>
              <p class="text-sm text-green-700 dark:text-green-300 mt-1">
                Your payment information is encrypted and secure. We use industry-standard security measures to protect your data.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { usePaymentsStore } from '@/stores/payments'
import { useUIStore } from '@/stores/ui'
import { Utils } from '@/utils'
import PaymentMethodSelector from './PaymentMethodSelector.vue'

const props = defineProps({
  bookingData: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['payment-success', 'payment-failed', 'cancel'])

const paymentsStore = usePaymentsStore()
const uiStore = useUIStore()

const selectedPaymentMethod = ref('')
const selectedMethod = ref(null)

const onPaymentMethodSelected = (method) => {
  selectedMethod.value = method
}

const processPayment = async () => {
  if (!selectedPaymentMethod.value) {
    uiStore.showError('Please select a payment method')
    return
  }

  if (!props.bookingData.booking_id) {
    uiStore.showError('Invalid booking data')
    return
  }

  try {
    const paymentData = {
      booking_id: props.bookingData.booking_id,
      amount: props.bookingData.amount,
      payment_method: selectedMethod.value?.code || selectedPaymentMethod.value
    }

    const result = await paymentsStore.createPayment(paymentData)
    
    if (result.success) {
      // Check if we have a snap_url for redirect
      if (result.data.snap_url) {
        // Redirect to payment gateway
        window.open(result.data.snap_url, '_blank')
        
        // Show success message and emit event
        uiStore.showSuccess('Payment initiated successfully. Please complete the payment in the new window.')
        emit('payment-success', {
          payment: result.data,
          booking_id: props.bookingData.booking_id
        })
      } else {
        // Handle other payment types
        emit('payment-success', {
          payment: result.data,
          booking_id: props.bookingData.booking_id
        })
      }
    } else {
      console.error('Payment processing error:', result.message)
      emit('payment-failed', {
        error: result.message,
        booking_id: props.bookingData.booking_id
      })
    }
  } catch (error) {
    console.error('Payment processing error:', error)
    emit('payment-failed', {
      error: error.message,
      booking_id: props.bookingData.booking_id
    })
  }
}
</script>

<style scoped>
.card {
  @apply bg-white dark:bg-gray-800 rounded-lg shadow-lg;
}

.card-body {
  @apply p-6;
}

.btn {
  @apply px-4 py-2 rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-outline {
  @apply border border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed;
  @apply dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700;
}
</style>
