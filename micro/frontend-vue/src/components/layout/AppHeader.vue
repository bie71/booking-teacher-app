<template>
  <header class="fixed top-0 left-0 right-0 z-50 bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
    <nav class="container mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <div class="flex items-center">
          <router-link to="/" class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-sm">JL</span>
            </div>
            <span class="text-xl font-bold text-gray-900 dark:text-white">JapanLearn</span>
          </router-link>
        </div>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-8">
          <!--
            Navigation links with icons.  Each link now includes a small
            SVG before the label to provide a consistent visual language
            across the application.  The icons are inline so they don't
            require additional dependencies and match the heroicon
            aesthetic used elsewhere in the app.
          -->
          <router-link
            to="/"
            class="inline-flex items-center text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 text-sm font-medium transition-colors"
            :class="{ 'text-blue-600 dark:text-blue-400': $route.name === 'Home' }"
          >
            <!-- Home icon -->
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9.75l9-7 9 7v9a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18.75v-9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 22V12h6v10" />
            </svg>
            <span>Home</span>
          </router-link>
          <router-link
            to="/teachers"
            class="inline-flex items-center text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 text-sm font-medium transition-colors"
            :class="{ 'text-blue-600 dark:text-blue-400': $route.name === 'Teachers' }"
          >
            <!-- Search/teachers icon -->
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35M11 19a8 8 0 100-16 8 8 0 000 16z" />
            </svg>
            <span>Find Teachers</span>
          </router-link>
          
          <!-- Authenticated Navigation -->
          <template v-if="authStore.isAuthenticated">
          <router-link
            :to="dashboardRoute"
            class="inline-flex items-center text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 text-sm font-medium transition-colors"
            :class="{ 'text-blue-600 dark:text-blue-400': isDashboardActive }"
          >
            <!-- Dashboard/grid icon -->
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h7v7H3V4zm0 9h7v7H3v-7zm9-9h9v7h-9V4zm0 9h9v7h-9v-7z" />
            </svg>
            <span>Dashboard</span>
          </router-link>
            <router-link
              to="/bookings"
              class="inline-flex items-center text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 text-sm font-medium transition-colors"
              :class="{ 'text-blue-600 dark:text-blue-400': $route.name === 'Bookings' }"
            >
              <!-- Calendar/bookings icon -->
              <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span>My Bookings</span>
            </router-link>
          </template>
        </div>

        <!-- Right Side Actions -->
        <div class="flex items-center space-x-4">
          <!-- Theme Toggle -->
          <button
            @click="uiStore.toggleTheme"
            class="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
            title="Toggle theme"
          >
            <svg v-if="uiStore.currentTheme === 'light'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </button>

          <!-- Authentication Actions -->
          <div v-if="!authStore.isAuthenticated" class="hidden md:flex items-center space-x-2">
            <router-link
              to="/login"
              class="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 text-sm font-medium transition-colors"
            >
              Login
            </router-link>
            <router-link
              to="/register"
              class="btn btn-primary text-sm"
            >
              Sign Up
            </router-link>
          </div>

          <!-- User Menu -->
          <div v-else class="relative">
            <button
              @click="showUserMenu = !showUserMenu"
              class="flex items-center space-x-2 text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <div class="w-8 h-8 rounded-full overflow-hidden flex items-center justify-center bg-blue-600">
                <template v-if="authStore.currentUser?.profile_image">
                  <img 
                    :src="authStore.currentUser.profile_image" 
                    alt="User Avatar" 
                    class="w-full h-full object-cover"
                  />
                </template>
                <template v-else>
                  <span class="text-white text-sm font-medium">
                    {{ authStore.currentUser?.name?.charAt(0).toUpperCase() }}
                  </span>
                </template>
              </div>
              <span class="hidden md:block text-sm font-medium">
                {{ authStore.currentUser?.name }}
              </span>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>

            <!-- User Dropdown -->
            <transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div
                v-if="showUserMenu"
                class="dropdown-menu"
                @click="showUserMenu = false"
              >
                <router-link to="/profile" class="dropdown-item">
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  Profile
                </router-link>
                <router-link to="/dashboard" class="dropdown-item">
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  Dashboard
                </router-link>
                <router-link to="/bookings" class="dropdown-item">
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  My Bookings
                </router-link>
                <hr class="my-1 border-gray-200 dark:border-gray-600">
                <button @click="handleLogout" class="dropdown-item w-full text-left text-red-600 dark:text-red-400">
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                  Logout
                </button>
              </div>
            </transition>
          </div>

          <!-- Mobile Menu Button -->
          <button
            @click="uiStore.toggleSidebar"
            class="md:hidden p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
      </div>
    </nav>

    <!-- Mobile Sidebar -->
    <transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="transform -translate-x-full"
      enter-to-class="transform translate-x-0"
      leave-active-class="transition ease-in duration-300"
      leave-from-class="transform translate-x-0"
      leave-to-class="transform -translate-x-full"
    >
      <div
        v-if="uiStore.isSidebarOpen"
        class="md:hidden fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-gray-800 shadow-lg"
      >
        <div class="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <div class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-sm">JL</span>
            </div>
            <span class="text-xl font-bold text-gray-900 dark:text-white">JapanLearn</span>
          </div>
          <button
            @click="uiStore.closeSidebar"
            class="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <nav class="p-4 space-y-2">
          <router-link
            to="/"
            @click="uiStore.closeSidebar"
            class="flex items-center px-3 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
          >
            <!-- Home icon -->
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9.75l9-7 9 7v9a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18.75v-9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 22V12h6v10" />
            </svg>
            <span>Home</span>
          </router-link>
          <router-link
            to="/teachers"
            @click="uiStore.closeSidebar"
            class="flex items-center px-3 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
          >
            <!-- Search/teachers icon -->
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35M11 19a8 8 0 100-16 8 8 0 000 16z" />
            </svg>
            <span>Find Teachers</span>
          </router-link>

          <template v-if="authStore.isAuthenticated">
            <router-link
              to="/dashboard"
              @click="uiStore.closeSidebar"
              class="flex items-center px-3 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
            >
              <!-- Dashboard/grid icon -->
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h7v7H3V4zm0 9h7v7H3v-7zm9-9h9v7h-9V4zm0 9h9v7h-9v-7z" />
              </svg>
              <span>Dashboard</span>
            </router-link>
            <router-link
              to="/bookings"
              @click="uiStore.closeSidebar"
              class="flex items-center px-3 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
            >
              <!-- Calendar/bookings icon -->
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span>My Bookings</span>
            </router-link>
            <router-link
              to="/profile"
              @click="uiStore.closeSidebar"
              class="flex items-center px-3 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
            >
              <!-- User/profile icon -->
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              <span>Profile</span>
            </router-link>
            <button @click="doLogout" class="px-3 py-2 rounded-md bg-rose-600 text-white hover:bg-rose-700">Logout</button>
          </template>

          <template v-else>
            <router-link
              to="/login"
              @click="uiStore.closeSidebar"
              class="flex items-center px-3 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition-colors"
            >
              <!-- Login icon (arrow in box) -->
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M12 5l7 7-7 7M5 5h2a2 2 0 012 2v2" />
              </svg>
              <span>Login</span>
            </router-link>
            <router-link
              to="/register"
              @click="uiStore.closeSidebar"
              class="flex items-center px-3 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-md transition-colors text-center"
            >
              <!-- User add icon -->
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9a3 3 0 11-6 0 3 3 0 016 0zM4 15a4 4 0 018 0v5H4v-5zm13 3v-2m0 0v-2m0 2h2m-2 0h-2" />
              </svg>
              <span>Sign Up</span>
            </router-link>
          </template>
        </nav>
      </div>
    </transition>

    <!-- Mobile Sidebar Backdrop -->
    <div
      v-if="uiStore.isSidebarOpen"
      @click="uiStore.closeSidebar"
      class="md:hidden fixed inset-0 bg-black bg-opacity-50 z-40"
    ></div>
  </header>
</template>

<script setup>
import { authService } from '@/services/api'
import { useRouter } from 'vue-router'
const doLogout = async()=>{ await authService.logout(); router.push('/login') }
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'

const authStore = useAuthStore()
const uiStore = useUIStore()
const showUserMenu = ref(false)
const router = useRouter()

// Computed properties for role-based dashboard routing
const dashboardRoute = computed(() => {
  if (!authStore.isAuthenticated) return '/login'
  
  switch (authStore.userRole) {
    case 'admin':
      return '/admin/dashboard'
    case 'teacher':
      return '/teacher/dashboard'
    case 'user':
    default:
      return '/student/dashboard'
  }
})

const isDashboardActive = computed(() => {
  const currentRoute = uiStore.currentRoute?.name || ''
  return ['StudentDashboard', 'TeacherDashboard', 'AdminDashboard'].includes(currentRoute)
})

const handleLogout = async () => {
  showUserMenu.value = false
  uiStore.closeSidebar()
  
  const confirmed = await uiStore.confirm('Are you sure you want to logout?', 'Logout')
  if (confirmed) {
    await authStore.logout()
    uiStore.showSuccess('Successfully logged out')
    router.push('/login')
  }
}

// Close user menu when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
