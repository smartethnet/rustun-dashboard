import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('../components/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/topology'
      },
      {
        path: 'topology',
        name: 'Topology',
        component: () => import('../views/Topology.vue'),
      },
      {
        path: 'ai',
        name: 'AIChat',
        component: () => import('../views/AIChat.vue'),
      }
    ]
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const username = localStorage.getItem('username')
  const password = localStorage.getItem('password')
  const isAuthenticated = !!(username && password)
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
