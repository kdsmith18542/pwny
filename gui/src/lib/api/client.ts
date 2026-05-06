import axios from "axios"

const API_BASE_URL = "http://127.0.0.1:31337/api/v1"

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
})

export interface APIResponse<T> {
  status: string
  data: T
  error?: {
    error: string
    detail?: string
  }
}

export const modulesApi = {
  list: async (type?: string) => {
    const resp = await apiClient.get<APIResponse<any[]>>("/modules", { params: { type } })
    return resp.data.data
  },
  get: async (name: string) => {
    const resp = await apiClient.get<APIResponse<any>>(`/modules/${name}`)
    return resp.data.data
  },
  validate: async (name: string, options: Record<string, any>) => {
    const resp = await apiClient.post<APIResponse<any>>(`/modules/${name}/validate`, { options })
    return resp.data.data
  },
  run: async (name: string, options: Record<string, any>) => {
    const resp = await apiClient.post<APIResponse<{ job_id: string }>>(`/modules/${name}/run`, { options })
    return resp.data.data
  },
}

export const sessionsApi = {
  list: async () => {
    const resp = await apiClient.get<APIResponse<any[]>>("/sessions")
    return resp.data.data
  },
  get: async (id: string) => {
    const resp = await apiClient.get<APIResponse<any>>(`/sessions/${id}`)
    return resp.data.data
  },
  close: async (id: string) => {
    const resp = await apiClient.delete<APIResponse<any>>(`/sessions/${id}`)
    return resp.data.data
  },
}

export const payloadApi = {
  generate: async (payload: any) => {
    const resp = await apiClient.post<APIResponse<{ payload: string; size: number }>>("/payload/generate", payload)
    return resp.data.data
  },
}

export const workspaceApi = {
  list: async () => {
    const resp = await apiClient.get<APIResponse<any[]>>("/workspaces")
    return resp.data.data
  },
  create: async (name: string, description?: string) => {
    const resp = await apiClient.post<APIResponse<any>>("/workspaces", { name, description })
    return resp.data.data
  },
  getHosts: async (id: string) => {
    const resp = await apiClient.get<APIResponse<any[]>>(`/workspaces/${id}/hosts`)
    return resp.data.data
  },
}
