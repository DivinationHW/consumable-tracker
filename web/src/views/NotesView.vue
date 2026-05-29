<template>
  <div>
    <el-button type="primary" style="margin-bottom:16px" @click="add">新建便签</el-button>
    <el-row :gutter="16">
      <el-col v-for="n in notes" :key="n.id" :span="8" style="margin-bottom:16px">
        <el-card>
          <template #header>
            <div style="display:flex;justify-content:space-between;align-items:center">
              <span style="font-weight:bold">{{ n.title }}</span>
              <div>
                <el-button size="small" link @click="edit(n)">编辑</el-button>
                <el-button size="small" type="danger" link @click="doDelete(n.id)">删除</el-button>
              </div>
            </div>
          </template>
          <div style="white-space:pre-wrap;font-size:13px">{{ n.content }}</div>
          <div style="margin-top:8px;font-size:12px;color:#999">{{ n.updated_at || n.created_at }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-empty v-if="!notes.length" description="暂无便签" />

    <el-dialog v-model="dialog" :title="editing ? '编辑便签' : '新建便签'" width="500px">
      <el-form :model="form" label-width="60px">
        <el-form-item label="标题"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="内容"><el-input v-model="form.content" type="textarea" :rows="6" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog=false">取消</el-button>
        <el-button type="primary" @click="save">{{ editing ? '更新' : '创建' }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listNotes, createNote, updateNote, deleteNote } from '@/api/notes'

const notes = ref<any[]>([])
const dialog = ref(false)
const editing = ref(false)
const form = ref({ title: '', content: '' })

async function load() { notes.value = (await listNotes()).data }
function add() { editing.value = false; form.value = { title: '', content: '' }; dialog.value = true }
function edit(n: any) { editing.value = true; form.value = { title: n.title, content: n.content }; dialog.value = true; form.value._id = n.id }
async function save() {
  if (editing.value) { await updateNote(form.value._id, form.value) } else { await createNote(form.value) }
  ElMessage.success(editing.value ? '已更新' : '已创建')
  dialog.value = false; load()
}
async function doDelete(id: number) {
  await ElMessageBox.confirm('确认删除？')
  await deleteNote(id); ElMessage.success('已删除'); load()
}
onMounted(load)
</script>
