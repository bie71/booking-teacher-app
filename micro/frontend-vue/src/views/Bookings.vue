<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="container py-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          My Bookings
        </h1>
        <p class="text-gray-600 dark:text-gray-400">
          Manage your lesson bookings and schedule
        </p>
      </div>

      <!-- Filters -->
      <div class="card mb-8">
        <div class="card-body">
          <div class="flex flex-col md:flex-row md:items-center md:gap-4 gap-3">
            <select
              v-model="selectedStatus"
              @change="handleStatusFilter"
              class="form-select"
            >
              <option value="">All Status</option>
              <option
                v-for="status in bookingsStore.bookingStatusOptions"
                :key="status.value"
                :value="status.value"
              >
                {{ status.label }}
              </option>
            </select>
            
            <input
              v-model="dateFrom"
              @change="handleDateFilter"
              type="date"
              class="form-input"
              placeholder="From date"
            />
            
            <input
              v-model="dateTo"
              @change="handleDateFilter"
              type="date"
              class="form-input"
              placeholder="To date"
            />
            
            <button
              @click="clearFilters"
              class="btn btn-primary text-sm"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>

      <!-- Bookings List -->
      <div v-if="bookingsStore.isLoading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="card">
          <div class="card-body">
            <div class="skeleton h-6 w-1/4 mb-2"></div>
            <div class="skeleton h-4 w-1/2 mb-2"></div>
            <div class="skeleton h-4 w-1/3"></div>
          </div>
        </div>
      </div>

      <div v-else-if="displayedBookings.length === 0" class="text-center py-12 min-h-[400px] flex flex-col items-center justify-center">
        <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
          No bookings found
        </h3>
        <p class="text-gray-600 dark:text-gray-400 mb-4">
          You haven't booked any lessons yet
        </p>
        <router-link to="/teachers" class="btn btn-primary">
          Book Your First Lesson
        </router-link>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="booking in displayedBookings"
          :key="booking.id"
          class="card hover:shadow-md transition-shadow dark:hover:shadow-lg dark:bg-gray-800 dark:hover:bg-gray-700"
        >
          <div class="card-body">
            <div class="flex flex-col md:flex-row md:items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-4 mb-2">
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                    {{ booking.schedule?.teacher?.name }}
                  </h3>
                  <span :class="getStatusBadgeClass(booking.status)" class="badge">
                    {{ Utils.capitalize(booking.status) }}
                  </span>
                </div>
                
                <div class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                  <div class="flex items-center">
                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    {{ Utils.formatDate(booking.schedule?.date) }}
                  </div>
                  <div class="flex items-center">
                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ booking.schedule?.start_time }} - {{ booking.schedule?.end_time }}
                  </div>
                  <div v-if="booking.note" class="flex items-start">
                    <svg class="w-4 h-4 mr-2 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    {{ booking.note }}
                  </div>
                </div>
              </div>
              
              <div class="flex items-center space-x-2 mt-4 md:mt-0">
                <!-- Reschedule booking icon button -->
                <button
                  v-if="bookingsStore.canRescheduleBooking(booking)"
                  @click="rescheduleBooking(booking)"
                  class="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800"
                  title="Reschedule"
                >
                  <!-- Icon: calendar with arrow (reschedule) -->
                  <svg
                    class="w-5 h-5 text-blue-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 13h-3v-3"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M17 17l-2-2m0 0l-2-2m2 2v3"
                    />
                  </svg>
                </button>
                <!-- Cancel booking icon button -->
                <button
                  v-if="bookingsStore.canCancelBooking(booking)"
                  @click="cancelBooking(booking)"
                  class="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800"
                  title="Cancel"
                >
                  <!-- Icon: X circle -->
                  <svg
                    class="w-5 h-5 text-red-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
                <!-- View details icon button -->
                <button
                  @click="openDetailModal(booking.id)"
                  class="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800"
                  title="View Details"
                >
                  <!-- Icon: eye -->
                  <svg
                    class="w-5 h-5 text-blue-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                </button>
                <!-- Admin-only edit status icon button -->
                <div v-if="authStore.currentUser?.role === 'admin'">
                  <button
                    @click="openEditStatusModal(booking)"
                    class="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800"
                    title="Edit Status"
                  >
                    <!-- Icon: pencil / edit -->
                    <svg
                      class="w-5 h-5 text-yellow-600"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15.232 5.232l3.536 3.536M9 11l3.536 3.536m-3.536-3.536L3 16.071V20h3.929l6.465-6.465m0 0l1.415-1.414a2 2 0 10-2.828-2.828l-1.415 1.414"
                      />
                    </svg>
                  </button>
                </div>
                <!-- Delete booking icon button -->
                <button
                  @click="deleteBooking(booking)"
                  class="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800"
                  title="Delete"
                >
                  <!-- Icon: trash -->
                  <svg
                    class="w-5 h-5 text-red-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5-4h4m-4 0a1 1 0 00-1 1v2h6V4a1 1 0 00-1-1m-4 0h4"
                    />
                  </svg>
                </button>
              </div>

            </div>
          </div>
        </div>

        <!-- Pagination -->
        <!-- Show pagination controls when more than one page exists.  Use
             totalPages > 1 instead of comparing totals to itemsPerPage
             directly, as this covers both client‑side and server‑side
             pagination scenarios. -->
        <div v-if="totalPages > 1" class="mt-8 flex flex-col items-center">
          <div class="flex items-center space-x-2">
            <button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage <= 1"
              class="btn btn-outline disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              Previous
            </button>
            
            <div class="flex items-center space-x-1">
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="goToPage(page)"
                :class="[
                  'btn btn-sm',
                  page === currentPage ? 'btn-primary' : 'btn-outline'
                ]"
              >
                {{ page }}
              </button>
            </div>
            
            <button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage >= totalPages"
              class="btn btn-outline disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
          
          <div class="mt-4 text-sm text-gray-600 dark:text-gray-400">
            Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to 
            {{ Math.min(currentPage * itemsPerPage, totalItems) }} of 
            {{ totalItems }} bookings
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Detail Booking -->
  <BookingDetailModal
    v-if="showDetailModal"
    :booking="selectedBooking"
    @close="selectedBooking = null;  showDetailModal = false"
  />

  <EditStatusModal
    v-if="showEditStatusModal"
    :booking="selectedBooking"
    @close="closeEditStatusModal"
    @updated="handleStatusUpdated"
  />



