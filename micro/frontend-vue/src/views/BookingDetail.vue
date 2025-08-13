<template>
  <div class="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-950">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <!-- Header -->
      <div class="flex items-start justify-between mb-8">
        <div>
          <h1 class="text-3xl font-semibold tracking-tight text-gray-900 dark:text-white">Booking Detail</h1>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">View payment, schedule, and participant information</p>
        </div>
        <div class="flex items-center gap-3">
          <router-link :to="{ name: 'BookingInvoice', params: { id: route.params.id } }"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 shadow-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M14 2H6a2 2 0 0 0-2 2v16l4-2 4 2 4-2 4 2V8l-6-6z"/></svg>
            Invoice
          </router-link>
          <router-link to="/admin/bookings"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-gray-200 text-gray-900 hover:bg-gray-300 dark:bg-gray-800 dark:text-gray-100 dark:hover:bg-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M15.41 7.41 14 6l-6 6 6 6 1.41-1.41L10.83 12z"/></svg>
            Back
          </router-link>
        </div>
      </div>

      <!-- Top summary card -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="md:col-span-2 bg-white dark:bg-gray-900 rounded-2xl shadow p-6 border border-gray-100 dark:border-gray-800">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-full overflow-hidden ring-1 ring-gray-200 dark:ring-gray-700">
                <img :src="userAvatar" alt="user avatar" class="h-10 w-10 object-cover" />
              </div>
              <div>
                <div class="text-xs uppercase tracking-wider text-gray-500">User</div>
                <div class="font-medium text-gray-900 dark:text-gray-100">{{ booking?.user?.name }}</div>
                <div class="text-sm text-gray-500">{{ booking?.user?.email }}</div>
              </div>
            </div>
            <span :class="['inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium', statusBadgeClass(booking?.status)]">
              {{ (booking?.status || '-').toUpperCase() }}
            </span>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-2">
            
<div class="flex items-start justify-between gap-6 py-2 border-b last:border-0 border-gray-100 dark:border-gray-800">
  <div class="text-sm text-gray-500 dark:text-gray-400">Booking ID</div>
  <div class="text-sm font-medium text-gray-900 dark:text-gray-100 text-right break-words">{{ booking?.id ?? '-' }}</div>
</div>
            
<div class="flex items-start justify-between gap-6 py-2 border-b last:border-0 border-gray-100 dark:border-gray-800">
  <div class="text-sm text-gray-500 dark:text-gray-400">Created At</div>
  <div class="text-sm font-medium text-gray-900 dark:text-gray-100 text-right break-words">{{ formatDate(booking?.created_at) ?? '-' }}</div>
</div>
            
<div class="flex items-start justify-between gap-6 py-2 border-b last:border-0 border-gray-100 dark:border-gray-800">
  <div class="text-sm text-gray-500 dark:text-gray-400">Total Price</div>
  <div class="text-sm font-medium text-gray-900 dark:text-gray-100 text-right break-words">{{ formatCurrency(booking?.total_price) ?? '-' }}</div>
</div>
            
<div class="flex items-start justify-between gap-6 py-2 border-b last:border-0 border-gray-100 dark:border-gray-800">
  <div class="text-sm text-gray-500 dark:text-gray-400">Note</div>
  <div class="text-sm font-medium text-gray-900 dark:text-gray-100 text-right break-words">{{ booking?.note ?? '-' }}</div>
