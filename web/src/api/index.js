import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    const msg = err.response?.data?.error || err.message
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export function login(username, password) {
  return api.post('/auth/login', { username, password }).then(r => r.data)
}

export function changePassword(password, newPassword) {
  return api.put('/auth/password', { password, new_password: newPassword }).then(r => r.data)
}

export function getUsers() { return api.get('/users').then(r => r.data) }
export function createUser(data) { return api.post('/users', data).then(r => r.data) }
export function updateUser(id, data) { return api.put(`/users/${id}`, data).then(r => r.data) }
export function deleteUser(id) { return api.delete(`/users/${id}`).then(r => r.data) }

export function getConsumables(params) { return api.get('/consumables', { params }).then(r => r.data) }
export function createConsumable(data) { return api.post('/consumables', data).then(r => r.data) }
export function updateConsumable(id, data) { return api.put(`/consumables/${id}`, data).then(r => r.data) }
export function deleteConsumable(id) { return api.delete(`/consumables/${id}`).then(r => r.data) }

export function getOffices() { return api.get('/offices').then(r => r.data) }
export function createOffice(data) { return api.post('/offices', data).then(r => r.data) }
export function updateOffice(id, data) { return api.put(`/offices/${id}`, data).then(r => r.data) }
export function deleteOffice(id) { return api.delete(`/offices/${id}`).then(r => r.data) }

export function getRecords(params) { return api.get('/records', { params }).then(r => r.data) }
export function createRecord(data) { return api.post('/records', data).then(r => r.data) }
export function updateRecord(id, data) { return api.put(`/records/${id}`, data).then(r => r.data) }
export function deleteRecord(id) { return api.delete(`/records/${id}`).then(r => r.data) }

export function getStats(params) { return api.get('/stats/summary', { params }).then(r => r.data) }

export function getNotes(params) { return api.get('/notes', { params }).then(r => r.data) }
export function createNote(data) { return api.post('/notes', data).then(r => r.data) }
export function updateNote(id, data) { return api.put(`/notes/${id}`, data).then(r => r.data) }
export function deleteNote(id) { return api.delete(`/notes/${id}`).then(r => r.data) }

export function getProblemTypes() { return api.get('/problem-types').then(r => r.data) }
export function createProblemType(data) { return api.post('/problem-types', data).then(r => r.data) }
export function updateProblemType(id, data) { return api.put(`/problem-types/${id}`, data).then(r => r.data) }
export function deleteProblemType(id) { return api.delete(`/problem-types/${id}`).then(r => r.data) }
export function sortProblemTypes(data) { return api.put('/problem-types/sort', data).then(r => r.data) }

export function getQRCodes() { return api.get('/qrcodes').then(r => r.data) }
export function createQRCode(data) { return api.post('/qrcodes', data).then(r => r.data) }
export function updateQRCode(id, data) { return api.put(`/qrcodes/${id}`, data).then(r => r.data) }
export function deleteQRCode(id) { return api.delete(`/qrcodes/${id}`).then(r => r.data) }
export function getQRCodeImage(code) { return api.get(`/qrcodes/${code}/image`).then(r => r.data) }
export function getDeviceModels() { return api.get('/device-models').then(r => r.data) }

export function getBackupConfig() { return api.get('/backup/config').then(r => r.data) }
export function saveBackupConfig(data) { return api.put('/backup/config', data).then(r => r.data) }
export function getBackupList() { return api.get('/backup/list').then(r => r.data) }
export function createBackup() { return api.post('/backup/now').then(r => r.data) }
export function downloadBackup(name) { return api.get(`/backup/download/${name}`, { responseType: 'blob' }).then(r => r.data) }
export function restoreBackup(name) { return api.post(`/backup/restore/${name}`).then(r => r.data) }
export function deleteBackup(name) { return api.delete(`/backup/${name}`).then(r => r.data) }

export function getTickets(params) { return api.get('/tickets', { params }).then(r => r.data) }
export function updateTicketStatus(id, status) { return api.post(`/tickets/${id}/status`, { status }).then(r => r.data) }
export function completeTicket(id, data) { return api.post(`/tickets/${id}/complete`, data).then(r => r.data) }
export function deleteTicket(id) { return api.delete(`/tickets/${id}`).then(r => r.data) }

export default api
