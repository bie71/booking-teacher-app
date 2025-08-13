<template>
  <div class="payment-method-selector">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
      Select Payment Method
    </h3>

    <!-- Loading -->
    <div v-if="paymentsStore.isLoading" class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div v-for="i in 4" :key="i" class="skeleton"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="paymentsStore.activePaymentMethods.length === 0" class="text-center py-8">
      <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
      </svg>
      <p class="text-gray-600 dark:text-gray-400">
        No payment methods available
      </p>
    </div>

    <!-- Payment Methods -->
    <div v-else>
      <!-- Cards grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="method in paginatedMethods"
          :key="method.code"
          @click="selectMethod(method)"
          :class="['payment-method-card', selectedMethod?.code === method.code ? 'selected' : '']"
        >
          <div class="flex items-center space-x-3">
            <div class="payment-method-icon">
              <component :is="getPaymentIcon(method.code)" />
            </div>
            <div>
              <h4 class="font-md text-gray-900 dark:text-white">
                {{ method.name }}
              </h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                {{ getPaymentDescription(method.code) }}
              </p>
            </div>
          </div>
          <input
            type="radio"
            :value="method.code"
            v-model="selectedMethodCode"
            class="payment-method-radio"
          />
        </div>
      </div>

      <!-- Pagination controls for payment methods -->
      <div v-if="totalPages > 1" class="mt-4 flex items-center justify-between">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage <= 1"
          class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>
        <div class="flex items-center space-x-1">
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="goToPage(page)"
            :class="['px-3 py-1 border rounded', page === currentPage ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100']"
          >
            {{ page }}
          </button>
        </div>
        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Selected Info -->
    <div v-if="selectedMethod" class="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
      <div class="flex items-start space-x-3">
        <svg class="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div>
          <h5 class="text-sm font-medium text-blue-900 dark:text-blue-100">
            {{ selectedMethod.name }} Selected
          </h5>
          <p class="text-xs text-blue-700 dark:text-blue-300 mt-1">
            {{ getPaymentInstructions(selectedMethod.code) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// Import heroicons for payment method icons.  These components are
// resolved automatically in the template when returned by
// getPaymentIcon.  We import only the icons we need to avoid
// increasing bundle size significantly.
import { CreditCardIcon, BanknotesIcon, DevicePhoneMobileIcon, BuildingLibraryIcon, BuildingStorefrontIcon } from '@heroicons/vue/24/solid'
import { ref, computed, onMounted, watch } from 'vue'
import { usePaymentsStore } from '@/stores/payments'
import { useUIStore } from '@/stores/ui'


const uiStore = useUIStore()

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'method-selected'])

const paymentsStore = usePaymentsStore()
const selectedMethodCode = ref(props.modelValue)

// Pagination state for displaying payment methods in pages of 10 items.
// When there are more than 10 active payment methods, pagination controls
// will appear. Otherwise, the full list is shown without pagination.
const limit = ref(10)
const currentPage = ref(1)

// Compute the total number of available (active) payment methods
const totalItems = computed(() => paymentsStore.activePaymentMethods.length)

// Determine how many pages are needed based on the limit
const totalPages = computed(() => {
  return totalItems.value > 0 ? Math.ceil(totalItems.value / limit.value) : 1
})

// Slice the active payment methods array to obtain only the items for the
// current page. If the backend returns more methods than the limit, only
// those within the range are displayed.
const paginatedMethods = computed(() => {
  const start = (currentPage.value - 1) * limit.value
  const end = start + limit.value
  return paymentsStore.activePaymentMethods.slice(start, end)
})

// Determine which page numbers to show in the pagination controls. This
// implementation shows up to five page buttons centred around the current
// page when possible. At the edges it will shift to show the first or
// last pages appropriately.
const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = start + maxVisible - 1
  if (end > totalPages.value) {
    end = totalPages.value
    start = Math.max(1, end - maxVisible + 1)
  }
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Navigate to a specific page number. This function ensures the page
// number stays within valid bounds. The component does not need to
// refetch payment methods since the full list is already loaded on mount.
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
}

// Watch for changes to the total number of pages. If the current page
// exceeds the new total, adjust it to the last valid page.
watch(totalPages, (newTotal) => {
  if (currentPage.value > newTotal) {
    currentPage.value = newTotal
  }
})

const selectedMethod = computed(() => {
  return paymentsStore.getPaymentMethodByCode(selectedMethodCode.value)
})

const selectMethod = (method) => {
  selectedMethodCode.value = method.code
}

