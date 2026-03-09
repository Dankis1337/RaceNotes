import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue'), meta: { guest: true } },
  { path: '/register', name: 'Register', component: () => import('../views/RegisterView.vue'), meta: { guest: true } },
  { path: '/', name: 'Races', component: () => import('../views/RacesListView.vue') },
  { path: '/races/new', name: 'RaceCreate', component: () => import('../views/RaceFormView.vue') },
  { path: '/races/:id', name: 'RaceDetail', component: () => import('../views/RaceDetailView.vue') },
  { path: '/races/:id/edit', name: 'RaceEdit', component: () => import('../views/RaceFormView.vue') },
  { path: '/setups', name: 'Setups', component: () => import('../views/SetupsListView.vue') },
  { path: '/setups/new', name: 'SetupCreate', component: () => import('../views/SetupFormView.vue') },
  { path: '/setups/:id', name: 'SetupDetail', component: () => import('../views/SetupDetailView.vue') },
  { path: '/setups/:id/edit', name: 'SetupEdit', component: () => import('../views/SetupFormView.vue') },
  { path: '/calculator', name: 'Calculator', component: () => import('../views/CalculatorView.vue') },
  { path: '/profile', name: 'Profile', component: () => import('../views/ProfileView.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (!to.meta.guest && !token) {
    return { name: 'Login' }
  }
  if (to.meta.guest && token) {
    return { name: 'Races' }
  }
})

export default router
