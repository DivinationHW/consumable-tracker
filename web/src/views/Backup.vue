<template>
  <div class="page-container">
    <div class="page-header">
      <h2>数据备份</h2>
      <div>
        <el-button type="primary" @click="handleCreateBackup" :loading="creating">立即备份</el-button>
      </div>
    </div>

    <el-card style="margin-bottom:20px;">
      <el-form :model="config" label-position="left" inline>
        <el-form-item label="自动备份">
          <el-switch v-model="autoBackup" @change="saveConfig" />
        </el-form-item>
        <el-form-item label="保留天数">
          <el-input-number v-model="config.keep_days" :min="1" :max="365" @change="saveConfig" />
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="items" stripe v-loading="loading">
      <el-table-column prop="name" label="文件名" min-width="300" />
      <el-table-column prop="size" label="大小" width="100">
        <template #default="{ row }">{{ (row.size / 1024).toFixed(1) }} KB</template>
      </el-table-column>
      <el-table-column prop="date" label="备份时间" width="160" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="handleDownload(row.name)" v-if="auth.isAdmin">下载</el-button>
          <el-button size="small" type="warning" @click="handleRestore(row.name)" v-if="auth.isAdmin">恢复</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.name)" v-if="auth.isAdmin">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getBackupConfig, saveBackupConfig, getBackupList, createBackup, downloadBackup, restoreBackup, deleteBackup } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuthStore()
const items = ref([])
const loading = ref(false)
const creating = ref(false)
const config = ref({ keep_days: 180 })
const autoBackup = ref(true)

onMounted(loadData)

async function loadData() {
  try { config.value = await getBackupConfig() } catch {}
  loadList()
}

async function loadList() {
  loading.value = true
  try { items.value = await getBackupList() } finally { loading.value = false }
}

async function saveConfig() {
  await saveBackupConfig({ keep_days: config.value.keep_days })
  ElMessage.success('配置已保存')
}

async function handleCreateBackup() {
  creating.value = true
  try {
    await createBackup()
    ElMessage.success('备份已创建')
    loadList()
  } finally { creating.value = false }
}

async function handleDownload(name) {
  try {
    const blob = await downloadBackup(name)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = name
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {}
}

async function handleRestore(name) {
  try {
    await ElMessageBox.confirm(`确定从备份文件 ${name} 恢复？此操作将覆盖当前数据。`)
    await restoreBackup(name)
    ElMessage.success('恢复已启动')
  } catch {}
}

async function handleDelete(name) {
  try {
    await ElMessageBox.confirm('确定删除该备份文件？')
    await deleteBackup(name)
    ElMessage.success('删除成功')
    loadList()
  } catch {}
}
</script>
