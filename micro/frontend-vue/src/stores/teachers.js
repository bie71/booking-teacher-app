import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { teacherService } from '@/services/api'
import { TEACHER_LEVELS } from '@/config'

export const useTeachersStore = defineStore('teachers', () => {
  // State
  const teachers = ref([])
  const currentTeacher = ref(null)
  const schedules = ref([])
  const loading = ref(false)
  const filters = ref({
    search: '',
    level: '',
    minPrice: null,
    maxPrice: null,
    availability: null
  })
  const pagination = ref({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0
  })

  // Getters
  const filteredTeachers = computed(() => {
    let filtered = teachers.value

    // Search filter
    if (filters.value.search) {
      const searchTerm = filters.value.search.toLowerCase()
      filtered = filtered.filter(teacher => {
        const name = (teacher.Name || teacher.name || '').toLowerCase()
        const bio = (teacher.Bio || teacher.bio || '').toLowerCase()
        return name.includes(searchTerm) || bio.includes(searchTerm)
      })
    }

    // Level filter
    if (filters.value.level) {
      filtered = filtered.filter(teacher => {
        const level = teacher.LanguageLevel || teacher.language_level
        return level === filters.value.level
      })
    }

    // Price range filter
    if (filters.value.minPrice !== null) {
      filtered = filtered.filter(teacher => {
        const price = teacher.PricePerHour || teacher.price_per_hour || 0
        return price >= filters.value.minPrice
      })
    }

    if (filters.value.maxPrice !== null) {
      filtered = filtered.filter(teacher => {
        const price = teacher.PricePerHour || teacher.price_per_hour || 0
        return price <= filters.value.maxPrice
      })
    }

    return filtered
  })

  const teacherLevels = computed(() => {
    return Object.values(TEACHER_LEVELS).map(level => ({
      value: level,
      label: level.charAt(0).toUpperCase() + level.slice(1)
    }))
  })

  const priceRange = computed(() => {
    if (teachers.value.length === 0) return { min: 0, max: 100 }
    
    const prices = teachers.value.map(t => t.PricePerHour || t.price_per_hour || 0).filter(p => p > 0)
    if (prices.length === 0) return { min: 0, max: 100 }
    
    return {
      min: Math.min(...prices),
      max: Math.max(...prices)
    }
  })

  const isLoading = computed(() => loading.value)

  // Fallback pagination derivation when API doesn't include meta/pagination
  function derivePagination({page, limit, itemsLength, incoming}) {
    // Try common shapes
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
    // Heuristic: if itemsLength == limit, assume there might be another page
    const baseTotal = (page - 1) * limit + itemsLength + (itemsLength === limit ? limit : 0)
    return {
      page,
      limit,
      total: baseTotal,
      totalPages: Math.max(1, Math.ceil(baseTotal / limit)),
    }
  }

  // Actions
  const fetchTeachers = async (params = {}) => {
    try {
      loading.value = true
    const requestedPage = params.page ?? pagination.value.page
    const requestedLimit = params.limit ?? pagination.value.limit


      const response = await teacherService.getTeachers({
         page: requestedPage,
         limit: requestedLimit,
        ...params
      })

      // Handle the response structure from your backend
      teachers.value = response.data || response.teachers || response
      
      if (response.pagination) {
        pagination.value = { 
          page: response.pagination.current_page || pagination.value.page,
          limit: response.pagination.limit || pagination.value.limit,
          total: response.pagination.total_data || pagination.value.total,
          totalPages: response.pagination.total_page
        }
      }

      return {
        success: true,
        data: teachers.value
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchTeacher = async (id) => {
    try {
      loading.value = true
      const response = await teacherService.getTeacher(id)
      
      currentTeacher.value = response.teacher || response.data || response
      
      return {
        success: true,
        data: currentTeacher.value
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchTeacherSchedules = async (teacherId, params = {}) => {
    try {
      loading.value = true
      const response = await teacherService.getTeacherSchedules(teacherId, params)
      
      schedules.value = response.schedules || response.data || response
      
      return {
        success: true,
        data: schedules.value
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const createTeacher = async (teacherData) => {
    try {
      loading.value = true
      const response = await teacherService.createTeacher(teacherData)
      
      const newTeacher = response.teacher || response.data || response
      teachers.value.push(newTeacher)
      
      return {
        success: true,
        data: newTeacher,
        message: 'Teacher created successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const updateTeacher = async (id, teacherData) => {
    try {
      loading.value = true
      const response = await teacherService.updateTeacher(id, teacherData)
      
      const updatedTeacher = response.teacher || response.data || response
      
      // Update in teachers array
      const index = teachers.value.findIndex(t => t.id === id)
      if (index !== -1) {
        teachers.value[index] = updatedTeacher
      }
      
      // Update current teacher if it's the same
      if (currentTeacher.value?.id === id) {
        currentTeacher.value = updatedTeacher
      }
      
      return {
        success: true,
        data: updatedTeacher,
        message: 'Teacher updated successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const deleteTeacher = async (id) => {
    try {
      loading.value = true
      await teacherService.deleteTeacher(id)
      
      // Remove from teachers array
      teachers.value = teachers.value.filter(t => t.id !== id)
      
      // Clear current teacher if it's the same
      if (currentTeacher.value?.id === id) {
        currentTeacher.value = null
      }
      
      return {
        success: true,
        message: 'Teacher deleted successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const createSchedule = async (scheduleData) => {
    try {
      loading.value = true
      const response = await teacherService.createSchedule(scheduleData)
      
      const newSchedule = response.schedule || response.data || response
      schedules.value.push(newSchedule)
      
      return {
        success: true,
        data: newSchedule,
        message: 'Schedule created successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const updateSchedule = async (id, scheduleData) => {
    try {
      loading.value = true
      const response = await teacherService.updateSchedule(id, scheduleData)
      
      const updatedSchedule = response.schedule || response.data || response
      
      // Update in schedules array
      const index = schedules.value.findIndex(s => s.id === id)
      if (index !== -1) {
        schedules.value[index] = updatedSchedule
      }
      
      return {
        success: true,
        data: updatedSchedule,
        message: 'Schedule updated successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const cancelSchedule = async (id) => {
    try {
      loading.value = true
      await teacherService.cancelSchedule(id)
      
      // Update schedule status in array
      const index = schedules.value.findIndex(s => s.id === id)
      if (index !== -1) {
        schedules.value[index].status = 'cancelled'
      }
      
      return {
        success: true,
        message: 'Schedule cancelled successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const uploadTeacherImage = async (file) => {
    try {
      loading.value = true
      const response = await teacherService.uploadTeacherImage(file)
      
      return {
        success: true,
        data: response,
        message: 'Image uploaded successfully'
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // Filter actions
  const setSearchFilter = (search) => {
    filters.value.search = search
  }

  const setLevelFilter = (level) => {
    filters.value.level = level
  }

  const setPriceFilter = (minPrice, maxPrice) => {
    filters.value.minPrice = minPrice
    filters.value.maxPrice = maxPrice
  }

  const clearFilters = () => {
    filters.value = {
      search: '',
      level: '',
      minPrice: null,
      maxPrice: null,
      availability: null
    }
  }

  // Pagination actions
  const setPage = (page) => {
    pagination.value.page = page
  }

  const setLimit = (limit) => {
    pagination.value.limit = limit
  }

  const resetPagination = () => {
    pagination.value.page = 1
  }

  return {
    // State
    teachers: computed(() => teachers.value),
    currentTeacher: computed(() => currentTeacher.value),
    schedules: computed(() => schedules.value),
    filters: computed(() => filters.value),
    pagination: computed(() => pagination.value),
    
    // Getters
    filteredTeachers,
    teacherLevels,
    priceRange,
    isLoading,
    
    // Actions
    fetchTeachers,
    fetchTeacher,
    fetchTeacherSchedules,
    createTeacher,
    updateTeacher,
    deleteTeacher,
    createSchedule,
    updateSchedule,
    cancelSchedule,
    uploadTeacherImage,
    
    // Filter actions
    setSearchFilter,
    setLevelFilter,
    setPriceFilter,
    clearFilters,
    
    // Pagination actions
    setPage,
    setLimit,
    resetPagination
  }
})
