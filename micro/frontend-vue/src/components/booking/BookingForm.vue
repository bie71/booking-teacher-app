<template>
  <div class="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-lg">
    <h2 class="text-2xl font-bold mb-6 text-gray-800">Book a Lesson</h2>
    
    <form @submit.prevent="handleBooking" class="space-y-6">
      <!-- Teacher Selection -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Select Teacher</label>
        <select 
          v-model="bookingForm.teacherId" 
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Choose a teacher</option>
          <option v-for="teacher in teachers" :key="teacher.id" :value="teacher.id">
            {{ teacher.name }} - {{ teacher.subject }}
          </option>
        </select>
      </div>

      <!-- Date Selection -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Select Date</label>
        <input 
          v-model="bookingForm.date" 
          type="date" 
          required
          :min="minDate"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <!-- Time Selection -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Select Time</label>
        <select 
          v-model="bookingForm.timeSlot" 
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Choose time slot</option>
          <option v-for="slot in availableSlots" :key="slot" :value="slot">
            {{ slot }}
          </option>
        </select>
      </div>

      <!-- Duration -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Duration (hours)</label>
        <select 
          v-model="bookingForm.duration" 
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="1">1 hour</option>
          <option value="1.5">1.5 hours</option>
          <option value="2">2 hours</option>
        </select>
      </div>

      <!-- Notes -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Notes (Optional)</label>
        <textarea 
          v-model="bookingForm.notes" 
          rows="3"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Any special requirements or notes..."
        ></textarea>
      </div>

      <!-- Price Display -->
      <div v-if="selectedTeacher" class="bg-gray-50 p-4 rounded-md">
        <div class="flex justify-between items-center">
          <span class="text-sm font-medium text-gray-700">Total Price:</span>
          <span class="text-xl font-bold text-green-600">${{ totalPrice }}</span>
        </div>
      </div>

      <!-- Submit Button -->
      <button 
        type="submit" 
        :disabled="loading || !bookingForm.teacherId || !bookingForm.date || !bookingForm.timeSlot"
        class="w-full bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition duration-200"
      >
        {{ loading ? 'Processing...' : 'Proceed to Payment' }}
      </button>
    </form>

    <!-- Error Message -->
    <div v-if="error" class="mt-4 p-3 bg-red-100 text-red-700 rounded-md">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTeachersStore } from '@/stores/teachers'
import { useBookingsStore } from '@/stores/bookings'
import { useRouter } from 'vue-router'

const router = useRouter()
const teachersStore = useTeachersStore()
const bookingsStore = useBookingsStore()

const loading = ref(false)
const error = ref('')

const bookingForm = ref({
  teacherId: '',
  date: '',
  timeSlot: '',
  duration: '1',
  notes: ''
})

const availableSlots = ref([
  '09:00', '10:00', '11:00', '13:00', '14:00', '15:00', '16:00', '17:00'
])

const teachers = computed(() => teachersStore.teachers)
const selectedTeacher = computed(() => 
  teachers.value.find(t => t.id === bookingForm.value.teacherId)
)

const totalPrice = computed(() => {
  if (!selectedTeacher.value) return 0
  return selectedTeacher.value.hourlyRate * parseFloat(bookingForm.value.duration)
})

const minDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

onMounted(async () => {
  await teachersStore.fetchTeachers()
})

const handleBooking = async () => {
  try {
    loading.value = true
    error.value = ''

    const bookingData = {
      teacherId: bookingForm.value.teacherId,
      studentId: JSON.parse(localStorage.getItem('user'))?.id,
      schedule: {
        date: bookingForm.value.date,
        startTime: bookingForm.value.timeSlot,
        duration: parseFloat(bookingForm.value.duration)
      },
      total_price: totalPrice.value,
      notes: bookingForm.value.notes
    }

    const response = await bookingsStore.createBooking(bookingData)
    
    // Redirect to payment page
    router.push(`/payment/${response.id}`)
  } catch (err) {
    error.value = err.message || 'Failed to create booking'
  } finally {
    loading.value = false
  }
}
</script>
