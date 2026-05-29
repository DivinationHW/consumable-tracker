<template>
  <div style="min-height:100vh;background:#f5f7fa;display:flex;align-items:center;justify-content:center;padding:20px;">
    <el-card style="width:460px;max-width:100%;">
      <h2 style="text-align:center;margin-bottom:24px;">设备报修</h2>
      <el-form v-if="!submitted" @submit.prevent="handleSubmit" label-position="top">
        <el-form-item label="房间号">
          <el-input :model-value="roomNumber" disabled />
        </el-form-item>
        <el-form-item label="设备类型">
          <el-input :model-value="deviceType" disabled />
        </el-form-item>
        <el-form-item label="问题类型">
          <el-select v-model="form.problem_type" style="width:100%">
            <el-option v-for="pt in problemTypes" :key="pt.id" :label="pt.name" :value="pt.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="问题描述（可选）">
          <el-input v-model="form.description" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item label="联系方式（可选）">
          <el-input v-model="form.contact" placeholder="手机号或微信号" maxlength="100" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="saving" style="width:100%">提交报修</el-button>
        </el-form-item>
      </el-form>
      <div v-else style="text-align:center;padding:20px 0;">
        <el-icon style="font-size:64px;color:#67c23a;"><CircleCheck /></el-icon>
        <h3 style="margin:16px 0 8px;">报修提交成功</h3>
        <p style="color:#909399;margin-bottom:16px;">工单编号：{{ ticketId }}</p>
        <el-button @click="checkStatus">查看进度</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const roomNumber = ref('')
const deviceType = ref('')
const problemTypes = ref([])
const officeId = ref(null)
const saving = ref(false)
const submitted = ref(false)
const ticketId = ref('')
const form = ref({ problem_type: '', description: '', contact: '' })

onMounted(async () => {
  const code = route.query.office
  if (!code) return
  try {
    const res = await axios.get(`/api/qrcodes/${code}/image`)
    roomNumber.value = res.data.room_number || '未知房间'
    deviceType.value = res.data.device_type || '未知设备'
  } catch {}
  try {
    const res = await axios.get(`/ticket/office/${officeId.value}/problem-types`)
    problemTypes.value = res.data.problem_types || []
    const defaultPt = problemTypes.value.find(p => p.is_default)
    if (defaultPt) form.value.problem_type = defaultPt.name
  } catch {}
})

function checkStatus() {
  router.push(`/ticket/${ticketId.value}`)
}

async function handleSubmit() {
  if (!form.value.problem_type) {
    ElMessage.warning('请选择问题类型')
    return
  }
  saving.value = true
  try {
    const res = await axios.post('/ticket', {
      office_id: officeId.value,
      problem_type: form.value.problem_type,
      description: form.value.description,
      contact: form.value.contact,
    })
    ticketId.value = res.data.id
    submitted.value = true
  } catch { saving.value = false }
}
</script>
