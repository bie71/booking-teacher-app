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
          Sign in to your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          Or
          <router-link to="/register" class="font-medium text-blue-600 hover:text-blue-500">
            create a new account
          </router-link>
        </p>
      </div>
      
      <form @submit.prevent="handleLogin" class="mt-8 space-y-6">
        <div class="space-y-4">
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
                placeholder="Enter your password"
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
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              v-model="form.rememberMe"
              type="checkbox"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-900 dark:text-gray-300">
              Remember me
            </label>
          </div>

          <div class="text-sm">
            <router-link to="/forgot-password" class="font-medium text-blue-600 hover:text-blue-500">
              Forgot your password?
            </router-link>
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="btn btn-primary w-full flex justify-center py-3"
          >
            <span v-if="authStore.isLoading" class="loading-spinner w-5 h-5 mr-2"></span>
            {{ authStore.isLoading ? 'Signing in...' : 'Sign in' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { handleApiCall, Utils } from '@/utils';

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const uiStore = useUIStore()

const showPassword = ref(false)
const form = reactive({
  email: '',
  password: '',
  rememberMe: false
})
const errors = reactive({})

const validateForm = () => {
  const newErrors = {}

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

  Object.assign(errors, newErrors)
  return Object.keys(newErrors).length === 0
}

  const handleLogin = async () => {
    // Clear previous errors
    Object.keys(errors).forEach(key => delete errors[key])

    if (!validateForm()) {
      return
    }

    const result = await handleApiCall(() => authStore.login(form.email, form.password));

    if (result && result.success) {
      // Redirect based on user role
      const userRole = authStore.userRole
      let dashboardRoute = '/dashboard'
      
      switch (userRole) {
        case 'admin':
          dashboardRoute = '/admin/dashboard'
          break
        case 'teacher':
          dashboardRoute = '/teacher/dashboard'
          break
        case 'user':
        default:
          dashboardRoute = '/student/dashboard'
          break
      }

      // Check if there's a redirect query parameter
      const redirectTo = route.query.redirect || dashboardRoute
      console.log('Redirecting to:', redirectTo)
      router.push(redirectTo)
    }
  }

onMounted(() => {
  uiStore.setPageTitle('Login')
})
</script>
