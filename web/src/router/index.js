import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/login', name: 'Login', component: () => import('@/views/Login.vue'), meta: { guest: true } },
  {
    path: '/',
    component: () => import('@/components/Layout.vue'),
    meta: { auth: true },
    children: [
      { path: '', name: 'Dashboard', component: () => import('@/views/Dashboard.vue') },
      { path: 'consumables', name: 'Consumables', component: () => import('@/views/Consumables.vue') },
      { path: 'offices', name: 'Offices', component: () => import('@/views/Offices.vue') },
      { path: 'records', name: 'Records', component: () => import('@/views/Records.vue') },
      { path: 'stats', name: 'Stats', component: () => import('@/views/Stats.vue') },
      { path: 'tickets', name: 'Tickets', component: () => import('@/views/Tickets.vue') },
      { path: 'qrcodes', name: 'QRCodes', component: () => import('@/views/QRCodes.vue') },
      { path: 'backup', name: 'Backup', component: () => import('@/views/Backup.vue') },
      { path: 'notes', name: 'Notes', component: () => import('@/views/Notes.vue') },
      { path: 'users', name: 'Users', component: () => import('@/views/Users.vue'), meta: { admin: true } },
      { path: 'settings', name: 'Settings', component: () => import('@/views/Settings.vue') },
    ],
  },
  { path: '/ticket', name: 'TicketForm', component: () => import('@/views/TicketForm.vue'), meta: { guest: true } },
  { path: '/ticket/:id', name: 'TicketStatus', component: () => import('@/views/TicketStatusView.vue'), meta: { guest: true } },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.token) return next('/login')
  if (to.meta.guest && auth.token) return next('/')
  next()
})

export default router
