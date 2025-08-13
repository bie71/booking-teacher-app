// API Configuration
export const API_CONFIG = {
  // Base URLs for different services
  // Base URLs for different services. These are mapped to the ports
  // exposed by the docker-compose configuration (8001–8004). If you
  // change the service ports in docker-compose.yml, update them here
  // accordingly. During local development without Docker, you may
  // continue using ports 8081–8084 as before.
  USER_SERVICE: 'http://localhost:8081/api/v1',
  ADMIN_SERVICE: 'http://localhost:8081/api/admin',
  TEACHER_SERVICE: 'http://localhost:8082/api/v1',
  BOOKING_SERVICE: 'http://localhost:8083/api/v1',
  PAYMENT_SERVICE: 'http://localhost:8084/api/v1',

  KEY_IMAGE_HERO: 'image_hero',
  
  // Default headers
  HEADERS: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  
  // Request timeout
  TIMEOUT: 10000,
  
  // Pagination
  DEFAULT_PAGE_SIZE: 10,
  
  // Local storage keys
  STORAGE_KEYS: {
    TOKEN: 'japanlearn_token',
    USER: 'japanlearn_user',
    REFRESH_TOKEN: 'japanlearn_refresh_token'
  }
}

// Application Configuration
export const APP_CONFIG = {
  // Application name
  NAME: 'JapanLearn',
  
  // Version
  VERSION: '1.0.0',
  
  // Default language
  DEFAULT_LANGUAGE: 'en',
  
  // Supported languages
  SUPPORTED_LANGUAGES: ['en', 'ja'],
  
  // Date/Time formats
  DATE_FORMAT: 'YYYY-MM-DD',
  TIME_FORMAT: 'HH:mm',
  DATETIME_FORMAT: 'YYYY-MM-DD HH:mm',
  
  // Price format
  CURRENCY: 'USD',
  CURRENCY_SYMBOL: '$',
  
  // Booking settings
  BOOKING: {
    MIN_ADVANCE_HOURS: 2, // Minimum hours in advance to book
    MAX_ADVANCE_DAYS: 30, // Maximum days in advance to book
    LESSON_DURATION: 60, // Default lesson duration in minutes
    CANCELLATION_HOURS: 24 // Hours before lesson to allow cancellation
  },
  
  // UI Settings
  UI: {
    TOAST_DURATION: 3000, // Toast notification duration in ms
    LOADING_DELAY: 300, // Delay before showing loading spinner
    ANIMATION_DURATION: 300 // Default animation duration
  }
}

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  SERVER_ERROR: 'Server error. Please try again later.',
  UNAUTHORIZED: 'Please log in to continue.',
  FORBIDDEN: 'You do not have permission to perform this action.',
  NOT_FOUND: 'The requested resource was not found.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  BOOKING_CONFLICT: 'This time slot is no longer available.',
  PAYMENT_FAILED: 'Payment failed. Please try again.',
  UNKNOWN_ERROR: 'An unexpected error occurred.'
}

// Success Messages
export const SUCCESS_MESSAGES = {
  LOGIN_SUCCESS: 'Successfully logged in!',
  REGISTER_SUCCESS: 'Account created successfully!',
  BOOKING_SUCCESS: 'Lesson booked successfully!',
  BOOKING_CANCELLED: 'Booking cancelled successfully.',
  BOOKING_RESCHEDULED: 'Booking rescheduled successfully.',
  PROFILE_UPDATED: 'Profile updated successfully.',
  PASSWORD_CHANGED: 'Password changed successfully.'
}

// Validation Rules
export const VALIDATION_RULES = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PASSWORD_MIN_LENGTH: 6,
  NAME_MIN_LENGTH: 2,
  PHONE: /^\+?[\d\s\-\(\)]+$/
}

// Teacher Levels
export const TEACHER_LEVELS = {
  BEGINNER: 'beginner',
  INTERMEDIATE: 'intermediate', 
  ADVANCED: 'advanced'
}

// Booking Status
export const BOOKING_STATUS = {
  PENDING: 'pending',
  PAID: 'paid',
  CANCELLED: 'cancelled',
  RESCHEDULED: 'rescheduled',
  COMPLETED: 'completed'
}

// Payment Status
export const PAYMENT_STATUS = {
  PENDING: 'pending',
  SETTLEMENT: 'settlement',
  FAILED: 'failed',
  CANCEL: 'cancel'
}

// Schedule Status
export const SCHEDULE_STATUS = {
  AVAILABLE: 'available',
  BOOKED: 'booked',
  CANCELLED: 'cancelled'
}

// User Roles
export const USER_ROLES = {
  USER: 'user',
  TEACHER: 'teacher',
  ADMIN: 'admin'
}
