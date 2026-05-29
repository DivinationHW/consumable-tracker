import api from './index'

export const login = (username: string, password: string) =>
  api.post('/login', { username, password })

export const changePassword = (oldPassword: string, newPassword: string) =>
  api.post('/change-password', { old_password: oldPassword, new_password: newPassword })

export const getMe = () => api.get('/me')
