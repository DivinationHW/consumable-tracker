<template>
  <div style="max-width:500px;margin:40px auto;padding:0 16px">
    <el-card>
      <template #header><h2 style="text-align:center">查询工单状态</h2></template>
      <el-form @keyup.enter="search">
        <el-form-item>
          <el-input v-model="id" placeholder="请输入工单编号" size="large">
            <template #append><el-button @click="search" :loading="loading">查询</el-button></template>
          </el-input>
        </el-form-item>
      </el-form>
      <div v-if="ticket">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="状态">
            <el-tag :type="ticket.status === 'completed' ? 'success' : ticket.status === 'processing' ? 'warning' : 'info'">
              {{ ticket.status === 'pending' ? '待处理' : ticket.status === 'processing' ? '处理中' : '已完成' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="办公室">{{ ticket.office_name }}</el-descriptions-item>
          <el-descriptions-item label="故障类型">{{ ticket.problem_type }}</el-descriptions-item>
          <el-descriptions-item label="描述">{{ ticket.description }}</el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ ticket.created_at }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <div v-if="notFound" style="text-align:center;color:#999">未找到该工单</div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { getPublicTicket } from '@/api/tickets'

const route = useRoute()
const id = ref((route.query.id as string) || '')
const ticket = ref<any>(null)
const loading = ref(false)
const notFound = ref(false)

async function search() {
  if (!id.value) return
  loading.value = true
  ticket.value = null
  notFound.value = false
  try {
    ticket.value = (await getPublicTicket(id.value)).data
  } catch {
    notFound.value = true
  } finally { loading.value = false }
}

if (route.query.id) search()
</script>
