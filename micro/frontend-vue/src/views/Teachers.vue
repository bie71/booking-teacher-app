<template>
  <div v-if="!storesInitialized" class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
    <div class="text-center">
      <h1 class="text-2xl font-bold text-gray-900 mb-4">Loading Teachers...</h1>
      <p class="text-gray-600">Please wait while we initialize the application.</p>
    </div>
  </div>
  
  <div v-else class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Header Section -->
    <section class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div class="container py-8">
        <div class="text-center mb-8">
          <h1 class="text-3xl md:text-4xl font-bold mb-4">
            Find Your Perfect Japanese Teacher
          </h1>
          <p class="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Browse our selection of qualified native Japanese teachers and book your personalized lesson
          </p>
        </div>

        <!-- Search and Filters -->
        <div class="max-w-4xl mx-auto dark:text-white dark:bg-gray-800">
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <!-- Search -->
            <div class="md:col-span-2">
              <div class="relative">
                <input
                  v-model="searchQuery"
                  @input="handleSearch"
                  type="text"
                  placeholder="Search teachers by name or bio..."
                  class="form-input pl-10"
                />
                <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>

            <!-- Level Filter -->
            <div>
              <select
                v-model="selectedLevel"
                @change="handleLevelFilter"
                class="form-select"
              >
                <option value="">All Levels</option>
                <option
                  v-for="level in (teachersStore?.teacherLevels || [])"
                  :key="level.value"
                  :value="level.value"
                >
                  {{ level.label }}
                </option>
              </select>
            </div>

            <!-- Sort -->
            <div>
              <select
                v-model="sortBy"
                @change="handleSort"
                class="form-select"
              >
                <option value="name">Sort by Name</option>
                <option value="price_low">Price: Low to High</option>
                <option value="price_high">Price: High to Low</option>
                <option value="level">Level</option>
              </select>
            </div>
          </div>

          <!-- Price Range Filter -->
          <div class="flex items-center space-x-4 mb-6">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Price Range:</span>
            <div class="flex items-center space-x-2">
              <input
                v-model.number="priceRange.min"
                @change="handlePriceFilter"
                type="number"
                placeholder="Min"
                class="form-input w-20"
                :min="teachersStore?.priceRange?.min || 0"
                :max="teachersStore?.priceRange?.max || 100"
              />
              <span class="text-gray-500">-</span>
              <input
                v-model.number="priceRange.max"
                @change="handlePriceFilter"
                type="number"
                placeholder="Max"
                class="form-input w-20"
                :min="teachersStore?.priceRange?.min || 0"
                :max="teachersStore?.priceRange?.max || 100"
              />
              <span class="text-sm text-gray-500 dark:text-gray-400">/hour</span>
            </div>
            <button
              @click="clearFilters"
              class="btn-link text-sm"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Teachers Grid -->
    <section class="section">
      <div class="container">
        <!-- Loading State -->
        <div v-if="teachersStore?.isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="i in 6" :key="i" class="card">
            <div class="card-body">
              <div class="skeleton h-48 w-full mb-4 rounded"></div>
              <div class="skeleton h-6 w-3/4 mb-2"></div>
              <div class="skeleton h-4 w-full mb-2"></div>
              <div class="skeleton h-4 w-2/3"></div>
            </div>
          </div>
        </div>

        <!-- No Results -->
        <div v-else-if="displayedTeachers.length === 0" class="text-center py-12">
          <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
            No teachers found
          </h3>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            Try adjusting your search criteria or filters
          </p>
          <button @click="clearFilters" class="btn btn-primary">
            Clear All Filters
          </button>
        </div>

        <!-- Teachers Grid -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="teacher in displayedTeachers"
            :key="teacher.ID || teacher.id"
            class="card hover:shadow-lg transition-all duration-300 cursor-pointer transform hover:-translate-y-1"
            @click="$router.push(`/teachers/${teacher.ID || teacher.id}`)"
          >
            <div class="card-body">
              <!-- Teacher Image -->
              <div class="relative mb-4">
                <div class="w-full h-48 bg-gray-200 dark:bg-gray-700 rounded-lg overflow-hidden">
                  <img
                    v-if="teacher.ProfileImage || teacher.profile_image"
                    :src="teacher.ProfileImage || teacher.profile_image"
                    :alt="teacher.Name || teacher.name"
                    class="w-full h-full object-cover"
                  />
                  <div v-else class="w-full h-full flex items-center justify-center text-4xl font-bold text-gray-500">
                    {{ (teacher.Name || teacher.name || '').charAt(0) }}
                  </div>
                </div>

                <!-- Favorite Button -->
                <div class="absolute top-2 left-2">
                  <button
                    @click.stop="toggleFavoriteTeacher(teacher)"
                    class="p-1 rounded-full bg-white bg-opacity-75 hover:bg-opacity-100 shadow"
                    title="Toggle favorite"
                  >
                    <svg
                      v-if="favoritesStore?.favoriteIds?.includes(teacher.ID || teacher.id)"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                      class="w-5 h-5 text-red-500"
                    >
                      <path d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 18.657l-6.828-6.829a4 4 0 010-5.656z" />
                    </svg>
                    <svg
                      v-else
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      class="w-5 h-5 text-gray-400 hover:text-red-500"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.172 5.172a4 4 0 015.656 0L12 8.343l3.172-3.171a4 4 0 115.656 5.656L12 21.657l-8.828-8.829a4 4 0 010-5.656z" />
                    </svg>
                  </button>
                </div>
                
                <!-- Online Status -->
                <div class="absolute top-2 right-2">
                  <div class="w-3 h-3 bg-green-400 rounded-full border-2 border-white"></div>
                </div>
              </div>

              <!-- Teacher Info -->
              <div class="space-y-3">
                <div>
                  <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-1">
                    {{ teacher.Name || teacher.name }}
                  </h3>
                  <div class="flex items-center space-x-2">
                    <span :class="getLevelBadgeClass(teacher.LanguageLevel || teacher.language_level)" class="badge text-xs">
                      {{ Utils.capitalize(teacher.LanguageLevel || teacher.language_level) }}
                    </span>
                    <span class="text-sm text-gray-500 dark:text-gray-400">
                      Native Speaker
                    </span>
                  </div>
                </div>

                <p class="text-gray-600 dark:text-gray-400 text-sm line-clamp-3">
                  {{ teacher.Bio || teacher.bio }}
                </p>

                <!-- Availability -->
                <div class="flex items-center text-sm text-gray-500 dark:text-gray-400">
                  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Available {{ teacher.AvailableStartTime || teacher.available_start_time }} - {{ teacher.AvailableEndTime || teacher.available_end_time }}
                </div>

                <!-- Price and Book Button -->
                <div class="flex items-center justify-between pt-2 border-t border-gray-200 dark:border-gray-700">
                  <div class="price-small">
                    {{ Utils.formatCurrency(teacher.PricePerHour || teacher.price_per_hour) }}/hour
                  </div>
                  <button
                    @click.stop="bookLesson(teacher)"
                    class="btn btn-primary text-sm px-4 py-2"
                  >
                    Book Lesson
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- Pagination controls -->
        <!-- Display pagination when more than one page exists.  We use
             pagination metadata from the store to compute the number of pages.
             If the backend does not support pagination, this will fall back
             to a single page. -->
        <div v-if="totalPages > 1" class="mt-12 flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0">
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
          <div class="text-sm text-gray-600 dark:text-gray-400">
            Showing {{ (currentPage - 1) * pagination.limit + 1 }} to
            {{ Math.min(currentPage * pagination.limit, totalItems) }} of
            {{ totalItems }} teachers
          </div>
        </div>
      </div>
    </section>

    <!-- Booking Modal -->
    <div v-if="showBookingModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
              Book Lesson with {{ selectedTeacher?.name }}
            </h2>
            <button
              @click="closeBookingModal"
              class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Booking Steps -->
          <div class="mb-6">
            <div class="flex items-center space-x-4">
              <div :class="['step', bookingStep >= 1 ? 'active' : '']">1</div>
              <div class="flex-1 h-1 bg-gray-200 dark:bg-gray-700">
                <div :class="['h-full bg-blue-600 transition-all', bookingStep >= 2 ? 'w-full' : 'w-0']"></div>
              </div>
              <div :class="['step', bookingStep >= 2 ? 'active' : '']">2</div>
              <div class="flex-1 h-1 bg-gray-200 dark:bg-gray-700">
                <div :class="['h-full bg-blue-600 transition-all', bookingStep >= 3 ? 'w-full' : 'w-0']"></div>
              </div>
              <div :class="['step', bookingStep >= 3 ? 'active' : '']">3</div>
            </div>
            <div class="flex justify-between mt-2 text-sm text-gray-600 dark:text-gray-400">
              <span>Select Schedule</span>
              <span>Booking Details</span>
              <span>Payment</span>
            </div>
          </div>

          <!-- Step 1: Schedule Selection -->
          <div v-if="bookingStep === 1" class="space-y-4">
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
                    ${{ selectedTeacher?.price_per_hour }}
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
          <div v-if="bookingStep === 2" class="space-y-4">
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
              
              <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
                <h4 class="font-medium mb-2">Booking Summary</h4>
                <div class="space-y-2 text-sm">
                  <div class="flex justify-between">
                    <span>Teacher:</span>
                    <span>{{ selectedTeacher?.name }}</span>
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
                    <span>${{ selectedTeacher?.price_per_hour }}</span>
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
          <div v-if="bookingStep === 3 && createdBooking">
            <PaymentProcessor
              :booking-data="{
                booking_id: createdBooking.id,
                teacher_name: selectedTeacher?.name,
                date: selectedSchedule?.date,
                start_time: selectedSchedule?.start_time,
                end_time: selectedSchedule?.end_time,
                amount: selectedTeacher?.price_per_hour * 100, // Convert to cents
                duration: 60
              }"
              @payment-success="handlePaymentSuccess"
              @payment-failed="handlePaymentFailed"
              @cancel="closeBookingModal"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { userService } from '@/services/api'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTeachersStore } from '@/stores/teachers'
