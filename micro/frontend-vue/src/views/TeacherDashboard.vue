<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="container py-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Teacher Dashboard
        </h1>
        <p class="text-gray-600 dark:text-gray-400">
          Manage your lessons, students, and schedule
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="dashboardStore.isLoading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="dashboardStore.error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-6">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-red-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-red-700 dark:text-red-300">{{ dashboardStore.error }}</span>
        </div>
      </div>

      <!-- Dashboard Content -->
      <div v-else-if="dashboardStore.dashboardData">
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div class="card">
            <div class="card-body text-center">
              <div class="text-3xl font-bold text-blue-600 dark:text-blue-400 mb-2">
                {{ dashboardStore.stats.upcomingBookings }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Upcoming Bookings
              </div>
            </div>
          </div>

          <div class="card">
            <div class="card-body text-center">
              <div class="text-3xl font-bold text-green-600 dark:text-green-400 mb-2">
                {{ dashboardStore.stats.totalStudents }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Total Students
              </div>
            </div>
          </div>

          <div class="card">
            <div class="card-body text-center">
              <div class="text-3xl font-bold text-purple-600 dark:text-purple-400 mb-2">
                {{ dashboardStore.stats.completedLessons }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Completed Lessons
              </div>
            </div>
          </div>

          <div class="card">
            <div class="card-body text-center">
              <div class="text-3xl font-bold text-yellow-600 dark:text-yellow-400 mb-2">
                ${{ dashboardStore.stats.totalEarnings.toFixed(2) }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Total Earnings
              </div>
            </div>
          </div>
        </div>

        <!-- Teacher Profile Card -->
        <div class="card mb-8">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Teacher Profile</h2>
          </div>
          <div class="card-body">
            <div v-if="dashboardStore.teacherProfile" class="flex items-center space-x-4">
              <div class="w-16 h-16 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                <span class="text-2xl font-bold text-blue-600 dark:text-blue-400">
                  {{ dashboardStore.teacherProfile.name?.charAt(0) }}
                </span>
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                  {{ dashboardStore.teacherProfile.name }}
                </h3>
                <p class="text-gray-600 dark:text-gray-400">
                  {{ dashboardStore.teacherProfile.email }}
                </p>
                <p v-if="dashboardStore.teacherProfile.specialization" class="text-sm text-gray-500 dark:text-gray-500">
                  {{ dashboardStore.teacherProfile.specialization }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- Upcoming Bookings -->
          <div class="card">
            <div class="card-header">
              <h2 class="text-xl font-semibold">Upcoming Bookings</h2>
            </div>
            <div class="card-body">
              <div v-if="dashboardStore.upcomingBookings.length === 0" class="text-center py-8">
                <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <p class="text-gray-500 dark:text-gray-400 mb-4">
                  No upcoming bookings
                </p>
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="booking in dashboardStore.upcomingBookings.slice(0, 5)"
                  :key="booking.id"
                  class="flex items-center space-x-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
                >
                  <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                    <span class="text-blue-600 dark:text-blue-400 font-semibold">
                      {{ booking.studentName?.charAt(0) }}
                    </span>
                  </div>
                  <div class="flex-1">
                    <h3 class="font-medium text-gray-900 dark:text-white">
                      {{ booking.studentName }}
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400">
                      {{ formatDate(booking.date) }} at {{ booking.startTime }}
                    </p>
                  </div>
                  <span :class="getStatusBadgeClass(booking.status)" class="badge text-xs">
                    {{ capitalize(booking.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Students -->
          <div class="card">
            <div class="card-header">
              <h2 class="text-xl font-semibold">Recent Students</h2>
            </div>
            <div class="card-body">
              <div v-if="dashboardStore.recentStudents.length === 0" class="text-center py-8">
                <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a6 6 0 00-9 5.197m0 0A13.935 13.935 0 0112 21a13.935 13.935 0 01-6.5-2.698" />
                </svg>
                <p class="text-gray-500 dark:text-gray-400">
                  No recent students
                </p>
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="student in dashboardStore.recentStudents.slice(0, 5)"
                  :key="student.id"
                  class="flex items-center space-x-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
                >
                  <div class="w-12 h-12 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center">
                    <span class="text-green-600 dark:text-green-400 font-semibold">
                      {{ student.name?.charAt(0) }}
                    </span>
                  </div>
                  <div class="flex-1">
                    <h3 class="font-medium text-gray-900 dark:text-white">
                      {{ student.name }}
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400">
                      {{ student.email }}
                    </p>
                    <p class="text-xs text-gray-500 dark:text-gray-500">
                      {{ student.totalLessons }} lessons • ${{ student.totalSpent.toFixed(2) }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Completed Lessons -->
        <div class="card mt-8">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Completed Lessons</h2>
          </div>
          <div class="card-body">
            <div v-if="dashboardStore.completedLessons.length === 0" class="text-center py-8">
              <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p class="text-gray-500 dark:text-gray-400">
                No completed lessons yet
              </p>
            </div>
            <div v-else class="space-y-4">
              <div
                v-for="lesson in dashboardStore.completedLessons.slice(0, 5)"
                :key="lesson.id"
                class="flex items-center space-x-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
              >
                <div class="w-12 h-12 bg-purple-100 dark:bg-purple-900 rounded-full flex items-center justify-center">
                  <span class="text-purple-600 dark:text-purple-400 font-semibold">
                    {{ lesson.studentName?.charAt(0) }}
                  </span>
                </div>
                <div class="flex-1">
                  <h3 class="font-medium text-gray-900 dark:text-white">
                    {{ lesson.studentName }}
                  </h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    {{ formatDate(lesson.date) }} • ${{ lesson.price.toFixed(2) }}
                  </p>
                </div>
                <span class="badge badge-success text-xs">
                  Completed
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-12">
        <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2-2v7a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <p class="text-gray-500 dark:text-gray-400">
          No dashboard data available
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTeacherDashboardStore } from '@/stores/teacherDashboard'
import { useUIStore } from '@/stores/ui'

const authStore = useAuthStore()
const dashboardStore = useTeacherDashboardStore()
const uiStore = useUIStore()

// Helper functions
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric', 
    year: 'numeric' 
  })
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
}

const getStatusBadgeClass = (status) => {
  const classes = {
    pending: 'badge badge-warning',
    confirmed: 'badge badge-info',
    completed: 'badge badge-success',
    cancelled: 'badge badge-error'
  }
  return classes[status] || 'badge badge-ghost'
}

const fetchTeacherDashboardData = async () => {
  if (authStore.currentUser?.id) {
    await dashboardStore.fetchTeacherDashboard(authStore.currentUser.id)
  } else {
    uiStore.showError('Error fetching teacher dashboard data')
  }
}

onMounted(async () => {
  uiStore.setPageTitle('Teacher Dashboard')
  try { await dashboardStore.fetchTeacherDashboard(authStore.currentUser?.id || 0) } catch(e) { /* ignore */ }
  fetchTeacherDashboardData()
})
</script>

<style scoped>
.card {
  @apply bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700;
}

.card-header {
  @apply px-6 py-4 border-b border-gray-200 dark:border-gray-700;
}

.card-body {
  @apply p-6;
}

.badge {
  @apply inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium;
}

.badge-success {
  @apply bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200;
}

.badge-warning {
  @apply bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200;
}

.badge-info {
  @apply bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200;
}

.badge-error {
  @apply bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200;
}

.badge-ghost {
  @apply bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200;
}
</style>
