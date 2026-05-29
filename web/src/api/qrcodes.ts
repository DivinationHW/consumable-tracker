import api from './index'

export const listQRCodes = () => api.get('/qrcodes')
export const createQRCode = (data: any) => api.post('/qrcodes', data)
export const generateBulk = (data: any) => api.post('/qrcodes/bulk', data)
export const deleteQRCode = (id: number) => api.delete(`/qrcodes/${id}`)
export const getQRImage = (code: string) => `/api/qrcodes/${code}/image`
export const getQRImageBase64 = (code: string) => api.get(`/qrcodes/${code}/image-base64`)
export const getQRPrintPage = () => '/api/qrcodes/print'