const getPaymentIcon = (code) => {
  // Map of payment codes to their corresponding imported icon
  // components.  If a code is not found, default to CreditCardIcon.
  const iconMap = {
    'credit_card': CreditCardIcon,
    'bank_transfer': BanknotesIcon,
    'gopay': DevicePhoneMobileIcon,
    'ovo': DevicePhoneMobileIcon,
    'dana': DevicePhoneMobileIcon,
    'shopeepay': DevicePhoneMobileIcon,
    'linkaja': DevicePhoneMobileIcon,
    'bca_va': BuildingLibraryIcon,
    'bni_va': BuildingLibraryIcon,
    'bri_va': BuildingLibraryIcon,
    'mandiri_va': BuildingLibraryIcon,
    'permata_va': BuildingLibraryIcon,
    'indomaret': BuildingStorefrontIcon,
    'alfamart': BuildingStorefrontIcon
  }
  return iconMap[code] || CreditCardIcon
}

const getPaymentDescription = (code) => {
  const descriptions = {
    'credit_card': 'Pay with your credit or debit card',
    'bank_transfer': 'Direct bank transfer',
    'gopay': 'Pay with GoPay e-wallet',
    'ovo': 'Pay with OVO e-wallet',
    'dana': 'Pay with DANA e-wallet',
    'shopeepay': 'Pay with ShopeePay e-wallet',
    'linkaja': 'Pay with LinkAja e-wallet',
    'bca_va': 'BCA Virtual Account',
    'bni_va': 'BNI Virtual Account',
    'bri_va': 'BRI Virtual Account',
    'mandiri_va': 'Mandiri Virtual Account',
    'permata_va': 'Permata Virtual Account',
    'indomaret': 'Pay at Indomaret stores',
    'alfamart': 'Pay at Alfamart stores'
  }
  
  return descriptions[code] || 'Payment method'
}

const getPaymentInstructions = (code) => {
  const instructions = {
    'credit_card': 'You will be redirected to enter your card details securely.',
    'bank_transfer': 'You will receive bank account details to complete the transfer.',
    'gopay': 'You will be redirected to GoPay to complete the payment.',
    'ovo': 'You will be redirected to OVO to complete the payment.',
    'dana': 'You will be redirected to DANA to complete the payment.',
    'shopeepay': 'You will be redirected to ShopeePay to complete the payment.',
    'linkaja': 'You will be redirected to LinkAja to complete the payment.',
    'bca_va': 'You will receive a BCA Virtual Account number to complete the payment.',
    'bni_va': 'You will receive a BNI Virtual Account number to complete the payment.',
    'bri_va': 'You will receive a BRI Virtual Account number to complete the payment.',
    'mandiri_va': 'You will receive a Mandiri Virtual Account number to complete the payment.',
    'permata_va': 'You will receive a Permata Virtual Account number to complete the payment.',
    'indomaret': 'You will receive a payment code to pay at any Indomaret store.',
    'alfamart': 'You will receive a payment code to pay at any Alfamart store.'
  }
  
  return instructions[code] || 'Follow the instructions to complete your payment.'
}

// Watch for changes and emit events
watch(selectedMethodCode, (newValue) => {
  emit('update:modelValue', newValue)
  if (newValue) {
    const method = paymentsStore.getPaymentMethodByCode(newValue)
    emit('method-selected', method)
  }
})

// Watch for prop changes
watch(() => props.modelValue, (newValue) => {
  selectedMethodCode.value = newValue
})

onMounted(async () => {
  // Fetch payment methods when component mounts
  const result = await paymentsStore.fetchPaymentMethods()
  if (!result.success) {
    uiStore.showError('Failed to load payment methods')
    console.error('Failed to load payment methods:', result.message)

  }
})
</script>

<style scoped>
.payment-method-card {
  @apply border border-gray-200 dark:border-gray-700 rounded-lg p-1.5 flex items-center justify-between transition-all cursor-pointer;
}
.payment-method-card:hover {
  @apply border-blue-300 dark:border-blue-600 shadow-sm;
}
.payment-method-card.selected {
  @apply border-blue-500 bg-blue-50 dark:bg-blue-900/20;
}
.payment-method-icon {
  @apply flex items-center justify-center w-7 h-7 bg-gray-100 dark:bg-gray-800 rounded-md;
}
.payment-method-icon svg {
  @apply w-3.5 h-3.5 text-gray-700 dark:text-gray-300;
}
.payment-method-radio {
  @apply w-3.5 h-3.5 text-blue-600 border-gray-300 focus:ring-blue-500;
}
.skeleton {
  @apply animate-pulse bg-gray-200 dark:bg-gray-700 rounded h-8;
}
</style>
