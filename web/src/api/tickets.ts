import api from './index'

export const listTickets = (params?: any) => api.get('/tickets', { params })
export const createPublicTicket = (data: any, qr?: string) => api.post('/public/ticket', data, { params: { qr } })
export const getTicket = (id: string) => api.get(`/tickets/${id}`)
export const getPublicTicket = (id: string) => api.get(`/public/ticket/${id}`)
export const processTicket = (id: string, data: any) => api.post(`/tickets/${id}/process`, data)
