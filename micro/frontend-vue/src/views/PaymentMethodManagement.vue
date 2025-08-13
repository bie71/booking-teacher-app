<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <!-- Header -->
    <div class="bg-white shadow dark:bg-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Payment Method Management</h1>
          <div class="flex items-center space-x-4">
            <router-link 
              to="/admin/dashboard"
              class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
            >
              Back to Dashboard
            </router-link>
            <button 
              @click="openAddModal"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 "
            >
              Add Payment Method
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white p-6 rounded-lg shadow dark:bg-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Total Methods</h3>
          <p class="text-3xl font-bold text-blue-600">{{ pagination?.total }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow dark:bg-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Active Methods</h3>
          <p class="text-3xl font-bold text-green-600">{{ activeCount }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow dark:bg-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Inactive Methods</h3>
          <p class="text-3xl font-bold text-red-600">{{ inactiveCount }}</p>
        </div>
      </div>

      <!-- Search and Filter -->
      <div class="bg-white p-6 rounded-lg shadow mb-6 dark:bg-gray-800">
        <div class="flex flex-col md:flex-row gap-4">
          <div class="flex-1">
            <input 
              v-model="searchQuery"
              type="text"
              placeholder="Search payment methods..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <select 
            v-model="filterStatus"
            class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <!-- Payment Methods Table -->
      <div class="bg-white rounded-lg shadow overflow-hidden dark:bg-gray-800">
        <div v-if="loading" class="p-8 text-center">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p class="mt-2 text-gray-600 dark:text-gray-400">Loading payment methods...</p>
        </div>

        <div v-else-if="error" class="p-8 text-center">
          <p class="text-red-600 dark:text-red-400">{{ error }}</p>
          <button 
            @click="loadPaymentMethods"
            class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Retry
          </button>
        </div>

        <div v-else-if="filteredPaymentMethods.length === 0" class="p-8 text-center">
          <p class="text-gray-600 dark:text-gray-400">No payment methods found</p>
        </div>

        <div v-else>
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-white uppercase tracking-wider">
                  Name
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-white uppercase tracking-wider">
                  Code
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-white uppercase tracking-wider">
                  Status
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-white uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200 dark:bg-gray-800">
              <!-- Loop through paginated payment methods rather than the full filtered list -->
              <tr v-for="method in responseData" :key="method.id" class="hover:bg-gray-50 dark:hover:bg-gray-700">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white flex items-center space-x-2">
                  <!-- Icon before payment method name -->
                  <svg v-if="['credit_card'].includes(method.code)" class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <rect x="2" y="5" width="20" height="14" rx="2" ry="2" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <line x1="2" y1="10" x2="22" y2="10" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <svg v-else-if="['bank_transfer'].includes(method.code)" class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M5 6h14a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2z" />
                  </svg>
                  <svg v-else-if="['gopay','ovo','dana','shopeepay','linkaja'].includes(method.code)" class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg v-else-if="['bca_va','bni_va','bri_va','mandiri_va','permata_va'].includes(method.code)" class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 21v-8a2 2 0 012-2h12a2 2 0 012 2v8" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v7" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 9h14" />
                  </svg>
                  <svg v-else-if="['indomaret','alfamart'].includes(method.code)" class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 21V7a2 2 0 012-2h14a2 2 0 012 2v14" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 3v4H8V3" />
                  </svg>
                  <!-- Default icon if none matched -->
                  <svg v-else class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <rect x="2" y="5" width="20" height="14" rx="2" ry="2" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <line x1="2" y1="10" x2="22" y2="10" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <span>{{ method.name }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                  {{ method.code }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="method.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                  >
                    {{ method.active ? 'Active' : 'Inactive' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium flex items-center space-x-2">
                  <!-- Edit icon -->
                  <button
                    @click="openEditModal(method)"
                    class="text-blue-600 hover:text-blue-900 p-1"
                    title="Edit"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536M9 11l3.536 3.536m-3.536-3.536L3 16.071V20h3.929l6.465-6.465m0 0l1.415-1.414a2 2 0 10-2.828-2.828l-1.415 1.414" />
                    </svg>
                  </button>
                  <!-- Delete icon -->
                  <button
                    @click="openDeleteModal(method)"
                    class="text-red-600 hover:text-red-900 p-1"
                    title="Delete"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                  <!-- Toggle status icon: if active show deactivate (cross), else show activate (check) -->
                    <button
                      @click="toggleStatus(method)"
                      :class="method.active ? 'text-orange-600 hover:text-orange-900' : 'text-green-600 hover:text-green-900'"
                      class="p-1"
                      :title="method.active ? 'Deactivate' : 'Activate'"
                    >
                      <svg v-if="method.active" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                      </svg>
                    </button>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Pagination controls for payment methods -->
          <div v-if="pagination?.totalPages > 0" >
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
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, totalItems) }} of {{ totalItems }} methods
            </div> -->
             <Pagination
            :page="pagination?.page || page || 1"
            :limit="pagination?.limit || limit || 10"
            :total="pagination?.total || total || 0"
              @update:page="goToPage"
              />
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50 dark:bg-gray-900 dark:bg-opacity-20 dark:text-white dark:shadow-sm">
      <div class="bg-white rounded-lg p-6 w-full max-w-md dark:bg-gray-800">
        <h2 class="text-xl font-bold mb-4">{{ modalMode === 'add' ? 'Add Payment Method' : 'Edit Payment Method' }}</h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-white">Name</label>
            <input 
              v-model="formData.name"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-white">Code</label>
            <input 
              v-model="formData.code"
              type="text"
              required
              :disabled="modalMode === 'edit'"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-white">Status</label>
            <select 
              v-model="formData.active"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            >
              <option :value="true">Active</option>
              <option :value="false">Inactive</option>
            </select>
          </div>
          <div class="flex justify-end space-x-3">
            <button 
              type="button"
              @click="closeModal"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 dark:text-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600"
            >
              Cancel
            </button>
            <button 
              type="submit"
              :disabled="loading"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 dark:bg-blue-500 dark:hover:bg-blue-600"
            >
              {{ loading ? 'Processing...' : (modalMode === 'add' ? 'Add' : 'Update') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Confirm Delete</h2>
        <p class="text-gray-700 mb-4">
          Are you sure you want to delete the payment method "{{ selectedMethod?.name }}"? This action cannot be undone.
        </p>
        <div class="flex justify-end space-x-3">
          <button 
            @click="closeDeleteModal"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
          >
            Cancel
          </button>
          <button 
            @click="deletePaymentMethod"
            :disabled="loading"
            class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50"
          >
            {{ loading ? 'Deleting...' : 'Delete' }}
          </button>
        </div>
      </div>
    </div>
  </div>

  
  
</template>


<script setup>
import Pagination from '@/components/ui/Pagination.vue'
import { ref, computed, onMounted, watch } from 'vue'
import { usePaymentsStore } from '@/stores/payments'

const paymentsStore = usePaymentsStore()

// State
const searchQuery = ref('')
const filterStatus = ref('')
const loading = ref(false)
const error = ref(null)
const showModal = ref(false)
const showDeleteModal = ref(false)
const modalMode = ref('add')
const selectedMethod = ref(null)
const formData = ref({
  name: '',
  code: '',
  active: true
})
const pagination = ref(null)
const responseData = ref([])
const responseDataFiltered = ref([])
// Methods
const loadPaymentMethods = async () => {
  loading.value = true
  error.value = null
  
  try {
   const response = await paymentsStore.fetchPaymentMethods({
      page: currentPage.value,
      limit: perPage.value
      })

      pagination.value = response.pagination
      responseData.value = response.data
  } catch (err) {
    error.value = err.message || 'Failed to load payment methods'
  } finally {
    loading.value = false
  }
}
const loadPaymentMethodsWithoutPagination = async () => {
  loading.value = true
  error.value = null
  
  try {
   const response = await paymentsStore.fetchPaymentMethods()

     responseDataFiltered.value = response.data
  } catch (err) {
    error.value = err.message || 'Failed to load payment methods'
  } finally {
    loading.value = false
  }
}

const currentPage = ref(1)
const perPage = ref(10)

// When search query or filter status changes, reset to first page
watch([searchQuery, filterStatus], () => {
  currentPage.value = 1
})

const filteredPaymentMethods = computed(() => {
  let filtered = responseDataFiltered.value
  
  if (searchQuery.value) {
    filtered = filtered.filter(method => 
      method.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      method.code.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  if (filterStatus.value) {
    const isActive = filterStatus.value === 'active'
    filtered = filtered.filter(method => method.active === isActive)
  }
  
  return filtered
})


// Navigate to a specific page.  Ensures the page number stays within
// bounds.  Since the full list is loaded on mount, there is no need
// to refetch from the server when switching pages here.  If in the
// future server‑side pagination is added, you could call
// paymentsStore.fetchPaymentMethods({ page: newPage, limit: itemsPerPage.value }).
const goToPage = (page) => {
  if (page < 1 || page > pagination.value.totalPages) return
  currentPage.value = page
  loadPaymentMethods()
}

const activeCount = computed(() => responseDataFiltered.value.filter(m => m.active).length)
const inactiveCount = computed(() => responseDataFiltered.value.filter(m => !m.active).length)


const openAddModal = () => {
  modalMode.value = 'add'
  formData.value = { name: '', code: '', active: true }
  showModal.value = true
}

const openEditModal = (method) => {
  modalMode.value = 'edit'
  selectedMethod.value = method
  formData.value = {
    name: method.name,
    code: method.code,
    active: method.active
  }
  showModal.value = true
}

const openDeleteModal = (method) => {
  selectedMethod.value = method
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedMethod.value = null
  formData.value = { name: '', code: '', active: true }
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  selectedMethod.value = null
}

const handleSubmit = async () => {
  loading.value = true
  
  try {
    if (modalMode.value === 'add') {
      await paymentsStore.createPaymentMethod(formData.value)
    } else {
      await paymentsStore.updatePaymentMethod(selectedMethod.value.id, formData.value)
    }
    
    closeModal()
    await loadPaymentMethods()
  } catch (err) {
    error.value = err.message || 'Failed to save payment method'
  } finally {
    loading.value = false
  }
}

const deletePaymentMethod = async () => {
  if (!selectedMethod.value) return
  
  loading.value = true
  
  try {
    await paymentsStore.deletePaymentMethod(selectedMethod.value.id)
    closeDeleteModal()
    await loadPaymentMethods()
  } catch (err) {
    error.value = err.message || 'Failed to delete payment method'
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (method) => {
  try {
    await paymentsStore.updatePaymentMethodStatus(method.code, !method.active)
    await loadPaymentMethods()
  } catch (err) {
    error.value = err.message || 'Failed to update status'
  }
}

onMounted(() => {
  loadPaymentMethods()
  loadPaymentMethodsWithoutPagination()
})
</script>
