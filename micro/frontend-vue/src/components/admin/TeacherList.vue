<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900 dark:text-white">Teacher Management</h2>
      <!-- <button 
        @click="$emit('add-teacher')"
        class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
      >
        Add New Teacher
      </button> -->
    </div>

    <!-- Search and Filters -->
    <div class="mb-4 flex space-x-4">
      <div class="flex-1">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search teachers..."
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
        />
      </div>
      <select
        v-model="selectedLevel"
        class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
      >
        <option value="">All Levels</option>
        <option value="beginner">Beginner</option>
        <option value="intermediate">Intermediate</option>
        <option value="advanced">Advanced</option>
      </select>
    </div>

    <!-- Teachers Table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
              Teacher
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
              Level
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
              Price/Hour
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
              Available Schedules
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
          <tr v-for="teacher in paginatedTeachers" :key="teacher.id || teacher.ID" class="hover:bg-gray-50 dark:hover:bg-gray-700">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="flex-shrink-0 h-10 w-10">
                  <img 
                    v-if="teacher.ProfileImage || teacher.profile_image"
                    :src="teacher.ProfileImage || teacher.profile_image" 
                    :alt="teacher.Name || teacher.name"
                    class="h-10 w-10 rounded-full"
                  />
                  <div v-else class="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
                    <span class="text-sm font-medium text-gray-500">
                      {{ (teacher.Name || teacher.name || '').charAt(0).toUpperCase() }}
                    </span>
                  </div>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ teacher.Name || teacher.name }}
                  </div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">
                    {{ teacher.Email || teacher.email }}
                  </div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="getLevelBadgeClass(teacher.LanguageLevel || teacher.language_level)">
                {{ (teacher.LanguageLevel || teacher.language_level)?.toUpperCase() }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
              <span class="text-green-600 font-semibold">
                ${{ teacher.PricePerHour || teacher.price_per_hour }}/hour
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ teacher.available_start_time }} - {{ teacher.available_end_time }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium flex items-center space-x-2">
              <button 
                @click="$emit('edit-teacher', teacher)"
                class="text-blue-600 hover:text-blue-900 p-1"
                title="Edit"
              >
                <!-- Pencil icon -->
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536M9 11l3 3L19.071 7.071a2 2 0 10-2.828-2.828L9 11zm0 0L5 15l-.707.707a1 1 0 00-.293.707V19a1 1 0 001 1h3.586a1 1 0 00.707-.293L15 13"></path></svg>
              </button>
              <button 
                @click="openScheduleModal(teacher)"
                class="text-green-600 hover:text-green-900 p-1"
                title="Manage Schedule"
              >
                <!-- Calendar icon -->
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
              </button>
              <button 
                @click="deleteTeacher(teacher.id || teacher.ID)"
                class="text-red-600 hover:text-red-900 p-1"
                title="Delete"
              >
                <!-- Trash icon -->
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6h18M8 6v12a2 2 0 002 2h4a2 2 0 002-2V6m-9 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path></svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination Controls -->
    <!--
      Show pagination only when there is more than one page.  When using
      server‑side pagination, the total number of pages is derived from
      the pagination metadata.  When falling back to client‑side
      pagination (no total returned), totalPages is computed from the
      length of the filtered teachers array.  The condition below uses
      totalPages > 1 rather than comparing total and limit to ensure
      pagination appears whenever multiple pages exist.
    -->
    <div v-if="totalPages > 0" >
      <!-- <div class="flex items-center space-x-2">
        <button 
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage <= 1"
          class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>
        <button 
          v-for="page in visiblePages" 
          :key="page"
          @click="goToPage(page)"
          :class="['px-3 py-1 border rounded', page === currentPage ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100']"
        >
          {{ page }}
        </button>
        <button 
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div> -->
      <!-- <div class="text-sm text-gray-600 dark:text-gray-400">
        <template v-if="pagination.total">
          Showing {{ (currentPage - 1) * pagination.limit + 1 }} to {{ Math.min(currentPage * pagination.limit, pagination.total) }} of {{ pagination.total }} teachers
        
        <Pagination
          :page="currentPage"
          :limit="pagination.limit"
          :total="pagination.total"
          @update:page="goToPage"
        />
      </template>

        <template v-else>
          Showing {{ (currentPage - 1) * pagination.limit + 1 }} to {{ Math.min(currentPage * pagination.limit, filteredTeachers.length) }} of {{ filteredTeachers.length }} teachers
        
        <Pagination
          :page="currentPage"
          :limit="pagination.limit"
          :total="pagination.total"
          @update:page="goToPage"
        />
      </template>

      </div> -->
      <Pagination
      :page="currentPage"
      :limit="pagination.limit"
      :total="pagination.total"
      @update:page="goToPage"
    />
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-4">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && filteredTeachers.length === 0" class="text-center py-8">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No teachers found</h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        Get started by creating a new teacher.
      </p>
    </div>

    <!-- Schedule Modal -->
    <ScheduleTeacherModal 
      v-if="showScheduleModal"
      :teacher="selectedTeacher"
      @close="closeScheduleModal"
    />
  </div>

  
</template>


<script>
import { useTeachersStore } from '@/stores/teachers'
import { useUIStore } from '@/stores/ui'
import ScheduleTeacherModal from '@/modals/ScheduleTeacherModal.vue'
import { userService } from '@/services/api'
import Pagination from '@/components/ui/Pagination.vue'

export default {
  name: 'TeacherList',
  components: {
    ScheduleTeacherModal, Pagination
  },
  emits: ['edit-teacher', 'add-teacher', 'manageSchedule'],
  data() {
    return {
      searchQuery: '',
      selectedLevel: '',
      selectedTeacher: null,
      showScheduleModal: false,
      teachersStore: null,
      uiStore: null
    }
  },
  computed: {
    teachers() {
      return this.teachersStore?.teachers || []
    },
    filteredTeachers() {
      let filtered = this.teachers
      
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        filtered = filtered.filter(teacher => 
          (teacher.Name || teacher.name)?.toLowerCase().includes(query) ||
          (teacher.Email || teacher.email)?.toLowerCase().includes(query)
        )
      }
      
      if (this.selectedLevel) {
        filtered = filtered.filter(teacher => 
          (teacher.LanguageLevel || teacher.language_level) === this.selectedLevel
        )
      }
      
      return filtered
    },
    loading() {
      return this.teachersStore?.isLoading || false
    },

    /**
     * Pagination information from the teachers store. If the backend
     * provides pagination metadata, we use it. Otherwise we fall back
     * to client‑side calculation based on the filtered teachers array.
     */
    pagination() {
      return this.teachersStore?.pagination || { page: 1, limit: 10, total: this.filteredTeachers.length }
    },

    /**
     * Current page number. Defaults to 1.
     */
    currentPage() {
      return this.pagination.page || 1
    },

    /**
     * Total number of pages based on total items and limit. If the
     * backend does not return total, we compute it from the length of
     * filteredTeachers.
     */
    totalPages() {
      const totalItems = this.pagination.total || this.filteredTeachers.length
      const perPage = this.pagination.limit || 10
      const pages = Math.ceil(totalItems / perPage)
      return pages > 0 ? pages : 1
    },

    /**
     * Slice of teachers to display on the current page when the backend
     * does not handle pagination. If backend pagination is active
     * (pagination.total is defined), we simply return the filteredTeachers
     * array (assuming the backend already paginated the response).
     */
    paginatedTeachers() {
      if (this.pagination.total) {
        return this.filteredTeachers
      }
      const perPage = this.pagination.limit || 10
      const start = (this.currentPage - 1) * perPage
      const end = start + perPage
      return this.filteredTeachers.slice(start, end)
    },

    /**
     * Compute a list of visible page numbers around the current page.
     * This helps to create a compact pagination control.
     */
    visiblePages() {
      const pages = []
      const current = this.currentPage
      const total = this.totalPages
      const start = Math.max(1, current - 2)
      const end = Math.min(total, current + 2)
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      return pages
    }
  },
  async mounted() {
    this.teachersStore = useTeachersStore()
    this.uiStore = useUIStore()
    
    // Load teachers
    await this.teachersStore.fetchTeachers()
  },
  methods: {
    getLevelBadgeClass(level) {
      const classes = {
        'beginner': 'bg-blue-100 text-blue-800',
        'intermediate': 'bg-yellow-100 text-yellow-800',
        'advanced': 'bg-red-100 text-red-800'
      }
      return classes[level] || 'bg-gray-100 text-gray-800'
    },
    openScheduleModal(teacher) {
      this.selectedTeacher = teacher
      this.showScheduleModal = true
    },
    closeScheduleModal() {
      this.showScheduleModal = false
      this.selectedTeacher = null
    },
    async deleteTeacher(teacherId) {
      if (confirm('Are you sure you want to delete this teacher?')) {
        try {
          await this.teachersStore.deleteTeacher(teacherId)
          this.uiStore.showSuccess('Teacher deleted successfully')
          // Log admin delete teacher activity
          try {
            await userService.logActivity({
              action: 'AdminDeleteTeacher',
              description: `Deleted teacher with ID ${teacherId}`
            })
          } catch (error) {
            console.warn('Failed to log admin delete teacher activity:', error)
          }
          await this.teachersStore.fetchTeachers()
        } catch (error) {
          this.uiStore.showError('Failed to delete teacher')
        }
      }
    },

    /**
     * Navigate to a specific page in the teacher list. If the backend
     * provides pagination metadata (via this.pagination.total), we call
     * the teachersStore to fetch the requested page. Otherwise, we
     * update the local pagination state in the store and allow the
     * computed property paginatedTeachers to slice the results. A
     * safeguard ensures the page number stays within valid bounds.
     *
     * @param {Number} page The page number to navigate to
     */
    async goToPage(page) {
      // Ensure the requested page is within range
      if (page < 1 || page > this.totalPages) {
        return
      }

      // Update the page in the teachers store
      this.teachersStore.setPage(page)

      // If the backend handles pagination (pagination.total is defined),
      // fetch the teachers for the new page. Otherwise, slicing will
      // occur in paginatedTeachers and no additional request is necessary.
      if (this.pagination.total) {
        await this.teachersStore.fetchTeachers()
      }
    }
  }
}
</script>
