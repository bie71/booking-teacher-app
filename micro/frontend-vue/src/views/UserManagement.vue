<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <!-- Header -->
    <div class="bg-white shadow dark:bg-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">User Management</h1>
          <div class="flex items-center space-x-4">
            <router-link 
              to="/admin/dashboard"
              class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
            >
              Back to Dashboard
            </router-link>
            <button 
              @click="showAddUser = true"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Add User
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filters -->
      <div class="bg-white rounded-lg shadow mb-6 dark:bg-gray-800 p-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Search</label>
            <input 
              v-model="searchQuery"
              type="text"
              placeholder="Search users..."
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Role</label>
            <select 
              v-model="selectedRole"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            >
              <option value="">All Roles</option>
              <option value="admin">Admin</option>
              <option value="teacher">Teacher</option>
              <option value="user">User</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Per Page</label>
            <select 
              v-model="perPage"
              @change="loadUsers"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            >
              <option value="10">10</option>
              <option value="25">25</option>
              <option value="50">50</option>
              <option value="100">100</option>
            </select>
          </div>
          <div class="flex items-end justify-between">
            <button 
              @click="loadUsers"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Apply Filters
            </button>
            <button 
              @click="clearFilters"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>

      <!-- Users Table -->
      <div class="bg-white rounded-lg shadow dark:bg-gray-800">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Profile Picture
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Name
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Email
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Role
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Created At
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
              <tr v-for="user in users" :key="user.id" class="hover:bg-gray-100 dark:hover:bg-gray-700">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                  <div class="w-10 h-10 rounded-xl overflow-hidden flex items-center justify-center bg-blue-600">
                    <template v-if="user?.profile_image">
                      <img 
                        :src="user.profile_image" 
                        alt="User Avatar" 
                        class="w-full h-full object-cover"
                      />
                    
  
                      </template>

                    <template v-else>
                      <span class="text-white text-sm font-medium">
                        {{ user?.name?.charAt(0).toUpperCase() }}
                      </span>
                    
 
                      </template>

                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                  {{ user.name }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                  {{ user.email }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getRoleBadgeClass(user.role)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                    {{ user.role }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                  {{ formatDate(user.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium flex items-center space-x-2">
                  <!-- Edit button with pencil icon -->
                  <button 
                    @click="editUser(user)"
                    class="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 p-1"
                    title="Edit"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536M9 11l3 3L19.071 7.071a2 2 0 10-2.828-2.828L9 11zm0 0L5 15l-.707.707a1 1 0 00-.293.707V19a1 1 0 001 1h3.586a1 1 0 00.707-.293L15 13" />
                    </svg>
                  </button>
                  <!-- Delete button with trash icon -->
                  <button 
                    @click="deleteUser(user)"
                    class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 p-1"
                    title="Delete"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6h18M8 6v12a2 2 0 002 2h4a2 2 0 002-2V6m-9 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
                    </svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        
        <!-- <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 dark:bg-gray-800 dark:border-gray-700 sm:px-6">
          <div class="flex-1 flex justify-between sm:hidden">
            <button 
              @click="previousPage"
              :disabled="currentPage === 1"
              class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              Previous
            </button>
            <button 
              @click="nextPage"
              :disabled="currentPage === totalPages"
              class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              Next
            </button>
          </div>
          <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                Showing {{ (currentPage - 1) * perPage + 1 }} to {{ Math.min(currentPage * perPage, totalUsers) }} of {{ totalUsers }} results
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                <button 
                  @click="previousPage"
                  :disabled="currentPage === 1"
                  class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                >
                  Previous
                </button>
                <button 
                  v-for="page in visiblePages"
                  :key="page"
                  @click="goToPage(page)"
                  :class="page === currentPage ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
                  class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                >
                  {{ currentPage }}
                </button>
                <button 
                  @click="nextPage"
                  :disabled="currentPage === totalPages"
                  class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                >
                  Next
                </button>
              </nav>
            </div>
          </div>
        </div> -->
        <div v-if="totalPages > 0">
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
              Showing {{ (currentPage - 1) * pagination.limit + 1 }} to {{ Math.min(currentPage * pagination.limit, pagination.total) }} of {{ pagination.total }} users
            
            <Pagination
              :page="currentPage"
              :limit="pagination.limit"
              :total="pagination.total"
              @update:page="goToPage"
            />
          </template>

            <template v-else>
              Showing {{ (currentPage - 1) * pagination.limit + 1 }} to {{ Math.min(currentPage * pagination.limit, pagination.length) }} of {{ pagination.length }} users
            
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
              @update:page="(p)=>{ currentPage=p; loadUsers && loadUsers() }"
            />
    </div>
  </div>
</div>

    <!-- Add/Edit User Modal -->
    <UserForm 
      v-if="showAddUser || editingUser"
      :user="editingUser"
      :mode="editingUser ? 'edit' : 'create'"
      @close="closeModal"
      @user-saved="handleUserSaved"
    />
  </div>

  <!-- <Pagination
    :page="currentPage"
    :limit="perPage"
    :total="totalUsers"
    @update:page="(p)=>{ currentPage=p; loadUsers && loadUsers() }"
  /> -->
</template>


<script setup>
import Pagination from '@/components/ui/Pagination.vue'
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { userService } from '@/services/api'
import { useUIStore } from '@/stores/ui'
import UserForm from '@/components/admin/UserForm.vue'

const uiStore = useUIStore()
const router = useRouter()

const users = ref([])
const currentPage = ref(1)
const perPage = ref(10)
const totalUsers = ref(0)
const totalPages = ref(0)
const searchQuery = ref('')
const selectedRole = ref('')
const showAddUser = ref(false)
const editingUser = ref(null)
const pagination = ref({
  page: 1,
  limit: 10,
  total: 0
})

const loadUsers = async () => {
  try {
    const response = await userService.getUsers({
      page: currentPage.value,
      limit: perPage.value,
      search: searchQuery.value,
      role: selectedRole.value
    })
    users.value = response.data.data
    pagination.value = response.data.pagination
    totalUsers.value = response.data.pagination.total
    totalPages.value = Math.ceil(totalUsers.value / perPage.value)
  } catch (error) {
    console.error('Failed to load users:', error)
    uiStore.showError('Failed to load users')
  }
}

const clearFilters = () => {
  searchQuery.value = ''
  selectedRole.value = ''
  perPage.value = 10
  loadUsers()
}

const editUser = (user) => {
  editingUser.value = user
}

const deleteUser = async (user) => {
  if (confirm(`Are you sure you want to delete ${user.name}?`)) {
    try {
      await userService.deleteUser(user.id)
      uiStore.showSuccess('User deleted successfully')
      // Log admin deletion activity for the current user
      try {
        await userService.logActivity({
          action: 'AdminDeleteUser',
          description: `Deleted user ${user.name}`
        })
      } catch (error) {
        console.warn('Failed to log admin delete user activity:', error)
      }
      loadUsers()
    } catch (error) {
      console.error('Failed to delete user:', error)
      uiStore.showError('Failed to delete user')
    }
  }
}

const closeModal = () => {
  showAddUser.value = false
  editingUser.value = null
}

const handleUserSaved = () => {
  closeModal()
  loadUsers()
}

const getRoleBadgeClass = (role) => {
  const classes = {
    admin: 'bg-red-100 text-red-800',
    teacher: 'bg-blue-100 text-blue-800',
    user: 'bg-green-100 text-green-800'
  }
  return classes[role] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadUsers()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadUsers()
  }
}

const goToPage = (page) => {
  currentPage.value = page
  loadUsers()
}

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

onMounted(() => {
  loadUsers()
})
</script>
