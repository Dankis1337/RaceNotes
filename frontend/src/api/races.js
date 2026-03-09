import api from './axios'

export const getRaces = (params) => api.get('/races', { params })
export const getRace = (id) => api.get(`/races/${id}`)
export const createRace = (data) => api.post('/races', data)
export const updateRace = (id, data) => api.put(`/races/${id}`, data)
export const deleteRace = (id) => api.delete(`/races/${id}`)
