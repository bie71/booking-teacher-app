<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section 
      class="relative text-white overflow-hidden"
      :class="!heroImage ? 'bg-gradient-to-br from-blue-400 via-purple-600 to-blue-600' : ''"
      :style="heroBackgroundStyle"
    >
      <!-- Background overlay for better text readability -->
      <div class="absolute inset-0 bg-black bg-opacity-40"></div>
      
      <div class="relative container mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div class="text-center">
          <h1 class="text-4xl md:text-6xl font-bold mb-6">
            Learn Japanese with
            <span class="text-gradient-japanese block mt-2">Native Teachers</span>
          </h1>
          <p class="text-xl md:text-2xl mb-8 text-blue-100 max-w-3xl mx-auto">
            Master Japanese through personalized lessons with experienced native speakers. 
            Book your first lesson today and start your journey to fluency.
          </p>
          <p class="text-japanese text-lg mb-8 text-blue-200">
            日本語を学ぼう！ネイティブの先生と一緒に上達しましょう。
          </p>
          
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <router-link
              to="/teachers"
              class="btn btn-primary bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
            >
              Find Teachers
            </router-link>
            <router-link
              v-if="!authStore.isAuthenticated"
              to="/register"
              class="btn btn-outline border-white text-white hover:bg-white hover:text-blue-600 text-lg px-8 py-3"
            >
              Get Started Free
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="section bg-white dark:bg-gray-900">
      <div class="container">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold mb-4">
            Why Choose JapanLearn?
          </h2>
          <p class="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Experience the best way to learn Japanese with our comprehensive platform
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Native Teachers</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Learn from qualified native Japanese speakers with years of teaching experience
            </p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Flexible Scheduling</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Book lessons at your convenience with our easy-to-use scheduling system
            </p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-purple-100 dark:bg-purple-900 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Personalized Learning</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Customized lessons tailored to your level and learning goals
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Featured Teachers Section -->
    <section class="section bg-gray-50 dark:bg-gray-800">
      <div class="container">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold mb-4">
            Meet Our Featured Teachers
          </h2>
          <p class="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Discover experienced Japanese teachers ready to help you achieve your language goals
          </p>
        </div>

        <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="i in 3" :key="i" class="card">
            <div class="card-body">
              <div class="skeleton h-32 w-32 rounded-full mx-auto mb-4"></div>
              <div class="skeleton h-6 w-3/4 mx-auto mb-2"></div>
              <div class="skeleton h-4 w-full mb-2"></div>
              <div class="skeleton h-4 w-2/3 mx-auto"></div>
            </div>
          </div>
        </div>

        <div v-else-if="featuredTeachers.length > 0" class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div
            v-for="teacher in featuredTeachers"
            :key="teacher.ID || teacher.id"
            class="card hover:shadow-lg transition-shadow cursor-pointer"
            @click="$router.push(`/teachers/${teacher.ID || teacher.id}`)"
          >
            <div class="card-body text-center">
              <div class="w-32 h-32 bg-gray-200 dark:bg-gray-700 rounded-full mx-auto mb-4 overflow-hidden">
                <img
                  v-if="teacher.ProfileImage || teacher.profile_image"
                  :src="teacher.ProfileImage || teacher.profile_image"
                  :alt="teacher.Name || teacher.name"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-2xl font-bold text-gray-500">
                  {{ (teacher?.Name || teacher?.name || '').charAt(0) }}
                </div>
              </div>
              
              <h3 class="text-xl font-semibold mb-2">{{ teacher?.Name || teacher?.name }}</h3>
              <div class="flex items-center justify-center mb-2">
                <span :class="getLevelBadgeClass(teacher.LanguageLevel || teacher.language_level)" class="badge">
                  {{ Utils.capitalize(teacher.LanguageLevel || teacher.language_level) }}
                </span>
              </div>
              <p class="text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
                {{ teacher.Bio || teacher.bio }}
              </p>
              <div class="price-small">
                {{ Utils.formatCurrency(teacher.PricePerHour || teacher.price_per_hour) }}/hour
              </div>
            </div>
          </div>
        </div>

        <div class="text-center mt-12">
          <router-link to="/teachers" class="btn btn-primary">
            View All Teachers
          </router-link>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="section bg-white dark:bg-gray-900">
      <div class="container">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold mb-4">
            How It Works
          </h2>
          <p class="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Start learning Japanese in just a few simple steps
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-2xl font-bold">
              1
            </div>
            <h3 class="text-lg font-semibold mb-2">Browse Teachers</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Explore our selection of qualified Japanese teachers
            </p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-green-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-2xl font-bold">
              2
            </div>
            <h3 class="text-lg font-semibold mb-2">Book a Lesson</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Choose your preferred time and book your lesson
            </p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-purple-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-2xl font-bold">
              3
            </div>
            <h3 class="text-lg font-semibold mb-2">Learn Online</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Join your personalized lesson and start learning
            </p>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 bg-yellow-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-2xl font-bold">
              4
            </div>
            <h3 class="text-lg font-semibold mb-2">Track Progress</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Monitor your improvement and continue learning
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="section bg-blue-600 text-white">
      <div class="container text-center">
        <h2 class="text-3xl md:text-4xl font-bold mb-4">
          Ready to Start Learning Japanese?
        </h2>
        <p class="text-xl mb-8 text-blue-100 max-w-2xl mx-auto">
          Join thousands of students who are already learning Japanese with our expert teachers
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <router-link
            to="/teachers"
            class="btn btn-primary bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
          >
            Find Your Teacher
          </router-link>
          <router-link
            v-if="!authStore.isAuthenticated"
            to="/register"
            class="btn btn-outline border-white text-white hover:bg-white hover:text-blue-600 text-lg px-8 py-3"
          >
            Sign Up Free
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTeachersStore } from '@/stores/teachers'
import { useUIStore } from '@/stores/ui'
import { Utils } from '@/utils';
import { userService } from '@/services/api'
import { API_CONFIG } from '@/config'

