
<template>
  <div class="min-h-screen bg-white text-gray-900 print:bg-white p-8">
    <div class="max-w-3xl mx-auto">
      <div class="flex items-center justify-between mb-6 no-print">
        <h1 class="text-2xl font-bold">Booking Invoice</h1>
        <div class="space-x-2">
          <button @click="window.print()" class="px-4 py-2 rounded-md bg-indigo-600 text-white">Print</button>
          <router-link to="/admin/bookings" class="px-4 py-2 rounded-md bg-gray-200">Back</router-link>
        </div>
      </div>

      <div class="border rounded-lg p-6 shadow">
        <div class="flex justify-between mb-4">
          <div>
            <div class="font-semibold text-lg">JapanLearn</div>
            <div class="text-sm text-gray-500">Japanese Teacher Booking</div>
          </div>
          <div class="text-right">
            <div class="font-semibold">Invoice #{{ booking?.id }}</div>
            <div class="text-sm text-gray-500">Date: {{ formatDate(booking?.created_at) }}</div>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-6 mb-6">
          <div>
            <div class="text-sm text-gray-500 mb-1">Bill To</div>
            <div class="font-medium">{{ booking?.user?.name }}</div>
            <div class="text-sm text-gray-600">{{ booking?.user?.email }}</div>
          </div>
          <div>
            <div class="text-sm text-gray-500 mb-1">Teacher</div>
            <div class="font-medium">{{ booking?.schedule?.teacher?.name }}</div>
            <div class="text-sm text-gray-600">Schedule ID: {{ booking?.schedule?.id }}</div>
          </div>
        </div>

        <table class="w-full text-sm mb-6">
          <thead>
            <tr class="border-b">
              <th class="text-left py-2">Description</th>
              <th class="text-right py-2">Amount</th>
            </tr>
          </thead>
          <tbody>
            <tr class="border-b">
              <td class="py-2">
                Lesson ({{ booking?.schedule?.start_time }} - {{ booking?.schedule?.end_time }})<br>
                <div class="text-gray-500">
                  Status: 
                  <span :class="statusBadgeClass(booking?.status)">
                    {{ booking?.status }}
                  </span>
                </div>
              </td>
              <td class="py-2 text-right">${{ formatAmount(booking?.total_price) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr>
              <td class="py-2 text-right font-semibold">Total</td>
              <td class="py-2 text-right font-semibold">${{ formatAmount(booking?.total_price) }}</td>
            </tr>
          </tfoot>
        </table>

       <div class="text-sm text-gray-500">
        Payment:
        {{ booking?.payment?.payment_method || '-' }}
        (
          <span :class="paymentBadgeClass(booking?.payment?.status || 'unpaid')">
            {{ (booking?.payment?.status || 'unpaid') }}
          </span>
        )
      </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { bookingService } from '@/services/api'

const route = useRoute()
const booking = ref(null)

function formatAmount(a){ return (a ?? 0).toLocaleString() }
function formatDate(s){ if(!s) return '-'; try{ return new Date(s).toLocaleString() }catch{return s} }

const paymentBadgeClass = (s) => {
  switch ((s || '').toLowerCase()) {
    case 'settlement': return 'bg-emerald-100 text-emerald-800 ring-1 ring-emerald-200'
    case 'failed':
    case 'cancel': return 'bg-rose-100 text-rose-800 ring-1 ring-rose-200'
    case 'pending': return 'bg-amber-100 text-amber-800 ring-1 ring-amber-200'
    default: return 'bg-gray-100 text-gray-800 ring-1 ring-gray-200'
  }
}

const statusBadgeClass = (s) => {
  switch ((s || '').toLowerCase()) {
    case 'paid': return 'bg-emerald-100 text-emerald-800 ring-1 ring-emerald-200'
    case 'completed': return 'bg-blue-100 text-blue-800 ring-1 ring-blue-200'
    case 'pending': return 'bg-amber-100 text-amber-800 ring-1 ring-amber-200'
    case 'cancelled': return 'bg-rose-100 text-rose-800 ring-1 ring-rose-200'
    default: return 'bg-gray-100 text-gray-800 ring-1 ring-gray-200'
  }
}

onMounted(async () => {
  const id = route.params.id
  const res = await bookingService.getBookingDetail(id)
  booking.value = res?.data || res?.booking || res
  // auto-print when opened in new tab
  setTimeout(()=> window.print(), 300)
})
</script>

<style>
@media print {
  .no-print { display: none !important; }
  .print\:bg-white { background-color: white !important; }
}
</style>
