<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" style="background: #304156; color: #fff; display: flex; flex-direction: column;">
      <div style="padding: 20px; font-size: 18px; font-weight: 700; border-bottom: 1px solid rgba(255,255,255,0.1);">
        耗材管理系统
      </div>
      <el-menu :default-active="route.path" router background-color="#304156" text-color="#bfcbd9" active-text-color="#409eff" style="flex:1; border: none;">
        <el-menu-item index="/"><el-icon><DataAnalysis /></el-icon>仪表盘</el-menu-item>
        <el-menu-item index="/consumables"><el-icon><Goods /></el-icon>耗材管理</el-menu-item>
        <el-menu-item index="/offices"><el-icon><OfficeBuilding /></el-icon>科室管理</el-menu-item>
        <el-menu-item index="/records"><el-icon><Document /></el-icon>使用记录</el-menu-item>
        <el-menu-item index="/stats"><el-icon><Histogram /></el-icon>数据统计</el-menu-item>
        <el-menu-item index="/tickets"><el-icon><Warning /></el-icon>报修工单</el-menu-item>
        <el-menu-item index="/qrcodes"><el-icon><Grid /></el-icon>二维码管理</el-menu-item>
        <el-menu-item index="/backup"><el-icon><Folder /></el-icon>数据备份</el-menu-item>
        <el-menu-item index="/notes"><el-icon><EditPen /></el-icon>备注管理</el-menu-item>
        <el-menu-item index="/users" v-if="auth.isAdmin"><el-icon><User /></el-icon>用户管理</el-menu-item>
        <el-menu-item index="/settings"><el-icon><Setting /></el-icon>系统设置</el-menu-item>
      </el-menu>
      <div style="padding: 16px 20px; border-top: 1px solid rgba(255,255,255,0.1); font-size: 13px; display: flex; justify-content: space-between; align-items: center;">
        <span>{{ auth.role === 'admin' ? '管理员' : '只读用户' }}</span>
        <el-button size="small" type="danger" plain @click="logout">退出</el-button>
      </div>
    </el-aside>
    <el-main style="background: #f5f7fa; overflow-y: auto;">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>
