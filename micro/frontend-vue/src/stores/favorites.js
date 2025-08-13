import { defineStore } from 'pinia'
import { ref } from 'vue'
import { userService, teacherService } from '@/services/api'

// The favorites store manages the list of teacher IDs that the
// authenticated user has marked as favorites and optionally caches
// detailed teacher objects for convenience. It provides actions to
// fetch favorites from the backend and to toggle a favorite state for
// a given teacher.
export const useFavoritesStore = defineStore('favorites', () => {
  // State: list of favorite teacher IDs and a map of detailed teacher objects
  const favoriteIds = ref([])
  const favoriteTeachers = ref([])
  const isLoading = ref(false)

  // Fetch favorite teacher IDs from the backend and optionally load
  // detailed teacher data. If loadDetails is true, teacher details
  // will be fetched from the teacher service for each ID. This is
  // useful when displaying a favorites list in the UI.
  const fetchFavorites = async (loadDetails = false) => {
    try {
      isLoading.value = true
      const response = await userService.getFavoriteTeachers()
      favoriteIds.value = response.data || []
      if (loadDetails) {
        await loadFavoriteTeacherDetails()
      }
    } catch (error) {
      console.error('Failed to fetch favorites:', error)
      favoriteIds.value = []
      favoriteTeachers.value = []
    } finally {
      isLoading.value = false
    }
  }

  // Load detailed teacher objects for each favorite ID. This helper
  // makes individual requests to the teacher service. It replaces the
  // favoriteTeachers array with the results.
  const loadFavoriteTeacherDetails = async () => {
    favoriteTeachers.value = []
    const ids = favoriteIds.value || []
    for (const id of ids) {
      try {
        const res = await teacherService.getTeacher(id)
        const teacher = res.teacher || res.data || res
        favoriteTeachers.value.push(teacher)
      } catch (error) {
        console.warn(`Failed to fetch teacher ${id} details`, error)
      }
    }
  }

  // Toggle a teacher's favorite status. If the teacher is already
  // favorited, this will remove it; otherwise it will add it. After
  // toggling, the favorites list is refreshed.
  const toggleFavorite = async (teacherId) => {
    const isFavorite = favoriteIds.value.includes(teacherId)
    try {
      await userService.toggleFavoriteTeacher({ teacher_id: teacherId, favorite: !isFavorite })
      // Refresh the favorites list after the update
      await fetchFavorites()
    } catch (error) {
      console.error('Failed to toggle favorite:', error)
    }
  }

  return {
    favoriteIds,
    favoriteTeachers,
    isLoading,
    fetchFavorites,
    toggleFavorite,
    loadFavoriteTeacherDetails
  }
})