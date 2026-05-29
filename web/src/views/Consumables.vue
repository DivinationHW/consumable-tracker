<template>
  <div class="page-container">
    <div class="page-header">
      <h2>耗材管理</h2>
      <el-button type="primary" @click="showDialog = true" v-if="auth.isAdmin"><el-icon><Plus /></el-icon>新增耗材</el-button>
    </div>
    <div class="filters">
      <el-input v-model="search" placeholder="搜索耗材名称" clearable @input="loadData" />
    </div>
    <el-table :data="items" stripe v-loading="loading" @row-dblclick="editItem">
      <el-table-column prop="name" label="耗材名称" min-width="160" />
      <el-table-column prop="category" label="分类" width="120" />
      <el-table-column prop="unit" label="单位" width="80" />
      <el-table-column prop="stock" label="库存" width="100" sortable />
      <el-table-column prop="threshold" label="预警值" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.stock <= row.threshold ? 'danger' : 'success'" size="small">
            {{ row.stock <= row.threshold ? '需补货' : '充足' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" v-if="auth.isAdmin">
        <template #default="{ row }">
          <el-button size="small" @click="editItem(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showDialog" :title="editing ? '编辑耗材' : '新增耗材'" width="500px">
      <el-form :model="form" label-position="top">
        <el-form-item label="耗材名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="分类"><el-input v-model="form.category" /></el-form-item>
        <el-form-item label="单位"><el-input v-model="form.unit" /></el-form-item>
        <el-form-item label="库存"><el-input-number v-model="form.stock" :min="0" /></el-form-item>
        <el-form-item label="预警值"><el-input-number v-model="form.threshold" :min="0" /></el-form-item>
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
import { getConsumables, createConsumable, updateConsumable, deleteConsumable } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuthStore()
const items = ref([])
const loading = ref(false)
const search = ref('')
const showDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = ref({ name: '', category: '', unit: '', stock: 0, threshold: 0 })

onMounted(loadData)

async function loadData() {
  loading.value = true
  try { items.value = await getConsumables({ search: search.value }) } finally { loading.value = false }
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
      await updateConsumable(form.value.id, form.value)
    } else {
      await createConsumable(form.value)
    }
    showDialog.value = false
    ElMessage.success('保存成功')
    loadData()
  } finally { saving.value = false }
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该耗材？')
    await deleteConsumable(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}
</script>