</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useBookingsStore } from '@/stores/bookings'
import { useUIStore } from '@/stores/ui'
import { handleApiCall, Utils } from '@/utils'
import BookingDetailModal from '@/components/BookingDetailModal.vue'
import EditStatusModal from '@/components/EditStatusModal.vue'

const authStore = useAuthStore()
const bookingsStore = useBookingsStore()
const uiStore = useUIStore()

const selectedStatus = ref('')
const dateFrom = ref('')
const dateTo = ref('')
const selectedBooking = ref(null)
const showDetailModal = ref(false)
const showEditStatusModal = ref(false)

// Client-side pagination state. If the backend provides pagination metadata
// (e.g. when the admin fetches bookings), these values can be updated
// accordingly. However, for user bookings where the backend returns a full
// list, we compute pagination locally.
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Filtered bookings from the store
const allFilteredBookings = computed(() => bookingsStore.filteredBookings)

// Total items based on the filtered list
const totalItems = computed(() => allFilteredBookings.value.length)

// Total pages derived from totalItems and itemsPerPage
const totalPages = computed(() => {
  // If the store has a pagination total greater than zero, use it to
  // compute the number of pages.  Otherwise derive from the local
  // filtered list length.  This allows the component to support both
  // server‑side and client‑side pagination seamlessly.
  if (bookingsStore.pagination?.total && bookingsStore.pagination.total > 0) {
    return Math.ceil(bookingsStore.pagination.total / itemsPerPage.value)
  }
  return totalItems.value > 0 ? Math.ceil(totalItems.value / itemsPerPage.value) : 1
})

// Slice the filtered bookings for the current page
const displayedBookings = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return allFilteredBookings.value.slice(start, end)
})

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


