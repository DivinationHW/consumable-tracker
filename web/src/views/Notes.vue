<template>
  <div class="page-container">
    <div class="page-header">
      <h2>备注管理</h2>
      <el-button type="primary" @click="showDialog = true" v-if="auth.isAdmin"><el-icon><Plus /></el-icon>新增备注</el-button>
    </div>
    <el-table :data="items" stripe v-loading="loading">
      <el-table-column prop="room_number" label="房间号" width="100" />
      <el-table-column prop="content" label="备注内容" min-width="300" show-overflow-tooltip />
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column label="操作" width="160" v-if="auth.isAdmin">
        <template #default="{ row }">
          <el-button size="small" @click="editItem(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showDialog" :title="editing ? '编辑备注' : '新增备注'" width="500px">
      <el-form :model="form" label-position="top">
        <el-form-item label="科室">
          <el-select v-model="form.office_id" filterable style="width:100%">
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注内容"><el-input v-model="form.content" type="textarea" :rows="4" /></el-form-item>
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
import { getNotes, createNote, updateNote, deleteNote, getOffices } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuthStore()
const items = ref([])
const offices = ref([])
const loading = ref(false)
const showDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = ref({ office_id: null, content: '' })

onMounted(async () => {
  offices.value = await getOffices()
  loadData()
})

async function loadData() {
  loading.value = true
  try { items.value = await getNotes() } finally { loading.value = false }
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
      await updateNote(form.value.id, form.value)
    } else {
      await createNote(form.value)
    }
    showDialog.value = false
    ElMessage.success('保存成功')
    loadData()
  } finally { saving.value = false }
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该备注？')
    await deleteNote(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}
</script>
