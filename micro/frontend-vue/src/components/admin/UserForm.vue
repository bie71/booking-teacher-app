<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="flex items-center justify-center min-h-screen p-4">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full dark:bg-gray-800">
        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ mode === 'create' ? 'Add New User' : 'Edit User' }}
          </h3>
        </div>

        <form @submit.prevent="handleSubmit" class="px-6 py-4">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Name
              </label>
              <input 
                v-model="form.name"
                type="text"
                required
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Email
              </label>
              <input 
                v-model="form.email"
                type="email"
                required
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>

            <div v-if="mode === 'create'">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Password
              </label>
              <input 
                v-model="form.password"
                type="password"
                required
                minlength="6"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Role
              </label>
              <select 
                v-model="form.role"
                required
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              >
                <option value="user">User</option>
                <option value="teacher">Teacher</option>
                <option value="admin">Admin</option>
              </select>
            </div>

            <!-- Profile image uploader (only shown in edit mode). Allows admins to upload
                 a new profile picture for the user. The uploaded image will be
                 displayed as a preview. -->
            <div v-if="mode === 'edit'">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Profile Image
              </label>
              <input
                type="file"
                accept="image/*"
                @change="handleFileChange"
                class="mt-1 block w-full text-sm text-gray-600 dark:text-gray-300 file:mr-4 file:py-2 file:px-4 file:border file:border-gray-300 file:rounded-md file:bg-gray-50 dark:file:bg-gray-700 dark:file:border-gray-600"
              />
              <!-- Show preview if an image has been uploaded or exists -->
              <div v-if="profileImageUrl" class="mt-2">
                <img :src="profileImageUrl" alt="Profile Image Preview" class="w-16 h-16 rounded-full object-cover" />
              </div>
            </div>
          </div>

          <div class="mt-6 flex justify-end space-x-3">
            <button 
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
            >
              Cancel
            </button>
            <button 
              type="submit"
              :disabled="loading"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {{ loading ? 'Saving...' : (mode === 'create' ? 'Create User' : 'Update User') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { userService } from '@/services/api'
import { useUIStore } from '@/stores/ui'

const props = defineProps({
  user: {
    type: Object,
    default: null
  },
  mode: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close', 'user-saved'])

const uiStore = useUIStore()
const loading = ref(false)
const form = ref({
  name: '',
  email: '',
  password: '',
  role: 'user'
})

// Hold the uploaded profile image URL. When editing a user this will be
// pre-filled with the current profile image. When a new file is selected
// this value will be updated after uploading the file.
const profileImageUrl = ref('')

// Handle file selection and upload the image to the server. Uses the
// userService.uploadUserImage helper which returns the image URL on
// success. If the upload fails, show an error toast.
async function handleFileChange(event) {
  const file = event.target.files && event.target.files[0]
  if (!file) return
  try {
    const response = await userService.uploadUserImage(file)
    profileImageUrl.value = response.file_url || response.path || response
    uiStore.showSuccess('Image uploaded successfully')
  } catch (error) {
    console.error('Failed to upload image:', error)
    uiStore.showError('Failed to upload image')
  }
}

const resetForm = () => {
  form.value = {
    name: '',
    email: '',
    password: '',
    role: 'user'
  }
  profileImageUrl.value = ''
}

const loadUserData = () => {
  if (props.user && props.mode === 'edit') {
    form.value = {
      name: props.user.name,
      email: props.user.email,
      password: '',
      role: props.user.role
    }
    // Initialize the profile image URL with the existing image, if any
    profileImageUrl.value = props.user.profile_image || ''
  }
}

const handleSubmit = async () => {
  loading.value = true
  
  try {
    if (props.mode === 'create') {
      await userService.createUser(form.value)
      uiStore.showSuccess('User created successfully')
      // Log admin create user activity
      try {
        await userService.logActivity({
          action: 'AdminCreateUser',
          description: `Created user ${form.value.name}`
        })
      } catch (error) {
        console.warn('Failed to log admin create user activity:', error)
      }
    } else {
      // Build the payload for updating the user. Include the
      // profile_image field only if a new image has been uploaded.
      const payload = {
        name: form.value.name,
        email: form.value.email,
        role: form.value.role
      }
      if (profileImageUrl.value) {
        payload.profile_image = profileImageUrl.value
      }
      await userService.updateUser(props.user.id, payload)
      uiStore.showSuccess('User updated successfully')
      // Log admin update user activity
      try {
        await userService.logActivity({
          action: 'AdminUpdateUser',
          description: `Updated user ${form.value.name}`
        })
      } catch (error) {
        console.warn('Failed to log admin update user activity:', error)
      }
    }

    emit('user-saved')
    resetForm()
  } catch (error) {
    console.error('Failed to save user:', error)
    uiStore.showError(error.response?.data?.error ||'Failed to save user')
  } finally {
    loading.value = false
  }
}

watch(() => props.user, loadUserData)
onMounted(loadUserData)
</script>
