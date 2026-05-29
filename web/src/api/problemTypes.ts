import api from './index'

export const listProblemTypes = (deviceType?: string) => api.get('/problem-types', { params: { device_type: deviceType } })
export const createProblemType = (data: any) => api.post('/problem-types', data)
export const updateProblemType = (id: number, data: any) => api.put(`/problem-types/${id}`, data)
export const deleteProblemType = (id: number) => api.delete(`/problem-types/${id}`)
export const getDeviceTypes = () => api.get('/device-types')