const authStore = useAuthStore()
const teachersStore = useTeachersStore()
const uiStore = useUIStore()

const loading = ref(true)
const featuredTeachers = ref([])
const heroImage = ref(null)

// Computed property for hero background style
const heroBackgroundStyle = computed(() => {
  if (heroImage.value) {
    return {
      backgroundImage: `url(${heroImage.value})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      backgroundRepeat: 'no-repeat'
    }
  }
  return {}
})
const fetchHeroImage = async () => {
  try {
    const response = await userService.getImageHero({ key: API_CONFIG.KEY_IMAGE_HERO })
    if (response && response.image_url) {
      heroImage.value = response.image_url
    } else {
      console.warn('No hero image found, fallback to default.')
    }
  } catch (error) {
    console.error('Error fetching hero image:', error)
  }
}
// Local API call handler
const handleApiCall = async (apiCall) => {
  try {
    const response = await apiCall()
    if (response && response.message && uiStore?.showSuccess) {
      uiStore.showSuccess(response.message)
    }
    return response
  } catch (error) {
    let message = 'An unexpected error occurred'
    
    if (error?.response?.data?.message) {
      message = error.response.data.message
    } else if (error?.message) {
      message = error.message
    }
    
    if (uiStore?.showError) {
      uiStore.showError(message)
    }
    return null
  }
}

const getLevelBadgeClass = (level) => {
  return `level-${level}`
}

const fetchFeaturedTeachers = async () => {
  loading.value = true
  const result = await handleApiCall(() => teachersStore.fetchTeachers());
  if (result && result.success) {
    featuredTeachers.value = result.data.slice(0, 3)
  }
  loading.value = false
}

onMounted(() => {
  uiStore.setPageTitle('Home')
  fetchFeaturedTeachers()
  fetchHeroImage()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
