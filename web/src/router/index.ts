import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'Login', component: () => import('@/views/LoginView.vue') },
    {
      path: '/',
      component: () => import('@/views/Layout.vue'),
      redirect: '/home',
      children: [
        { path: 'home', name: 'Home', component: () => import('@/views/HomeView.vue') },
        { path: 'records', name: 'Records', component: () => import('@/views/RecordsView.vue') },
        { path: 'stats', name: 'Stats', component: () => import('@/views/StatsView.vue') },
        { path: 'notes', name: 'Notes', component: () => import('@/views/NotesView.vue') },
        { path: 'tickets', name: 'Tickets', component: () => import('@/views/TicketListView.vue') },
        { path: 'qrcodes', name: 'QRCodes', component: () => import('@/views/QRCodeManageView.vue') },
        { path: 'settings', name: 'Settings', component: () => import('@/views/SettingsView.vue') },
        { path: 'backup', name: 'Backup', component: () => import('@/views/BackupView.vue') },
      ],
    },
    { path: '/public/ticket/new', name: 'TicketSubmit', component: () => import('@/views/TicketSubmitView.vue') },
    { path: '/public/ticket/status', name: 'TicketStatus', component: () => import('@/views/TicketStatusView.vue') },
  ],
})

router.beforeEach((to, _from) => {
  const token = localStorage.getItem('token')
  if (to.name !== 'Login' && to.name !== 'TicketSubmit' && to.name !== 'TicketStatus' && !token) {
    return { name: 'Login' }
  }
})

export default router