import { useBookingsStore } from '@/stores/bookings'
import { useUIStore } from '@/stores/ui'
import { useFavoritesStore } from '@/stores/favorites'
import { Utils } from '@/utils';
import PaymentProcessor from '@/components/PaymentProcessor.vue'

const router = useRouter()

// Initialize stores with error handling
let authStore = null
let teachersStore = null
let bookingsStore = null
let uiStore = null
let favoritesStore = null
const storesInitialized = ref(false)

try {
  authStore = useAuthStore()
  teachersStore = useTeachersStore()
  bookingsStore = useBookingsStore()
  uiStore = useUIStore()
  favoritesStore = useFavoritesStore()
  storesInitialized.value = true
} catch (error) {
  console.error('Failed to initialize stores in Teachers view:', error)
}

// Reactive data
const searchQuery = ref('')
const selectedLevel = ref('')
const sortBy = ref('name')
const priceRange = ref({
  min: null,
  max: null
})

// Computed properties
const displayedTeachers = computed(() => {
  if (!storesInitialized.value || !teachersStore?.filteredTeachers) {
    return []
  }
  
  let teachers = [...teachersStore.filteredTeachers]

  // Apply sorting
  switch (sortBy.value) {
    case 'price_low':
      teachers.sort((a, b) => (a.PricePerHour || a.price_per_hour) - (b.PricePerHour || b.price_per_hour))
      break
    case 'price_high':
      teachers.sort((a, b) => (b.PricePerHour || b.price_per_hour) - (a.PricePerHour || a.price_per_hour))
      break
    case 'level':
      teachers.sort((a, b) => (a.LanguageLevel || a.language_level).localeCompare(b.LanguageLevel || b.language_level))
      break
    default:
      teachers.sort((a, b) => (a.Name || a.name).localeCompare(b.Name || b.name))
  }

  return teachers
})