</div>
          </div>
        </div>

        <div class="bg-white dark:bg-gray-900 rounded-2xl shadow p-6 border border-gray-100 dark:border-gray-800">
          <div class="flex items-center gap-3 mb-4">
            <div class="h-10 w-10 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2a10 10 0 1 0 .001 20.001A10 10 0 0 0 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>
            </div>
            <div>
              <div class="text-xs uppercase tracking-wider text-gray-500">Payment</div>
              <div class="font-medium text-gray-900 dark:text-gray-100">{{ booking?.payment?.payment_method || '-' }}</div>
              <div class="flex items-center gap-2 text-sm text-gray-500">
                Status:
              <span
                :class="['inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium', paymentBadgeClass(booking?.payment?.status)]"
              >
                {{ (booking?.payment?.status || '-').toUpperCase() }}
              </span>
            </div>
            </div>
          </div>
          <div class="text-xs text-gray-500">* Payment details synced from payment service.</div>
        </div>
      </div>

      <!-- Details grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Schedule & Teacher -->
        <div class="lg:col-span-2 bg-white dark:bg-gray-900 rounded-2xl shadow p-6 border border-gray-100 dark:border-gray-800">
          <div class="flex items-center gap-3 mb-4">
            <div class="h-10 w-10 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor"><path d="M19 3H5a2 2 0 0 0-2 2v14l4-4h12a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2z"/></svg>
            </div>
            <div>
              <div class="text-xs uppercase tracking-wider text-gray-500">Schedule</div>
              <div class="font-medium text-gray-900 dark:text-gray-100">#{{ booking?.schedule?.id }}</div>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4">
            <div class="flex items-center gap-4 p-3 rounded-xl bg-gray-50 dark:bg-gray-800/60">
              <div class="h-12 w-12 rounded-full overflow-hidden ring-1 ring-gray-200 dark:ring-gray-700">
                <img :src="teacherAvatar" alt="teacher avatar" class="h-12 w-12 object-cover" />
              </div>
              <div>
                <div class="font-medium text-gray-900 dark:text-gray-100">{{ booking?.schedule?.teacher?.name }}</div>
                <div class="text-xs text-gray-500">Rate: {{ booking?.schedule?.teacher?.price_per_hour ? formatCurrency(booking?.schedule?.teacher?.price_per_hour) + ' / hour' : '-' }}</div>
              </div>
            </div>
            
<div class="flex items-start justify-between gap-6 py-2 border-b last:border-0 border-gray-100 dark:border-gray-800">
  <div class="text-sm text-gray-500 dark:text-gray-400">Time</div>
  <div class="text-sm font-medium text-gray-900 dark:text-gray-100 text-right break-words">{{ `${booking?.schedule?.start_time || '-'} - ${booking?.schedule?.end_time || '-'}` ?? '-' }}</div>
</div>
            
<div class="flex items-start justify-between gap-6 py-2 border-b last:border-0 border-gray-100 dark:border-gray-800">
  <div class="text-sm text-gray-500 dark:text-gray-400">Teacher Rate</div>
  <div class="text-sm font-medium text-gray-900 dark:text-gray-100 text-right break-words">{{ booking?.schedule?.teacher?.price_per_hour ? formatCurrency(booking?.schedule?.teacher?.price_per_hour) + ' / hour' : '-' ?? '-' }}</div>
</div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="bg-white dark:bg-gray-900 rounded-2xl shadow p-6 border border-gray-100 dark:border-gray-800">
          <div class="text-sm font-semibold text-gray-900 dark:text-gray-100 mb-4">Actions</div>
          <div class="grid grid-cols-1 gap-3">
            <button @click="changeStatus('paid')" class="w-full px-4 py-2 rounded-lg bg-emerald-600 text-white hover:bg-emerald-700">Mark as Paid</button>
            <button @click="changeStatus('completed')" class="w-full px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700">Mark as Completed</button>
            <button @click="changeStatus('cancelled')" class="w-full px-4 py-2 rounded-lg bg-rose-600 text-white hover:bg-rose-700">Cancel Booking</button>
          </div>
          <p class="text-xs text-gray-500 mt-3">* Actions affect booking status immediately.</p>
        </div>
      </div>
    
      <!-- Timeline -->
      <div class="mt-8 bg-white dark:bg-gray-900 rounded-2xl shadow p-6 border border-gray-100 dark:border-gray-800">
        <div class="flex items-center gap-3 mb-4">
          <div class="h-10 w-10 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor"><path d="M19 3h-1V1h-2v2H8V1H6v2H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2zM5 19V9h14v10H5z"/></svg>
          </div>
          <div>
            <div class="text-xs uppercase tracking-wider text-gray-500">Timeline</div>
            <div class="font-medium text-gray-900 dark:text-gray-100">Booking Progress</div>
          </div>
        </div>
        
<ol class="relative border-s border-gray-200 dark:border-gray-700 ml-3">
  <li v-for="(it, idx) in timelineItems" :key="idx" class="mb-6 ms-4">
    <span :class="[
      'absolute -start-2.5 flex h-5 w-5 items-center justify-center rounded-full ring-2 ring-white dark:ring-gray-900',
      it.state === 'done' ? 'bg-emerald-500' : it.state === 'error' ? 'bg-rose-500' : it.state === 'progress' ? 'bg-indigo-500' : 'bg-gray-400'
    ]">
      <svg class="h-3 w-3 text-white" viewBox="0 0 20 20" fill="currentColor">
        <path v-if="it.state==='done'" fill-rule="evenodd" d="M16.707 5.293a1 1 0 00-1.414 0L8 12.586 4.707 9.293a1 1 0 10-1.414 1.414l4 4a1 1 0 001.414 0l8-8a1 1 0 000-1.414z" clip-rule="evenodd"/>
        <circle v-else cx="10" cy="10" r="3" />
      </svg>
    </span>
    <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{ it.title }}</h3>
    <p class="text-sm text-gray-500 dark:text-gray-400">{{ it.desc }}</p>
  </li>
