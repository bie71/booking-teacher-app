<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <!-- Header -->
    <div class="bg-white shadow dark:bg-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Admin Dashboard</h1>
          <div class="flex items-center space-x-4">
            <button 
              @click="showHeroUpload = !showHeroUpload"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Upload Hero Image
            </button>
          <button 
            @click="deleteHeroImage"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
          >
            Delete Hero Image
          </button>
          <router-link 
            to="/admin/users"
            class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 dark:text-white"
          >
            User Management
          </router-link>
          <router-link 
            to="/admin/payment-methods"
            class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 dark:text-white"
          >
            Payment Methods
          </router-link>
          
            <button 
              @click="$router.push({ name: 'BookingManagement' })"
              class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
            >
              Booking Management
            </button>
          </div>
        </div>

      </div>
    </div>

    <!-- Hero Upload Modal -->
    <div v-if="showHeroUpload" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg max-w-2xl w-full mx-4 dark:bg-gray-800 dark:text-white">
        <AdminHeroUpload @hero-updated="handleHeroUpdated" @close="showHeroUpload = false" />
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white p-6 rounded-lg shadow dark:bg-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Total Users</h3>
          <p class="text-3xl font-bold text-blue-600">{{ stats.totalUsers }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow  dark:bg-gray-800 ">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Total Teachers</h3>
          <p class="text-3xl font-bold text-green-600 dark:text-white">{{ stats.totalTeachers }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow  dark:bg-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Total Bookings</h3>
          <p class="text-3xl font-bold text-purple-600">{{ stats.totalBookings }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow  dark:bg-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Revenue</h3>
          <p class="text-3xl font-bold text-yellow-600">${{ stats.totalRevenue }}</p>
        </div>
      </div>

      <!-- Teacher Management Section -->
      <div class="bg-white rounded-lg shadow mb-8  dark:bg-gray-800">
        <div class="border-b border-gray-200 dark:border-gray-700">
          <nav class="flex space-x-8 px-6" aria-label="Tabs">
            <button 
              @click="activeTab = 'teachers'"
              :class="activeTab === 'teachers' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
              class="whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
            >
              Teachers
            </button>
            <button 
              @click="activeTab = 'add-teacher'"
              :class="activeTab === 'add-teacher' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
              class="whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
            >
              Add Teacher
            </button>
            <!-- <button 
              @click="activeTab = 'schedules'"
              :class="activeTab === 'schedules' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
              class="whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
            >
              Schedules
            </button> -->
          </nav>
        </div>

        <div class="p-6">
          <!-- Teachers Tab -->
          <div v-if="activeTab === 'teachers'">
            <TeacherList 
             @edit-teacher="handleEditTeacher"
             @manage-schedule="handleManageSchedule"
            />
            <EditTeacherModal
              v-if="showEditModal"
              :teacher="selectedTeacher"
              @close="showEditModal = false"
            />

            <ScheduleTeacherModal
              v-if="showScheduleModal"
              :teacher-id="selectedTeacher?.id"
              @close="showScheduleModal = false"
            />

          </div>

          <!-- Add Teacher Tab -->
          <div v-if="activeTab === 'add-teacher'">
            <TeacherForm mode="create" @teacher-created="handleTeacherCreated" />
          </div>

          <!-- Schedules Tab -->
          <!-- <div v-if="activeTab === 'schedules'">
            <TeacherSchedule 
              @add-schedule="handleAddSchedule"
              @edit-schedule="handleEditSchedule"
              @delete-schedule="handleDeleteSchedule"
            />
          </div> -->
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white p-6 rounded-lg shadow dark:bg-gray-800">
          <h3 class="text-lg font-semibold mb-4">Quick Actions</h3>
          <div class="space-y-3">
            <router-link 
              to="/bookings" 
              class="block w-full px-4 py-2 bg-blue-600 text-white text-center rounded-md hover:bg-blue-700 dark:text-white"
            >
              View Bookings
            </router-link>
            <router-link 
              to="/profile" 
              class="block w-full px-4 py-2 bg-green-600 text-white text-center rounded-md hover:bg-green-700 dark:text-white"
            >
              Manage Profile
            </router-link>
          
          </div>
        </div>


        <div class="bg-white p-6 rounded-lg shadow dark:bg-gray-800">
          <!-- Upcoming Lessons Section -->
          <h3 class="text-lg font-semibold mb-4">Upcoming Lessons</h3>
          <div v-if="upcomingLessons.length === 0" class="text-sm text-gray-600 dark:text-gray-400 mb-4">
            No upcoming lessons scheduled
          </div>
          <div v-else class="space-y-3 mb-6">
            <div v-for="lesson in upcomingLessons" :key="lesson.id" class="flex items-center space-x-3">
              <div class="w-2 h-2 bg-green-500 rounded-full"></div>
              <div class="flex-1">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ lesson.schedule?.teacher?.name }}
                </span>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  {{ Utils.formatDate(lesson.schedule?.date) }} at {{ lesson.schedule?.start_time }}
                </div>
              </div>
              <span class="text-xs text-gray-500 dark:text-gray-400 ml-2">
                {{ Utils.capitalize(lesson.status) }}
              </span>
            </div>
          </div>

          <!-- Recent Activity Section -->
          <h3 class="text-lg font-semibold mb-4">Recent Activity</h3>
          <div v-if="recentActivities.length === 0" class="text-sm text-gray-600 dark:text-gray-400">
            No recent activities recorded
          </div>
          <div v-else class="space-y-3">
            <div v-for="activity in recentActivities" :key="activity.id" class="flex items-center space-x-3">
              <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
              <span class="text-sm text-gray-600 dark:text-gray-400 flex-1">{{ activity.description }}</span>
              <span class="text-xs text-gray-500 dark:text-gray-400 ml-2">
                {{ Utils.formatDateTime(activity.created_at) }}
              </span>
            </div>
          
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBookingsStore } from '@/stores/bookings'
import { userService, bookingService } from '@/services/api'
import { Utils } from '@/utils'
import { useUIStore } from '@/stores/ui'
import { API_CONFIG } from '@/config'

