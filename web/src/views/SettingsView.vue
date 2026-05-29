<template>
  <el-card>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="耗材管理" name="consumables">
        <el-button type="primary" size="small" style="margin-bottom:8px" @click="addConsumable">新增</el-button>
        <el-table :data="consumables" stripe>
          <el-table-column prop="name" label="名称" width="300" show-overflow-tooltip />
          <el-table-column prop="unit" label="单位" width="80" />
          <el-table-column label="操作" width="140">
            <template #default="{ row }">
              <el-button size="small" link @click="editConsumable(row)">编辑</el-button>
              <el-button v-if="!row.is_default" size="small" type="danger" link @click="delConsumable(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="办公室管理" name="offices">
        <el-button type="primary" size="small" style="margin-bottom:8px" @click="addOffice">新增</el-button>
        <el-table :data="offices" stripe>
          <el-table-column prop="room_number" label="房间号" width="120" />
          <el-table-column prop="device_type" label="设备类型" width="120" />
          <el-table-column prop="device_model" label="设备型号" show-overflow-tooltip />
          <el-table-column label="操作" width="140">
            <template #default="{ row }">
              <el-button size="small" link @click="editOffice(row)">编辑</el-button>
              <el-button size="small" type="danger" link @click="delOffice(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="用户管理" name="users">
        <el-button type="primary" size="small" style="margin-bottom:8px" @click="addUser">新增</el-button>
        <el-table :data="users" stripe>
          <el-table-column prop="username" label="用户名" width="150" />
          <el-table-column prop="role" label="角色" width="100">
            <template #default="{ row }">{{ row.role === 'admin' ? '管理员' : '只读' }}</template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="160" />
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button size="small" link @click="editUser(row)">编辑</el-button>
              <el-button size="small" type="danger" link @click="delUser(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="故障类型" name="problemTypes">
        <el-button type="primary" size="small" style="margin-bottom:8px" @click="addProblemType">新增</el-button>
        <el-table :data="problemTypes" stripe>
          <el-table-column prop="device_type" label="设备类型" width="120" />
          <el-table-column prop="name" label="故障名称" width="200" />
          <el-table-column prop="sort_order" label="排序" width="70" />
          <el-table-column label="操作" width="140">
            <template #default="{ row }">
              <el-button size="small" link @click="editProblemType(row)">编辑</el-button>
              <el-button v-if="!row.is_default" size="small" type="danger" link @click="delProblemType(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </el-card>

  <el-dialog v-model="dialog" :title="dialogTitle" width="450px">
    <el-form :model="dialogForm" label-width="80px">
      <el-form-item v-if="dialogType==='consumables'" label="名称"><el-input v-model="dialogForm.name" /></el-form-item>
      <el-form-item v-if="dialogType==='consumables'" label="单位"><el-input v-model="dialogForm.unit" /></el-form-item>
      <el-form-item v-if="dialogType==='offices'" label="房间号"><el-input v-model="dialogForm.room_number" /></el-form-item>
      <el-form-item v-if="dialogType==='offices'" label="设备类型"><el-input v-model="dialogForm.device_type" /></el-form-item>
      <el-form-item v-if="dialogType==='offices'" label="设备型号"><el-input v-model="dialogForm.device_model" /></el-form-item>
      <el-form-item v-if="dialogType==='users'" label="用户名"><el-input v-model="dialogForm.username" /></el-form-item>
      <el-form-item v-if="dialogType==='users'" label="密码"><el-input v-model="dialogForm.password" type="password" show-password /></el-form-item>
      <el-form-item v-if="dialogType==='users'" label="角色"><el-select v-model="dialogForm.role"><el-option label="管理员" value="admin" /><el-option label="只读" value="readonly" /></el-select></el-form-item>
      <el-form-item v-if="dialogType==='problemTypes'" label="设备类型"><el-input v-model="dialogForm.device_type" /></el-form-item>
      <el-form-item v-if="dialogType==='problemTypes'" label="故障名称"><el-input v-model="dialogForm.name" /></el-form-item>
      <el-form-item v-if="dialogType==='problemTypes'" label="排序"><el-input-number v-model="dialogForm.sort_order" :min="0" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog=false">取消</el-button>
      <el-button type="primary" @click="saveDialog">{{ dialogEdit ? '更新' : '创建' }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listConsumables, createConsumable, updateConsumable, deleteConsumable } from '@/api/consumables'
