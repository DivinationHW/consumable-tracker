<template>
  <el-container style="height:100vh">
    <el-aside width="200px">
      <el-menu :default-active="route.path" router style="height:100%">
        <el-menu-item index="/home"><el-icon><HomeFilled /></el-icon>首页录入</el-menu-item>
        <el-menu-item index="/records"><el-icon><List /></el-icon>记录查询</el-menu-item>
        <el-menu-item index="/stats"><el-icon><DataAnalysis /></el-icon>统计报表</el-menu-item>
        <el-menu-item v-if="auth.isAdmin" index="/notes"><el-icon><EditPen /></el-icon>便签</el-menu-item>
        <el-menu-item index="/tickets"><el-icon><Ticket /></el-icon>工单管理</el-menu-item>
        <el-menu-item index="/qrcodes"><el-icon><Camera /></el-icon>二维码管理</el-menu-item>
        <el-menu-item index="/settings"><el-icon><Setting /></el-icon>设置</el-menu-item>
        <el-menu-item v-if="auth.isAdmin" index="/backup"><el-icon><Folder /></el-icon>备份管理</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="display:flex;align-items:center;justify-content:space-between;border-bottom:1px solid #eee">
        <span style="font-size:16px;font-weight:bold">耗材使用管理系统</span>
        <el-dropdown @command="handleCommand">
          <span style="cursor:pointer">{{ auth.user?.username }} ({{ auth.user?.role }}) <el-icon><ArrowDown /></el-icon></span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="password">修改密码</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main style="background:#f5f7fa">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
  <el-dialog v-model="pwDialog" title="修改密码" width="400px">
    <el-form :model="pwForm">
      <el-form-item label="原密码"><el-input v-model="pwForm.old" type="password" show-password /></el-form-item>
      <el-form-item label="新密码"><el-input v-model="pwForm.newPw" type="password" show-password /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="pwDialog=false">取消</el-button>
      <el-button type="primary" @click="changePw">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { changePassword } from '@/api/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const pwDialog = ref(false)
const pwForm = reactive({ old: '', newPw: '' })

function handleCommand(cmd: string) {
  if (cmd === 'logout') { auth.logout(); router.push('/login') }
  if (cmd === 'password') pwDialog.value = true
}

async function changePw() {
  await changePassword(pwForm.old, pwForm.newPw)
  ElMessage.success('密码修改成功')
  pwDialog.value = false
  pwForm.old = ''
  pwForm.newPw = ''
}
</script>
