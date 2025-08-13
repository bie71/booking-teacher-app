<template bg-white p-4 dark:bg-gray-800>
  <Modal :show="true" @close="$emit('close')">
    <template #header>
      <h2 class="text-lg font-semibold">Edit Booking Status</h2>
    </template>

    <template #body>
      <form @submit.prevent="submit">
        <div class="mb-4 bg-white p-4 rounded shadow dark:bg-gray-800">
          <label for="status" class="block text-sm font-medium mb-1">Status</label>
          <select v-model="newStatus" id="status" class="select select-bordered dark:bg-gray-800">
            <option disabled value="">Select status</option>
            <option value="pending">Pending</option>
            <option value="paid">Paid</option>
            <option value="cancelled">Cancelled</option>
            <option value="rescheduled">Rescheduled</option>
            <option value="completed">Completed</option>
          </select>
        </div>

        <div class="flex justify-end space-x-2 m-2">
          <button type="button" class="btn btn-secondary dark:text-white" @click="$emit('close')">Cancel</button>
          <button type="submit" class="btn btn-primary">Update</button>
        </div>
      </form>
    </template>
  </Modal>
</template>

<script setup>
import Modal from '@/components/Modal.vue'
import { ref } from 'vue'
import { useUIStore } from '@/stores/ui'
import { useBookingsStore } from '@/stores/bookings'

const props = defineProps({
  booking: Object
})

const emit = defineEmits(['close', 'updated'])

const newStatus = ref(props.booking.status)
const uiStore = useUIStore()
const bookingsStore = useBookingsStore()

const submit = async () => {
  const result = await bookingsStore.changeStatus(props.booking.id, newStatus.value)
  if (result && result.success) {
    uiStore.showSuccess('Status updated successfully!')
    emit('updated')
  } else {
    uiStore.showError('Failed to update status.')
  }
}
</script>
