import axios from 'axios'

const API_URL = '/api'

interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'moderator' | 'member'
  created_at: string
}

interface AuthResponse {
  token: string
  user: User
}

interface LoginRequest {
  username: string
  password: string
}

interface RegisterRequest {
  username: string
  email: string
  password: string
  code: string
}

class AuthService {
  getToken(): string | null {
    return localStorage.getItem('token')
  }

  setToken(token: string): void {
    localStorage.setItem('token', token)
  }

  removeToken(): void {
    localStorage.removeItem('token')
  }

  getApi() {
    const token = this.getToken()
    return axios.create({
      baseURL: API_URL,
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await axios.post(`${API_URL}/auth/login`, data)
    this.setToken(response.data.token)
    return response.data
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await axios.post(`${API_URL}/auth/register`, data)
    this.setToken(response.data.token)
    return response.data
  }

  async getCurrentUser(): Promise<User> {
    const api = this.getApi()
    const response = await api.get('/auth/me')
    return response.data
  }

  logout(): void {
    this.removeToken()
  }

  isAuthenticated(): boolean {
    return !!this.getToken()
  }
}

export const authService = new AuthService()
