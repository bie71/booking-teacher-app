import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { APP_CONFIG } from '@/config'

export const useUIStore = defineStore('ui', () => {
  // State
  const notifications = ref([])
  const loading = ref(false)
  const sidebarOpen = ref(false)
  const modals = ref({
    login: false,
    register: false,
    booking: false,
    profile: false
  })
  const theme = ref(localStorage.getItem('theme') || 'light')

  // Getters
  const isLoading = computed(() => loading.value)
  const isSidebarOpen = computed(() => sidebarOpen.value)
  const activeNotifications = computed(() => notifications.value)
  const currentTheme = computed(() => theme.value)

  // Actions
  const showNotification = (message, type = 'info', duration = APP_CONFIG.UI.TOAST_DURATION) => {
    const id = Date.now().toString()
    const notification = {
      id,
      message,
      type, // 'success', 'error', 'warning', 'info'
      duration,
      timestamp: new Date()
    }

    notifications.value.push(notification)

    // Auto remove notification after duration
    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id)
      }, duration)
    }

    return id
  }

  const removeNotification = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  const clearNotifications = () => {
    notifications.value = []
  }

  // Specific notification types
  const showSuccess = (message, duration) => {
    return showNotification(message, 'success', duration)
  }

  const showError = (message, duration = 6000) => {
    // Longer duration for errors to ensure users see them
    return showNotification(message, 'error', duration)
  }

  const showWarning = (message, duration) => {
    return showNotification(message, 'warning', duration)
  }

  const showInfo = (message, duration) => {
    return showNotification(message, 'info', duration)
  }

  // Enhanced error handling methods
  const showNetworkError = (retryCallback = null) => {
    const message = retryCallback 
      ? 'Network error occurred. Click to retry.' 
      : 'Network error. Please check your connection.'
    
    const id = showError(message, 8000)
    
    if (retryCallback) {
      // Store retry callback for potential use
      const notification = notifications.value.find(n => n.id === id)
      if (notification) {
        notification.retryCallback = retryCallback
      }
    }
    
    return id
  }

  const showValidationErrors = (errors) => {
    if (typeof errors === 'object' && errors !== null) {
      Object.entries(errors).forEach(([field, message]) => {
        showError(`${field}: ${message}`, 5000)
      })
    } else if (typeof errors === 'string') {
      showError(errors, 5000)
    }
  }

  const showFormError = (message) => {
    return showError(`Form Error: ${message}`, 5000)
  }

  const showApiError = (error, context = '') => {
    let message = 'An unexpected error occurred'
    
    if (error?.message) {
      message = error.message
    } else if (error?.response?.data?.message) {
      message = error.response.data.message
    } else if (typeof error === 'string') {
      message = error
    }
    
    if (context) {
      message = `${context}: ${message}`
    }
    
    return showError(message, 6000)
  }

  // Loading state
  const setLoading = (isLoading) => {
    loading.value = isLoading
  }

  const startLoading = () => {
    loading.value = true
  }

  const stopLoading = () => {
    loading.value = false
  }

  // Sidebar
  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const openSidebar = () => {
    sidebarOpen.value = true
  }

  const closeSidebar = () => {
    sidebarOpen.value = false
  }

  // Modals
  const openModal = (modalName) => {
    if (modals.value.hasOwnProperty(modalName)) {
      modals.value[modalName] = true
    }
  }

  const closeModal = (modalName) => {
    if (modals.value.hasOwnProperty(modalName)) {
      modals.value[modalName] = false
    }
  }

  const closeAllModals = () => {
    Object.keys(modals.value).forEach(key => {
      modals.value[key] = false
    })
  }

  const isModalOpen = (modalName) => {
    return modals.value[modalName] || false
  }

  // Theme
  const setTheme = (newTheme) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    
    // Apply theme to document
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  const toggleTheme = () => {
    const newTheme = theme.value === 'light' ? 'dark' : 'light'
    setTheme(newTheme)
  }

  // Initialize theme
  const initializeTheme = () => {
    if (theme.value === 'dark') {
      document.documentElement.classList.add('dark')
    }
  }

  // Utility functions
  const scrollToTop = () => {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  }

  const scrollToElement = (elementId, offset = 0) => {
    const element = document.getElementById(elementId)
    if (element) {
      const top = element.offsetTop - offset
      window.scrollTo({
        top: top,
        behavior: 'smooth'
      })
    }
  }

  // Page title management
  const setPageTitle = (title) => {
    document.title = title ? `${title} - ${APP_CONFIG.NAME}` : APP_CONFIG.NAME
  }

  // Meta description
  const setMetaDescription = (description) => {
    const metaDescription = document.querySelector('meta[name="description"]')
    if (metaDescription) {
      metaDescription.setAttribute('content', description)
    }
  }

  // Confirmation dialog
  const confirm = (message, title = 'Confirm') => {
    // TODO: Replace with a non-blocking modal component for better UX.
    return new Promise((resolve) => {
      const result = window.confirm(`${title}\n\n${message}`)
      resolve(result)
    })
  }

  // Alert dialog
  const alert = (message, title = 'Alert') => {
    // TODO: Replace with a non-blocking modal component for better UX.
    return new Promise((resolve) => {
      window.alert(`${title}\n\n${message}`)
      resolve()
    })
  }

  // Copy to clipboard
  const copyToClipboard = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      showSuccess('Copied to clipboard')
      return true
    } catch (error) {
      showError('Failed to copy to clipboard')
      return false
    }
  }

  // Share functionality
  const share = async (data) => {
    try {
      if (navigator.share) {
        await navigator.share(data)
        return true
      } else {
        // Fallback to copying URL
        await copyToClipboard(data.url || window.location.href)
        return true
      }
    } catch (error) {
      showError('Failed to share')
      return false
    }
  }

  // Keyboard shortcuts
  const handleKeyboardShortcut = (event) => {
    // Ctrl/Cmd + K for search
    if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
      event.preventDefault()
      // Trigger search modal or focus search input
      const searchInput = document.querySelector('input[type="search"]')
      if (searchInput) {
        searchInput.focus()
      }
    }

    // Escape to close modals
    if (event.key === 'Escape') {
      closeAllModals()
      closeSidebar()
    }
  }

  // Initialize keyboard shortcuts
  const initializeKeyboardShortcuts = () => {
    document.addEventListener('keydown', handleKeyboardShortcut)
  }

  // Cleanup keyboard shortcuts
  const cleanupKeyboardShortcuts = () => {
    document.removeEventListener('keydown', handleKeyboardShortcut)
  }

  return {
    // Getters
    isLoading,
    isSidebarOpen,
    activeNotifications,
    currentTheme,
    modals: computed(() => modals.value), // Expose modals state

    // --- Actions ---

    // Notification Actions
    showNotification,
    removeNotification,
    clearNotifications,
    showSuccess,
    showError,
    showWarning,
    showInfo,

    // Error Handling Actions
    showNetworkError,
    showValidationErrors,
    showFormError,
    showApiError,

    // Loading Actions
    setLoading,
    startLoading,
    stopLoading,

    // Sidebar Actions
    toggleSidebar,
    openSidebar,
    closeSidebar,

    // Modal Actions
    openModal,
    closeModal,
    closeAllModals,
    isModalOpen,

    // Theme Actions
    setTheme,
    toggleTheme,
    initializeTheme,

    // Utility Actions
    scrollToTop,
    scrollToElement,
    setPageTitle,
    setMetaDescription,
    confirm,
    alert,
    copyToClipboard,
    share,
    initializeKeyboardShortcuts,
    cleanupKeyboardShortcuts
  }
})