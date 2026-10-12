import axios from 'axios'

const API_URL = '/api'

// Custom error class for API errors
export class ApiError extends Error {
  status: number
  details?: string

  constructor(message: string, status: number, details?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.details = details
  }
}

export interface Node {
  id: string
  name: string
  slug: string
  description: string
  rules: string
  host_id: string
  created_at: string
  updated_at: string
}

export interface Post {
  id: string
  node_id: string
  author_id: string
  title: string
  body: string
  url?: string
  vote_count: number
  author: {
    id: string
    username: string
    email: string
    role: string
  }
  created_at: string
  updated_at: string
}

export interface InviteCode {
  id: string
  code: string
  max_uses: number
  uses: number
  expires_at?: string
  role: string
  created_at: string
}

class ApiService {
  private getApi() {
    const token = localStorage.getItem('token')
    return axios.create({
      baseURL: API_URL,
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
  }

  private async request<T>(fn: () => Promise<any>): Promise<T> {
    try {
      const response = await fn()
      return response.data
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        const status = error.response.status
        const message = error.response.data?.error || 'Request failed'
        throw new ApiError(message, status, error.response.data?.details)
      }
      throw new ApiError('Network error', 500)
    }
  }

  // Nodes
  async getNodes(): Promise<Node[]> {
    return this.request<Node[]>(() => this.getApi().get('/nodes'))
  }

  async getNode(slug: string): Promise<Node> {
    return this.request<Node>(() => this.getApi().get(`/nodes/${slug}`))
  }

  async createNode(data: { name: string; description?: string; rules?: string }): Promise<Node> {
    return this.request<Node>(() => this.getApi().post('/nodes', data))
  }

  // Posts
  async getPosts(slug: string, page = 1, limit = 20): Promise<Post[]> {
    return this.request<Post[]>(() => this.getApi().get(`/nodes/${slug}/posts?page=${page}&limit=${limit}`))
  }

  async getPost(id: string): Promise<Post> {
    return this.request<Post>(() => this.getApi().get(`/posts/${id}`))
  }

  async createPost(slug: string, data: { title: string; body?: string; url?: string; tags?: string[] }): Promise<Post> {
    return this.request<Post>(() => this.getApi().post(`/nodes/${slug}/posts`, data))
  }

  async votePost(id: string, vote: 'up' | 'down'): Promise<void> {
    return this.request<void>(() => this.getApi().post(`/posts/${id}/vote`, { vote }))
  }

  // Invites
  async createInvite(data: { max_uses?: number; expires?: string; role?: string }): Promise<InviteCode> {
    return this.request<InviteCode>(() => this.getApi().post('/invites', data))
  }

  async getInvites(): Promise<InviteCode[]> {
    return this.request<InviteCode[]>(() => this.getApi().get('/invites'))
  }

  async revokeInvite(id: string): Promise<void> {
    return this.request<void>(() => this.getApi().delete(`/invites/${id}`))
  }
}

export const apiService = new ApiService()
