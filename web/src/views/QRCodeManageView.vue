<template>
  <div>
    <el-card style="margin-bottom:16px">
      <template #header><span>生成二维码</span></template>
      <el-form :inline="true" :model="genForm" label-width="auto">
        <el-form-item label="办公室"><el-select v-model="genForm.office_id" clearable placeholder="可选" style="width:160px">
          <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
        </el-select></el-form-item>
        <el-form-item label="设备类型"><el-input v-model="genForm.device_type" placeholder="如: printer" style="width:140px" /></el-form-item>
        <el-form-item label="型号"><el-input v-model="genForm.device_model" placeholder="可选" style="width:180px" /></el-form-item>
        <el-form-item>
          <el-button type="primary" @click="generateSingle">生成单张</el-button>
          <el-button @click="generateBulk">批量生成</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>二维码列表</span>
          <el-button size="small" @click="print">打印全部</el-button>
        </div>
      </template>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column label="二维码" width="90">
          <template #default="{ row }">
            <el-image :src="`/api/qrcodes/${row.code}/image`" style="width:60px;height:60px" />
          </template>
        </el-table-column>
        <el-table-column prop="code" label="码值" width="100" />
        <el-table-column prop="office_name" label="办公室" width="120" />
        <el-table-column prop="device_type" label="设备类型" width="100" />
        <el-table-column prop="device_model" label="型号" width="160" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button size="small" link @click="download(row.code)">下载</el-button>
            <el-button size="small" type="danger" link @click="doDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listOffices } from '@/api/offices'
import { listQRCodes, createQRCode, generateBulk, deleteQRCode } from '@/api/qrcodes'

const loading = ref(false)
const offices = ref<any[]>([])
const list = ref<any[]>([])
const genForm = reactive({ office_id: undefined as number | undefined, device_type: '', device_model: '' })

async function load() {
  const [o, q] = await Promise.all([listOffices(), listQRCodes()])
  offices.value = o.data
  list.value = q.data
}
async function generateSingle() {
  await createQRCode(genForm)
  ElMessage.success('已生成')
  load()
}
async function generateBulk() {
  await ElMessageBox.prompt('生成数量', '批量生成', { inputValue: '10', inputPattern: /^\d+$/ })
  const { value } = await ElMessageBox.prompt('生成数量', '批量生成', { inputValue: '10' })
  await generateBulk({ count: parseInt(value), ...genForm })
  ElMessage.success(`已生成 ${value} 个`)
  load()
}
function download(code: string) {
  window.open(`/api/qrcodes/${code}/image`, '_blank')
}
async function doDelete(id: number) {
  await ElMessageBox.confirm('确认删除？')
  await deleteQRCode(id)
  ElMessage.success('已删除')
  load()
}
function print() {
  window.open('/api/qrcodes/print', '_blank')
}
onMounted(load)
</script>
