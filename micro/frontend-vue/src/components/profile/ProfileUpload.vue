<template>
  <div class="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
    <h3 class="text-lg font-semibold mb-4">Upload Profile Picture</h3>
    
    <div class="space-y-4">
      <!-- Current Profile Picture -->
      <div v-if="currentImage" class="text-center">
        <img 
          :src="currentImage" 
          alt="Current Profile" 
          class="w-32 h-32 rounded-full mx-auto object-cover border-4 border-gray-200"
        />
      </div>
      
      <!-- Default Avatar -->
      <div v-else class="text-center">
        <div class="w-32 h-32 rounded-full mx-auto bg-gray-300 flex items-center justify-center">
          <svg class="w-16 h-16 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
          </svg>
        </div>
      </div>

      <!-- Upload Button -->
      <div class="text-center">
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          @change="handleFileSelect"
          class="hidden"
        />
        <button
          @click="triggerFileInput"
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition duration-200"
        >
          Choose Image
        </button>
      </div>

      <!-- Preview -->
      <div v-if="previewImage" class="text-center">
        <img 
          :src="previewImage" 
          alt="Preview" 
          class="w-32 h-32 rounded-full mx-auto object-cover border-4 border-blue-200"
        />
        <div class="mt-4 space-x-2">
          <button
            @click="uploadImage"
            :disabled="uploading"
            class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:bg-gray-400"
          >
            {{ uploading ? 'Uploading...' : 'Upload' }}
          </button>
          <button
            @click="clearPreview"
            class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
          >
            Cancel
          </button>
        </div>
      </div>

      <!-- Messages -->
      <div v-if="successMessage" class="p-3 bg-green-100 text-green-700 rounded-md text-sm">
        {{ successMessage }}
      </div>
      <div v-if="errorMessage" class="p-3 bg-red-100 text-red-700 rounded-md text-sm">
        {{ errorMessage }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  userType: {
    type: String,
    required: true,
    validator: (value) => ['user', 'teacher'].includes(value)
  },
  userId: {
    type: String,
    required: true
  }
})

const authStore = useAuthStore()
const fileInput = ref(null)
const previewImage = ref(null)
const currentImage = ref(null)
const uploading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

const selectedFile = ref(null)

onMounted(async () => {
  await loadCurrentProfile()
})

const loadCurrentProfile = async () => {
  try {
    const response = await authStore.getUserProfile(props.userId)
    currentImage.value = response.profileImage || response.avatar || null
  } catch (error) {
    console.error('Error loading profile:', error)
  }
}

const triggerFileInput = () => {
  fileInput.value.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file && file.type.startsWith('image/')) {
    if (file.size > 5 * 1024 * 1024) {
      errorMessage.value = 'File size must be less than 5MB'
      return
    }
    
    selectedFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      previewImage.value = e.target.result
    }
    reader.readAsDataURL(file)
  } else {
    errorMessage.value = 'Please select a valid image file'
  }
}

const uploadImage = async () => {
  if (!selectedFile.value) return

  uploading.value = true
  successMessage.value = ''
  errorMessage.value = ''

  try {
    const formData = new FormData()
    formData.append('image', selectedFile.value)
    formData.append('userId', props.userId)
    formData.append('userType', props.userType)

    let response
    if (props.userType === 'teacher') {
      response = await authStore.uploadTeacherProfile(formData)
    } else {
      response = await authStore.uploadUserProfile(formData)
    }
    
    successMessage.value = 'Profile picture updated successfully!'
    currentImage.value = response.url
    clearPreview()
    
    // Emit event to parent
    emit('profile-updated', response.url)
  } catch (error) {
    errorMessage.value = error.message || 'Failed to upload profile picture'
  } finally {
    uploading.value = false
  }
}

const clearPreview = () => {
  previewImage.value = null
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const emit = defineEmits(['profile-updated'])
</script>
