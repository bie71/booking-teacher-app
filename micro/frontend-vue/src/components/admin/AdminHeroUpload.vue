<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 dark:bg-gray-900 dark:text-white dark:bg-opacity-0">
    <div class="max-w-4xl mx-auto p-6 bg-white rounded-lg shadow-lg relative max-h-[90vh] overflow-y-auto dark:bg-gray-800">
      <!-- Close Button -->
      <button
        @click="$emit('close')"
        class="absolute top-4 right-4 text-gray-500 hover:text-gray-700 text-2xl font-bold dark:text-gray-400"
      >
        &times;
      </button>
      
      <h2 class="text-2xl font-bold mb-6 text-gray-800 dark:text-white">Upload Hero Image</h2>
      
      <div class="space-y-6">
        <!-- Current Hero Display -->
        <div v-if="currentHero" class="mb-6">
          <h3 class="text-lg font-semibold mb-3">Current Hero Image</h3>
          <img :src="currentHero" alt="Current Hero" class="w-full h-64 object-cover rounded-lg" />

        </div>
        
        <!-- Upload Form -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2 dark:text-white">
            Select New Hero Image
          </label>
          <input
            type="file"
            @change="handleFileSelect"
            accept="image/*"
            class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
          />
        </div>
        
        <!-- Preview -->
        <div v-if="previewUrl" class="mb-6">
          <h3 class="text-lg font-semibold mb-3">Preview</h3>
          <img :src="previewUrl" alt="Preview" class="w-full h-64 object-cover rounded-lg" />
        </div>
        
        <!-- Upload Button -->
        <div class="flex justify-end space-x-4">
          <button
            @click="$emit('close')"
            class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
          >
            Cancel
          </button>
          <button
            @click="uploadHero"
            :disabled="!selectedFile || uploading"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            {{ uploading ? 'Uploading...' : 'Upload' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { userService } from '@/services/api'
import { API_CONFIG } from '@/config'
import { useUIStore } from '@/stores/ui'

const uiStore = useUIStore()

export default {
  name: 'AdminHeroUpload',
  props: {
    currentHero: {
      type: String,
      default: null
    }
  },
  emits: ['close', 'uploaded'],
  setup(props, { emit }) {
    const selectedFile = ref(null)
    const previewUrl = ref(null)
    const uploading = ref(false)

    const handleFileSelect = (event) => {
      const file = event.target.files[0]
      if (file) {
        selectedFile.value = file
        previewUrl.value = URL.createObjectURL(file)
      }
    }

const uploadHero = async () => {
  if (!selectedFile.value) return

  uploading.value = true

  try {
    const uploadResponse = await userService.uploadImageHero(selectedFile.value)
    const imageUrl = uploadResponse.file_url 

    console.log('Uploaded image URL:', imageUrl)

    await userService.saveImageHero({ image_url: imageUrl, key_image: API_CONFIG.KEY_IMAGE_HERO })

    uiStore.showSuccess('Image uploaded successfully!')
    emit('uploaded')
    emit('close')

  } catch (error) {
    console.error('Error uploading hero image:', error)
    uiStore.showError('Error uploading hero image')
  } finally {
    uploading.value = false
  }
}

    return {
      selectedFile,
      previewUrl,
      uploading,
      handleFileSelect,
      uploadHero
    }
  }
}
</script>
