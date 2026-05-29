import api from './index'

export interface RecordQuery {
  office_id?: string
  consumable_id?: string
  date_from?: string
  date_to?: string
  page?: number
  page_size?: number
}

export const listRecords = (params: RecordQuery) => api.get('/records', { params })
export const createRecord = (data: any) => api.post('/records', data)
export const updateRecord = (id: number, data: any) => api.put(`/records/${id}`, data)
export const deleteRecord = (id: number) => api.delete(`/records/${id}`)
export const exportRecords = (params: RecordQuery) => api.get('/records/export', { params, responseType: 'blob' })
