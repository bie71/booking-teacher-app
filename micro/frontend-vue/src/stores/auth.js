import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userService } from '@/services/api'
import { API_CONFIG, SUCCESS_MESSAGES, ERROR_MESSAGES } from '@/config'

export const useAuthStore = defineStore('auth', () => {
  // Helper function to safely parse JSON from localStorage
  const safeParseJSON = (key, fallback = null) => {
    try {
      const item = localStorage.getItem(key)
      if (item === null || item === undefined || item === 'undefined') {
        return fallback
      }
      return JSON.parse(item)
    } catch (error) {
      console.warn(`Failed to parse localStorage item "${key}":`, error)
      return fallback
    }
  }

  // State
  const token = ref(localStorage.getItem(API_CONFIG.STORAGE_KEYS.TOKEN))
  const user = ref(safeParseJSON(API_CONFIG.STORAGE_KEYS.USER, null))
  const refreshTokenValue = ref(localStorage.getItem(API_CONFIG.STORAGE_KEYS.REFRESH_TOKEN))
  const loading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!(token.value && user.value))
  const currentUser = computed(() => user.value)
  const userRole = computed(() => user.value?.role || null)
  const isLoading = computed(() => loading.value)

  // Actions
  const setAuthData = (tokenValue, userData, refreshToken = null) => {
    token.value = tokenValue
    user.value = userData
    refreshTokenValue.value = refreshToken

    // Store in localStorage with error handling
    try {
      if (tokenValue) {
        localStorage.setItem(API_CONFIG.STORAGE_KEYS.TOKEN, tokenValue)
      }
      if (userData) {
        localStorage.setItem(API_CONFIG.STORAGE_KEYS.USER, JSON.stringify(userData))
      }
      if (refreshToken) {
        localStorage.setItem(API_CONFIG.STORAGE_KEYS.REFRESH_TOKEN, refreshToken)
      }
    } catch (error) {
      console.warn('Failed to save auth data to localStorage:', error)
    }
  }

  const clearAuthData = () => {
    token.value = null
    user.value = null
    refreshTokenValue.value = null

    // Remove from localStorage with error handling
    try {
      localStorage.removeItem(API_CONFIG.STORAGE_KEYS.TOKEN)
      localStorage.removeItem(API_CONFIG.STORAGE_KEYS.USER)
      localStorage.removeItem(API_CONFIG.STORAGE_KEYS.REFRESH_TOKEN)
    } catch (error) {
      console.warn('Failed to clear auth data from localStorage:', error)
    }
  }

  const login = async (email, password) => {
    try {
      loading.value = true
      
      const response = await userService.login({
        email,
        password
      })

      setAuthData(response.token, response.user, response.refresh_token)
      
      return {
        success: true,
        message: SUCCESS_MESSAGES.LOGIN_SUCCESS,
        user: response.user
      }
    } catch (error) {
      throw error

    } finally {
      loading.value = false
    }
  }

  const register = async (name, email, password) => {
    try {
      loading.value = true
      
      const response = await userService.register({
        name,
        email,
        password,
        role: 'user'
      })

      // Auto-login after successful registration
      if (response.token && response.user) {
        setAuthData(response.token, response.user, response.refresh_token)
      }
      
      return {
        success: true,
        message: SUCCESS_MESSAGES.REGISTER_SUCCESS,
        user: response.user
      }
    } catch (error) {
      throw error

    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      // Optional: Call logout endpoint to invalidate token on server
      // await userService.logout()
    } catch (error) {
      console.warn('Logout endpoint error:', error)
    } finally {
      clearAuthData()
    }
  }

  const getUserProfile = async () => {
    try {
      if (!isAuthenticated.value) {
        throw new Error(ERROR_MESSAGES.UNAUTHORIZED)
      }

      const response = await userService.getProfile()
      
      // Update stored user data
      user.value = response.user || response
      try {
        localStorage.setItem(API_CONFIG.STORAGE_KEYS.USER, JSON.stringify(user.value))
      } catch (error) {
        console.warn('Failed to update user data in localStorage:', error)
      }
      
      return {
        success: true,
        user: user.value
      }
    } catch (error) {
      if (error.status === 401) {
        clearAuthData()
      }
      return {
        success: false,
        message: error.message
      }
    }
  }

  const updateProfile = async (profileData) => {
    try {
      if (!isAuthenticated.value) {
        throw new Error(ERROR_MESSAGES.UNAUTHORIZED)
      }

      loading.value = true
      const response = await userService.updateProfile(profileData)
      
      // Update stored user data
      user.value = { ...user.value, ...response.user }
      try {
        localStorage.setItem(API_CONFIG.STORAGE_KEYS.USER, JSON.stringify(user.value))
      } catch (error) {
        console.warn('Failed to update user profile in localStorage:', error)
      }
      
      return {
        success: true,
        message: SUCCESS_MESSAGES.PROFILE_UPDATED,
        user: user.value
      }
    } catch (error) {
      if (error.status === 401) {
        clearAuthData()
      }
      return {
        success: false,
        message: error.message || ERROR_MESSAGES.NETWORK_ERROR
      }
    } finally {
      loading.value = false
    }
  }

  const changePassword = async (currentPassword, newPassword) => {
    try {
      if (!isAuthenticated.value) {
        throw new Error(ERROR_MESSAGES.UNAUTHORIZED)
      }

      loading.value = true
      await userService.changePassword({
        current_password: currentPassword,
        new_password: newPassword
      })
      
      return {
        success: true,
        message: SUCCESS_MESSAGES.PASSWORD_CHANGED
      }
    } catch (error) {
      if (error.status === 401) {
        clearAuthData()
      }
      return {
        success: false,
        message: error.message || ERROR_MESSAGES.NETWORK_ERROR
      }
    } finally {
      loading.value = false
    }
  }

  const refreshToken = async () => {
    try {
      if (!refreshTokenValue.value) {
        throw new Error('No refresh token available')
      }

      const response = await userService.refreshToken(refreshTokenValue.value)
      
      // Update tokens
      setAuthData(response.token, user.value, response.refresh_token)
      
      return {
        success: true,
        token: response.token
      }
    } catch (error) {
      clearAuthData()
      return {
        success: false,
        message: error.message
      }
    }
  }

  const resetPassword = async (resetData) => {
    try {
      loading.value = true
      await userService.resetPassword(resetData)
      
      return {
        success: true,
        message: 'Password reset email sent successfully'
      }
    } catch (error) {
      throw error

    } finally {
      loading.value = false
    }
  }

  const getStatistics = async () => {
    try {
      loading.value = true
      const response = await userService.getStatistics()
      return response
    } catch (error) {
      console.error('Error fetching statistics:', error)
      throw error
    } finally {
      loading.value = false
    }
  }



  // Initialize auth state
  const initializeAuth = async () => {
    if (isAuthenticated.value) {
      // Verify token is still valid
      const result = await getUserProfile()
      if (!result.success) {
        clearAuthData()
      }
    }
  }

  return {
    // State
    token: computed(() => token.value),
    user: computed(() => user.value),
    loading: computed(() => loading.value),
    
    // Getters
    isAuthenticated,
    currentUser,
    userRole,
    isLoading,
    
    // Actions
    login,
    register,
    logout,
    getUserProfile,
    updateProfile,
    changePassword,
    refreshToken,
    resetPassword,
    initializeAuth,
    setAuthData,
    clearAuthData,
    getStatistics
  }
})
