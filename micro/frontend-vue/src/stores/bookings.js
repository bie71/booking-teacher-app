import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { bookingService } from '@/services/api'
import { BOOKING_STATUS, SUCCESS_MESSAGES } from '@/config'

export const useBookingsStore = defineStore('bookings', () => {
  // State
  const bookings = ref([])
  const currentBooking = ref(null)
  const loading = ref(false)
  const filters = ref({
    status: '',
    dateFrom: null,
    dateTo: null,
    teacherId: null
  })
  const pagination = ref({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0
  })

  // Helper function to ensure bookings is always an array
  const normalizeBookings = (data) => {
    if (Array.isArray(data)) {
      return data
    } else if (data && typeof data === 'object') {
      // Handle different response structures
      if (Array.isArray(data.bookings)) {
        return data.bookings
      }
      if (Array.isArray(data.data)) {
        return data.data
      }
    }
    return []
  }

  // Getters
  const filteredBookings = computed(() => {
    let filtered = Array.isArray(bookings.value) ? bookings.value : []

    // Status filter
    if (filters.value.status) {
      filtered = filtered.filter(booking => 
        booking?.status === filters.value.status
      )
    }

    // Date range filter
    if (filters.value.dateFrom) {
      filtered = filtered.filter(booking => {
        const bookingDate = booking?.schedule?.date ? new Date(booking.schedule.date) : null
        return bookingDate && bookingDate >= new Date(filters.value.dateFrom)
      })
    }

    if (filters.value.dateTo) {
      filtered = filtered.filter(booking => {
        const bookingDate = booking?.schedule?.date ? new Date(booking.schedule.date) : null
        return bookingDate && bookingDate <= new Date(filters.value.dateTo)
      })
    }

    // Teacher filter
    if (filters.value.teacherId) {
      filtered = filtered.filter(booking => 
        booking?.schedule?.teacher?.id === filters.value.teacherId
      )
    }

    return filtered
  })

  const upcomingBookings = computed(() => {
    const now = new Date()
    const allBookings = Array.isArray(bookings.value) ? bookings.value : []

    return allBookings.filter(booking => {
      if (!booking?.schedule?.date) return false
      const bookingDate = new Date(booking.schedule.date)
      return bookingDate > now && booking.status === BOOKING_STATUS.PAID
    }).sort((a, b) => new Date(a.schedule.date) - new Date(b.schedule.date))
  })

  const pastBookings = computed(() => {
    const now = new Date()
    const allBookings = Array.isArray(bookings.value) ? bookings.value : []

    return allBookings.filter(booking => {
      if (!booking?.schedule?.date) return false
      // const bookingDate = new Date(booking.schedule.date)
      // return bookingDate <= now || booking.status === BOOKING_STATUS.COMPLETED
      return booking.status === BOOKING_STATUS.COMPLETED
    }).sort((a, b) => new Date(b.schedule.date) - new Date(a.schedule.date))
  })

  const pendingBookings = computed(() => {
    const allBookings = Array.isArray(bookings.value) ? bookings.value : []
    return allBookings.filter(booking => 
      booking?.status === BOOKING_STATUS.PENDING
    )
  })

  const cancelledBookings = computed(() => {
    const allBookings = Array.isArray(bookings.value) ? bookings.value : []
    return allBookings.filter(booking => 
      booking?.status === BOOKING_STATUS.CANCELLED
    )
  })

  const bookingStatusOptions = computed(() => {
    return Object.values(BOOKING_STATUS).map(status => ({
      value: status,
      label: status.charAt(0).toUpperCase() + status.slice(1)
    }))
  })

  const isLoading = computed(() => loading.value)

  function derivePagination({page, limit, itemsLength, incoming}) {
    const shapes = [
      incoming?.pagination,
      incoming?.meta?.pagination,
      incoming?.data?.pagination,
      incoming?.meta,
      incoming
    ].filter(Boolean)

    for (const p of shapes) {
      const cur = {
        page: Number(p.current_page ?? p.page ?? page ?? 1),
        limit: Number(p.per_page ?? p.limit ?? limit ?? 10),
        total: Number(p.total ?? p.total_data ?? p.totalCount ?? p.count ?? 0),
        totalPages: Number(p.total_pages ?? p.total_page ?? p.pages ?? 0),
      }
      if (cur.total || cur.totalPages) {
        if (!cur.total && cur.totalPages) cur.total = cur.totalPages * cur.limit
        if (!cur.totalPages && cur.total) cur.totalPages = Math.max(1, Math.ceil(cur.total / cur.limit))
        return cur
      }
    }
    const baseTotal = (page - 1) * limit + itemsLength + (itemsLength === limit ? limit : 0)
    return {
      page,
      limit,
      total: baseTotal,
      totalPages: Math.max(1, Math.ceil(baseTotal / limit)),
    }
  }

  // Actions
  const fetchBookings = async (params = {}) => {
    try {
      loading.value = true
      const response = await bookingService.getBookings({
        page: pagination.value.page,
        limit: pagination.value.limit,
        ...params
      })

      bookings.value = normalizeBookings(response.bookings || response.data || response)
      
      if (response.pagination) {
        pagination.value = { ...pagination.value, ...response.pagination }
      }

      return {
        success: true,
        data: bookings.value
      }
    } catch (error) {
      console.error('Error fetching bookings:', error)
      bookings.value = []
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchBookingsByUser = async (userId) => {
    try {
      loading.value = true
      const response = await bookingService.getBookingsByUser(userId)
      
      bookings.value = normalizeBookings(response.bookings || response.data || response)

      if (response.data?.pagination) {
        pagination.value = { 
          page: response.data.pagination.current_page, 
          limit: response.data.pagination.limit, 
          total: response.data.pagination.total ,
          totalPages: response.data.pagination.total_pages
         }
      }

      return {
        success: true,
        data: bookings.value
      }
    } catch (error) {
      console.error('Error fetching user bookings:', error)
      bookings.value = []
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchBooking = async (id) => {
    try {
      loading.value = true
      const response = await bookingService.getBooking(id)
      
      currentBooking.value = response.booking || response.data || response
      
      return {
        success: true,
        data: currentBooking.value
      }
    } catch (error) {
      console.error('Error fetching booking:', error)
      currentBooking.value = null
      throw error
    } finally {
      loading.value = false
    }
  }


   const bookingDetail = async (id) => {
    try {
      loading.value = true
      const response = await bookingService.getBookingDetail(id)
      
      currentBooking.value = response.booking || response.data || response
      return {
        success: true,
        data: currentBooking.value
      }
    } catch (error) {
      console.error('Error fetching booking:', error)
      currentBooking.value = null
      throw error
    } finally {
      loading.value = false
    }
  }


  const createBooking = async (bookingData) => {
    try {
      loading.value = true
      const response = await bookingService.createBooking(bookingData)
      
      const newBooking = response.booking || response.data || response
      if (newBooking && typeof newBooking === 'object') {
        bookings.value.unshift(newBooking)
      }
      
      return {
        success: true,
        data: newBooking,
        message: SUCCESS_MESSAGES.BOOKING_SUCCESS
      }
    } catch (error) {
      console.error('Error creating booking:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const rescheduleBooking = async (id, rescheduleData) => {
    try {
      loading.value = true
      const response = await bookingService.rescheduleBooking(id, rescheduleData)
      
      const updatedBooking = response.booking || response.data || response
      
      const index = bookings.value.findIndex(b => b?.id === id)
      if (index !== -1 && updatedBooking && typeof updatedBooking === 'object') {
        bookings.value[index] = updatedBooking
      }
      
      if (currentBooking.value?.id === id) {  
        currentBooking.value = updatedBooking
      }
      
      return {
        success: true,
        data: updatedBooking,
        message: SUCCESS_MESSAGES.BOOKING_RESCHEDULED
      }
    } catch (error) {
      console.error('Error rescheduling booking:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const cancelBooking = async (id) => {
    try {
      loading.value = true
      const response = await bookingService.cancelBooking(id)
      
      const updatedBooking = response.booking || response.data || response
      
      const index = bookings.value.findIndex(b => b?.id === id)
      if (index !== -1) {
        bookings.value[index] = { ...bookings.value[index], status: BOOKING_STATUS.CANCELLED }
      }
      
      if (currentBooking.value?.id === id) {
        currentBooking.value = { ...currentBooking.value, status: BOOKING_STATUS.CANCELLED }
      }
      
      return {
        success: true,
        data: updatedBooking,
        message: SUCCESS_MESSAGES.BOOKING_CANCELLED
      }
    } catch (error) {
      console.error('Error cancelling booking:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const updateBookingStatus = async (id, status) => {
    try {
      loading.value = true
      const response = await bookingService.updateBookingStatus(id, { status })
      
      const updatedBooking = response.booking || response.data || response
      
      const index = bookings.value.findIndex(b => b?.id === id)
      if (index !== -1) {
        bookings.value[index] = { ...bookings.value[index], status }
      }
      
      if (currentBooking.value?.id === id) {
        currentBooking.value = { ...currentBooking.value, status }
      }
      
      return {
        success: true,
        data: updatedBooking,
        message: 'Booking status updated successfully'
      }
    } catch (error) {
      console.error('Error updating booking status:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const changeStatus = async (id, statusData) => {
    try {
      loading.value = true
      const response = await bookingService.changeStatus(id, statusData)
      
      return {
        success: true,
        data: response,
        message: 'Booking status changed successfully'
      }
    } catch (error) {
      console.error('Error changing booking status:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const deleteBooking = async (id) => {
    try {
      loading.value = true
      const response = await bookingService.deleteBooking(id)
      return response
    } catch (error) {
      console.error('Error deleting booking:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
     

  const uploadBookingImage = async (file) => {
    try {
      loading.value = true
      const response = await bookingService.uploadBookingImage(file)
      
      return {
        success: true,
        data: response,
        message: 'Image uploaded successfully'
      }
    } catch (error) {
      console.error('Error uploading booking image:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // Filter actions
  const setStatusFilter = (status) => {
    filters.value.status = status
  }

  const setDateRangeFilter = (dateFrom, dateTo) => {
    filters.value.dateFrom = dateFrom
    filters.value.dateTo = dateTo
  }

  const setTeacherFilter = (teacherId) => {
    filters.value.teacherId = teacherId
  }

  const clearFilters = () => {
    filters.value = {
      status: '',
      dateFrom: null,
      dateTo: null,
      teacherId: null
    }
  }

  // Utility actions
  const getBookingById = (id) => {
    return bookings.value.find(booking => booking?.id === id)
  }

  const canCancelBooking = (booking) => {
    if (!booking || !booking.schedule?.date || !booking.schedule?.start_time) return false
    
    const bookingDateTime = new Date(`${booking.schedule.date} ${booking.schedule.start_time}`)
    const now = new Date()
    const hoursUntilBooking = (bookingDateTime - now) / (1000 * 60 * 60)
    
    return hoursUntilBooking >= 24 && booking.status === BOOKING_STATUS.PAID
  }

  const canRescheduleBooking = (booking) => {
    if (!booking || !booking.schedule?.date || !booking.schedule?.start_time) return false
    
    const bookingDateTime = new Date(`${booking.schedule.date} ${booking.schedule.start_time}`)
    const now = new Date()
    const hoursUntilBooking = (bookingDateTime - now) / (1000 * 60 * 60)
    
    return hoursUntilBooking >= 24 && booking.status === BOOKING_STATUS.PAID
  }

  return {
    // State
    bookings: computed(() => bookings.value),
    currentBooking: computed(() => currentBooking.value),
    filters: computed(() => filters.value),
    pagination: computed(() => pagination.value),
    
    // Getters
    filteredBookings,
    upcomingBookings,
    pastBookings,
    pendingBookings,
    cancelledBookings,
    bookingStatusOptions,
    isLoading,
    
    // Actions
    fetchBookings,
    fetchBookingsByUser,
    fetchBooking,
    createBooking,
    rescheduleBooking,
    cancelBooking,
    updateBookingStatus,
    uploadBookingImage,
    bookingDetail,
    changeStatus,
    deleteBooking,
    
    // Filter actions
    setStatusFilter,
    setDateRangeFilter,
    setTeacherFilter,
    clearFilters,
    
    // Utility actions
    getBookingById,
    canCancelBooking,
    canRescheduleBooking
  }
})
