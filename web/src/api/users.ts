import api from './index'

export const listUsers = () => api.get('/users')
export const createUser = (data: { username: string; password: string; role: string }) => api.post('/users', data)
export const updateUser = (id: number, data: { username?: string; password?: string; role?: string }) => api.put(`/users/${id}`, data)
export const deleteUser = (id: number) => api.delete(`/users/${id}`)
