
import axios from "axios"

const api = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_URL
})

api.interceptors.request.use((config) => {
  //   const credentials = <IAuthCredentials>JSON.parse(sessionStorage.getItem("credentials")!)
  //
  // if (credentials !== null && credentials !== undefined) {
  //   config.headers.Authorization = `${credentials.token}`
  // }

  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response !== undefined && error.response.status === 401) {
      // sessionStorage.clear()
    }

    return Promise.reject(error)
  }
)

export default api
