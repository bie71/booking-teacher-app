<template>
  <div class="container mx-auto py-10 px-4 max-w-md">
    <h1 class="text-3xl font-bold mb-6 text-center">Reset Password</h1>
    <div v-if="tokenChecked && !tokenValid" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      The reset link is invalid or has expired.
    </div>
    <div v-else-if="tokenChecked && tokenValid">
      <div v-if="successMessage" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
        {{ successMessage }}
      </div>
      <div v-if="errorMessage" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
        {{ errorMessage }}
      </div>
      <form @submit.prevent="submit" class="space-y-4 bg-white p-6 rounded shadow">
        <div>
          <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
          <input
            v-model="newPassword"
            type="password"
            id="newPassword"
            required
            class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            placeholder="Enter new password"
          />
        </div>
        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
          <input
            v-model="confirmPassword"
            type="password"
            id="confirmPassword"
            required
            class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            placeholder="Confirm new password"
          />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none"
        >
          {{ loading ? 'Resetting...' : 'Reset Password' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { userService } from '@/services/api'

const route = useRoute()
const router = useRouter()

const tokenValid = ref(false)
const tokenChecked = ref(false)
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

const token = ref('')

onMounted(async () => {
  token.value = route.query.token || ''
  if (!token.value) {
    tokenValid.value = false
    tokenChecked.value = true
    return
  }
  try {
    await userService.verifyResetToken(token.value)
    tokenValid.value = true
  } catch (error) {
    tokenValid.value = false
  } finally {
    tokenChecked.value = true
  }
})

async function submit() {
  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = 'Passwords do not match.'
    return
  }
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await userService.resetPassword({ token: token.value, new_password: newPassword.value })
    successMessage.value = 'Your password has been reset successfully. You can now log in.'
    // Optionally redirect to login after a short delay
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (error) {
    errorMessage.value = error.message || 'Failed to reset password. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>