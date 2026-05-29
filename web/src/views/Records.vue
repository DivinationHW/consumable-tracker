<template>
  <div class="page-container">
    <div class="page-header">
      <h2>使用记录</h2>
      <div>
        <el-button @click="exportExcel">导出Excel</el-button>
        <el-button type="primary" @click="showDialog = true" v-if="auth.isAdmin"><el-icon><Plus /></el-icon>新增记录</el-button>
      </div>
    </div>
    <div class="filters">
      <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" @change="loadData" />
      <el-select v-model="filterOffice" placeholder="选择科室" clearable @change="loadData">
        <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
      </el-select>
      <el-select v-model="filterConsumable" placeholder="选择耗材" clearable @change="loadData">
        <el-option v-for="c in consumables" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
    </div>
    <el-table :data="items" stripe v-loading="loading">
      <el-table-column prop="usage_date" label="日期" width="120" />
      <el-table-column prop="room_number" label="房间号" width="100" />
      <el-table-column prop="consumable_name" label="耗材" width="140" />
      <el-table-column prop="quantity" label="数量" width="80" sortable />
      <el-table-column prop="unit" label="单位" width="60" />
      <el-table-column prop="note" label="备注" min-width="160" show-overflow-tooltip />
      <el-table-column label="操作" width="160" v-if="auth.isAdmin">
        <template #default="{ row }">
          <el-button size="small" @click="editItem(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showDialog" :title="editing ? '编辑记录' : '新增记录'" width="500px">
      <el-form :model="form" label-position="top">
        <el-form-item label="科室">
          <el-select v-model="form.office_id" filterable style="width:100%">
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="耗材">
          <el-select v-model="form.consumable_id" filterable style="width:100%">
            <el-option v-for="c in consumables" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量"><el-input-number v-model="form.quantity" :min="1" /></el-form-item>
        <el-form-item label="日期"><el-date-picker v-model="form.usage_date" type="date" style="width:100%" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.note" type="textarea" :rows="2" /></el-form-item>
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
import { getRecords, createRecord, updateRecord, deleteRecord, getOffices, getConsumables } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as XLSX from 'xlsx'

const auth = useAuthStore()
const items = ref([])
const offices = ref([])
const consumables = ref([])
const loading = ref(false)
const dateRange = ref([])
const filterOffice = ref(null)
const filterConsumable = ref(null)
const showDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = ref({ office_id: null, consumable_id: null, quantity: 1, usage_date: new Date(), note: '' })

onMounted(async () => {
  offices.value = await getOffices()
  consumables.value = await getConsumables()
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const params = {}
    if (dateRange.value && dateRange.value[0]) {
      params.start_date = dateRange.value[0].toISOString().slice(0,10)
      params.end_date = dateRange.value[1].toISOString().slice(0,10)
    }
    if (filterOffice.value) params.office_id = filterOffice.value
    if (filterConsumable.value) params.consumable_id = filterConsumable.value
    items.value = await getRecords(params)
  } finally { loading.value = false }
}

function editItem(row) {
  if (!auth.isAdmin) return
  form.value = { ...row, usage_date: new Date(row.usage_date) }
  editing.value = true
  showDialog.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const payload = { ...form.value, usage_date: form.value.usage_date.toISOString().slice(0,10) }
    if (editing.value) {
      await updateRecord(form.value.id, payload)
    } else {
      await createRecord(payload)
    }
    showDialog.value = false
    ElMessage.success('保存成功')
    loadData()
  } finally { saving.value = false }
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定删除该记录？')
    await deleteRecord(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}

function exportExcel() {
  const ws = XLSX.utils.json_to_sheet(items.value.map(r => ({
    日期: r.usage_date, 房间号: r.room_number, 耗材: r.consumable_name, 数量: r.quantity, 单位: r.unit, 备注: r.note
  })))
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '使用记录')
  XLSX.writeFile(wb, `使用记录_${new Date().toISOString().slice(0,10)}.xlsx`)
}
</script>