import { listOffices, createOffice, updateOffice, deleteOffice } from '@/api/offices'
import { listUsers } from '@/api/users'
import * as userApi from '@/api/users'
import { listProblemTypes, createProblemType, updateProblemType, deleteProblemType } from '@/api/problemTypes'

const activeTab = ref('consumables')
const consumables = ref<any[]>([])
const offices = ref<any[]>([])
const users = ref<any[]>([])
const problemTypes = ref<any[]>([])

const dialog = ref(false)
const dialogType = ref('')
const dialogEdit = ref(false)
const dialogTitle = ref('')
const dialogForm = reactive<any>({})

async function loadAll() {
  const [c, o, u, p] = await Promise.all([
    listConsumables(), listOffices(),
    listUsers ? listUsers() : Promise.resolve({ data: [] }),
    listProblemTypes(),
  ])
  consumables.value = c.data
  offices.value = o.data
  users.value = (u as any)?.data || []
  problemTypes.value = p.data
}

function openDialog(type: string, edit: boolean, data?: any) {
  dialogType.value = type
  dialogEdit.value = edit
  const titles: any = { consumables: '耗材', offices: '办公室', users: '用户', problemTypes: '故障类型' }
  dialogTitle.value = (edit ? '编辑' : '新增') + titles[type]
  dialogForm.value = edit ? { ...data } : {}
  if (!edit) {
    if (type === 'consumables') dialogForm.unit = '个'
    if (type === 'problemTypes') dialogForm.sort_order = 0
  }
  dialog.value = true
}

function addConsumable() { openDialog('consumables', false) }
function editConsumable(row: any) { openDialog('consumables', true, row) }
function addOffice() { openDialog('offices', false) }
function editOffice(row: any) { openDialog('offices', true, row) }
function addUser() { openDialog('users', false) }
function editUser(row: any) { openDialog('users', true, row) }
function addProblemType() { openDialog('problemTypes', false) }
function editProblemType(row: any) { openDialog('problemTypes', true, row) }

async function saveDialog() {
  const type = dialogType.value
  const isEdit = dialogEdit.value
  let fn: any
  if (type === 'consumables') fn = isEdit ? updateConsumable(dialogForm.id, dialogForm) : createConsumable(dialogForm)
  else if (type === 'offices') fn = isEdit ? updateOffice(dialogForm.id, dialogForm) : createOffice(dialogForm)
  else if (type === 'users') fn = isEdit ? userApi.updateUser(dialogForm.id, dialogForm) : userApi.createUser(dialogForm)
  else if (type === 'problemTypes') fn = isEdit ? updateProblemType(dialogForm.id, dialogForm) : createProblemType(dialogForm)
  await fn
  ElMessage.success(isEdit ? '已更新' : '已创建')
  dialog.value = false
  loadAll()
}

async function delConsumable(id: number) { await ElMessageBox.confirm('确认删除？'); await deleteConsumable(id); ElMessage.success('已删除'); loadAll() }
async function delOffice(id: number) { await ElMessageBox.confirm('确认删除？'); await deleteOffice(id); ElMessage.success('已删除'); loadAll() }
async function delUser(id: number) { await ElMessageBox.confirm('确认删除？'); await userApi.deleteUser(id); ElMessage.success('已删除'); loadAll() }
async function delProblemType(id: number) { await ElMessageBox.confirm('确认删除？'); await deleteProblemType(id); ElMessage.success('已删除'); loadAll() }

onMounted(loadAll)
</script>
