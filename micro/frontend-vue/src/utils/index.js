import { useUIStore } from '@/stores/ui';
import dayjs from 'dayjs'
import { APP_CONFIG, VALIDATION_RULES } from '@/config'
import { ApiError } from '@/services/api'

export async function handleApiCall(apiCall) {
  const uiStore = useUIStore();
  try {
    const response = await apiCall();
    if (response && response.message) {
      uiStore.showSuccess(response.message);
    }
    return response;
  } catch (error) {

    let message = 'An unexpected error occurred';

    if (error instanceof ApiError) {
      const data = error.data;
      message =
        typeof data === 'string' ? data :
        data?.message || data?.error || error.message || message;
    } else {
      const data = error.response?.data;
      message =
        typeof data === 'string' ? data :
        data?.message || data?.error || error.message || message;
    }

    uiStore.showError(message);
    return null;
  }
}




// Utility Functions
export const Utils = {
  // Format currency
  formatCurrency: (amount, currency = APP_CONFIG.CURRENCY) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(amount)
  },
  
  // Format date
  formatDate: (date, format = APP_CONFIG.DATE_FORMAT) => {
    if (!date) return ''
    return dayjs(date).format(format)
  },
  
  // Format time
  formatTime: (time) => {
    if (!time) return ''
    return time.substring(0, 5) // HH:mm format
  },
  
  // Format datetime
  formatDateTime: (datetime) => {
    if (!datetime) return ''
    return dayjs(datetime).format(APP_CONFIG.DATETIME_FORMAT)
  },
  
  // Validate email
  isValidEmail: (email) => {
    return VALIDATION_RULES.EMAIL.test(email)
  },
  
  // Validate password
  isValidPassword: (password) => {
    return password && password.length >= VALIDATION_RULES.PASSWORD_MIN_LENGTH
  },
  
  // Generate unique ID
  generateId: () => {
    return Date.now().toString(36) + Math.random().toString(36).substr(2)
  },
  
  // Debounce function
  debounce: (func, wait) => {
    let timeout
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout)
        func(...args)
      }
      clearTimeout(timeout)
      timeout = setTimeout(later, wait)
    }
  },
  
  // Throttle function
  throttle: (func, limit) => {
    let inThrottle
    return function() {
      const args = arguments
      const context = this
      if (!inThrottle) {
        func.apply(context, args)
        inThrottle = true
        setTimeout(() => inThrottle = false, limit)
      }
    }
  },
  
  // Deep clone object
  deepClone: (obj) => {
    return JSON.parse(JSON.stringify(obj))
  },
  
  // Check if object is empty
  isEmpty: (obj) => {
    return Object.keys(obj).length === 0
  },
  
  // Capitalize first letter
  capitalize: (str) => {
    return str.charAt(0).toUpperCase() + str.slice(1)
  },
  
  // Truncate text
  truncate: (text, length = 100) => {
    if (text.length <= length) return text
    return text.substring(0, length) + '...'
  },
  
  // Get query parameter
  getQueryParam: (param) => {
    const urlParams = new URLSearchParams(window.location.search)
    return urlParams.get(param)
  },
  
  // Scroll to element
  scrollTo: (elementId, offset = 0) => {
    const element = document.getElementById(elementId)
    if (element) {
      const top = element.offsetTop - offset
      window.scrollTo({
        top: top,
        behavior: 'smooth'
      })
    }
  },
  
  // Check if element is in viewport
  isInViewport: (element) => {
    const rect = element.getBoundingClientRect()
    return (
      rect.top >= 0 &&
      rect.left >= 0 &&
      rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
      rect.right <= (window.innerWidth || document.documentElement.clientWidth)
    )
  },

  // Format Japanese text
  formatJapanese: (text) => {
    if (!text) return ''
    return text.trim()
  },

  // Get relative time
  getRelativeTime: (date) => {
    return dayjs(date).fromNow()
  },

  // Check if date is today
  isToday: (date) => {
    return dayjs(date).isSame(dayjs(), 'day')
  },

  // Check if date is in the future
  isFuture: (date) => {
    return dayjs(date).isAfter(dayjs())
  },

  // Get time slots for a day
  getTimeSlots: (startTime, endTime, duration = 60) => {
    const slots = []
    let current = dayjs(`2000-01-01 ${startTime}`)
    const end = dayjs(`2000-01-01 ${endTime}`)
    
    while (current.isBefore(end)) {
      const next = current.add(duration, 'minute')
      if (next.isAfter(end)) break
      
      slots.push({
        start: current.format('HH:mm'),
        end: next.format('HH:mm'),
        label: `${current.format('HH:mm')} - ${next.format('HH:mm')}`
      })
      
      current = next
    }
    
    return slots
  },

  // Price range formatter
  formatPriceRange: (min, max) => {
    if (min === max) {
      return Utils.formatCurrency(min)
    }
    return `${Utils.formatCurrency(min)} - ${Utils.formatCurrency(max)}`
  },

  // Get level badge color
  getLevelColor: (level) => {
    const colors = {
      beginner: 'bg-green-100 text-green-800',
      intermediate: 'bg-yellow-100 text-yellow-800',
      advanced: 'bg-red-100 text-red-800'
    }
    return colors[level] || 'bg-gray-100 text-gray-800'
  },

  // Get status badge color
  getStatusColor: (status) => {
    const colors = {
      available: 'bg-green-100 text-green-800',
      booked: 'bg-blue-100 text-blue-800',
      cancelled: 'bg-red-100 text-red-800',
      pending: 'bg-yellow-100 text-yellow-800',
      paid: 'bg-green-100 text-green-800',
      completed: 'bg-gray-100 text-gray-800'
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }
}

export default Utils
