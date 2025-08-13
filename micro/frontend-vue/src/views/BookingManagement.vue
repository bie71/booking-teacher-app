
<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Booking Management</h1>
        <div class="flex gap-2">
          <input type="date" v-model="filters.start_date" @change="fetchBookings" class="form-input" />
          <input type="date" v-model="filters.end_date" @change="fetchBookings" class="form-input" />
          <button @click="exportCSV" class="btn btn-secondary">Export CSV</button>
          <input v-model="filters.q" @keyup.enter="fetchBookings" placeholder="Search user/teacher/payment" class="form-input w-64" />
          <select v-model="filters.status" @change="fetchBookings" class="form-select">
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="paid">Paid</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-700">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">User</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Teacher</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Schedule</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Payment</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider">Action</th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="b in bookings" :key="b.id" class="hover:bg-gray-100 dark:hover:bg-gray-700">
              <td class="px-6 py-4 whitespace-nowrap text-sm">{{ b.id }}</td>
              <td class="px-6 py-4 text-sm">
                <div class="font-medium">{{ b.user?.name || '-' }}</div>
                <div class="text-gray-500 text-xs">{{ b.user?.email }}</div>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="font-medium">{{ b.schedule?.teacher?.name || '-' }}</div>
                <div class="text-gray-500 text-xs">Booking Total: ${{ b.total_price }}</div>
                <div class="text-gray-500 text-xs">Rate: ${{ b.schedule?.teacher?.price || b.schedule?.teacher?.price_per_hour || "-" }}/hour</div>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="font-medium">{{ b.schedule?.start_time }} - {{ b.schedule?.end_time }}</div>
                <div class="text-gray-500 text-xs">ID: {{ b.schedule?.id }}</div>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="font-medium">{{ b.payment?.payment_method || '-' }}</div>
                <div class="text-gray-500 text-xs">{{ b.payment?.status || '-' }}</div>
              </td>
              <td class="px-6 py-4 text-sm">
                <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                  :class="statusBadge(b.status)">{{ b.status }}</span>
              </td>
              <td class="px-6 py-4 text-right text-sm">
                <div class="flex gap-2 justify-end">
                  <button @click="viewDetail(b.id)" class="btn btn-secondary">Detail</button>
                  <button @click="printInvoice(b.id)" class="btn btn-secondary">Print</button>
                  <button v-if="b.status==='pending'" @click="markPaid(b)" class="btn btn-primary">Mark Paid</button>
                  <button v-if="b.status!=='cancelled'" @click="cancel(b)" class="btn btn-danger">Cancel</button>
                  <button v-if="b.status==='paid'" @click="complete(b)" class="btn btn-success">Complete</button>
                </div>
              </td>
            </tr>
            <tr v-if="bookings.length===0">
              <td colspan="7" class="px-6 py-6 text-center text-gray-500">No bookings.</td>
            </tr>
          </tbody>
        </table>
        <Pagination :page="page" :limit="limit" :total="pagination.total" @update:page="(p)=>{ page=p; fetchBookings() }" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Pagination from '@/components/ui/Pagination.vue'
import { bookingService } from '@/services/api'

const router = useRouter()
const bookings = ref([])
const pagination = ref({ page: 1, limit: 10, total: 0, total_pages: 1 })
const page = ref(1)
const limit = ref(10)
const filters = ref({ status: '', q: '', start_date: '', end_date: '' })

function statusBadge(status){
  switch(status){
    case 'paid': return 'bg-green-100 text-green-800'
    case 'pending': return 'bg-yellow-100 text-yellow-800'
    case 'completed': return 'bg-blue-100 text-blue-800'
    case 'cancelled': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

async function fetchBookings(){
  const params = { page: page.value, limit: limit.value, status: filters.value.status, q: filters.value.q, start_date: filters.value.start_date, end_date: filters.value.end_date }
  const res = await bookingService.getAdminBookings(params)
  bookings.value = res.bookings || res.data || []
  pagination.value = res.pagination || { page: page.value, limit: limit.value, total: res.total || 0, total_pages: res.total_pages || 1 }
}

function viewDetail(id){
  router.push({ name: 'BookingDetail', params: { id } })
}
async function markPaid(b){ await bookingService.changeStatus(b.id, 'paid'); fetchBookings() }
async function cancel(b){ await bookingService.changeStatus(b.id, 'cancelled'); fetchBookings() }
async function complete(b){ await bookingService.changeStatus(b.id, 'completed'); fetchBookings() }

function toCSV(rows){
  if(!rows || !rows.length) return ''
  const header = ['ID','User','Email','Teacher','ScheduleID','Start','End','PaymentMethod','PaymentStatus','Status','CreatedAt']
  const esc = (v)=>(''+(v??'')).replace(/"/g,'""')
  const lines = [header.join(',')]
  for(const b of rows){
    lines.push([
      b.id,
      esc(b.user?.name),
      esc(b.user?.email),
      esc(b.schedule?.teacher?.name),
      b.schedule?.id,
      esc(b.schedule?.start_time),
      esc(b.schedule?.end_time),
      esc(b.payment?.payment_method),
      esc(b.payment?.status),
      esc(b.status),
      esc(b.created_at)
    ].join(','))
  }
  return lines.join('\n')
}
function downloadBlob(content, filename, type='text/csv'){
  const blob = new Blob([content], {type})
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a'); a.href=url; a.download=filename; a.click()
  URL.revokeObjectURL(url)
}
function exportCSV(){
  const csv = toCSV(bookings.value)
  const today = new Date().toISOString().slice(0,10)
  downloadBlob(csv, `bookings_${today}.csv`)
}
function printInvoice(id){
  window.open(`/admin/booking/${id}/invoice`, '_blank')
}

onMounted(fetchBookings)
</script>

<style scoped>
.btn{ @apply px-3 py-1.5 rounded-md text-white shadow }
.btn-primary{ @apply bg-blue-600 hover:bg-blue-700 }
.btn-secondary{ @apply bg-gray-600 hover:bg-gray-700 }
.btn-danger{ @apply bg-red-600 hover:bg-red-700 }
.btn-success{ @apply bg-green-600 hover:bg-green-700 }
.form-input{ @apply bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md px-3 py-2 text-sm text-gray-900 dark:text-gray-100 }
.form-select{ @apply bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md px-3 py-2 text-sm text-gray-900 dark:text-gray-100 }
</style>
