<template>
  <div class="page-container">
    <div class="page-header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="showDialog = true"><el-icon><Plus /></el-icon>新增用户</el-button>
    </div>
    <el-table :data="items" stripe v-loading="loading">
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
            {{ row.role === 'admin' ? '管理员' : '只读用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="editItem(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)" :disabled="row.id === currentUserId">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showDialog" :title="editing ? '编辑用户' : '新增用户'" width="500px">
      <el-form :model="form" label-position="top">
        <el-form-item label="用户名"><el-input v-model="form.username" /></el-form-item>
        <el-form-item label="密码" v-if="!editing"><el-input v-model="form.password" type="password" show-password /></el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role" style="width:100%">
            <el-option label="管理员" value="admin" />
            <el-option label="只读用户" value="readonly" />
          </el-select>
        </el-form-item>
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
import { getUsers, createUser, updateUser, deleteUser } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const items = ref([])
const loading = ref(false)
const showDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const currentUserId = ref(parseInt(localStorage.getItem('user_id') || '0'))
const form = ref({ username: '', password: '', role: 'readonly' })

onMounted(loadData)

async function loadData() {
  loading.value = true
  try { items.value = await getUsers() } finally { loading.value = false }
}

function editItem(row) {
  form.value = { id: row.id, username: row.username, role: row.role, password: '' }
  editing.value = true
  showDialog.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editing.value) {
      const payload = { username: form.value.username, role: form.value.role }
      if (form.value.password) payload.password = form.value.password
      await updateUser(form.value.id, payload)
    } else {
      await createUser(form.value)
    }
    showDialog.value = false
    ElMessage.success('保存成功')
    loadData()
  } finally { saving.value = false }
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该用户？')
    await deleteUser(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}
</script>
