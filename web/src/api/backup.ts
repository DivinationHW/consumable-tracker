import api from './index'

export const listBackups = () => api.get('/backups')
export const createBackup = () => api.post('/backups')
export const restoreBackup = (filename: string) => api.post(`/backups/${encodeURIComponent(filename)}/restore`)
export const deleteBackup = (filename: string) => api.delete(`/backups/${encodeURIComponent(filename)}`)