import TeacherList from '@/components/admin/TeacherList.vue'
import TeacherForm from '@/components/admin/TeacherForm.vue'
import EditTeacherModal from '@/modals/EditTeacherModal.vue'
import ScheduleTeacherModal from '@/modals/ScheduleTeacherModal.vue'

const uiStore = useUIStore()
import AdminHeroUpload from '@/components/admin/AdminHeroUpload.vue'

const router = useRouter()
const authStore = useAuthStore()
const bookingsStore = useBookingsStore()

const showHeroUpload = ref(false)
const activeTab = ref('teachers')

const stats = ref({
  totalUsers: 0,
  totalTeachers: 0,
  totalBookings: 0,
  totalRevenue: 0
})

// Recent activities for the current user.  Initially empty and populated on mount.
const recentActivities = ref([])

// Upcoming lessons for the current user.  Initially empty.
const upcomingLessons = ref([])

// Fetch the recent activity logs for the admin (current user)
const fetchRecentActivities = async () => {
  try {
    const res = await userService.getRecentActivity({ limit: 5 })
    recentActivities.value = res.activities || res.data || res || []
  } catch (error) {
    console.error('Failed to fetch recent activities', error)
  }
}

// Fetch upcoming lessons for the admin (current user)
const fetchUpcomingLessons = async () => {
  if (!authStore.currentUser) return
  try {
    const res = await bookingService.getUpcomingLessons(authStore.currentUser.id, { limit: 5 })
    // Response might wrap data under a `data` field
    upcomingLessons.value = res.data || res || []
  } catch (error) {
    console.error('Failed to fetch upcoming lessons', error)
  }
}



const selectedTeacher = ref(null)
const showEditModal = ref(false)
const showScheduleModal = ref(false)

const handleEditTeacher = (teacher) => {
  selectedTeacher.value = teacher
  showEditModal.value = true
}

const handleManageSchedule = (teacher) => {
  selectedTeacher.value = teacher
  showScheduleModal.value = true
}



const handleHeroUpdated = (url) => {
  showHeroUpload.value = false
  // Update hero image in dashboard
  console.log('Hero updated:', url)
}

const handleTeacherCreated = () => {
  activeTab.value = 'teachers'
}

const deleteHeroImage = async () => {
  try {
    uiStore.confirm('Are you sure you want to delete the hero image?', 'Delete Hero Image')

    await userService.deleteImageHero({ key: API_CONFIG.KEY_IMAGE_HERO })
    uiStore.showSuccess('Hero image deleted successfully!')
  } catch (error) {
    console.error('Failed to delete hero image:', error)
    uiStore.showError('Failed to delete hero image')
  }
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}


onMounted(async () => {
  // Load dashboard stats
  await bookingsStore.fetchBookings()

 const rawStats = await authStore.getStatistics()

  if (rawStats) {
    stats.value = {
      totalUsers: rawStats.data.total_users,
      totalTeachers: rawStats.data.total_teachers,
      totalBookings: rawStats.data.total_bookings,
      totalRevenue: rawStats.data.total_revenue
    }
  }

  // Load recent activities and upcoming lessons
  await fetchRecentActivities()
  await fetchUpcomingLessons()
})

</script>
