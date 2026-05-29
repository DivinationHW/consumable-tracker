import api from './index'

export const listNotes = () => api.get('/notes')
export const createNote = (data: any) => api.post('/notes', data)
export const updateNote = (id: number, data: any) => api.put(`/notes/${id}`, data)
export const deleteNote = (id: number) => api.delete(`/notes/${id}`)
