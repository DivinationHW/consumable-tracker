<template>
  <div class="page-container">
    <div class="page-header">
      <h2>二维码管理</h2>
      <div>
        <el-button @click="showGenerate = true" v-if="auth.isAdmin">生成二维码</el-button>
      </div>
    </div>
    <el-table :data="items" stripe v-loading="loading" @row-dblclick="editItem">
      <el-table-column prop="code" label="二维码编码" width="160" />
      <el-table-column prop="room_number" label="绑定房间" width="100">
        <template #default="{ row }">{{ row.room_number || '-' }}</template>
      </el-table-column>
      <el-table-column prop="device_type" label="设备类型" width="120" />
      <el-table-column prop="device_model" label="设备型号" width="150" />
      <el-table-column label="配置状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_configured ? 'success' : 'info'" size="small">
            {{ row.is_configured ? '已配置' : '未配置' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column label="操作" width="200" v-if="auth.isAdmin">
        <template #default="{ row }">
          <el-button size="small" @click="editItem(row)">配置</el-button>
          <el-button size="small" @click="showImage(row.code)">二维码</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showGenerate" title="生成二维码" width="400px">
      <el-form label-position="top">
        <el-form-item label="生成数量">
          <el-input-number v-model="generateCount" :min="1" :max="50" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showGenerate = false">取消</el-button>
        <el-button type="primary" @click="handleGenerate" :loading="generating">生成</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEdit" :title="'配置二维码 - ' + editForm.code" width="500px">
      <el-form :model="editForm" label-position="top">
        <el-form-item label="绑定科室">
          <el-select v-model="editForm.office_id" filterable clearable style="width:100%">
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number + ' - ' + o.department" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备类型"><el-input v-model="editForm.device_type" /></el-form-item>
        <el-form-item label="设备型号">
          <el-autocomplete v-model="editForm.device_model" :fetch-suggestions="queryModels" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEdit = false">取消</el-button>
        <el-button type="primary" @click="handleSaveEdit" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showQR" title="二维码图片" width="360px" align-center>
      <div style="text-align:center;">
        <img :src="qrImage" style="width:256px;height:256px;" v-if="qrImage" />
        <div style="margin-top:12px;word-break:break-all;font-size:12px;color:#909399;">{{ qrDataUrl }}</div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getQRCodes, createQRCode, updateQRCode, deleteQRCode, getQRCodeImage, getDeviceModels, getOffices } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuthStore()
const items = ref([])
const offices = ref([])
const models = ref([])
const loading = ref(false)
const showGenerate = ref(false)
const generating = ref(false)
const generateCount = ref(10)
const showEdit = ref(false)
const saving = ref(false)
const editForm = ref({ id: null, code: '', office_id: null, device_type: '', device_model: '' })
const showQR = ref(false)
const qrImage = ref('')
const qrDataUrl = ref('')

onMounted(async () => {
  offices.value = await getOffices()
  models.value = (await getDeviceModels()).models || []
  loadData()
})

async function loadData() {
  loading.value = true
  try { items.value = await getQRCodes() } finally { loading.value = false }
}

async function handleGenerate() {
  generating.value = true
  try {
    await createQRCode({ count: generateCount.value })
    showGenerate.value = false
    ElMessage.success('生成成功')
    loadData()
  } finally { generating.value = false }
}

function editItem(row) {
  if (!auth.isAdmin) return
  editForm.value = { id: row.id, code: row.code, office_id: row.office_id, device_type: row.device_type || '', device_model: row.device_model || '' }
  showEdit.value = true
}

async function handleSaveEdit() {
  saving.value = true
  try {
    await updateQRCode(editForm.value.id, { office_id: editForm.value.office_id, device_type: editForm.value.device_type, device_model: editForm.value.device_model })
    showEdit.value = false
    ElMessage.success('保存成功')
    loadData()
  } finally { saving.value = false }
}

async function showImage(code) {
  const data = await getQRCodeImage(code)
  qrImage.value = data.image
  qrDataUrl.value = data.data_url
  showQR.value = true
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该二维码？')
    await deleteQRCode(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}

function queryModels(query, cb) {
  const filtered = models.value.filter(m => m.toLowerCase().includes(query.toLowerCase()))
  cb(filtered.map(m => ({ value: m })))
}
</script>
