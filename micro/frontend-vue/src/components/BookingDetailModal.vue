<template class="bg-white p-4 dark:bg-gray-800">
  <Modal :show="true" @close="$emit('close')" >
    <template #header >
      <h2 class="text-xl font-bold">Booking Detail</h2>
    </template>

    <template #body>
      <table class="table-auto w-full text-sm border border-gray-200 dark:border-gray-700 dark:bg-gray-800">
        <tbody>

          <!-- Status -->
          <tr>
            <td class="font-medium px-4 py-2 w-40 border-b">Status</td>
            <td class="px-4 py-2 border-b">
              <span :class="`inline-block text-xs font-semibold px-2 py-1 rounded ${getStatusBadgeClass(booking.status)}`">
                {{ booking.status }}
              </span>
            </td>
          </tr>

          <!-- User -->
          <tr v-if="booking.user">
            <td class="font-medium px-4 py-2 border-b">User</td>
            <td class="px-4 py-2 border-b flex items-center gap-3">
              <img :src="booking.user.profile_image" class="w-10 h-10 rounded-xl object-cover" />
              <div>
                <p class="font-medium">{{ booking.user.name }}</p>
                <p class="text-gray-500 text-xs dark:text-gray-400">{{ booking.user.email }}</p>
              </div>
            </td>
          </tr>

          <!-- Schedule -->
          <tr v-if="booking.schedule">
            <td class="font-medium px-4 py-2 border-b dark:text-white">Schedule</td>
            <td class="px-4 py-2 border-b dark:text-white" >
              <p><strong>Date:</strong> {{ formatDateOnly(booking.schedule.date) }} </p>
              <p><strong>Time:</strong> {{ booking.schedule.start_time }} - {{ booking.schedule.end_time }}</p>
              <p><strong>Total:</strong> ${{ booking.total_price }}</p>
            </td>
          </tr>

          <!-- Teacher -->
          <tr v-if="booking.schedule && booking.schedule.teacher">
            <td class="font-medium px-4 py-2 border-b">Teacher</td>
            <td class="px-4 py-2 border-b">
              <p class="font-semibold">{{ booking.schedule.teacher.name }}</p>
              <p class="text-sm text-gray-500 italic dark:text-gray-400">{{ booking.schedule.teacher.bio }}</p>
            </td>
          </tr>

          <!-- Payment -->
          <tr v-if="booking.payment">
            <td class="font-medium px-4 py-2 border-b">Payment</td>
            <td class="px-4 py-2 border-b space-y-1">
              <p><strong>Transaction ID:</strong> {{ booking.payment.midtrans_transaction_id || '-' }}</p>
              <p><strong>Status:</strong> {{ booking.payment.status || '-' }}</p>
              <p><strong>Method:</strong> {{ booking.payment.payment_method || '-' }}</p>
              <p><strong>Amount:</strong> ${{ booking.payment.amount }}</p>
            </td>
          </tr>

          <!-- Notes -->
          <tr v-if="booking.note">
            <td class="font-medium px-4 py-2 border-b">Note</td>
            <td class="px-4 py-2 border-b text-gray-700 dark:text-gray-300">{{ booking.note }}</td>
          </tr>

          <!-- Created At -->
          <tr>
            <td class="font-medium px-4 py-2">Booked At</td>
            <td class="px-4 py-2 text-sm text-gray-500 dark:text-gray-400">{{ formatDate(booking.created_at) }}</td>
          </tr>

        </tbody>
      </table>
    </template>

    <template #footer>
      <button class="btn btn-primary dark:text-white" @click="$emit('close')">Close</button>
    </template>
  </Modal>
</template>

<script setup>
import Modal from '@/components/Modal.vue'

const props = defineProps({
  booking: Object
})

const emit = defineEmits(['close'])

const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'pending':
      return 'bg-yellow-100 text-yellow-800 border border-yellow-300'
    case 'paid':
      return 'bg-blue-100 text-blue-800 border border-blue-300'
    case 'cancelled':
      return 'bg-red-100 text-red-800 border border-red-300'
    case 'rescheduled':
      return 'bg-purple-100 text-purple-800 border border-purple-300'
    case 'completed':
      return 'bg-green-100 text-green-800 border border-green-300'
    default:
      return 'bg-gray-100 text-gray-800 border border-gray-300'
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const formatDateOnly = (date)  => {
  return new Date(date).toLocaleDateString('en-GB', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}


</script>
