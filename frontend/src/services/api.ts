import axios from 'axios'

const API_URL = '/api'

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

  // Nodes
  async getNodes(): Promise<Node[]> {
    const api = this.getApi()
    const response = await api.get('/nodes')
    return response.data
  }

  async getNode(slug: string): Promise<Node> {
    const api = this.getApi()
    const response = await api.get(`/nodes/${slug}`)
    return response.data
  }

  async createNode(data: { name: string; description?: string; rules?: string }): Promise<Node> {
    const api = this.getApi()
    const response = await api.post('/nodes', data)
    return response.data
  }

  // Posts
  async getPosts(slug: string, page = 1, limit = 20): Promise<Post[]> {
    const api = this.getApi()
    const response = await api.get(`/nodes/${slug}/posts?page=${page}&limit=${limit}`)
    return response.data
  }

  async getPost(id: string): Promise<Post> {
    const api = this.getApi()
    const response = await api.get(`/posts/${id}`)
    return response.data
  }

  async createPost(slug: string, data: { title: string; body?: string; url?: string; tags?: string[] }): Promise<Post> {
    const api = this.getApi()
    const response = await api.post(`/nodes/${slug}/posts`, data)
    return response.data
  }

  async votePost(id: string, vote: 'up' | 'down'): Promise<void> {
    const api = this.getApi()
    await api.post(`/posts/${id}/vote`, { vote })
  }

  // Invites
  async createInvite(data: { max_uses?: number; expires?: string; role?: string }): Promise<InviteCode> {
    const api = this.getApi()
    const response = await api.post('/invites', data)
    return response.data
  }

  async getInvites(): Promise<InviteCode[]> {
    const api = this.getApi()
    const response = await api.get('/invites')
    return response.data
  }

  async revokeInvite(id: string): Promise<void> {
    const api = this.getApi()
    await api.delete(`/invites/${id}`)
  }
}

export const apiService = new ApiService()
