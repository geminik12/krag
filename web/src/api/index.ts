import axios, { InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'

const api = axios.create({
  timeout: 10000
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use((response: AxiosResponse) => response, (error: AxiosError) => {
  if (error.response && error.response.status === 401) {
    localStorage.removeItem('token')
    window.location.href = '/login'
  }
  return Promise.reject(error)
})

export default api
