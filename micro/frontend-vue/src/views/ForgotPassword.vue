<template>
  <div class="container mx-auto mt-10 py-10 px-4 max-w-md">
    <h1 class="text-3xl font-bold mb-6 text-center">Forgot Password</h1>
    <div v-if="successMessage" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
      {{ successMessage }}
    </div>
    <div v-if="errorMessage" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {{ errorMessage }}
    </div>
    <form @submit.prevent="submit" class="space-y-4 bg-white p-6 rounded shadow">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
        <input
          v-model="email"
          type="email"
          id="email"
          required
          class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
          placeholder="Enter your email"
        />
      </div>
      <button
        type="submit"
        :disabled="loading"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none"
      >
        {{ loading ? 'Sending...' : 'Send Reset Link' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { userService } from '@/services/api'

const email = ref('')
const loading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

async function submit() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await userService.forgotPassword(email.value)
    successMessage.value = 'If the email exists in our system, a password reset link has been sent.'
    email.value = ''
  } catch (error) {
    errorMessage.value = error.message || 'Something went wrong. Please try again later.'
  } finally {
    loading.value = false
  }
}
</script>