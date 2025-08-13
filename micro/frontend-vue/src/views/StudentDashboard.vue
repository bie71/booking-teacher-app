<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="container py-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Welcome back, {{ authStore.currentUser?.name }}!
        </h1>
        <p class="text-gray-600 dark:text-gray-400">
          Here's your learning progress and upcoming lessons
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <!-- Stats Cards -->
        <div class="card">
          <div class="card-body text-center">
            <div class="text-3xl font-bold text-blue-600 dark:text-blue-400 mb-2">
              {{ upcomingLessons.length }}
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Upcoming Lessons
            </div>
          </div>
        </div>

        <div class="card">
          <div class="card-body text-center">
            <div class="text-3xl font-bold text-green-600 dark:text-green-400 mb-2">
              {{ completedLessons.length }}
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Completed Lessons
            </div>
          </div>
        </div>

        <div class="card">
          <div class="card-body text-center">
            <div class="text-3xl font-bold text-purple-600 dark:text-purple-400 mb-2">
              {{ totalHours }}
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Total Hours
            </div>
          </div>
        </div>

        <div class="card">
          <div class="card-body text-center">
            <div class="text-3xl font-bold text-yellow-600 dark:text-yellow-400 mb-2">
              {{ favoriteTeachers.length }}
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Favorite Teachers
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Upcoming Lessons -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Upcoming Lessons</h2>
          </div>
          <div class="card-body">
            <div v-if="upcomingLessons.length === 0" class="text-center py-8">
              <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <p class="text-gray-500 dark:text-gray-400 mb-4">
                No upcoming lessons scheduled
              </p>
              <router-link to="/teachers" class="btn btn-primary">
                Book a Lesson
              </router-link>
            </div>
            <div v-else class="space-y-4">
              <div
                v-for="lesson in upcomingLessons.slice(0, 3)"
                :key="lesson.id"
                class="flex items-center space-x-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
              >
                <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                  <span class="text-blue-600 dark:text-blue-400 font-semibold">
                    {{ lesson.schedule?.teacher?.name?.charAt(0) }}
                  </span>
                </div>
                <div class="flex-1">
                  <h3 class="font-medium text-gray-900 dark:text-white">
                    {{ lesson.schedule?.teacher?.name }}
                  </h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    {{ Utils.formatDate(lesson.schedule?.date) }} at {{ lesson.schedule?.start_time }}
                  </p>
                </div>
                <span :class="getStatusBadgeClass(lesson.status)" class="badge text-xs">
                  {{ Utils.capitalize(lesson.status) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Recent Activity</h2>
          </div>
          <div class="card-body">
            <!-- If there are no recent activities, show a placeholder -->
            <div v-if="recentActivities.length === 0" class="text-center py-8">
              <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <p class="text-gray-500 dark:text-gray-400">
                No recent activities recorded
              </p>
            </div>
            <!-- Otherwise display recent activities -->
            <div v-else class="space-y-3">
              <div v-for="activity in recentActivities" :key="activity.id" class="flex items-center justify-between p-3 rounded-lg bg-gray-50 dark:bg-gray-700">
                <span class="text-sm text-gray-700 dark:text-gray-300 flex-1">
                  {{ activity.description }}
                </span>
                <span class="text-xs text-gray-500 dark:text-gray-400 ml-4 whitespace-nowrap">
                  {{ Utils.formatDateTime(activity.created_at) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useBookingsStore } from '@/stores/bookings'
import { useUIStore } from '@/stores/ui'
import { handleApiCall, Utils } from '@/utils';
import { userService, bookingService } from '@/services/api'

const authStore = useAuthStore()
const bookingsStore = useBookingsStore()
const uiStore = useUIStore()

const upcomingLessons = computed(() => bookingsStore.upcomingBookings)
const completedLessons = computed(() => bookingsStore.pastBookings)
const totalHours = computed(() => completedLessons.value.length)
import { useFavoritesStore } from '@/stores/favorites'

// Access the favorites store to manage and retrieve favorite teachers
const favoritesStore = useFavoritesStore()
// Computed property that reflects the current list of favorite teacher IDs
const favoriteTeachers = computed(() => favoritesStore.favoriteIds)

// Recent activity logs for the current user
const recentActivities = ref([])

// Fetch recent activity logs for the current user
const fetchRecentActivities = async () => {
  if (!authStore.currentUser) return
  try {
    const result = await userService.getRecentActivity({ limit: 5 })
    // The API may return data under different keys; attempt to normalize
    recentActivities.value =
      result.activities || result.data || result || []
  } catch (error) {
    console.error('Failed to fetch recent activities', error)
  }
}

const getStatusBadgeClass = (status) => {
  return `status-${status}`
}

const fetchDashboardData = async () => {
  if (authStore.currentUser) {
    await handleApiCall(() => bookingsStore.fetchBookingsByUser(authStore.currentUser.id));
  }

  // Always fetch recent activities when loading the dashboard
  await fetchRecentActivities()

  // Fetch the user's favorite teachers so the count reflects any changes
  await favoritesStore.fetchFavorites()
}

onMounted(() => {
  uiStore.setPageTitle('Student Dashboard')
  fetchDashboardData()
})
</script>
