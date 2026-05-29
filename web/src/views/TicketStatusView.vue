<template>
  <div style="min-height:100vh;background:#f5f7fa;display:flex;align-items:center;justify-content:center;padding:20px;">
    <el-card style="width:460px;max-width:100%;">
      <h2 style="text-align:center;margin-bottom:24px;">工单进度查询</h2>
      <div v-if="ticket.id">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="工单编号">{{ ticket.id }}</el-descriptions-item>
          <el-descriptions-item label="房间号">{{ ticket.room_number }}</el-descriptions-item>
          <el-descriptions-item label="问题类型">{{ ticket.problem_type }}</el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ ticket.created_at }}</el-descriptions-item>
          <el-descriptions-item label="当前状态">
            <el-tag :class="'tag-' + ticket.status" size="large">
              {{ {pending:'待处理',processing:'处理中',completed:'已完成'}[ticket.status] }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <el-input v-else v-model="searchId" placeholder="输入工单编号" style="margin-bottom:16px;" />
      <el-button v-if="!ticket.id" type="primary" @click="searchTicket" style="width:100%">查询</el-button>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const ticket = ref({})
const searchId = ref(route.params.id || '')

onMounted(() => {
  if (searchId.value) searchTicket()
})

async function searchTicket() {
  if (!searchId.value) return
  try {
    const res = await axios.get(`/ticket/${searchId.value}`)
    ticket.value = res.data
  } catch {}
}
</script>
