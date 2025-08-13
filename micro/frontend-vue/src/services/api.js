import axios from 'axios'
import { API_CONFIG, ERROR_MESSAGES } from '@/config'
import { useAuthStore } from '@/stores/auth'

// Create axios instance
const apiClient = axios.create({
  timeout: API_CONFIG.TIMEOUT,
  headers: API_CONFIG.HEADERS
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const authStore = useAuthStore()
    
    if (error.response?.status === 401 && authStore.token) {
      // Try to refresh token
      try {
        await authStore.refreshToken()
        // Retry the original request
        return apiClient.request(error.config)
      } catch (refreshError) {
        authStore.logout()
        throw new ApiError(ERROR_MESSAGES.UNAUTHORIZED, 401)
      }
    }
    
    const message = error.response?.data?.message || getErrorMessage(error.response?.status)
    throw new ApiError(message, error.response?.status || 0, error.response?.data)
  }
)

// Custom API Error class
export class ApiError extends Error {
  constructor(message, status = 0, data = null) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

// Get error message based on status code
function getErrorMessage(status) {
  switch (status) {
    case 400:
      return ERROR_MESSAGES.VALIDATION_ERROR
    case 401:
      return ERROR_MESSAGES.UNAUTHORIZED
    case 403:
      return ERROR_MESSAGES.FORBIDDEN
    case 404:
      return ERROR_MESSAGES.NOT_FOUND
    case 409:
      return ERROR_MESSAGES.BOOKING_CONFLICT
    case 500:
    case 502:
    case 503:
    case 504:
      return ERROR_MESSAGES.SERVER_ERROR
    default:
      return ERROR_MESSAGES.UNKNOWN_ERROR
  }
}

// User Service
export const userService = {
  async register(userData) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/register`, userData)
    return response.data
  },

  async login(credentials) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/login`, credentials)
    return response.data
  },

  async getUsers(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.ADMIN_SERVICE}/users`, { params })
    return response
  },

  async createUser(userData) {
    const response = await apiClient.post(`${API_CONFIG.ADMIN_SERVICE}/users`, userData)
    return response.data
  },

  async getUser(id) {
    const response = await apiClient.get(`${API_CONFIG.ADMIN_SERVICE}/users/${id}`)
    return response.data
  },

  async updateUser(id, userData) {
    const response = await apiClient.put(`${API_CONFIG.ADMIN_SERVICE}/users/${id}`, userData)
    return response.data
  },

  async deleteUser(id) {
    const response = await apiClient.delete(`${API_CONFIG.ADMIN_SERVICE}/users/${id}`)
    return response.data
  },

  async getProfile() {
    const response = await apiClient.get(`${API_CONFIG.USER_SERVICE}/me`)
    return response.data
  },

  async updateProfile(profileData) {
    const response = await apiClient.put(`${API_CONFIG.USER_SERVICE}/profile`, profileData)
    return response.data
  },

  async changePassword(passwordData) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/change-password`, passwordData)
    return response.data
  },

  async resetPassword(resetData) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/reset-password`, resetData)
    return response.data
  },

  /**
   * Request a password reset email. This will trigger the backend to
   * generate a reset token and send a reset link to the specified
   * email address. The response is generic to avoid revealing whether
   * the email exists in the system.
   *
   * @param {string} email The email address to send the reset link to.
   */
  async forgotPassword(email) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/forgot-password`, { email })
    return response.data
  },

  /**
   * Verify the validity of a password reset token. The backend will
   * return an error if the token is invalid or expired. This is used
   * to pre-validate tokens before showing the reset form.
   *
   * @param {string} token The reset token to verify.
   */
  async verifyResetToken(token) {
    const response = await apiClient.get(`${API_CONFIG.USER_SERVICE}/verify-reset-token`, { params: { token } })
    return response.data
  },

  async refreshToken(refreshToken) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/refresh`, {
      refresh_token: refreshToken
    })
    return response.data
  },

  async uploadUserImage(file) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/upload-image`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },
  async uploadImageHero(file) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await apiClient.post(`${API_CONFIG.ADMIN_SERVICE}/upload-hero-image`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },
  async saveImageHero(requestUpload) {
    const response = await apiClient.post(`${API_CONFIG.ADMIN_SERVICE}/hero-image`, requestUpload)
    return response.data
  },
  async getImageHero(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.USER_SERVICE}/hero-image`, {
      params
    })
    return response.data
  },
  async deleteImageHero(params = {}) {
    const response = await apiClient.delete(`${API_CONFIG.ADMIN_SERVICE}/hero-image`, {
      params
    })
    return response.data
  },

  async getStatistics() {
    const response = await apiClient.get(`${API_CONFIG.ADMIN_SERVICE}/dashboard-stats`)
    return response.data
  }

  ,
  /**
   * Log a user activity. This endpoint sends an action and optional
   * description to the backend, which records the activity for the
   * currently authenticated user.
   *
   * @param {Object} payload An object containing `action` and
   *                         optional `description` fields.
   */
  async logActivity(payload) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/activity`, payload)
    return response.data
  },

  /**
   * Retrieve recent activities for the current user. The `limit`
   * parameter controls how many log entries are returned.
   *
   * @param {Object} params Query parameters, such as { limit: 5 }.
   */
  async getRecentActivity(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.USER_SERVICE}/activity/recent`, { params })
    return response.data
  }

  ,
  /**
   * Retrieve the IDs of teachers that the current user has favorited.
   * Returns an array of teacher IDs. If the user has no favorites,
   * returns an empty array.
   */
  async getFavoriteTeachers() {
    const response = await apiClient.get(`${API_CONFIG.USER_SERVICE}/favorites`)
    return response.data
  },

  /**
   * Add or remove a teacher from the current user's favorites list.
   * Pass an object containing `teacher_id` and an optional `favorite`
   * boolean. If `favorite` is omitted or true, the teacher will be
   * added; if false, the teacher will be removed. Returns a message
   * indicating success.
   *
   * @param {Object} payload { teacher_id: number, favorite?: boolean }
   */
  async toggleFavoriteTeacher(payload) {
    const response = await apiClient.post(`${API_CONFIG.USER_SERVICE}/favorites`, payload)
    return response.data
  }

}

// Teacher Service
export const teacherService = {
  async getMe(){ const r = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/teachers/me`); return r.data },
  async getTeachers(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/teachers`, { params })
    return response.data
  },

  async getTeacher(id) {
    const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/teachers/${id}`)
    return response.data
  },

  async createTeacher(teacherData) {
    const response = await apiClient.post(`${API_CONFIG.TEACHER_SERVICE}/teachers`, teacherData)
    return response.data
  },

  async updateTeacher(id, teacherData) {
    const response = await apiClient.put(`${API_CONFIG.TEACHER_SERVICE}/teachers/${id}`, teacherData)
    return response.data
  },

  async deleteTeacher(id) {
    const response = await apiClient.delete(`${API_CONFIG.TEACHER_SERVICE}/teachers/${id}`)
    return response.data
  },

  async getTeacherSchedules(teacherId, params = {}) {
    const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/schedule/teacher/${teacherId}`, { params })
    return response.data
  },

  async getSchedules(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/schedules`, { params })
    return response.data
  },

  async getSchedule(id) {
    const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/schedule/${id}`)
    return response.data
  },

  async createSchedule(scheduleData) {
    const response = await apiClient.post(`${API_CONFIG.TEACHER_SERVICE}/schedule`, scheduleData)
    return response.data
  },

  async updateSchedule(id, scheduleData) {
    const response = await apiClient.put(`${API_CONFIG.TEACHER_SERVICE}/schedule/${id}`, scheduleData)
    return response.data
  },

  async updateScheduleStatus(statusData) {
    const response = await apiClient.put(`${API_CONFIG.TEACHER_SERVICE}/schedule-status`, statusData)
    return response.data
  },

  async cancelSchedule(id) {
    const response = await apiClient.put(`${API_CONFIG.TEACHER_SERVICE}/cancel-schedule/${id}`)
    return response.data
  },

  async filterSchedulesByTeacher(filterData) {
    const response = await apiClient.post(`${API_CONFIG.TEACHER_SERVICE}/schedule/filter-by-teacher`, filterData)
    return response.data
  },

  async getBatchScheduleDetails(scheduleIds) {
    const response = await apiClient.post(`${API_CONFIG.TEACHER_SERVICE}/schedule/batch-detail`, {
      schedule_ids: scheduleIds
    })
    return response.data
  },

  async uploadTeacherImage(file) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await apiClient.post(`${API_CONFIG.TEACHER_SERVICE}/upload-image`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  async getTeacherDashboard(teacherId) {
    const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/teachers/dashboard/${teacherId}`)
    return response.data
  },
  
    // Delete schedule
    async deleteSchedule(scheduleId) {
      try {
        const response = await apiClient.delete(`${API_CONFIG.TEACHER_SERVICE}/schedule/${scheduleId}`)
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete schedule')
      }
    },
  
    // Get schedule by ID
    async getScheduleById(scheduleId) {
      try {
        const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/schedule/${scheduleId}`)
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch schedule')
      }
    },
  
    // Get all schedules with pagination
    async getAllSchedules(params = {}) {
      try {
        const response = await apiClient.get(`${API_CONFIG.TEACHER_SERVICE}/schedules`, { params })
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch schedules')
      }
    }
}
// Booking Service
export const bookingService = {
  async getAdminBookings(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.BOOKING_SERVICE}/admin/bookings`, { params })
    return response.data
  }
