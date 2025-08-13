import { defineStore } from 'pinia'
import { teacherService } from '@/services/api'
import { useAuthStore } from '@/stores/auth'

export const useTeacherDashboardStore = defineStore('teacherDashboard', {
  state: () => ({
    dashboardData: null,
    isLoading: false,
    error: null
  }),

  getters: {
    teacherProfile: (state) => state.dashboardData?.teacherProfile || null,
    stats: (state) => state.dashboardData?.stats || {
      totalStudents: 0,
      upcomingBookings: 0,
      completedLessons: 0,
      totalEarnings: 0
    },
    upcomingBookings: (state) => state.dashboardData?.upcomingBookings || [],
    recentStudents: (state) => state.dashboardData?.recentStudents || [],
    completedLessons: (state) => state.dashboardData?.completedLessons || []
  },

  actions: {
    async fetchTeacherDashboard(teacherId) {
      this.isLoading = true
      this.error = null
      
      try {
        const response = await teacherService.getTeacherDashboard(teacherId)
        this.dashboardData = response.data
      } catch (error) {
        this.error = error.message || 'Failed to load dashboard data'
        console.error('Error fetching teacher dashboard:', error)
      } finally {
        this.isLoading = false
      }
    },

    async refreshDashboard() {
      const authStore = useAuthStore()
      if (authStore.currentUser?.id) {
        await this.fetchTeacherDashboard(authStore.currentUser.id)
      }
    },

    clearDashboard() {
      this.dashboardData = null
      this.error = null
    }
  }
})
