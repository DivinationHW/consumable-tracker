<template>
  <div class="page-container">
    <div class="page-header"><h2>系统设置</h2></div>
    <el-card style="max-width:600px;">
      <el-form label-position="top">
        <el-form-item label="当前角色">
          <el-tag :type="auth.isAdmin ? 'danger' : 'info'" size="large">
            {{ auth.isAdmin ? '管理员' : '只读用户' }}
          </el-tag>
        </el-form-item>
        <el-divider />
        <el-form-item label="当前密码">
          <el-input v-model="form.password" type="password" show-password placeholder="输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="form.new_password" type="password" show-password placeholder="输入新密码" />
        </el-form-item>
        <el-form-item label="确认新密码">
          <el-input v-model="form.confirm_password" type="password" show-password placeholder="再次输入新密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleChangePassword" :loading="saving">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { changePassword } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const auth = useAuthStore()
const saving = ref(false)
const form = ref({ password: '', new_password: '', confirm_password: '' })

async function handleChangePassword() {
  if (!form.value.password || !form.value.new_password) {
    ElMessage.warning('请填写完整')
    return
  }
  if (form.value.new_password !== form.value.confirm_password) {
    ElMessage.warning('两次密码不一致')
    return
  }
  if (form.value.new_password.length < 6) {
    ElMessage.warning('密码至少6位')
    return
  }
  saving.value = true
  try {
    await changePassword(form.value.password, form.value.new_password)
    ElMessage.success('密码修改成功')
    form.value = { password: '', new_password: '', confirm_password: '' }
  } finally { saving.value = false }
}
</script>
