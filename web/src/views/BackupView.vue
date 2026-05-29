<template>
  <div>
    <el-button type="primary" style="margin-bottom:16px" :loading="creating" @click="doCreate">创建备份</el-button>
    <el-table :data="backups" stripe v-loading="loading">
      <el-table-column prop="filename" label="文件名" />
      <el-table-column prop="size" label="大小" width="120">
        <template #default="{ row }">{{ (row.size / 1024).toFixed(1) }} KB</template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="warning" link @click="doRestore(row.filename)">恢复</el-button>
          <el-button size="small" type="danger" link @click="doDelete(row.filename)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listBackups, createBackup, restoreBackup, deleteBackup } from '@/api/backup'

const loading = ref(false)
const creating = ref(false)
const backups = ref<any[]>([])

async function load() {
  loading.value = true
  try { backups.value = (await listBackups()).data } finally { loading.value = false }
}

async function doCreate() {
  creating.value = true
  try {
    await createBackup()
    ElMessage.success('备份成功')
    load()
  } finally { creating.value = false }
}

async function doRestore(filename: string) {
  try {
    await ElMessageBox.confirm(`确认从 ${filename} 恢复？需要重启服务器`, '警告', { confirmButtonText: '恢复', type: 'warning' })
    await restoreBackup(filename)
    ElMessage.warning(`已从 ${filename} 恢复，请重启服务器`)
  } catch { }
}

async function doDelete(filename: string) {
  await ElMessageBox.confirm('确认删除？')
  await deleteBackup(filename)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>
