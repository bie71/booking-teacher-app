<template>
  <div class="payment-method-management">
    <div class="bg-white shadow rounded-lg p-6">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800">Payment Method Management</h2>
        <button 
          @click="showAddModal = true"
          class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
        >
          Add Payment Method
        </button>
      </div>

      <!-- Search and Filter -->
      <div class="mb-4 flex gap-4">
        <div class="flex-1">
          <input 
            v-model="searchQuery"
            type="text"
            placeholder="Search payment methods..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <select 
          v-model="filterStatus"
          class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading payment methods...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-8">
        <div class="text-red-600">
          <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <p class="mt-2">{{ error }}</p>
          <button 
            @click="loadPaymentMethods"
            class="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>

      <!-- Payment Methods Table -->
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
                Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
                Code
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-white">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200 dark:bg-gray-800">
            <tr v-for="method in filteredPaymentMethods" :key="method.code">
              <!-- Name with icon -->
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 flex items-center">
                <!-- Generic credit card icon to visually represent payment methods -->
                <svg class="w-5 h-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.25 7.5A2.25 2.25 0 014.5 5.25h15a2.25 2.25 0 012.25 2.25v9a2.25 2.25 0 01-2.25 2.25h-15A2.25 2.25 0 012.25 16.5v-9z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.25 9.75h19.5" />
                </svg>
                {{ method.name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
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
                <!-- Toggle status button using icons -->
                <button 
                  @click="toggleStatus(method)"
                  :class="method.active ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'"
                  class="text-white p-1 rounded"
                  :title="method.active ? 'Deactivate' : 'Activate'"
                >
                  <!-- Deactivate icon (X) -->
                  <svg v-if="method.active" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  <!-- Activate icon (check) -->
                  <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="!loading && !error && totalPages > 1" class="mt-4 flex justify-between items-center">
        <!-- <div class="text-sm text-gray-700">
          Showing {{ (currentPage - 1) * limit + 1 }} to {{ Math.min(currentPage * limit, totalItems) }} of {{ totalItems }} results
        </div>
        <div class="flex space-x-2">
          <button 
            @click="previousPage"
            :disabled="currentPage === 1"
            class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button 
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div> -->
          <Pagination
          :page="currentPage"
          :limit="pagination.limit"
          :total="pagination.total"
          @update:page="goToPage"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { usePaymentsStore } from '@/stores/payments'
import Pagination from '@/components/ui/Pagination.vue'

export default {
  name: 'PaymentMethodManagement',
  setup() {
    const paymentStore = usePaymentsStore()
    const searchQuery = ref('')
    const filterStatus = ref('')
    const loading = ref(false)
    const error = ref(null)
    const currentPage = ref(1)
    const limit = ref(10)

    const paymentMethods = computed(() => paymentStore.paymentMethods)
    const totalItems = computed(() => paymentStore.totalItems)
    const totalPages = computed(() => Math.ceil(totalItems.value / limit.value))

    const filteredPaymentMethods = computed(() => {
      let filtered = paymentMethods.value
      
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

    const loadPaymentMethods = async () => {
      loading.value = true
      error.value = null
      
      try {
        await paymentStore.fetchPaymentMethods({
          page: currentPage.value,
          limit: limit.value
        })
      } catch (err) {
        error.value = err.message || 'Failed to load payment methods'
      } finally {
        loading.value = false
      }
    }

    const toggleStatus = async (method) => {
      try {
        await paymentStore.updatePaymentMethodStatus(method.code, !method.active)
        await loadPaymentMethods()
      } catch (err) {
        error.value = err.message || 'Failed to update status'
      }
    }

    const previousPage = () => {
      if (currentPage.value > 1) {
        currentPage.value--
        loadPaymentMethods()
      }
    }

    const nextPage = () => {
      if (currentPage.value < totalPages.value) {
        currentPage.value++
        loadPaymentMethods()
      }
    }

    onMounted(() => {
      loadPaymentMethods()
    })

    return {
      searchQuery,
      filterStatus,
      loading,
      error,
      currentPage,
      limit,
      filteredPaymentMethods,
      totalItems,
      totalPages,
      loadPaymentMethods,
      toggleStatus,
      previousPage,
      nextPage
    }
  }
}
</script>

<style scoped>
.payment-method-management {
  @apply p-6;
}
</style>
