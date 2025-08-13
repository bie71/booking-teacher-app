<template>
  <div>
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900">
        {{ mode === 'create' ? 'Add New Teacher' : 'Edit Teacher' }}
      </h2>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Basic Information -->
      <div class="bg-white shadow rounded-lg p-6 dark:bg-gray-800">
        <h3 class="text-lg font-medium text-gray-900 mb-4 dark:text-white">Basic Information</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-white">Full Name *</label>
            <input 
              v-model="form.name"
              type="text"
              id="name"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          
          <!-- <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Email *</label>
            <input 
              v-model="form.email"
              type="email"
              id="email"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div> -->
          
          <!-- <div>
            <label for="phone" class="block text-sm font-medium text-gray-700">Phone Number</label>
            <input 
              v-model="form.phone"
              type="tel"
              id="phone"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div> -->
          
          <div>
            <label for="language_level" class="block text-sm font-medium text-gray-700 dark:text-white">Language Level *</label>
            <select 
              v-model="form.language_level"
              id="language_level"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            >
              <option value="">Select Level</option>
              <option value="beginner">Beginner</option>
              <option value="intermediate">Intermediate</option>
              <option value="advanced">Advanced</option>
              <option value="native">Native</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Profile Details -->
      <div class="bg-white shadow rounded-lg p-6 dark:bg-gray-800">
        <h3 class="text-lg font-medium text-gray-900 mb-4 dark:text-white">Profile Details</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="bio" class="block text-sm font-medium text-gray-700 dark:text-white">Bio</label>
            <textarea 
              v-model="form.bio"
              id="bio"
              rows="3"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            ></textarea>
          </div>
          
          <div>
            <label for="price_per_hour" class="block text-sm font-medium text-gray-700 dark:text-white">Price per Hour *</label>
            <input 
              v-model.number="form.price_per_hour"
              type="number"
              id="price_per_hour"
              min="0"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>

          <label for="available_start_time" class="block text-sm font-medium text-gray-700 dark:text-white">Available Start Time *</label>
          <input 
            v-model="form.available_start_time"
            type="time"
            id="available_start_time"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          />
        </div>

        <div>
          <label for="available_end_time" class="block text-sm font-medium text-gray-700 dark:text-white">Available End Time *</label>
          <input 
            v-model="form.available_end_time"
            type="time"
            id="available_end_time"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          />
        </div>
          <!-- <div>
            <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
            <select 
              v-model="form.status"
              id="status"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </div> -->
        </div>
      </div>

      <!-- Profile Image -->
      <div class="bg-white shadow rounded-lg p-6 dark:bg-gray-800">
        <h3 class="text-lg font-medium text-gray-900 mb-4 dark:text-white">Profile Image</h3>
        <div>
          <label for="profile_image" class="block text-sm font-medium text-gray-700 dark:text-white">Profile Image</label>
          <div class="mt-1 flex items-center">
            <div v-if="previewImage" class="mr-4">
              <img :src="previewImage" alt="Preview" class="h-20 w-20 rounded-full object-cover" />
            </div>
            <input 
              type="file"
              id="profile_image"
              accept="image/*"
              @change="handleImageUpload"
              class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-medium file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 dark:file:bg-gray-700 dark:file:text-gray-300 dark:hover:file:bg-gray-600 dark:hover:file:text-gray-200"
            />
          </div>
        </div>
      </div>

      <!-- Submit Button -->
      <div class="flex justify-end space-x-3">
        <button 
          type="button"
          @click="$emit('cancel')"
          class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
        >
          Cancel
        </button>
        <button 
          type="submit"
          :disabled="isSubmitting"
          class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
        >
          {{ isSubmitting ? 'Saving...' : (mode === 'create' ? 'Create Teacher' : 'Update Teacher') }}
        </button>
      </div>
    </form>

    <!-- Loading Overlay -->
    <div v-if="isSubmitting" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-4 rounded-lg">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-2 text-sm text-gray-600">Saving teacher...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useTeachersStore } from '@/stores/teachers'
import { useUIStore } from '@/stores/ui'
import { userService } from '@/services/api'

const props = defineProps({
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  },
  teacher: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

const teachersStore = useTeachersStore()
const uiStore = useUIStore()

const isSubmitting = ref(false)
const previewImage = ref('')
const form = ref({
  name: '',
  email: '',
  phone: '',
  language_level: '',
  bio: '',
  price_per_hour: 0,
  status: 'active',
  profile_image: '',
  available_start_time: '',
  available_end_time: ''
})


// Initialize form based on mode
const initializeForm = () => {
  if (props.mode === 'edit' && props.teacher) {
    form.value = {
      name: props.teacher.name || '',
      email: props.teacher.email || '',
      phone: props.teacher.phone || '',
      language_level: props.teacher.language_level || '',
      bio: props.teacher.bio || '',
      price_per_hour: props.teacher.price_per_hour || 0,
      status: props.teacher.status || 'active',
      profile_image: props.teacher.profile_image || '',
      available_start_time: props.teacher.available_start_time || '',
      available_end_time: props.teacher.available_end_time || ''
    }
  } else {
    form.value = {
      name: '',
      email: '',
      phone: '',
      language_level: '',
      bio: '',
      price_per_hour: 0,
      status: 'active',
      profile_image: '',
      available_start_time: '',
      available_end_time: ''
    }
  }
}

// Methods
const handleImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      previewImage.value = e.target.result
    }
    reader.readAsDataURL(file)
    
    // Upload image
    teachersStore.uploadTeacherImage(file).then((response) => {
      form.value.profile_image = response.file_url || response.data?.file_url
    }).catch((error) => {
      uiStore.showError('Failed to upload image')
    })
  }
}

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    
    const teacherData = {
      ...form.value,
      price_per_hour: Number(form.value.price_per_hour)
    }

    if (props.mode === 'create') {
      const result = await teachersStore.createTeacher(teacherData)
      emit('submit', result.data)
      uiStore.showSuccess('Teacher created successfully')
      // Log admin create teacher activity
      try {
        await userService.logActivity({
          action: 'AdminCreateTeacher',
          description: `Created teacher ${teacherData.name}`
        })
      } catch (error) {
        console.warn('Failed to log admin create teacher activity:', error)
      }
    } else {
      const result = await teachersStore.updateTeacher(props.teacher.id, teacherData)
      emit('submit', result.data)
      uiStore.showSuccess('Teacher updated successfully')
      // Log admin update teacher activity
      try {
        await userService.logActivity({
          action: 'AdminUpdateTeacher',
          description: `Updated teacher ${teacherData.name}`
        })
      } catch (error) {
        console.warn('Failed to log admin update teacher activity:', error)
      }
    }
    
    // Reset form
    initializeForm()
    previewImage.value = ''
  } catch (error) {
    uiStore.showError(error.message || 'Failed to save teacher')
  } finally {
    isSubmitting.value = false
  }
}

// Initialize form on mount
onMounted(() => {
  initializeForm()
})

// Watch for teacher prop changes
watch(() => props.teacher, () => {
  initializeForm()
}, { immediate: true })
</script>