,
  
  async createBooking(bookingData) {
    const response = await apiClient.post(`${API_CONFIG.BOOKING_SERVICE}/bookings`, bookingData)
    return response.data
  },

  async getBookings(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.BOOKING_SERVICE}/bookings`, { params })
    return response.data
  },

  async getBooking(id) {
    const response = await apiClient.get(`${API_CONFIG.BOOKING_SERVICE}/booking/${id}`)
    return response.data
  },
    async getBookingDetail(id) {
    const response = await apiClient.get(`${API_CONFIG.BOOKING_SERVICE}/booking-detail/${id}`)
    return response.data
  },

  async getBookingsByUser(userId) {
    const response = await apiClient.get(`${API_CONFIG.BOOKING_SERVICE}/bookings/user/${userId}`)
    return response.data
  },

  /**
   * Fetch upcoming lessons for a user.
   *
   * Retrieves bookings with status "paid" and future schedule dates for the
   * specified user.  An optional `limit` parameter can be provided to
   * restrict the number of records returned.  Defaults to 5 if not
   * specified.
   *
   * @param {Number|String} userId The ID of the user
   * @param {Object} params Optional query parameters (e.g. { limit: 5 })
   * @returns {Promise<Object>} The response data containing upcoming bookings
   */
  async getUpcomingLessons(userId, params = {}) {
    const response = await apiClient.get(`${API_CONFIG.BOOKING_SERVICE}/bookings/user/${userId}/upcoming-lessons`, { params })
    return response.data
  },

  async rescheduleBooking(id, rescheduleData) {
    const response = await apiClient.post(`${API_CONFIG.BOOKING_SERVICE}/bookings/${id}/reschedule`, rescheduleData)
    return response.data
  },

  async cancelBooking(id) {
    const response = await apiClient.post(`${API_CONFIG.BOOKING_SERVICE}/bookings/${id}/cancel`)
    return response.data
  },

  async updateBookingStatus(id, statusData) {
    const response = await apiClient.put(`${API_CONFIG.BOOKING_SERVICE}/bookings/${id}/status`, statusData)
    return response.data
  },
  async changeStatus(id, statusData) {
    const response = await apiClient.put(`${API_CONFIG.BOOKING_SERVICE}/booking-change/${id}/status/${statusData}`)
    return response.data
  },

  async deleteBooking(id) {
    const response = await apiClient.delete(`${API_CONFIG.BOOKING_SERVICE}/booking/${id}`)
    return response.data
  },

  async uploadBookingImage(file) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await apiClient.post(`${API_CONFIG.BOOKING_SERVICE}/upload-image`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  }
}
// Payment Service
export const paymentService = {
  async createPayment(paymentData) {
    const response = await apiClient.post(`${API_CONFIG.PAYMENT_SERVICE}/payments`, paymentData)
    return response.data
  },

  async getPayments(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.PAYMENT_SERVICE}/payments`, { params })
    return response.data
  },

  async getPayment(id) {
    const response = await apiClient.get(`${API_CONFIG.PAYMENT_SERVICE}/payment/${id}`)
    return response.data
  },

  async handlePaymentCallback(callbackData) {
    const response = await apiClient.post(`${API_CONFIG.PAYMENT_SERVICE}/payments/callback`, callbackData)
    return response.data
  },

  async createPaymentMethod(methodData) {
    const response = await apiClient.post(`${API_CONFIG.PAYMENT_SERVICE}/admin/payment-methods/`, methodData)
    return response.data
  },
  async deletePaymentMethod(id) {
    const response = await apiClient.delete(`${API_CONFIG.PAYMENT_SERVICE}/admin/payment-methods/${id}`)
    return response.data
  },
   async updatePaymentMethod(id, methodData) {
    const response = await apiClient.put(`${API_CONFIG.PAYMENT_SERVICE}/admin/payment-methods/${id}`, methodData)
    return response.data
  },

  async updatePaymentMethodStatus(code, active) {
    const response = await apiClient.post(`${API_CONFIG.PAYMENT_SERVICE}/payment-methods/status`, { code, active })
    return response.data
  },

  async getPaymentMethods(params = {}) {
    // Public endpoint for retrieving available payment methods. Optional
    // query parameters are passed through but may be ignored by the
    // backend. This still allows consumers to request specific pages
    // or limits if the API supports it without breaking existing
    // functionality.
    const response = await apiClient.get(`${API_CONFIG.PAYMENT_SERVICE}/payment-methods/`, { params })
    return response.data
  },

  /**
   * Retrieve payment methods from the admin endpoint. This variant
   * returns paginated results along with pagination metadata. Only
   * authenticated users with admin privileges should call this. The
   * params argument should include page and limit.
   */
  async getAdminPaymentMethods(params = {}) {
    const response = await apiClient.get(`${API_CONFIG.PAYMENT_SERVICE}/admin/payment-methods/`, { params })
    return response.data
  }

}


export const authService = {
  async logout () {
    try {
      await apiClient.post(`${API_CONFIG.BOOKING_SERVICE}/auth/logout`)
    } catch (e) {
      console.error('Logout error', e)
    }
    try { localStorage.removeItem('token') } catch (e) {}
    if (apiClient?.defaults?.headers) {
      delete apiClient.defaults.headers.common['Authorization']
    }
  }
}

export default apiClient
