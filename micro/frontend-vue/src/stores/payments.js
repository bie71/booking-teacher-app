import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { paymentService } from '@/services/api'
import { PAYMENT_STATUS, API_CONFIG, ERROR_MESSAGES } from '@/config'

export const usePaymentsStore = defineStore('payments', () => {
  // State
  const payments = ref([])
  const paymentMethods = ref([])
  const currentPayment = ref(null)
  const loading = ref(false)
  const processingPayment = ref(false)
  const pagination = ref({
    page: 1,
    limit: 10,
    total: 0
  })

  // Separate pagination state for payment methods. This object
  // tracks the current page, limit and total count of payment
  // methods when using the admin endpoint to retrieve a paginated
  // list. For public payment methods (used by PaymentMethodSelector),
  // the total will be computed from the length of the returned
  // array.
  const methodsPagination = ref({
    page: 1,
    limit: 10,
    total: 0
  })

  // Getters
  const activePaymentMethods = computed(() => {
    return paymentMethods.value.filter(method => method.active)
  })

  const pendingPayments = computed(() => {
    return payments.value.filter(payment => 
      payment.status === PAYMENT_STATUS.PENDING
    )
  })

  const successfulPayments = computed(() => {
    return payments.value.filter(payment => 
      payment.status === PAYMENT_STATUS.SETTLEMENT
    )
  })

  const failedPayments = computed(() => {
    return payments.value.filter(payment => 
      payment.status === PAYMENT_STATUS.FAILED
    )
  })

  const isLoading = computed(() => loading.value)
  const isProcessingPayment = computed(() => processingPayment.value)

  // Actions
  /**
   * Fetch payment methods from the backend. When params.page is
   * provided, the admin endpoint will be used to return a paginated
   * response along with metadata. Otherwise, the public endpoint is
   * used which returns all active methods at once. Pagination
   * information is stored in methodsPagination and can be consumed by
   * components such as PaymentMethodManagement.
   *
   * @param {Object} params Optional parameters (page, limit, etc.)
   */
  const fetchPaymentMethods = async (params = {}) => {
  loading.value = true
  try {
    const useAdmin = params.page != null // 0/1/2 valid; beda dengan truthy
    const res = useAdmin
      ? await paymentService.getAdminPaymentMethods(params)
      : await paymentService.getPaymentMethods(params)

    // Unwrap payload: dukung bentuk Axios (res.data) atau plain
    const payload = res?.data ?? res
    const dataArr =
      Array.isArray(payload?.data) ? payload.data
      : Array.isArray(payload) ? payload
      : [] // fallback kalau salah bentuk

    paymentMethods.value = dataArr

    if (useAdmin) {
      // Ambil pagination dari payload (axios: payload.pagination)
      const rawP = payload?.pagination ?? res?.pagination ?? null

      const page  = Number(rawP?.current_page ?? params.page ?? 1)
      const limit = Number(rawP?.limit ?? params.limit ?? 10)
      const total = Number(
        rawP?.total_items ?? rawP?.total ?? rawP?.total_count ?? dataArr.length ?? 0
      )

      methodsPagination.value = {
        page,
        limit,
        total,
        totalPages: rawP?.total_pages ?? rawP?.total_page ?? rawP?.pages ?? 0
      }

      // Safety: kalau page > totalPages, balik ke terakhir
      if (methodsPagination.value.page > methodsPagination.value.totalPages) {
        methodsPagination.value.page = methodsPagination.value.totalPages
      }
    } else {
      // Public endpoint: tidak paginated
      const total = dataArr.length || 0
      methodsPagination.value = {
        page: 1,
        limit: total || 10,
        total,
        totalPages: 1
      }
    }

    return { success: true, data: paymentMethods.value, pagination: methodsPagination.value }
  } catch (error) {
    console.error('Error fetching payment methods:', error)
    return {
      success: false,
      message: error?.message || ERROR_MESSAGES.NETWORK_ERROR
    }
  } finally {
    loading.value = false
  }
}


  // Create a new payment method
  const createPaymentMethod = async (methodData) => {
    try {
      loading.value = true
      const response = await paymentService.createPaymentMethod(methodData)

      const newMethod = response.data || response
      paymentMethods.value.unshift(newMethod.data)
      return {
        success: true,
        data: newMethod.data
      }
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Failed to create payment method'
      }
    } finally {
      loading.value = false
    }
  }

  // Update a payment method
  const updatePaymentMethod = async (id, updateData) => {
    try {
      loading.value = true
      const response = await paymentService.updatePaymentMethod(id, updateData)

      const updatedMethod = response.data || response
      const index = paymentMethods.value.findIndex(method => method.id === id)
      if (index !== -1) {
        paymentMethods.value[index] = { ...paymentMethods.value[index], ...updatedMethod.data }
      }
      return {
        success: true,
        data: updatedMethod.data
      }
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Failed to update payment method'
      }
    } finally {
      loading.value = false
    }
  }

  // Delete a payment method
  const deletePaymentMethod = async (id) => {
    try {
      loading.value = true
     await paymentService.deletePaymentMethod(id)

      return {
        success: true,
        message: 'Payment method deleted successfully'
      }
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Failed to delete payment method'
      }
    } finally {
      loading.value = false
    }
  }

  // Update payment method status
  const updatePaymentMethodStatus = async (code, active) => {
    try {
      loading.value = true
      const response = await paymentService.updatePaymentMethodStatus(code, active)

      const methodIndex = paymentMethods.value.findIndex(method => method.code === code)
      if (methodIndex !== -1) {
        paymentMethods.value[methodIndex].active = active
      }
      return {
        success: true,
        message: 'Payment method status updated'
      }
    } catch (error) {
      return {
        success: false,
        message: error.message || 'Failed to update status'
      }
    } finally {
      loading.value = false
    }
  }

  const fetchPayments = async (params = {}) => {
    try {
      loading.value = true
      const response = await paymentService.getPayments({
        page: pagination.value.page,
        limit: pagination.value.limit,
        ...params
      })

      payments.value = response.data || response
      
      if (response.pagination) {
        pagination.value = { ...pagination.value, ...response.pagination }
      }

      return {
        success: true,
        data: payments.value
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchPayment = async (id) => {
    try {
      loading.value = true
      const response = await paymentService.getPayment(id)
      
      currentPayment.value = response.data || response
      
      return {
        success: true,
        data: currentPayment.value
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  const createPayment = async (paymentData) => {
    try {
      processingPayment.value = true
      const response = await paymentService.createPayment(paymentData)
      
      const newPayment = response.data || response
      
      if (newPayment.id) {
        payments.value.unshift(newPayment)
      }
      
      return {
        success: true,
        data: newPayment,
        message: 'Payment created successfully'
      }
    } catch (error) {
      return {
        success: false,
        message: error.message || ERROR_MESSAGES.PAYMENT_FAILED
      }
    } finally {
      processingPayment.value = false
    }
  }

  const handlePaymentCallback = async (callbackData) => {
    try {
      const response = await paymentService.handlePaymentCallback(callbackData)
      
      const updatedPayment = response.payment || response.data || response
      
      const index = payments.value.findIndex(p => p.id === updatedPayment.id)
      if (index !== -1) {
        payments.value[index] = updatedPayment
      }
      
      if (currentPayment.value?.id === updatedPayment.id) {
        currentPayment.value = updatedPayment
      }
      
      return {
        success: true,
        data: updatedPayment,
        message: 'Payment status updated'
      }
    } catch (error) {
      throw error
    }
  }

  // Utility functions
  const getPaymentById = (id) => {
    return payments.value.find(payment => payment.id === id)
  }

  const getPaymentMethodByCode = (code) => {
    return paymentMethods.value.find(method => method.code === code)
  }

  const getPaymentMethodByName = (name) => {
    return paymentMethods.value.find(method => method.name.toLowerCase() === name.toLowerCase())
  }

  const formatPaymentAmount = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR'
    }).format(amount)
  }

  const getPaymentStatusColor = (status) => {
    switch (status) {
      case PAYMENT_STATUS.PENDING:
        return 'text-yellow-600 bg-yellow-100'
      case PAYMENT_STATUS.SETTLEMENT:
        return 'text-green-600 bg-green-100'
      case PAYMENT_STATUS.FAILED:
      case PAYMENT_STATUS.CANCEL:
        return 'text-red-600 bg-red-100'
      default:
        return 'text-gray-600 bg-gray-100'
    }
  }

  const isPaymentSuccessful = (payment) => {
    return payment.status === PAYMENT_STATUS.SETTLEMENT
  }

  const isPaymentPending = (payment) => {
    return payment.status === PAYMENT_STATUS.PENDING
  }

  const isPaymentFailed = (payment) => {
    return payment.status === PAYMENT_STATUS.FAILED || payment.status === PAYMENT_STATUS.CANCEL
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
    payments: computed(() => payments.value),
    paymentMethods: computed(() => paymentMethods.value),
    currentPayment: computed(() => currentPayment.value),
    pagination: computed(() => pagination.value),
    // Pagination state for payment methods (admin view)
    // Total number of payment methods available in the current paginated view
    totalItems: computed(() => methodsPagination.value.total),
    
    // Getters
    activePaymentMethods,
    pendingPayments,
    successfulPayments,
    failedPayments,
    isLoading,
    isProcessingPayment,
    
    // Actions
    fetchPaymentMethods,
    fetchPayments,
    fetchPayment,
    createPayment,
    handlePaymentCallback,
    updatePaymentMethodStatus,
    createPaymentMethod,
    deletePaymentMethod,
    updatePaymentMethod,
    
    // Utility functions
    getPaymentById,
    getPaymentMethodByCode,
    getPaymentMethodByName,
    formatPaymentAmount,
    getPaymentStatusColor,
    isPaymentSuccessful,
    isPaymentPending,
    isPaymentFailed,
    
    // Pagination actions
    setPage,
    setLimit,
    resetPagination
  }
})
