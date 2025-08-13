<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="container py-8">
      <div class="max-w-2xl mx-auto">
        <div class="mb-8">
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Profile Settings
          </h1>
          <p class="text-gray-600 dark:text-gray-400">
            Manage your account information and preferences
          </p>
        </div>

        <!-- Profile Information -->
        <div class="card mb-8">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Profile Information</h2>
          </div>
          <div class="card-body">
            <form @submit.prevent="updateProfile" class="space-y-6">
              <div class="flex items-center space-x-4 mb-4">
                <img
                  v-if="profileForm.profile_image"
                  :src="profileForm.profile_image"
                  alt="Profile Image"
                  class="w-20 h-20 rounded-full object-cover"
                />
                <div>
                  <label
                    for="profileImage"
                    class="cursor-pointer inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                  >
                    Upload Profile Image
                  </label>
                  <input
                    id="profileImage"
                    type="file"
                    accept="image/*"
                    class="hidden"
                    @change="handleFileUpload"
                  />
                </div>
              </div>

              <div>
                <label for="name" class="form-label">Full Name</label>
                <input
                  id="name"
                  v-model="profileForm.name"
                  type="text"
                  required
                  class="form-input"
                  :class="{ 'border-red-500': errors.name }"
                />
                <p v-if="errors.name" class="form-error">{{ errors.name }}</p>
              </div>

              <div>
                <label for="email" class="form-label">Email Address</label>
                <input
                  id="email"
                  v-model="profileForm.email"
                  type="email"
                  required
                  class="form-input"
                  :class="{ 'border-red-500': errors.email }"
                />
                <p v-if="errors.email" class="form-error">{{ errors.email }}</p>
              </div>

              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="authStore.isLoading"
                  class="btn btn-primary"
                >
                  <span v-if="authStore.isLoading" class="loading-spinner w-4 h-4 mr-2"></span>
                  {{ authStore.isLoading ? 'Updating...' : 'Update Profile' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Change Password -->
        <div class="card mb-8">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Change Password</h2>
          </div>
          <div class="card-body">
            <form @submit.prevent="changePassword" class="space-y-6">
              <div>
                <label for="currentPassword" class="form-label">Current Password</label>
                <input
                  id="currentPassword"
                  v-model="passwordForm.currentPassword"
                  type="password"
                  required
                  class="form-input"
                  :class="{ 'border-red-500': errors.currentPassword }"
                />
                <p v-if="errors.currentPassword" class="form-error">{{ errors.currentPassword }}</p>
              </div>

              <div>
                <label for="newPassword" class="form-label">New Password</label>
                <input
                  id="newPassword"
                  v-model="passwordForm.newPassword"
                  type="password"
                  required
                  class="form-input"
                  :class="{ 'border-red-500': errors.newPassword }"
                />
                <p v-if="errors.newPassword" class="form-error">{{ errors.newPassword }}</p>
              </div>

              <div>
                <label for="confirmPassword" class="form-label">Confirm New Password</label>
                <input
                  id="confirmPassword"
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  required
                  class="form-input"
                  :class="{ 'border-red-500': errors.confirmPassword }"
                />
                <p v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</p>
              </div>

              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="authStore.isLoading"
                  class="btn btn-primary"
                >
                  <span v-if="authStore.isLoading" class="loading-spinner w-4 h-4 mr-2"></span>
                  {{ authStore.isLoading ? 'Changing...' : 'Change Password' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Account Actions -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold">Account Actions</h2>
          </div>
          <div class="card-body">
            <div class="space-y-4">
              <div class="flex items-center justify-between p-4 bg-red-50 dark:bg-red-900/20 rounded-lg">
                <div>
                  <h3 class="font-medium text-red-900 dark:text-red-100">Delete Account</h3>
                  <p class="text-sm text-red-700 dark:text-red-300">
                    Permanently delete your account and all associated data
                  </p>
                </div>
                <button
                  @click="deleteAccount"
                  class="btn btn-danger"
                >
                  Delete Account
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { handleApiCall, Utils } from '@/utils';
import { userService } from '@/services/api'

const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUIStore()

const profileForm = reactive({
  name: '',
  email: '',
  profile_image: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const errors = reactive({})

const validateProfileForm = () => {
  const newErrors = {}

  if (!profileForm.name) {
    newErrors.name = 'Name is required'
  } else if (profileForm.name.length < 2) {
    newErrors.name = 'Name must be at least 2 characters'
  }

  if (!profileForm.email) {
    newErrors.email = 'Email is required'
  } else if (!Utils.isValidEmail(profileForm.email)) {
    newErrors.email = 'Please enter a valid email address'
  }

  Object.assign(errors, newErrors)
  return Object.keys(newErrors).length === 0
}

const validatePasswordForm = () => {
  const newErrors = {}

  if (!passwordForm.currentPassword) {
    newErrors.currentPassword = 'Current password is required'
  }

  if (!passwordForm.newPassword) {
    newErrors.newPassword = 'New password is required'
  } else if (!Utils.isValidPassword(passwordForm.newPassword)) {
    newErrors.newPassword = 'Password must be at least 6 characters'
  }

  if (!passwordForm.confirmPassword) {
    newErrors.confirmPassword = 'Please confirm your new password'
  } else if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    newErrors.confirmPassword = 'Passwords do not match'
  }

  Object.assign(errors, newErrors)
  return Object.keys(newErrors).length === 0
}

const updateProfile = async () => {
  // Clear previous errors
  Object.keys(errors).forEach(key => delete errors[key])

  if (!validateProfileForm()) {
    return
  }

  const result = await handleApiCall(() => authStore.updateProfile(profileForm))

  // If the profile was updated successfully, log an activity entry
  if (result && result.success) {
    try {
      await userService.logActivity({
        action: 'ProfileUpdated',
        description: 'Updated profile information'
      })
    } catch (error) {
      console.warn('Failed to log profile update activity:', error)
    }
  }
}

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  try {
    
    uiStore.showInfo('Uploading image...')
    const response = await userService.uploadUserImage(file)
    
      profileForm.profile_image = response.file_url
      uiStore.showSuccess('Image uploaded successfully!')
    
  } catch (error) {
    console.error('Upload failed', error)
    uiStore.showError('Failed to upload image')
  }
}

const changePassword = async () => {
  // Clear previous errors
  Object.keys(errors).forEach(key => delete errors[key])

  if (!validatePasswordForm()) {
    return
  }

  const result = await handleApiCall(() => authStore.changePassword(
    passwordForm.currentPassword,
    passwordForm.newPassword
  ));

  if (result && result.success) {
    // Clear form
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  }
}

const deleteAccount = async () => {
  const confirmed = await uiStore.confirm(
    'Are you sure you want to delete your account? This action cannot be undone and all your data will be permanently deleted.',
    'Delete Account'
  )

  if (confirmed) {
    uiStore.showError('Account deletion functionality coming soon!')
  }
}

const loadUserData = () => {
  if (authStore.currentUser) {
    profileForm.name = authStore.currentUser.name || ''
    profileForm.email = authStore.currentUser.email || ''
    profileForm.profile_image = authStore.currentUser.profile_image
  }
}

onMounted(() => {
  uiStore.setPageTitle('Profile Settings')
  loadUserData()
})
</script>