const handleStatusFilter = () => {
  bookingsStore.setStatusFilter(selectedStatus.value)
  resetToFirstPage()
}

const handleDateFilter = () => {
  bookingsStore.setDateRangeFilter(dateFrom.value, dateTo.value)
  resetToFirstPage()
}

const clearFilters = () => {
  selectedStatus.value = ''
  dateFrom.value = ''
  dateTo.value = ''
  bookingsStore.clearFilters()
  resetToFirstPage()
}

const resetToFirstPage = () => {
  if (currentPage.value !== 1) {
    goToPage(1)
  } else {
    // Always fetch bookings from the backend when filters change
    fetchBookings()
  }
}

const goToPage = (page) => {
  // Navigate to a given page. If using backend pagination, update the store
  // page and fetch. Otherwise, simply set the currentPage locally.
  if (page >= 1 && page <= totalPages.value) {
    // If the store pagination total is greater than zero, assume server-side pagination
    if (bookingsStore.pagination?.total && bookingsStore.pagination.total > 0) {
      bookingsStore.pagination.page = page
      fetchBookings()
    } else {
      currentPage.value = page
    }
  }
}

// Pagination logic for visible pages
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  const maxVisible = 5
  const pages = []
  
  if (total <= maxVisible) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    let start = Math.max(1, current - Math.floor(maxVisible / 2))
    let end = Math.min(total, start + maxVisible - 1)
    
    if (end - start < maxVisible - 1) {
      start = Math.max(1, end - maxVisible + 1)
    }
    
    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
  }
  
  return pages
})

const rescheduleBooking = async (booking) => {
  uiStore.showInfo('Reschedule functionality coming soon!')
}

const cancelBooking = async (booking) => {
  const confirmed = await uiStore.confirm(
    'Are you sure you want to cancel this booking? This action cannot be undone.',
    'Cancel Booking'
  )
  
  if (confirmed) {
    const result = await handleApiCall(() => bookingsStore.cancelBooking(booking.id));
    if (result && result.success) {
      uiStore.showSuccess(result.message)
    }
  }
}

const editStatusBooking = async (booking) => {
  const newStatus = await uiStore.prompt(
    'Enter new status for this booking:',
    booking.status
  )
  
  if (newStatus) {
    const result = await handleApiCall(() =>
      bookingsStore.changeStatus(booking.id, newStatus)
    )
    if (result && result.success) {
      uiStore.showSuccess('Booking status updated!')
      fetchBookings()
    }
  }
}

const deleteBooking = async (booking) => {
  const confirmed = await uiStore.confirm(
    'Are you sure you want to delete this booking? This cannot be undone.',
    'Delete Booking'
  )
  
  if (confirmed) {
   try {
     await handleApiCall(() =>
      bookingsStore.deleteBooking(booking.id)
    )
    await fetchBookings()
   }catch (error) {
     uiStore.showError('Failed to delete booking', error)
   }
  }
}



const viewBookingDetails = (booking) => {
  console.log(booking)
  selectedBooking.value = booking
}

const fetchBookings = async () => {
  if (authStore.currentUser) {
    await handleApiCall(() => bookingsStore.fetchBookingsByUser(authStore.currentUser.id));
  } else {
    uiStore.showError('Please login to view your bookings')
  }
}
const openDetailModal = async (id) => {
  try {
    const res = await bookingsStore.bookingDetail(id)
    if (res?.success) {
      selectedBooking.value = res.data
      showDetailModal.value = true
    }
    
  }catch (error) {
    uiStore.showError('Failed to load booking details ', error)
    
  }
}

const openEditStatusModal = (booking) => {
  selectedBooking.value = booking
  showEditStatusModal.value = true
}

const closeEditStatusModal = () => {
  showEditStatusModal.value = false
  selectedBooking.value = null
}

const handleStatusUpdated = () => {
  fetchBookings()
  closeEditStatusModal()
}


onMounted(() => {
  uiStore.setPageTitle('My Bookings')
  fetchBookings()
})
</script>
