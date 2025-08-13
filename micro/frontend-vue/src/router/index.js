import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'

// Import views
import Home from '@/views/Home.vue'
import Teachers from '@/views/Teachers.vue'
import TeacherDetail from '@/views/TeacherDetail.vue'
import Dashboard from '@/views/Dashboard.vue'
import BookingManagement from '@/views/BookingManagement.vue'
import BookingInvoice from '@/views/BookingInvoice.vue'
import BookingDetail from '@/views/BookingDetail.vue'
import Bookings from '@/views/Bookings.vue'
import Profile from '@/views/Profile.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import NotFound from '@/views/NotFound.vue'

const routes = [
  { path: '/admin/booking/:id', name: 'BookingDetail', component: BookingDetail, meta: { requiresAuth: true, admin: true, title: 'Booking Detail' } },
  { path: '/admin/booking/:id/invoice', name: 'BookingInvoice', component: BookingInvoice, meta: { requiresAuth: true, admin: true, title: 'Booking Invoice' } },
  { path: '/admin/bookings', name: 'BookingManagement', component: BookingManagement, meta: { requiresAuth: true, admin: true, title: 'Booking Management' } },
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'Home',
      description: 'Learn Japanese with native teachers - JapanLearn'
    }
  },
  {
    path: '/teachers',
    name: 'Teachers',
    component: Teachers,
    meta: {
      title: 'Find Teachers',
      description: 'Browse and find the perfect Japanese teacher for you'
    }
  },
  {
    path: '/teachers/:id',
    name: 'TeacherDetail',
    component: TeacherDetail,
    props: true,
    meta: {
      title: 'Teacher Profile',
      description: 'View teacher profile and book lessons'
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: {
      title: 'Login',
      description: 'Login to your JapanLearn account',
      guest: true
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: {
      title: 'Register',
      description: 'Create your JapanLearn account',
      guest: true
    }
  },
  // Password reset routes. Allow users to request a reset link and
  // recover their password using the token from their email.
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/ForgotPassword.vue'),
    meta: {
      title: 'Forgot Password',
      description: 'Request a password reset link',
      guest: true
    }
  },
  {
    path: '/recover-password',
    name: 'RecoverPassword',
    component: () => import('@/views/RecoverPassword.vue'),
    meta: {
      title: 'Reset Password',
      description: 'Set a new password using your reset token',
      guest: true
    }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      title: 'Dashboard',
      description: 'Your learning dashboard',
      requiresAuth: true
    }
  },
  {
    path: '/student/dashboard',
    name: 'StudentDashboard',
    component: () => import('@/views/StudentDashboard.vue'),
    meta: {
      title: 'Student Dashboard',
      description: 'Your learning dashboard',
      requiresAuth: true,
      role: 'user'
    }
  },
  {
    path: '/teacher/dashboard',
    name: 'TeacherDashboard',
    component: () => import('@/views/TeacherDashboard.vue'),
    meta: {
      title: 'Teacher Dashboard',
      description: 'Your teaching dashboard',
      requiresAuth: true,
      role: 'teacher'
    }
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: () => import('@/views/AdminDashboard.vue'),
    meta: {
      title: 'Admin Dashboard',
      description: 'System management dashboard',
      requiresAuth: true,
      role: 'admin'
    }
  },
  {
    path: '/admin/users',
    name: 'UserManagement',
    component: () => import('@/views/UserManagement.vue'),
    meta: {
      title: 'User Management',
      description: 'Manage system users',
      requiresAuth: true,
      role: 'admin'
    }
  },
  {
    path: '/admin/payment-methods',
    name: 'PaymentMethodManagement',
    component: () => import('@/views/PaymentMethodManagement.vue'),
    meta: {
      title: 'Payment Method Management',
      description: 'Manage payment methods',
      requiresAuth: true,
      role: 'admin'
    }
  },
  {
    path: '/bookings',
    name: 'Bookings',
    component: Bookings,
    meta: {
      title: 'My Bookings',
      description: 'Manage your lesson bookings',
      requiresAuth: true
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: {
      title: 'Profile',
      description: 'Manage your profile settings',
      requiresAuth: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: {
      title: 'Page Not Found'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    } else {
      return { top: 0 }
    }
  }
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const uiStore = useUIStore()

  // Set page title and meta description
  if (to.meta.title) {
    uiStore.setPageTitle(to.meta.title)
  }
  
  if (to.meta.description) {
    uiStore.setMetaDescription(to.meta.description)
  }

  // Check if route requires authentication
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      uiStore.showError('Please log in to access this page')
      next({
        name: 'Login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // Check role-based access
    if (to.meta.role) {
      const userRole = authStore.userRole
      if (userRole !== to.meta.role) {
        uiStore.showError('Access denied: insufficient permissions')
        
        // Redirect to appropriate dashboard based on role
        switch (userRole) {
          case 'admin':
            next({ name: 'AdminDashboard' })
            break
          case 'teacher':
            next({ name: 'TeacherDashboard' })
            break
          case 'user':
          default:
            next({ name: 'StudentDashboard' })
            break
        }
        return
      }
    }
  }

  // Check if route is for guests only (login/register)
  if (to.meta.guest && authStore.isAuthenticated) {
    // Redirect to appropriate dashboard based on role
    const userRole = authStore.userRole
    switch (userRole) {
      case 'admin':
        next({ name: 'AdminDashboard' })
        break
      case 'teacher':
        next({ name: 'TeacherDashboard' })
        break
      case 'user':
      default:
        next({ name: 'StudentDashboard' })
        break
    }
    return
  }

  // Close sidebar on route change (mobile)
  uiStore.closeSidebar()

  next()
})

router.afterEach((to, from) => {
  // Close any open modals after navigation
  const uiStore = useUIStore()
  uiStore.closeAllModals()
})

// Handle navigation errors
router.onError((error) => {
  console.error('Router error:', error)
  const uiStore = useUIStore()
  uiStore.showError('Navigation error occurred')
})

export default router
