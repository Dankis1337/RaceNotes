import api from './axios'

export const calculatePressure = (data) => api.post('/calculator/tire-pressure', data)