</ol>

      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { bookingService } from '@/services/api'

const route = useRoute()
const booking = ref(null)
const loading = ref(false)

function formatCurrency(n){ 
  const v = Number(n || 0)
  return '$' + v.toLocaleString() 
}
function formatDate(s){ if(!s) return '-'; try{ return new Date(s).toLocaleString() }catch{return s} }

async function fetchDetail(){
  loading.value = true
  try{
    const id = route.params.id
    const res = await bookingService.getBookingDetail(id)
    booking.value = res?.data || res?.booking || res
  } finally {
    loading.value = false
  }
}

async function changeStatus(status){
  if(!booking.value) return
  await bookingService.changeStatus(booking.value.id, status)
  await fetchDetail()
}

onMounted(fetchDetail)


const teacherInitials = (name)=>{
  if(!name) return '?';
  return name.split(' ').slice(0,2).map(s=>s[0]).join('').toUpperCase();
}

const timelineItems = computed(()=>{
  const items = []
  if(booking.value){
    items.push({ 
      title: 'Created', 
      desc: formatDate(booking.value.created_at), 
      state: 'done' 
    })
    // Payment
    const p = booking.value.payment
    if(p && p.status){
      items.push({
        title: 'Payment',
        desc: `${(p.payment_method||'-').toUpperCase()} • ${p.status.toUpperCase()}`,
        state: p.status === 'settlement' ? 'done' : (p.status === 'failed' || p.status === 'cancel') ? 'error' : 'progress'
      })
    } else {
      items.push({ title: 'Payment', desc: 'No payment recorded', state: 'idle' })
    }
    // Schedule
    items.push({
      title: 'Schedule',
      desc: `${booking.value.schedule?.start_time || '-'} → ${booking.value.schedule?.end_time || '-'}`,
      state: 'progress'
    })
    // Completion / Cancellation
    const st = (booking.value.status||'').toLowerCase()
    if(st === 'completed'){
      items.push({ title: 'Completed', desc: 'Lesson completed', state: 'done' })
    } else if(st === 'cancelled'){
      items.push({ title: 'Cancelled', desc: 'Booking cancelled', state: 'error' })
    } else {
      items.push({ title: 'Status', desc: (booking.value.status||'-').toUpperCase(), state: 'progress' })
    }
  }
  return items
})

const pickAvatar = (o)=> o?.avatar || o?.photo_url || o?.profile_url || o?.image_url || o?.picture || o?.profile_image || ''
const userAvatar = computed(()=> pickAvatar(booking.value?.user) || 'https://placehold.co/80x80?text=U')
const teacherAvatar = computed(()=> pickAvatar(booking.value?.schedule?.teacher) || 'https://placehold.co/96x96?text=T')

// BADGE HELPERS
const statusBadgeClass = (s) => {
  switch ((s || '').toLowerCase()) {
    case 'paid': return 'bg-emerald-100 text-emerald-800 ring-1 ring-emerald-200'
    case 'completed': return 'bg-blue-100 text-blue-800 ring-1 ring-blue-200'
    case 'pending': return 'bg-amber-100 text-amber-800 ring-1 ring-amber-200'
    case 'cancelled': return 'bg-rose-100 text-rose-800 ring-1 ring-rose-200'
    default: return 'bg-gray-100 text-gray-800 ring-1 ring-gray-200'
  }
}
const paymentBadgeClass = (s) => {
  switch ((s || '').toLowerCase()) {
    case 'settlement': return 'bg-emerald-100 text-emerald-800 ring-1 ring-emerald-200'
    case 'failed':
    case 'cancel': return 'bg-rose-100 text-rose-800 ring-1 ring-rose-200'
    case 'pending': return 'bg-amber-100 text-amber-800 ring-1 ring-amber-200'
    default: return 'bg-gray-100 text-gray-800 ring-1 ring-gray-200'
  }
}




</script>


