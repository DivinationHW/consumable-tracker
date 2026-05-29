import api from './index'

export const listConsumables = () => api.get('/consumables')
export const createConsumable = (data: { name: string; unit: string }) => api.post('/consumables', data)
export const updateConsumable = (id: number, data: { name: string; unit: string }) => api.put(`/consumables/${id}`, data)
export const deleteConsumable = (id: number) => api.delete(`/consumables/${id}`)
