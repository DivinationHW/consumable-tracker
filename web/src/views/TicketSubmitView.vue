<template>
  <div style="max-width:500px;margin:40px auto;padding:0 16px">
    <el-card>
      <template #header><h2 style="text-align:center">提交维修工单</h2></template>
      <el-form :model="form" label-width="80px" @keyup.enter="submit">
        <el-form-item label="房间号">
          <el-select v-model="form.office_id" filterable placeholder="选择办公室" style="width:100%">
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备类型">{{ form.device_type || '自动识别' }}</el-form-item>
        <el-form-item label="设备型号">{{ form.device_model || '自动识别' }}</el-form-item>
        <el-form-item label="故障类型">
          <el-radio-group v-model="form.problem_type">
            <el-radio-button v-for="p in problemTypes" :key="p.id" :value="p.name">{{ p.name }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="联系方式"><el-input v-model="form.contact" placeholder="电话/微信" /></el-form-item>
        <el-form-item><el-button type="primary" style="width:100%" :loading="loading" @click="submit">提交</el-button></el-form-item>
      </el-form>
      <div v-if="result" style="text-align:center;color:#67c23a">
        <p>提交成功！工单编号：</p>
        <p style="font-size:18px;font-weight:bold">{{ result }}</p>
        <p>请保存编号查询进度</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { listOffices } from '@/api/offices'
import { listProblemTypes } from '@/api/problemTypes'
import { createPublicTicket } from '@/api/tickets'
import { listQRCodes } from '@/api/qrcodes'

const route = useRoute()
const loading = ref(false)
const offices = ref<any[]>([])
const problemTypes = ref<any[]>([])
const result = ref('')
const form = reactive({
  office_id: undefined as number | undefined,
  device_type: '',
  device_model: '',
  problem_type: '',
  description: '',
  contact: '',
})

async function submit() {
  if (!form.office_id || !form.problem_type) {
    ElMessage.warning('请选择办公室和故障类型')
    return
  }
  loading.value = true
  try {
    const res = await createPublicTicket(form, route.query.qr as string)
    result.value = res.data.id
    ElMessage.success('提交成功')
  } finally { loading.value = false }
}

onMounted(async () => {
  const qr = route.query.qr as string
  const [o, p] = await Promise.all([listOffices(), listProblemTypes()])
  offices.value = o.data
  problemTypes.value = p.data

  if (qr) {
    const qrList = (await listQRCodes()).data
    const match = qrList.find((q: any) => q.code === qr)
    if (match && match.office_id) {
      form.office_id = match.office_id
      form.device_type = match.device_type
      form.device_model = match.device_model
    }
  }
})
</script>
