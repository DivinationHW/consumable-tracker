<template>
  <div>
    <el-card style="margin-bottom:16px">
      <template #header><span>快速录入</span></template>
      <el-form :inline="true" :model="form" label-width="auto">
        <el-form-item label="日期">
          <el-date-picker v-model="form.date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="办公室">
          <el-select v-model="form.office_id" filterable placeholder="选择办公室" style="width:180px">
            <el-option v-for="o in offices" :key="o.id" :label="o.room_number" :value="o.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="耗材名称">
          <el-select v-model="form.consumable_id" filterable placeholder="选择耗材" style="width:300px">
            <el-option v-for="c in consumables" :key="c.id" :label="c.name" :value="c.id">
              <span>{{ c.name }}</span>
              <span style="color:#999;margin-left:8px;font-size:12px">{{ c.unit }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="使用了">
          <el-input-number v-model="form.quantity" :min="1" :max="999" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.note" placeholder="可选" style="width:150px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit">录入</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>最近记录</span>
          <el-button size="small" @click="$router.push('/records')">查看全部</el-button>
        </div>
      </template>
      <el-table :data="recent" stripe v-loading="loading" max-height="500">
        <el-table-column prop="usage_date" label="日期" width="110" />
        <el-table-column prop="office_name" label="办公室" width="120" />
        <el-table-column prop="device_type" label="设备类型" width="100" />
        <el-table-column prop="device_model" label="型号" width="160" />
        <el-table-column prop="consumable_name" label="耗材" width="200" show-overflow-tooltip />
        <el-table-column prop="quantity" label="数量" width="70" />
        <el-table-column prop="note" label="备注" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listOffices } from '@/api/offices'
import { listConsumables } from '@/api/consumables'
import { listRecords, createRecord } from '@/api/records'

const loading = ref(false)
const offices = ref<any[]>([])
const consumables = ref<any[]>([])
const recent = ref<any[]>([])
const form = reactive({
  date: new Date().toISOString().slice(0, 10),
  office_id: undefined as number | undefined,
  consumable_id: undefined as number | undefined,
  quantity: 1,
  note: '',
})

async function load() {
  const [o, c, r] = await Promise.all([
    listOffices(),
    listConsumables(),
    listRecords({ page: 1, page_size: 20 }),
  ])
  offices.value = o.data
  consumables.value = c.data
  recent.value = r.data.data
}

async function submit() {
  if (!form.office_id || !form.consumable_id) {
    ElMessage.warning('请选择办公室和耗材')
    return
  }
  await createRecord(form)
  ElMessage.success('录入成功')
  form.quantity = 1
  form.note = ''
  const r = await listRecords({ page: 1, page_size: 20 })
  recent.value = r.data.data
}

onMounted(load)
</script>
