<template>
  <div style="height: 100vh; display: flex; align-items: center; justify-content: center; background: #f5f7fa;">
    <el-card style="width: 400px; padding: 20px;">
      <h2 style="text-align: center; margin-bottom: 24px;">耗材使用明细系统</h2>
      <el-form @submit.prevent="handleLogin" label-position="top">
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" show-password placeholder="请输入密码" @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading" style="width: 100%">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login as apiLogin } from '@/api'

const router = useRouter()
const auth = useAuthStore()
const username = ref('')
const password = ref('')
const loading = ref(false)

async function handleLogin() {
  if (!username.value || !password.value) return
  loading.value = true
  try {
    const data = await apiLogin(username.value, password.value)
    auth.setToken(data.token, data.role)
    router.push('/')
  } finally {
    loading.value = false
  }
}
</script>
