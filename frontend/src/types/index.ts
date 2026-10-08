export interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'moderator' | 'member'
  created_at: string
}

export interface Node {
  id: string
  slug: string
  name: string
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
  author: User
  created_at: string
  updated_at: string
}

export interface Comment {
  id: string
  post_id: string
  author_id: string
  parent_id?: string
  body: string
  vote_count: number
  author: User
  replies?: Comment[]
  created_at: string
  updated_at: string
}

export interface InviteCode {
  id: string
  code: string
  max_uses: number
  uses: number
  expires_at?: string
  role: 'member' | 'moderator'
  created_at: string
}

export interface Tag {
  id: string
  name: string
}

export interface Vote {
  id: string
  user_id: string
  post_id?: string
  comment_id?: string
  vote_type: 'up' | 'down'
  created_at: string
}

export interface ApiResponse<T> {
  data: T
  error?: string
}

export interface PaginationParams {
  page?: number
  limit?: number
  sort?: 'new' | 'hot' | 'top'
}
