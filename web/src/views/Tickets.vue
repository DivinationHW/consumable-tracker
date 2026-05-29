<template>
  <div class="page-container">
    <div class="page-header">
      <h2>报修工单</h2>
      <el-select v-model="filterStatus" placeholder="工单状态" clearable @change="loadData" style="width:140px">
        <el-option label="待处理" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
      </el-select>
    </div>
    <el-table :data="items" stripe v-loading="loading">
      <el-table-column prop="room_number" label="房间" width="80" />
      <el-table-column prop="device_type" label="设备类型" width="100" />
      <el-table-column prop="problem_type" label="问题类型" width="120" />
      <el-table-column prop="description" label="描述" min-width="160" show-overflow-tooltip />
      <el-table-column prop="contact" label="联系方式" width="140" />
      <el-table-column prop="created_at" label="提交时间" width="160" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :class="'tag-' + row.status" size="small">
            {{ {pending:'待处理',processing:'处理中',completed:'已完成'}[row.status] }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" v-if="auth.isAdmin">
        <template #default="{ row }">
          <template v-if="row.status === 'pending'">
            <el-button size="small" type="warning" @click="updateStatus(row.id, 'processing')">处理</el-button>
          </template>
          <template v-if="row.status === 'processing'">
            <el-button size="small" type="success" @click="openComplete(row)">完成</el-button>
          </template>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showComplete" title="完成工单" width="500px">
      <el-form :model="completeForm" label-position="top">
        <el-form-item label="消耗耗材"><el-input v-model="completeForm.consumable_used" placeholder="耗材名称" /></el-form-item>
        <el-form-item label="消耗数量"><el-input-number v-model="completeForm.consumable_quantity" :min="0" /></el-form-item>
        <el-form-item label="处理备注"><el-input v-model="completeForm.handle_note" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showComplete = false">取消</el-button>
        <el-button type="primary" @click="handleComplete" :loading="saving">确定完成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getTickets, updateTicketStatus, completeTicket, deleteTicket } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuthStore()
const items = ref([])
const loading = ref(false)
const filterStatus = ref('')
const showComplete = ref(false)
const saving = ref(false)
const completeForm = ref({ consumable_used: '', consumable_quantity: 0, handle_note: '' })
const currentTicketId = ref(null)

onMounted(loadData)

async function loadData() {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value) params.status = filterStatus.value
    items.value = await getTickets(params)
  } finally { loading.value = false }
}

async function updateStatus(id, status) {
  await updateTicketStatus(id, status)
  ElMessage.success('状态已更新')
  loadData()
}

function openComplete(row) {
  currentTicketId.value = row.id
  completeForm.value = { consumable_used: '', consumable_quantity: 0, handle_note: '' }
  showComplete.value = true
}

async function handleComplete() {
  saving.value = true
  try {
    await completeTicket(currentTicketId.value, completeForm.value)
    showComplete.value = false
    ElMessage.success('工单已完成')
    loadData()
  } finally { saving.value = false }
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该工单？')
    await deleteTicket(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}
</script>
