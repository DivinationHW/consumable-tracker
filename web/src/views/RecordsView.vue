<template>
  <div>
    <el-card style="margin-bottom:16px">
      <el-form :inline="true" :model="query" label-width="auto">
        <el-form-item label="办公室">
          <el-select v-model="query.office_id" clearable placeholder="全部" style="width:160px">
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="耗材">
          <el-select v-model="query.consumable_id" clearable placeholder="全部" style="width:200px">
            <el-option v-for="c in consumables" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始">
          <el-date-picker v-model="query.date_from" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="结束">
          <el-date-picker v-model="query.date_to" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="reset">重置</el-button>
          <el-button @click="doExport" :disabled="loading">导出Excel</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table :data="records" stripe v-loading="loading" @row-dblclick="edit">
        <el-table-column prop="usage_date" label="日期" width="110" sortable />
        <el-table-column prop="office_name" label="办公室" width="120" />
        <el-table-column prop="device_type" label="设备类型" width="100" />
        <el-table-column prop="device_model" label="型号" width="160" show-overflow-tooltip />
        <el-table-column prop="consumable_name" label="耗材" width="220" show-overflow-tooltip />
        <el-table-column prop="quantity" label="数量" width="70" />
        <el-table-column prop="note" label="备注" show-overflow-tooltip />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button v-if="auth.isAdmin" size="small" type="primary" link @click="edit(row)">编辑</el-button>
            <el-button v-if="auth.isAdmin" size="small" type="danger" link @click="doDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="query.page"
        :page-size="50"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top:16px;justify-content:center"
        @current-change="load"
      />
    </el-card>

    <el-dialog v-model="editDialog" title="编辑记录" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="日期"><el-date-picker v-model="editForm.usage_date" type="date" value-format="YYYY-MM-DD" /></el-form-item>
        <el-form-item label="办公室">
          <el-select v-model="editForm.office_id" filterable>
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="耗材">
          <el-select v-model="editForm.consumable_id" filterable>
            <el-option v-for="c in consumables" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量"><el-input-number v-model="editForm.quantity" :min="1" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="editForm.note" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialog=false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listOffices } from '@/api/offices'
import { listConsumables } from '@/api/consumables'
import { listRecords, updateRecord, deleteRecord, exportRecords } from '@/api/records'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const loading = ref(false)
const offices = ref<any[]>([])
const consumables = ref<any[]>([])
const records = ref<any[]>([])
const total = ref(0)
const editDialog = ref(false)
const editForm = reactive<any>({})
const query = reactive({
  office_id: undefined as number | undefined,
  consumable_id: undefined as number | undefined,
  date_from: undefined as string | undefined,
  date_to: undefined as string | undefined,
  page: 1,
})

async function load() {
  loading.value = true
  try {
    const res = await listRecords(query)
    records.value = res.data.data
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}

function search() { query.page = 1; load() }
function reset() {
  query.office_id = undefined
  query.consumable_id = undefined
  query.date_from = undefined
  query.date_to = undefined
  query.page = 1
  load()
}

function edit(row: any) {
  Object.assign(editForm, row)
  editDialog.value = true
}

async function saveEdit() {
  await updateRecord(editForm.id, editForm)
  ElMessage.success('更新成功')
  editDialog.value = false
  load()
}

async function doDelete(id: number) {
  await ElMessageBox.confirm('确认删除？')
  await deleteRecord(id)
  ElMessage.success('已删除')
  load()
}

async function doExport() {
  loading.value = true
  try {
    const res = await exportRecords(query)
    const url = URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a')
    a.href = url
    a.download = `耗材记录_${new Date().toISOString().slice(0,10)}.xlsx`
    a.click()
    URL.revokeObjectURL(url)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const [o, c] = await Promise.all([listOffices(), listConsumables()])
  offices.value = o.data
  consumables.value = c.data
  load()
})
</script>
