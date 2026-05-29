<template>
  <div>
    <el-card style="margin-bottom:16px">
      <el-form :inline="true" :model="query">
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部" style="width:120px">
            <el-option label="待处理" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="load">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <el-table :data="tickets" stripe v-loading="loading" @row-dblclick="showDetail">
      <el-table-column prop="created_at" label="提交时间" width="160" />
      <el-table-column prop="office_name" label="办公室" width="120" />
      <el-table-column prop="device_type" label="设备类型" width="100" />
      <el-table-column prop="problem_type" label="故障类型" width="120" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'completed' ? 'success' : row.status === 'processing' ? 'warning' : 'info'">
            {{ row.status === 'pending' ? '待处理' : row.status === 'processing' ? '处理中' : '已完成' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" link @click="showDetail(row)">详情</el-button>
          <el-button v-if="row.status !== 'completed' && auth.isAdmin" size="small" type="warning" link @click="process(row)">处理</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="detailDialog" title="工单详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="办公室">{{ detail.office_name }}</el-descriptions-item>
        <el-descriptions-item label="设备">{{ detail.device_type }} {{ detail.device_model }}</el-descriptions-item>
        <el-descriptions-item label="故障">{{ detail.problem_type }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detail.status === 'completed' ? 'success' : detail.status === 'processing' ? 'warning' : 'info'">{{ detail.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ detail.description }}</el-descriptions-item>
        <el-descriptions-item label="联系方式">{{ detail.contact }}</el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ detail.created_at }}</el-descriptions-item>
        <el-descriptions-item label="处理人" v-if="detail.handled_by_user">{{ detail.handled_by_user }}</el-descriptions-item>
        <el-descriptions-item label="处理备注" v-if="detail.handle_note" :span="2">{{ detail.handle_note }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog v-model="processDialog" title="处理工单" width="500px">
      <el-form :model="processForm" label-width="80px">
        <el-form-item label="处理结果">
          <el-select v-model="processForm.status">
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-form-item>
        <el-form-item label="耗材用量"><el-input v-model="processForm.consumable_used" placeholder="耗材名称" /></el-form-item>
        <el-form-item label="数量"><el-input-number v-model="processForm.consumable_quantity" :min="0" /></el-form-item>
        <el-form-item label="处理备注"><el-input v-model="processForm.handle_note" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="processDialog=false">取消</el-button>
        <el-button type="primary" @click="saveProcess">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listTickets, getTicket, processTicket } from '@/api/tickets'
import { useAuthStore } from '@/stores/auth'
import { useWebSocket } from '@/api/websocket'

const auth = useAuthStore()
const loading = ref(false)
const tickets = ref<any[]>([])
const detail = ref<any>({})
const detailDialog = ref(false)
const processDialog = ref(false)
const processForm = reactive({ status: 'completed', consumable_used: '', consumable_quantity: 0, handle_note: '' })
let processId = ''
const query = reactive({ status: '' as string })

async function load() {
  loading.value = true
  try { tickets.value = (await listTickets(query)).data } finally { loading.value = false }
}
async function showDetail(row: any) {
  const res = await getTicket(row.id)
  detail.value = res.data
  detailDialog.value = true
}
function process(row: any) {
  processId = row.id
  processForm.status = 'completed'
  processForm.consumable_used = ''
  processForm.consumable_quantity = 0
  processForm.handle_note = ''
  processDialog.value = true
}
async function saveProcess() {
  await processTicket(processId, processForm)
  ElMessage.success('处理成功')
  processDialog.value = false
  load()
}

const ws = useWebSocket()
ws.on('ticket_new', () => load())
ws.on('ticket_update', () => load())
ws.connect()

onMounted(load)
</script>
