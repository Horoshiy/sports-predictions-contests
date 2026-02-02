// gRPC-Web client configuration

class GrpcClient {
  private baseUrl: string

  constructor() {
    // Use Vite proxy in development, environment variable in production
    this.baseUrl = import.meta.env.VITE_API_URL || ''
  }

  getBaseUrl(): string {
    return this.baseUrl
  }

  // Common headers for all requests
  getHeaders(isFormData = false): Record<string, string> {
    const headers: Record<string, string> = {}

    // Don't set Content-Type for FormData, let browser set it with boundary
    if (!isFormData) {
      headers['Content-Type'] = 'application/json'
    }

    // Add auth token if available
    const token = localStorage.getItem('authToken')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    return headers
  }

  // Handle HTTP requests to gRPC-Gateway
  async request<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    path: string,
    body?: any,
    isFormData = false
  ): Promise<T> {
    const url = `${this.baseUrl}${path}`
    
    const config: RequestInit = {
      method,
      headers: this.getHeaders(isFormData),
    }

    if (body && method !== 'GET') {
      if (isFormData) {
        config.body = body // FormData object
      } else {
        config.body = JSON.stringify(body)
      }
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: 'Request failed',
          code: response.status,
          message: response.statusText,
        }))
        throw new Error(errorData.message || `HTTP ${response.status}`)
      }

      return await response.json()
    } catch (error) {
      console.error('gRPC request failed:', error)
      throw error
    }
  }

  // GET request
  async get<T>(path: string): Promise<T> {
    return this.request<T>('GET', path)
  }

  // POST request
  async post<T>(path: string, body?: any): Promise<T> {
    return this.request<T>('POST', path, body)
  }

  // PUT request
  async put<T>(path: string, body?: any): Promise<T> {
    return this.request<T>('PUT', path, body)
  }

  // DELETE request
  async delete<T>(path: string): Promise<T> {
    return this.request<T>('DELETE', path)
  }

  // POST FormData request (for file uploads)
  async postFormData<T>(path: string, formData: FormData): Promise<T> {
    return this.request<T>('POST', path, formData, true)
  }
}

// Singleton instance
export const grpcClient = new GrpcClient()
export default grpcClient