const hasMoreTeachers = computed(() => {
  if (!storesInitialized.value || !teachersStore?.pagination) {
    return false
  }
  return teachersStore.pagination.page * teachersStore.pagination.limit < teachersStore.pagination.total
})

// Local API call handler
const handleApiCall = async (apiCall) => {
  if (!storesInitialized.value || !uiStore) {
    console.warn('Stores not initialized, cannot handle API call')
    return null
  }

  try {
    const response = await apiCall()
    if (response && response.message && uiStore.showSuccess) {
      uiStore.showSuccess(response.message)
    }
    return response
  } catch (error) {
    let message = 'An unexpected error occurred'
    
    if (error?.response?.data?.message) {
      message = error.response.data.message
    } else if (error?.message) {
      message = error.message
    }
    
    if (uiStore?.showError) {
      uiStore.showError(message)
    }
    return null
  }
}

// Methods
const getLevelBadgeClass = (level) => {
  return `level-${level}`
}

const handleSearch = Utils.debounce(() => {
  if (storesInitialized.value && teachersStore?.setSearchFilter) {
    teachersStore.setSearchFilter(searchQuery.value)
  }
}, 300)

const handleLevelFilter = () => {
  if (storesInitialized.value && teachersStore?.setLevelFilter) {
    teachersStore.setLevelFilter(selectedLevel.value)
  }
}

