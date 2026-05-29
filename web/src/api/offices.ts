import api from './index'

export const listOffices = () => api.get('/offices')
export const createOffice = (data: { room_number: string; device_type: string; device_model: string }) => api.post('/offices', data)
export const updateOffice = (id: number, data: { room_number: string; device_type: string; device_model: string }) => api.put(`/offices/${id}`, data)
export const deleteOffice = (id: number) => api.delete(`/offices/${id}`)
export const getDeviceModels = () => api.get('/device-models')
