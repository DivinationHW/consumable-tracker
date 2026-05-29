<template>
  <div class="page-container">
    <div class="page-header">
      <h2>科室管理</h2>
      <el-button type="primary" @click="showDialog = true" v-if="auth.isAdmin"><el-icon><Plus /></el-icon>新增科室</el-button>
    </div>
    <el-table :data="items" stripe v-loading="loading" @row-dblclick="editItem">
      <el-table-column prop="room_number" label="房间号" width="120" />
      <el-table-column prop="department" label="所属科室" min-width="160" />
      <el-table-column prop="device_type" label="设备类型" width="120" />
      <el-table-column prop="device_model" label="设备型号" width="160" />
      <el-table-column label="操作" width="160" v-if="auth.isAdmin">
        <template #default="{ row }">
          <el-button size="small" @click="editItem(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showDialog" :title="editing ? '编辑科室' : '新增科室'" width="500px">
      <el-form :model="form" label-position="top">
        <el-form-item label="房间号"><el-input v-model="form.room_number" /></el-form-item>
        <el-form-item label="所属科室"><el-input v-model="form.department" /></el-form-item>
        <el-form-item label="设备类型"><el-input v-model="form.device_type" /></el-form-item>
        <el-form-item label="设备型号"><el-input v-model="form.device_model" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOffices, createOffice, updateOffice, deleteOffice } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuthStore()
const items = ref([])
const loading = ref(false)
const showDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = ref({ room_number: '', department: '', device_type: '', device_model: '' })

onMounted(loadData)

async function loadData() {
  loading.value = true
  try { items.value = await getOffices() } finally { loading.value = false }
}

function editItem(row) {
  if (!auth.isAdmin) return
  form.value = { ...row }
  editing.value = true
  showDialog.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editing.value) {
      await updateOffice(form.value.id, form.value)
    } else {
      await createOffice(form.value)
    }
    showDialog.value = false
    ElMessage.success('保存成功')
    loadData()
  } finally { saving.value = false }
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该科室？')
    await deleteOffice(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}
</script>