const handlePriceFilter = () => {
  if (storesInitialized.value && teachersStore?.setPriceFilter) {
    teachersStore.setPriceFilter(priceRange.value.min, priceRange.value.max)
  }
}

const handleSort = () => {
  // Sorting is handled in computed property
}

const clearFilters = () => {
  searchQuery.value = ''
  selectedLevel.value = ''
  priceRange.value = { min: null, max: null }
  sortBy.value = 'name'
  if (storesInitialized.value && teachersStore?.clearFilters) {
    teachersStore.clearFilters()
  }
}

const showBookingModal = ref(false)
const selectedTeacher = ref(null)
const selectedSchedule = ref(null)
const bookingStep = ref(1)
const bookingNotes = ref('')
const availableSchedules = ref([])
const createdBooking = ref(null)

const bookLesson = async (teacher) => {
  if (!storesInitialized.value) {
    console.warn('Stores not initialized, cannot book lesson')
    return
  }

  if (!authStore?.isAuthenticated) {
    if (uiStore?.showError) {
      uiStore.showError('Please log in to book a lesson')
    }
    router.push('/login')
    return
  }
  
  selectedTeacher.value = teacher
  showBookingModal.value = true
  bookingStep.value = 1
  
  // Fetch available schedules for the teacher
  if (teachersStore?.fetchTeacherSchedules) {
    const result = await handleApiCall(() => teachersStore.fetchTeacherSchedules(teacher.id));
    if (result && result.success) {
      availableSchedules.value = result.data.filter(schedule => schedule.status === 'available')
    } else if (uiStore?.showError) {
      uiStore.showError('Failed to load available schedules')
    }
  }
}

const selectSchedule = (schedule) => {
  selectedSchedule.value = schedule
}

const nextStep = () => {
  if (bookingStep.value < 3) {
    bookingStep.value++
  }
}

const prevStep = () => {
  if (bookingStep.value > 1) {
    bookingStep.value--
  }
}

const createBooking = async () => {
  if (!storesInitialized.value) {
    console.warn('Stores not initialized, cannot create booking')
    return
  }

  if (!selectedSchedule.value || !authStore?.currentUser) {
    if (uiStore?.showError) {
      uiStore.showError('Missing required booking information')
    }
    return
  }

  const bookingData = {
    user_id: authStore.currentUser.id,
    schedule_id: selectedSchedule.value.id,
    note: bookingNotes.value,
    total_price: selectedTeacher.value.price_per_hour
  }

  if (bookingsStore?.createBooking) {
    const result = await handleApiCall(() => bookingsStore.createBooking(bookingData));
    if (result && result.success) {
      createdBooking.value = result.data
      nextStep()

      // Log activity for booking creation.  Use the user service to
      // record that the user has created a booking.  The backend will
      // associate the log entry with the current authenticated user.
      try {
        await userService.logActivity({
          action: 'BookingCreated',
          description: `Created a booking for teacher ${selectedTeacher.value?.name || selectedTeacher.value?.Name || ''}`
        })
      } catch (error) {
        console.warn('Failed to log booking activity:', error)
      }
    }
  }
}

const handlePaymentSuccess = (data) => {
  if (uiStore?.showSuccess) {
    uiStore.showSuccess('Payment completed successfully!')
  }

  // Log activity for successful payment.  Notify the backend that the
  // user has completed a payment for a booking.  The backend will
  // associate this activity with the current user.
  userService.logActivity({
    action: 'PaymentCompleted',
    description: `Completed payment for booking ${createdBooking.value?.id || ''}`
  }).catch((err) => {
    console.warn('Failed to log payment activity:', err)
  })

  closeBookingModal()
  router.push('/bookings')
}

const handlePaymentFailed = (data) => {
  if (uiStore?.showError) {
    uiStore.showError('Payment failed. Please try again.')
  }
}

// Toggle favorite state for a teacher. When invoked, this will
// add or remove the teacher from the user's favorites via the
// favorites store. The click is stopped at the button level to
// prevent triggering the card click navigation.
const toggleFavoriteTeacher = async (teacher) => {
  if (!storesInitialized.value || !favoritesStore) return
  const teacherId = teacher.ID || teacher.id
  if (!teacherId) return
  if (!authStore?.isAuthenticated) {
    uiStore?.showError('Please log in to favorite a teacher')
    return
  }
  await favoritesStore.toggleFavorite(teacherId)
}

