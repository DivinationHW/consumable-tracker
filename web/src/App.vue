<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElNotification } from 'element-plus'

const auth = useAuthStore()
let ws = null

function connectWS() {
  if (!auth.token || auth.role !== 'admin') return
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const url = `${protocol}//${location.host}/ws?token=${auth.token}`
  ws = new WebSocket(url)
  ws.onmessage = (e) => {
    const data = JSON.parse(e.data)
    if (data.type === 'new_ticket') {
      ElNotification({ title: '新报修工单', message: `房间 ${data.ticket.room_number} - ${data.ticket.problem_type}`, type: 'warning', duration: 0 })
    }
  }
  ws.onclose = () => setTimeout(connectWS, 5000)
}

onMounted(() => {
  if (auth.token) connectWS()
})
</script>
