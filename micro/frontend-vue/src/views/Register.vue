<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <div class="flex justify-center">
          <div class="w-12 h-12 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-lg">JL</span>
          </div>
        </div>
        <h2 class="mt-6 text-center text-3xl font-bold text-gray-900 dark:text-white">
          Create your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          Or
          <router-link to="/login" class="font-medium text-blue-600 hover:text-blue-500">
            sign in to your existing account
          </router-link>
        </p>
      </div>
      
      <form @submit.prevent="handleRegister" class="mt-8 space-y-6">
        <div class="space-y-4">
          <div>
            <label for="name" class="form-label">Full Name</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.name }"
              placeholder="Enter your full name"
            />
            <p v-if="errors.name" class="form-error">{{ errors.name }}</p>
          </div>

          <div>
            <label for="email" class="form-label">Email address</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.email }"
              placeholder="Enter your email"
            />
            <p v-if="errors.email" class="form-error">{{ errors.email }}</p>
          </div>
          
          <div>
            <label for="password" class="form-label">Password</label>
            <div class="relative">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="form-input pr-10"
                :class="{ 'border-red-500': errors.password }"
                placeholder="Create a password"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg v-if="showPassword" class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
                </svg>
                <svg v-else class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </button>
            </div>
            <p v-if="errors.password" class="form-error">{{ errors.password }}</p>
            <div class="mt-1">
              <div class="flex items-center space-x-1 text-xs">
                <div :class="passwordStrength.length >= 6 ? 'text-green-600' : 'text-gray-400'">
                  ✓ At least 6 characters
                </div>
              </div>
            </div>
          </div>

          <div>
            <label for="confirmPassword" class="form-label">Confirm Password</label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.confirmPassword }"
              placeholder="Confirm your password"
            />
            <p v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</p>
          </div>
        </div>

        <div class="flex items-center">
          <input
            id="agree-terms"
            v-model="form.agreeTerms"
            type="checkbox"
            required
            class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label for="agree-terms" class="ml-2 block text-sm text-gray-900 dark:text-gray-300">
            I agree to the
            <router-link to="/terms" class="text-blue-600 hover:text-blue-500">Terms of Service</router-link>
            and
            <router-link to="/privacy" class="text-blue-600 hover:text-blue-500">Privacy Policy</router-link>
          </label>
        </div>

        <div>
          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="btn btn-primary w-full flex justify-center py-3"
          >
            <span v-if="authStore.isLoading" class="loading-spinner w-5 h-5 mr-2"></span>
            {{ authStore.isLoading ? 'Creating account...' : 'Create account' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { handleApiCall, Utils } from '@/utils';

const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUIStore()

const showPassword = ref(false)
const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeTerms: false
})
const errors = reactive({})

const passwordStrength = computed(() => {
  return {
    length: form.password.length
  }
})

const validateForm = () => {
  const newErrors = {}

  if (!form.name) {
    newErrors.name = 'Name is required'
  } else if (form.name.length < 2) {
    newErrors.name = 'Name must be at least 2 characters'
  }

  if (!form.email) {
    newErrors.email = 'Email is required'
  } else if (!Utils.isValidEmail(form.email)) {
    newErrors.email = 'Please enter a valid email address'
  }

  if (!form.password) {
    newErrors.password = 'Password is required'
  } else if (!Utils.isValidPassword(form.password)) {
    newErrors.password = 'Password must be at least 6 characters'
  }

  if (!form.confirmPassword) {
    newErrors.confirmPassword = 'Please confirm your password'
  } else if (form.password !== form.confirmPassword) {
    newErrors.confirmPassword = 'Passwords do not match'
  }

  if (!form.agreeTerms) {
    newErrors.agreeTerms = 'You must agree to the terms and conditions'
  }

  Object.assign(errors, newErrors)
  return Object.keys(newErrors).length === 0
}

const handleRegister = async () => {
  // Clear previous errors
  Object.keys(errors).forEach(key => delete errors[key])

  if (!validateForm()) {
    return
  }

  const result = await handleApiCall(() => authStore.register(form.name, form.email, form.password));

  if (result && result.success) {
    router.push('/dashboard')
  }
}

onMounted(() => {
  uiStore.setPageTitle('Register')
})
</script>