const closeBookingModal = () => {
  showBookingModal.value = false
  selectedTeacher.value = null
  selectedSchedule.value = null
  bookingStep.value = 1
  bookingNotes.value = ''
  availableSchedules.value = []
  createdBooking.value = null
}

const loadMoreTeachers = async () => {
  if (storesInitialized.value && teachersStore?.setPage && teachersStore?.pagination) {
    teachersStore.setPage(teachersStore.pagination.page + 1)
    await fetchTeachers()
  }
}

const fetchTeachers = async () => {
  if (storesInitialized.value && teachersStore?.fetchTeachers) {
    await handleApiCall(() => teachersStore.fetchTeachers());
  }
}
// Pagination for teacher listing.  We derive pagination information
// from the teachers store.  When the backend supplies total and limit
// values, we use them to compute the number of pages.  Otherwise,
// fall back to computing from the number of teachers loaded locally.
const pagination = computed(() => {
  return teachersStore?.pagination || { page: 1, limit: 10, total: 0 }
})

// Current page number used by pagination controls.  Defaults to 1.
const currentPage = computed(() => {
  return pagination.value.page || 1
})

// Total number of items available.  Prefer the total reported by the
// backend (if greater than zero) but fall back to the length of the
// teachers array when no total metadata is provided.  This allows
// pagination controls to appear when the backend returns a full list
// without pagination metadata.
const totalItems = computed(() => {
  if (pagination.value.total && pagination.value.total > 0) {
    return pagination.value.total
  }
  return Array.isArray(teachersStore?.teachers) ? teachersStore.teachers.length : 0
})

// Compute the total number of pages by dividing total items by the
// pagination limit.  Ensure at least one page is returned.
const totalPages = computed(() => {
  const limit = pagination.value.limit || 10
  const pages = Math.ceil(totalItems.value / limit)
  return pages > 0 ? pages : 1
})

// Generate a set of visible pages for pagination controls (max 5 pages)
const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  const total = totalPages.value
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = start + maxVisible - 1
  if (end > total) {
    end = total
    start = Math.max(1, end - maxVisible + 1)
  }
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Navigate to a specific page.  Update the store's page and then
// refetch teachers from the backend.  If the requested page is out of
// bounds, do nothing.  The teachers store will handle whether
// server‑side pagination is applied; if not, the computed properties
// above will slice the local array.
const goToPage = async (page) => {
  if (!storesInitialized.value || !teachersStore) return
  if (page < 1 || page > totalPages.value) return
  teachersStore?.setPage?.(page)

  if (pagination.value.total && pagination.value.total > 0 && teachersStore?.fetchTeachers) {
    await fetchTeachers()
  }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// Lifecycle
onMounted(async () => {
  if (!storesInitialized.value) {
    console.warn('Stores not initialized, skipping Teachers view initialization')
    return
  }

  if (uiStore?.setPageTitle) {
    uiStore.setPageTitle('Find Teachers')
  }
  
  await fetchTeachers()

  // Load the user's favorite teachers so the UI can show favorite
  // status on each teacher card. Do this after fetching teachers so
  // the favorites list is ready for rendering.
  if (authStore?.isAuthenticated && favoritesStore?.fetchFavorites) {
    await favoritesStore.fetchFavorites()
  }
  
  // Set initial price range
  if (teachersStore?.priceRange && teachersStore.priceRange.min !== teachersStore.priceRange.max) {
    priceRange.value = {
      min: teachersStore.priceRange.min,
      max: teachersStore.priceRange.max
    }
  }
})

// Watch for filter changes
watch([searchQuery, selectedLevel, priceRange], async () => {
  if (storesInitialized.value && teachersStore?.resetPagination) {
    teachersStore.resetPagination()
    // Refetch teachers for the first page when filters change
    await fetchTeachers()
  }
}, { deep: true })
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.step {
  @apply w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium;
  @apply bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-400;
}

.step.active {
  @apply bg-blue-600 text-white;
}

.form-select, .form-input, .form-textarea {
  @apply w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent;
  @apply dark:border-gray-600 dark:bg-gray-700 dark:text-white;
}

.btn {
  @apply px-4 py-2 rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-outline {
  @apply border border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-gray-500;
  @apply dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700;
}
</style>
